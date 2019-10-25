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

func (db *DB) Table(t interface{}) *DB {
	db.mod = t
	return db
}

// Select *: none input
// Select columns: []string
func (db *DB) Select(columns ...string) *DB {
	return db.select_(columns...)
}

func (db *DB) Where(where interface{}, args ...interface{}) *DB {
	return db.and(where, args...)
}

func (db *DB) And(where interface{}, args ...interface{}) *DB {
	return db.and(where, args...)
}

func (db *DB) Or(where interface{}, args ...interface{}) *DB {
	return db.or(where, args...)
}

func (db *DB) End(any ...interface{}) *DB {
	return db.end(any...)
}
