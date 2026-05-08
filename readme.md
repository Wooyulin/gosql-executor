# SQL Executor

A lightweight command-line SQL query tool supporting MySQL, Oracle, and PostgreSQL. Query results can be exported to CSV/JSON files.

## Features

- Interactive SQL input with multi-line query support (end with `;`)
- Support for MySQL, Oracle, and PostgreSQL
- Output to console (table format) and/or file (CSV/JSON)
- Arrow key history navigation, backspace editing
- Static binary (`CGO_ENABLED=0`), works across Linux distributions

## Quick Start

### 1. Download

Download the archive for your platform from [GitHub Releases](https://github.com/Wooyulin/gosql-executor/releases) and extract it.

### 2. Configuration

There are two ways to configure the database connection:

**Option A: Interactive Setup (Recommended)**

Run the program directly without creating a config file. It will guide you through the database connection setup:

```bash
./sql-executor
```

```
Config file not found. Starting interactive setup:
Database type [mysql/oracle/pgsql]: mysql
Host [127.0.0.1]: 192.168.1.100
Port [3306]:
Username: root
Password: ****
Database name: mydb

Config saved to ./config.yaml
```

**Option B: Manual Configuration**

Copy `config.yaml.example` to `config.yaml` and edit as needed:

```yaml
database:
  type: "mysql"          # mysql, oracle, pgsql
  username: "root"
  password: "password"
  dsn: "{username}:{password}@tcp(localhost:3306)/dbname"

output:
  directory: "./output"
  format: "csv"          # csv, json
  show_in_console: true  # display results in console
  save_to_file: true     # save results to file
```

DSN Examples:

| Database   | DSN Format                                                           |
|------------|----------------------------------------------------------------------|
| MySQL      | `{username}:{password}@tcp(host:3306)/dbname`                        |
| Oracle     | `oracle://{username}:{password}@host:1521/service_name`              |
| PostgreSQL | `postgres://{username}:{password}@host:5432/dbname?sslmode=disable`  |

### 3. Run

```bash
# Use default config file ./config.yaml (enters interactive setup on first run)
./sql-executor

# Specify config file
./sql-executor --config=/path/to/config.yaml
```

### Usage Example

```
SQL> SELECT * FROM users LIMIT 5;
id   |name   |email
---  |---    |---
1    |Alice  |alice@example.com
2    |Bob    |bob@example.com

[INFO] Query executed in 12.5ms, saved to: query_result_1735374074

SQL> exit
```

Multi-line query:

```
SQL> SELECT u.name, o.total
  -> FROM users u
  -> JOIN orders o ON u.id = o.user_id
  -> WHERE o.total > 100;
```

## Build from Source

Requires Go 1.22+.

```bash
# Windows
go build -o sql-executor.exe ./cmd

# Linux / macOS
go build -o sql-executor ./cmd

# Cross-compile for Linux (recommended, CGO_ENABLED=0 ensures broad compatibility)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sql-executor ./cmd
```

## Project Structure

```
cmd/main.go                 # Entry point
config/config.go            # Config parsing
config/interactive.go       # Interactive setup
internal/
  database/database.go      # Database interface
  database/mysql.go         # MySQL implementation
  database/oracle.go        # Oracle implementation
  database/postgres.go      # PostgreSQL implementation
  executor/executor.go      # SQL execution & interaction
  output/writer.go          # Result output (CSV/JSON/console)
pkg/logger/logger.go        # Logger
```

## License

[MIT](LICENSE)
