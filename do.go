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
	var query = db.sentence.String()

	var now = GetNowTime()
	var cost time.Duration
	var finishQueryAt time.Time

	if dbConfig.debugOn {
		defer func() {
			switch db.err {
			case nil:
				log.InfoflnN(3, "%s%v@%v", query, cost, finishQueryAt)
			case ErrNotFound:
				log.WarningflnN(3, "%s %s;%v@%v", query, db.err, cost, finishQueryAt)
			default:
				log.ErrorflnN(3, "%s %s;%v@%v", query, db.err, cost, finishQueryAt)
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
			rows, db.err = db.SqlTxDB.Query(query)
		} else {
			rows, db.err = db.SqlDB.Query(query)
		}

		cost = time.Since(now)
		finishQueryAt = time.Now()

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
		var scanSlice = holderType.Kind() == reflect.Slice
		var scanTable = holderType.Kind() == reflect.Struct && holderType.Name() != "Time"
		var scanMap = holderType.Kind() == reflect.Map

		switch {
		// scan to slice
		case scanSlice:
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
		case scanTable:
			for rows.Next() {
				db.err = scanRowsToTableValue(rows, columns, holderValue)
				if db.err != nil {
					return db
				}
			}

		// scan to map
		case scanMap:
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
			_, db.err = db.SqlTxDB.Exec(query)
		} else {
			_, db.err = db.SqlDB.Exec(query)
		}

		cost = time.Since(now)
		finishQueryAt = time.Now()
	}

	if db.err != nil && dbConfig.handleError != nil {
		dbConfig.handleError(db.err)
	}

	return db
}
