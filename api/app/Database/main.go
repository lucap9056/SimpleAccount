package Database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type API struct {
	db       *sql.DB
	saltSize int
}

type Config struct {
	SourceName string
}

func New(data *Config) (*API, error) {
	database := API{}
	connect, err := sql.Open("mysql", data.SourceName)
	if err != nil {
		return nil, err
	}
	connect.SetMaxOpenConns(100)
	connect.SetMaxIdleConns(10)
	database.db = connect
	database.saltSize = 16

	return &database, nil
}

func (ctr *API) Connect() *sql.DB {
	return ctr.db
}

func (ctr *API) Close() error {
	return ctr.db.Close()
}
