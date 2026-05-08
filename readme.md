# SQL Executor

一个轻量级的命令行 SQL 查询工具，支持 MySQL、Oracle、PostgreSQL，查询结果可导出为 CSV/JSON 文件。

## 功能

- 交互式 SQL 输入，支持多行查询（以 `;` 结束）
- 支持 MySQL、Oracle、PostgreSQL 三种数据库
- 查询结果输出到控制台（表格格式）和/或文件（CSV/JSON）
- 上下键回顾历史命令，支持退格编辑
- 纯静态编译（`CGO_ENABLED=0`），Linux 各发行版通用

## 快速开始

### 1. 下载

从 [GitHub Releases](https://github.com/Wooyulin/gosql-executor/releases) 下载对应平台的压缩包并解压。

### 2. 配置

将 `config.yaml.example` 复制为 `config.yaml`，按需修改：

```yaml
database:
  type: "mysql"          # mysql, oracle, pgsql
  username: "root"
  password: "password"
  dsn: "{username}:{password}@tcp(localhost:3306)/dbname"

output:
  directory: "./output"
  format: "csv"          # csv, json
  show_in_console: true  # 控制台显示结果
  save_to_file: true     # 结果保存到文件
```

DSN 示例：

| 数据库     | DSN 格式                                                             |
|-----------|----------------------------------------------------------------------|
| MySQL     | `{username}:{password}@tcp(host:3306)/dbname`                        |
| Oracle    | `oracle://{username}:{password}@host:1521/service_name`              |
| PostgreSQL| `postgres://{username}:{password}@host:5432/dbname?sslmode=disable`  |

### 3. 运行

```bash
# 使用默认配置文件 ./config.yaml
./sql-executor

# 指定配置文件
./sql-executor --config=/path/to/config.yaml
```

### 使用示例

```
SQL> SELECT * FROM users LIMIT 5;
id   |name   |email
---  |---    |---
1    |张三   |zhangsan@example.com
2    |李四   |lisi@example.com

[INFO] 查询执行完成，耗时: 12.5ms，结果已保存到: query_result_1735374074

SQL> exit
```

输入多行查询：

```
SQL> SELECT u.name, o.total
  -> FROM users u
  -> JOIN orders o ON u.id = o.user_id
  -> WHERE o.total > 100;
```

## 从源码编译

需要 Go 1.22+。

```bash
# Windows
go build -o sql-executor.exe ./cmd

# Linux / macOS
go build -o sql-executor ./cmd

# 交叉编译 Linux（推荐，CGO_ENABLED=0 保证各发行版通用）
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sql-executor ./cmd
```

## 项目结构

```
cmd/main.go                 # 入口
config/config.go            # 配置解析
internal/
  database/database.go      # 数据库接口
  database/mysql.go         # MySQL 实现
  database/oracle.go        # Oracle 实现
  database/postgres.go      # PostgreSQL 实现
  executor/executor.go      # SQL 执行与交互逻辑
  output/writer.go          # 结果输出（CSV/JSON/控制台）
pkg/logger/logger.go        # 日志
```

## License

[MIT](LICENSE)
