package orm

import (
	"database/sql"
)

var dbConfig connConfig

type DB struct {
	*sql.DB
	// mod is a table
	mod interface{}
	// sentence gen sql sentence
	sentence *sqlSentence
	// err returns any err
	err error
}

func (db *DB) copy() *DB {
	var d DB
	d.DB = db.DB
	d.err = db.err
	d.mod = db.mod
	d.sentence = new(sqlSentence)
	if db.sentence != nil {
		*d.sentence = *db.sentence
	}
	return &d
}
