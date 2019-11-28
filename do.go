package orm

import (
	"database/sql"
	"reflect"
	"time"

	"github.com/wizhodl/go-utils/log"
)

func (db *DB) do(any ...interface{}) *DB {
	var val = reflect.Indirect(reflect.ValueOf(db.sentence.mod))
	if val.Kind() != reflect.Struct {
		panic(ErrInvalidTable)
	}

	db.sentence.table(db.sentence.mod)
	db.sentence.String()

	var now = GetNowTime()
	var cost time.Duration

	if dbConfig.debugOn {
		defer func() {
			switch db.err {
			case nil:
				log.InfoflnN(3, db.sentence.raw)
			case ErrNotFound:
				log.WarningflnN(3, "%s %s;%v", db.sentence.raw, db.err, cost)
			default:
				log.ErrorflnN(3, "%s %s;%v", db.sentence.raw, db.err, cost)
			}
		}()
	}

	switch db.sentence.head.option {
	case optionSelect:
		if any == nil {
			db.err = ErrSelectQueryNeedDataHolder
			return db
		}

		var holder = any[0]

		var val = reflect.ValueOf(holder)
		var typ = reflect.TypeOf(holder)
		if val.Kind() != reflect.Ptr {
			db.err = ErrScanHolderMustBeValidPointer
			return db
		}

		var holderValue = val.Elem()
		var holderType = typ.Elem()

		if val.IsNil() {
			db.err = ErrScanHolderMustBeValidPointer
			return db
		}

		var rows *sql.Rows
		if db.isTxOn {
			rows, db.err = db.SqlTxDB.Query(db.sentence.raw)
		} else {
			rows, db.err = db.SqlDB.Query(db.sentence.raw)
		}

		cost = time.Since(now)

		if rows != nil {
			defer rows.Close()
		}

		if db.err != nil {
			return db
		}

		var columns []string
		columns, db.err = rows.Columns()
		if db.err != nil {
			return db
		}

		db.err = ErrNotFound

		switch holderType.Kind() {
		// scan to slice
		case reflect.Slice:
			if val.IsNil() {
				db.err = ErrScanHolderMustBeValidPointer
				return db
			}

			var baseType = holderType.Elem()
			var isBasePtr bool
			if isBasePtr = baseType.Kind() == reflect.Ptr; isBasePtr {
				baseType = baseType.Elem()
			}

			switch baseType.Kind() {
			case reflect.Struct:
				for rows.Next() {
					var table = reflect.New(baseType)
					var holder = table.Elem()
					db.err = scanRowsToTableValue(rows, columns, holder)
					if db.err != nil {
						return db
					}

					if isBasePtr {
						holderValue.Set(reflect.Append(holderValue, table))
					} else {
						holderValue.Set(reflect.Append(holderValue, holder))
					}
				}

			// TODO
			case reflect.Slice:

			default:
				db.err = ErrNotFound
				for rows.Next() {
					var column = reflect.New(baseType)
					var holder = column.Elem()
					db.err = rows.Scan(column.Interface())
					if db.err != nil {
						return db
					}

					if isBasePtr {
						holderValue.Set(reflect.Append(holderValue, column))
					} else {
						holderValue.Set(reflect.Append(holderValue, holder))
					}
				}
			}

		// scan to table struct
		case reflect.Struct:
			for rows.Next() {
				db.err = scanRowsToTableValue(rows, columns, holderValue)
				if db.err != nil {
					return db
				}
			}

		// scan to map
		case reflect.Map:
			initMap(holder)
			m, ok := (holder).(*(map[string]interface{}))
			if ok {
				db.err = scanRowsToMap(rows, *m)
			}

		default:
			if any != nil {
				db.err = scanRowsToAny(rows, any...)
			}
		}

	case optionInsert, optionUpdate:
		if db.isTxOn {
			_, db.err = db.SqlTxDB.Exec(db.sentence.raw)
		} else {
			_, db.err = db.SqlDB.Exec(db.sentence.raw)
		}

		cost = time.Since(now)
	}

	if db.err != nil && dbConfig.handleError != nil {
		dbConfig.handleError(db.err)
	}

	return db
}
