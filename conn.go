package orm

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func (c *connConfig) Open() (*DB, error) {
	db, err := sql.Open(c.driverName, c.loginString())
	if err != nil {
		return nil, err
	}
	dbConfig = *c
	return &DB{SqlDB: db, sentence: &sqlSentence{}}, nil
}
