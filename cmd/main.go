package main

import (
	"flag"
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
	configPath := flag.String("config", "./config.yaml", "配置文件路径")
	flag.Parse()

	// 检查配置文件是否存在
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Fatalf("配置文件不存在: %s\n请根据 config.yaml.example 创建配置文件", *configPath)
	}

	// 加载配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	logger := logger.NewLogger()

	// 初始化数据库连接
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		logger.Fatal("创建数据库连接失败", err)
	}

	if err := db.Connect(); err != nil {
		logger.Fatal("连接数据库失败", err)
	}
	defer db.Close()

	// 初始化输出写入器
	writer := output.NewWriter(&cfg.Output)

	// 创建执行器
	executor := executor.NewSQLExecutor(db, writer, logger)

	// 运行程序
	if err := executor.Run(); err != nil {
		logger.Fatal("程序执行失败", err)
	}
}
