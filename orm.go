package orm

import (
	"database/sql"
)

var dbConfig connConfig

type DB struct {
	*sql.DB
	// sentence gen sql sentence
	sentence *sqlSentence
	// err returns any err
	err error
}

func (db *DB) copy() *DB {
	var d DB
	d.DB = db.DB

	if db.sentence == nil {
		d.sentence = newSentence()
	} else {
		d.sentence = db.sentence.copy()
	}

	return &d
}

func (db *DB) Table(t interface{}) *DB {
	return db.table(t)
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

// Insert(&struct{})
// Insert([]*struct{})
// Insert(format string, args ...interface{})
func (db *DB) Insert(set interface{}, args ...interface{}) *DB {
	return db.insert(set, args...)
}

// Update(&struct{}, string, string...)
// Update([]*struct{})
// Update(map[string]interface{}, string, string ...)
// Update(format string, args ...interface{})
func (db *DB) Update(set interface{}, args ...interface{}) *DB {
	return db.update(set, args...)
}

func (db *DB) Do(any ...interface{}) error {
	return db.do(any...).err
}
