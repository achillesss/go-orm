package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/wizhodl/go-utils/log"
)

func (db *DB) do(any ...interface{}) *DB {
	var op = db.sentence.head.option
	var val = reflect.Indirect(reflect.ValueOf(db.sentence.mod))
	if op != optionRaw {
		if val.Kind() != reflect.Struct {
			panic(fmt.Sprintf("%s:%s\n", ErrInvalidTable, val.Kind()))
		}
		db.sentence.table(db.sentence.mod)
	}

	var query = db.sentence.String()
	var now = GetNowTime()
	var queryID string
	var caller = log.CallerLine(2)

	if dbConfig.startQueryMonitor != nil || dbConfig.endQueryMonitor != nil {
		queryID = uuid.New().String()
	}

	*db.StartCount++
	if dbConfig.startQueryMonitor != nil {
		go func() {
			startQueryChan <- &StartQuery{
				ID:      queryID,
				Query:   query,
				StartAt: now,
				Caller:  caller,
			}
		}()
	}

	var cost time.Duration
	var finishQueryAt time.Time

	defer func() {
		*db.EndCount++
		if dbConfig.endQueryMonitor == nil {
			return
		}
		go func() {
			endQueryChan <- &EndQuery{
				ID:     queryID,
				Query:  query,
				EndAt:  finishQueryAt,
				Cost:   cost,
				Error:  db.err,
				Caller: caller,
			}
		}()
	}()

	switch op {
	case optionInsert, optionUpdate, optionDelete, optionRaw:

		var r sql.Result
		if db.isTxOn {
			r, db.err = db.SqlTxDB.Exec(query)
		} else {
			r, db.err = db.SqlDB.Exec(query)
		}

		if op == optionInsert && db.err == nil && db.sentence.updateIDFunc != nil {
			db.sentence.updateIDFunc(r)
		}

		cost = time.Since(now)
		finishQueryAt = time.Now()

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
				return db
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
	}

	if db.err != nil && dbConfig.handleError != nil {
		dbConfig.handleError(db.err)
	}

	return db
}
