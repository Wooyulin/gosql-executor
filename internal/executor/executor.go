package executor

import (
	"fmt"
	"io"
	"sql-executor/internal/database"
	"sql-executor/internal/output"
	"sql-executor/pkg/logger"
	"strings"
	"time"

	"github.com/peterh/liner"
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
	state := liner.NewLiner()
	defer state.Close()

	state.SetCtrlCAborts(false) // Ctrl+C 仅取消当前输入，不退出程序

	for {
		var query string
		for {
			prompt := "SQL> "
			if query != "" {
				prompt = "  -> "
			}
			line, err := state.Prompt(prompt)
			if err != nil {
				if err == io.EOF {
					return nil
				}
				e.logger.Error("Failed to read input", err)
				query = ""
				break
			}

			if line == "exit" {
				return nil
			}

			query += line + " "
			query = strings.TrimSpace(query)
			if query == "" {
				continue
			}
			if line == "" || strings.HasSuffix(strings.TrimSpace(line), ";") {
				break
			}
		}

		if query = strings.TrimSpace(query); query == "" {
			continue
		}

		query = processQuery(query, e.db)

		start := time.Now()
		rows, err := e.db.Execute(query)
		if err != nil {
			e.logger.Error("Failed to execute SQL", err)
			continue
		}
		defer rows.Close()

		filename := fmt.Sprintf("query_result_%d", time.Now().Unix())
		if err := e.writer.Write(filename, rows); err != nil {
			e.logger.Error("Failed to write results", err)
			continue
		}

		duration := time.Since(start)
		if e.writer.(*output.FileWriter).Config().SaveToFile {
			e.logger.Info(fmt.Sprintf("Query executed in %v, saved to: %s", duration, filename))
		} else {
			e.logger.Info(fmt.Sprintf("Query executed in %v", duration))
		}

		// 将执行过的 SQL 加入历史，方便上下键回顾
		state.AppendHistory(query)
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
