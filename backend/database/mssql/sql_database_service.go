package mssql

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/microsoft/go-mssqldb"
)

type DBServiceInterface interface {
	Rebind(query string) string
	Select(dest interface{}, query string, args ...interface{}) error
	Close() error
}

type DBService struct {
	*sqlx.DB
}

func NewMSSQLDatabaseService(dsn string) (DBServiceInterface, error) {
	db, err := sqlx.Connect("sqlserver", dsn)
	if err != nil {
		return nil, err
	}
	return &DBService{db}, nil
}
