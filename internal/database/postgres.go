package database

import (
	"database/sql"
	"fmt"
	"sql-executor/config"

	_ "github.com/lib/pq"
)

type PostgresDatabase struct {
	BaseDatabase
}

func NewPostgresDatabase(cfg *config.DatabaseConfig) *PostgresDatabase {
	return &PostgresDatabase{
		BaseDatabase: BaseDatabase{
			Config: cfg,
		},
	}
}

func (p *PostgresDatabase) Connect() error {
	db, err := sql.Open("postgres", p.BaseDatabase.GetDSN())
	if err != nil {
		return fmt.Errorf("连接PostgreSQL数据库失败: %w", err)
	}

	p.BaseDatabase.DB = db
	return db.Ping()
}

func (p *PostgresDatabase) Execute(query string) (*sql.Rows, error) {
	return p.BaseDatabase.DB.Query(query)
}

func (p *PostgresDatabase) Close() error {
	if p.BaseDatabase.DB != nil {
		return p.BaseDatabase.DB.Close()
	}
	return nil
}
