package orm

import (
	"reflect"
)

func (db *DB) do(any ...interface{}) *DB {
	var val = reflect.Indirect(reflect.ValueOf(db.sentence.mod))
	if val.Kind() != reflect.Struct {
		db.err = errInvalidTable
		return db
	}

	db.sentence.table(db.sentence.mod)

	switch db.sentence.head.option {
	case optionSelect:
		if any == nil {
			db.err = errSelectQueryNeedDataHolder
			return db
		}

		db.doSelect(any[0])

	case optionInsert:
		db.doInsert()

	case optionUpdate:
		db.doUpdate()
	}

	return db
}
