package executor

import (
	"bufio"
	"fmt"
	"os"
	"sql-executor/internal/database"
	"sql-executor/internal/output"
	"sql-executor/pkg/logger"
	"strings"
	"time"
)

type SQLExecutor struct {
	db     database.Database
	writer output.Writer
	logger *logger.Logger
}

func NewSQLExecutor(db database.Database, writer output.Writer, logger *logger.Logger) *SQLExecutor {
	return &SQLExecutor{
		db:     db,
		writer: writer,
		logger: logger,
	}
}

func (e *SQLExecutor) Run() error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("SQL> ")
		var query string
		for scanner.Scan() {
			line := scanner.Text()
			if line == "exit" {
				return nil
			}

			query += line + " "
			if line == "" || line[len(line)-1] == ';' {
				break
			}
			fmt.Print("  -> ")
		}

		if query = strings.TrimSpace(query); query == "" {
			continue
		}

		query = processQuery(query, e.db)

		start := time.Now()
		rows, err := e.db.Execute(query)
		if err != nil {
			e.logger.Error("执行SQL失败", err)
			continue
		}
		defer rows.Close()

		filename := fmt.Sprintf("query_result_%d", time.Now().Unix())
		if err := e.writer.Write(filename, rows); err != nil {
			e.logger.Error("写入结果失败", err)
			continue
		}

		duration := time.Since(start)
		if e.writer.(*output.FileWriter).Config().SaveToFile {
			e.logger.Info(fmt.Sprintf("查询执行完成，耗时: %v，结果已保存到: %s", duration, filename))
		} else {
			e.logger.Info(fmt.Sprintf("查询执行完成，耗时: %v", duration))
		}
		fmt.Println()
	}
}

func processQuery(query string, db database.Database) string {
	query = strings.TrimSpace(query)
	query = strings.TrimSuffix(query, ";")

	switch db.(type) {
	case *database.OracleDatabase:
		return query
	default:
		return query + ";"
	}
}
