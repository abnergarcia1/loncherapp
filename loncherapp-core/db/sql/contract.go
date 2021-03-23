package sql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type StorageDB interface {
	Connect() error
	IsConnected() bool
	Disconnect()
	Execute(query string, args ...interface{}) (sql.Result, error)
	BeginTx() (*sqlx.Tx, error)
	ExecuteTx(tx *sqlx.Tx, query string, args [][]interface{}) error
	Query(model interface{}, query string, args ...interface{}) (err error)
	Select(model interface{}, query string, args ...interface{}) (err error)
	QueryModel(model interface{}, query string, args interface{}) (err error)
	QueryOne(model interface{}, query string, args ...interface{}) (err error)
}
