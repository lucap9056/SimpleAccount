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

	db, err := sql.Open("mysql", data.SourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	database.db = db
	database.saltSize = 16

	if err := database.checkAndCreateTable(); err != nil {
		return nil, err
	}

	return &database, nil
}

func (ctr *API) checkAndCreateTable() error {
	var table string
	err := ctr.db.QueryRow("SHOW TABLES LIKE 'User'").Scan(&table)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if table == "" {
		_, err := ctr.db.Exec(`
		CREATE TABLE User (
			id              INT             AUTO_INCREMENT,
			username        VARCHAR(32)     NOT NULL,
			email           VARCHAR(256)    NOT NULL,
			salt            VARCHAR(64)     NOT NULL,
			hash            VARCHAR(64)     NOT NULL,
			last_edit       TIMESTAMP       NOT NULL        DEFAULT CURRENT_TIMESTAMP(),
			register_time   TIMESTAMP       NOT NULL        DEFAULT CURRENT_TIMESTAMP(),
			deleted         TIMESTAMP,
			PRIMARY KEY (id)
		)`)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ctr *API) Connect() *sql.DB {
	return ctr.db
}

func (ctr *API) Close() error {
	return ctr.db.Close()
}
