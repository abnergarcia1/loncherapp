package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	configOptions    []Config
	DbErrNoDocuments = errors.New("No documents matched")
)

type Client struct {
	id          string
	isConnected bool
	connStr     string
	reconnects  int
	sync.RWMutex
	*sqlx.DB
}

func (c *Client) Connect() error {
	if c.DB != nil {
		return nil
	}

	db, err := sqlx.Connect("mysql", c.connStr)

	if err != nil {
		return err
	}

	c.Lock()
	c.DB = db
	c.isConnected = true
	c.Unlock()

	return nil
}

func (c *Client) IsConnected() bool {
	if c.DB != nil {
		err := c.DB.Ping()
		return err == nil
	}
	return false
}

func (c *Client) Disconnect() {
	if c.DB == nil {
		return
	}
	c.RLock()
	defer c.RUnlock()
	if c.isConnected {
		_ = c.DB.Close()
	}

}

func (c *Client) Select(model interface{}, query string, args ...interface{}) (err error) {
	if c.DB == nil {
		return fmt.Errorf("cannot request to empty or closed connection")
	}

	err = c.DB.Select(model, query, args...)
	if err != nil {
		return
	}
	return
}

func (c *Client) Query(model interface{}, query string, args ...interface{}) (err error) {
	if c.DB == nil {
		return fmt.Errorf("cannot request to empty or closed connection")
	}

	res, err := c.DB.Query(query, args...)
	if err != nil {
		return
	}

	if err = res.Err(); err != nil {
		return
	}

	for res.Next() {
		if err = res.Err(); err != nil {
			return err
		}
		err := res.Scan(model)
		if err != nil {
			return err
		}
	}
	if err = res.Err(); err != nil {
		return err
	}

	err = res.Close()

	return
}

func (c *Client) QueryModel(model interface{}, query string, arg interface{}) (err error) {
	if c.DB == nil {
		return fmt.Errorf("Cannot request to empty or closed connection")
	}

	res, err := c.DB.NamedQuery(query, arg)
	if err != nil {
		return
	}
	if err = res.Err(); err != nil {
		return
	}

	for res.Next() {
		if err = res.Err(); err != nil {
			return err
		}
		err := res.StructScan(model)
		if err != nil {
			return err
		}
	}
	if err = res.Err(); err != nil {
		return err
	}

	err = res.Close()

	return
}

func (c *Client) Execute(query string, args ...interface{}) (sql.Result, error) {
	if c.DB == nil {
		return nil, fmt.Errorf("Cannot request to empty or closed connection")
	}

	res, err := c.DB.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) BeginTx() (*sqlx.Tx, error) {
	if c.DB == nil {
		return nil, fmt.Errorf("Cannot request to empty or closed connection")
	}

	return c.DB.Beginx()
}

func (c *Client) ExecuteTx(tx *sqlx.Tx, query string, args [][]interface{}) error {
	if tx == nil {
		return fmt.Errorf("Cannot request to emtpy or closed connection")
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, values := range args {
		_, err = stmt.Exec(values...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) QueryOne(model interface{}, query string, args ...interface{}) (err error) {
	if c.DB == nil {
		return fmt.Errorf("Cannot request to empty or closed connection")
	}

	err = c.DB.Get(model, query, args...)
	if err != nil {
		return
	}

	return
}

type Config func(*Client)

func NewClient() StorageDB {
	if os.Getenv("TESTING") == "true" {
		//return NewStorageDbMock()
	}

	client := &Client{}

	for _, conf := range configOptions {
		conf(client)
	}

	return client
}

func Init(configs ...Config) {
	configOptions = configs
}

func SetClientID(id string) Config {
	return func(c *Client) {
		c.id = id
	}
}

func SetConnectionString(conn string) Config {
	return func(c *Client) {
		c.connStr = conn
	}
}
