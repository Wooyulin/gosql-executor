package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sql-executor/config"
	"sql-executor/internal/database"
	"sql-executor/internal/executor"
	"sql-executor/internal/output"
	"sql-executor/pkg/logger"
)

func main() {
	// 允许通过命令行参数指定配置文件路径
	configPath := flag.String("config", "./config.yaml", "config file path")
	flag.Parse()

	// 加载配置：文件存在则读取，否则进入交互式配置
	var cfg *config.Config
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		cfg, err = config.InteractiveSetup()
		if err != nil {
			log.Fatalf("Interactive setup failed: %v", err)
		}
		if err := config.SaveConfig(cfg, *configPath); err != nil {
			log.Fatalf("Failed to save config: %v", err)
		}
		fmt.Printf("\nConfig saved to %s\n", *configPath)
	} else {
		cfg, err = config.LoadConfig(*configPath)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
	}

	// 初始化日志
	logger := logger.NewLogger()

	// 初始化数据库连接
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to create database connection", err)
	}

	if err := db.Connect(); err != nil {
		logger.Fatal("Failed to connect to database", err)
	}
	defer db.Close()

	// 初始化输出写入器
	writer := output.NewWriter(&cfg.Output)

	// 创建执行器
	executor := executor.NewSQLExecutor(db, writer, logger)

	// 运行程序
	if err := executor.Run(); err != nil {
		logger.Fatal("Execution failed", err)
	}
}
