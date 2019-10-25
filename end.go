package orm

import (
	"reflect"
)

func (db *DB) end(any ...interface{}) *DB {
	var val = reflect.Indirect(reflect.ValueOf(db.mod))
	if val.Kind() != reflect.Struct {
		db.err = errInvalidTable
		return db
	}

	db.sentence.table(db.mod)

	switch db.sentence.head.option {
	case optionSelect:
		db.scan(any[0])
	}
	return db
}

func (db *DB) End(any ...interface{}) *DB {
	return db.end(any...)
}
