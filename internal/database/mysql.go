package database

import (
	"database/sql"
	"sql-executor/config"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLDatabase struct {
	BaseDatabase
}

func NewMySQLDatabase(cfg *config.DatabaseConfig) *MySQLDatabase {
	return &MySQLDatabase{
		BaseDatabase: BaseDatabase{
			Config: cfg,
		},
	}
}

func (m *MySQLDatabase) Connect() error {
	db, err := sql.Open("mysql", m.BaseDatabase.GetDSN())
	if err != nil {
		return err
	}

	m.BaseDatabase.DB = db
	return db.Ping()
}

func (m *MySQLDatabase) Execute(query string) (*sql.Rows, error) {
	return m.BaseDatabase.DB.Query(query)
}

func (m *MySQLDatabase) Close() error {
	if m.BaseDatabase.DB != nil {
		return m.BaseDatabase.DB.Close()
	}
	return nil
}
