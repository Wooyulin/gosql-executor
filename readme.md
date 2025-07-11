# SQL Executor

一个支持 MySQL、Oracle、PostgreSQL 的 SQL 执行工具。

## 编译说明

### Windows 版本编译（在 Windows 系统上运行）
SET GOOS=windows
SET GOARCH=amd64
go build -o sql-executor.exe cmd/main.go

# Linux 版本
SET GOOS=linux
SET GOARCH=amd64
go build -o sql-executor cmd/main.go



# 直接运行
./sql-executor.exe

# 或者指定配置文件
./sql-executor.exe --config=config.yaml