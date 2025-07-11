package database

import (
	"database/sql"
	"fmt"
	"sql-executor/config"

	_ "github.com/sijms/go-ora/v2"
)

type OracleDatabase struct {
	BaseDatabase
}

func NewOracleDatabase(cfg *config.DatabaseConfig) *OracleDatabase {
	return &OracleDatabase{
		BaseDatabase: BaseDatabase{
			Config: cfg,
		},
	}
}

func (o *OracleDatabase) Connect() error {
	db, err := sql.Open("oracle", o.BaseDatabase.GetDSN())
	if err != nil {
		return fmt.Errorf("连接Oracle数据库失败: %w", err)
	}

	o.BaseDatabase.DB = db
	return db.Ping()
}

func (o *OracleDatabase) Execute(query string) (*sql.Rows, error) {
	return o.BaseDatabase.DB.Query(query)
}

func (o *OracleDatabase) Close() error {
	if o.BaseDatabase.DB != nil {
		return o.BaseDatabase.DB.Close()
	}
	return nil
}
