package database

import (
	"database/sql"
	"fmt"
	"net/url"
	"sql-executor/config"
	"strings"
)

type Database interface {
	Connect() error
	Close() error
	Execute(query string) (*sql.Rows, error)
}

type BaseDatabase struct {
	DB     *sql.DB
	Config *config.DatabaseConfig
}

// 处理 DSN 字符串，替换用户名和密码
func (b *BaseDatabase) GetDSN() string {
	var username, password string

	switch b.Config.Type {
	case "oracle":
		username = url.QueryEscape(b.Config.Username)
		password = url.QueryEscape(b.Config.Password)
	default:
		username = b.Config.Username
		password = b.Config.Password
	}

	dsn := b.Config.DSN
	dsn = strings.ReplaceAll(dsn, "{username}", username)
	dsn = strings.ReplaceAll(dsn, "{password}", password)

	return dsn
}

func NewDatabase(cfg *config.DatabaseConfig) (Database, error) {
	switch cfg.Type {
	case "mysql":
		return NewMySQLDatabase(cfg), nil
	case "oracle":
		return NewOracleDatabase(cfg), nil
	case "pgsql":
		return NewPostgresDatabase(cfg), nil
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", cfg.Type)
	}
}
