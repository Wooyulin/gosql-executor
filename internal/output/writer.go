package output

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sql-executor/config"
	"text/tabwriter"
)

const (
	// 批量处理的行数
	batchSize = 1000
	// 最大文件大小 (1GB)
	maxFileSize = 1 << 30
)

type Writer interface {
	Write(filename string, rows *sql.Rows) error
}

type FileWriter struct {
	config *config.OutputConfig
}

func NewWriter(cfg *config.OutputConfig) Writer {
	return &FileWriter{
		config: cfg,
	}
}

func (w *FileWriter) Write(filename string, rows *sql.Rows) error {
	// 获取列名（在使用前先获取）
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("获取列名失败: %w", err)
	}

	// 准备存放数据的切片
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// 收集所有数据
	var allRows [][]interface{}
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("扫描行数据失败: %w", err)
		}

		// 复制当前行的值
		rowValues := make([]interface{}, len(values))
		for i, v := range values {
			rowValues[i] = v
		}
		allRows = append(allRows, rowValues)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("读取行数据失败: %w", err)
	}

	// 在命令行显示结果
	if w.config.ShowInConsole {
		if err := w.displayInConsole(columns, allRows); err != nil {
			return fmt.Errorf("显示结果失败: %w", err)
		}
	}

	// 如果不需要保存到文件，直接返回
	if !w.config.SaveToFile {
		return nil
	}

	// 写入文件
	if err := os.MkdirAll(w.config.Directory, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	ext := ".csv"
	if w.config.Format == "json" {
		ext = ".json"
	}

	fullPath := filepath.Join(w.config.Directory, filename+ext)

	if w.config.Format == "json" {
		return w.writeJSON(fullPath, columns, allRows)
	}
	return w.writeCSV(fullPath, columns, allRows)
}

// 修改显示函数以接受列和数据
func (w *FileWriter) displayInConsole(columns []string, rows [][]interface{}) error {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)

	// 写入表头
	for i, col := range columns {
		if i > 0 {
			fmt.Fprint(tw, "\t")
		}
		fmt.Fprint(tw, col)
	}
	fmt.Fprintln(tw)

	// 写入分隔线
	for i := range columns {
		if i > 0 {
			fmt.Fprint(tw, "\t")
		}
		fmt.Fprint(tw, "---")
	}
	fmt.Fprintln(tw)

	// 显示数据（限制显示前100行）
	maxRows := len(rows)
	if maxRows > 100 {
		maxRows = 100
	}

	for i := 0; i < maxRows; i++ {
		for j, v := range rows[i] {
			if j > 0 {
				fmt.Fprint(tw, "\t")
			}
			fmt.Fprint(tw, formatValue(v))
		}
		fmt.Fprintln(tw)
	}

	// 如果还有更多行
	if len(rows) > 100 {
		fmt.Fprintln(tw, "... 更多行已省略 ...")
	}

	tw.Flush()
	fmt.Println() // 添加一个空行
	return nil
}

func (w *FileWriter) writeCSV(filepath string, columns []string, rows [][]interface{}) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("创建CSV文件失败: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	if err := writer.Write(columns); err != nil {
		return fmt.Errorf("写入CSV表头失败: %w", err)
	}

	// 写入数据行
	for i, row := range rows {
		strRow := make([]string, len(row))
		for j, v := range row {
			strRow[j] = formatValue(v)
		}

		if err := writer.Write(strRow); err != nil {
			return fmt.Errorf("写入CSV数据行失败: %w", err)
		}

		// 每 batchSize 行刷新一次缓冲区
		if (i+1)%batchSize == 0 {
			writer.Flush()
			if fileInfo, err := file.Stat(); err == nil {
				if fileInfo.Size() > maxFileSize {
					return fmt.Errorf("文件大小超过限制(1GB)")
				}
			}
		}
	}

	return nil
}

func (w *FileWriter) writeJSON(filepath string, columns []string, rows [][]interface{}) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("创建JSON文件失败: %w", err)
	}
	defer file.Close()

	// 写入JSON数组开始标记
	if _, err := file.WriteString("[\n"); err != nil {
		return fmt.Errorf("写入JSON头部失败: %w", err)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("  ", "  ")

	// 写入数据
	for i, row := range rows {
		// 创建map存储行数据
		rowMap := make(map[string]interface{})
		for j, col := range columns {
			rowMap[col] = row[j]
		}

		// 添加逗号分隔符
		if i > 0 {
			if _, err := file.WriteString(",\n"); err != nil {
				return fmt.Errorf("写入JSON分隔符失败: %w", err)
			}
		}

		if err := encoder.Encode(rowMap); err != nil {
			return fmt.Errorf("写入JSON数据行失败: %w", err)
		}

		// 检查文件大小
		if (i+1)%batchSize == 0 {
			if fileInfo, err := file.Stat(); err == nil {
				if fileInfo.Size() > maxFileSize {
					return fmt.Errorf("文件大小超过限制(1GB)")
				}
			}
		}
	}

	// 写入JSON数组结束标记
	if _, err := file.WriteString("\n]"); err != nil {
		return fmt.Errorf("写入JSON尾部失败: %w", err)
	}

	return nil
}

// formatValue 格式化不同类型的值
func formatValue(v interface{}) string {
	if v == nil {
		return "NULL"
	}
	switch v := v.(type) {
	case []byte:
		return string(v)
	case string:
		if len(v) > 50 {
			return v[:47] + "..."
		}
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// 添加方法以获取配置
func (w *FileWriter) Config() *config.OutputConfig {
	return w.config
}
