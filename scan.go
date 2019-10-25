package orm

import (
	"database/sql"
	"reflect"

	"github.com/wizhodl/go-utils/log"
)

type scanFunc func(rows *sql.Rows, holder interface{}) error

// slice must be like &[]*struct{} or &[]struct{}
func scanRowsToTableSlice(rows *sql.Rows, holder interface{}) error {
	var val = reflect.Indirect(reflect.ValueOf(holder))
	var typ = val.Type()
	var baseType = typ.Elem()
	var isBasePtr bool
	if isBasePtr = baseType.Kind() == reflect.Ptr; isBasePtr {
		baseType = baseType.Elem()
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	for rows.Next() {
		var table = reflect.New(baseType)
		var holder = table.Elem()
		err = scanRowsToTableValue(rows, columns, holder)
		if err != nil {
			return err
		}

		if isBasePtr {
			val.Set(reflect.Append(val, table))
		} else {
			val.Set(reflect.Append(val, holder))
		}
	}

	return nil
}

func scanRowsToTable(rows *sql.Rows, holder interface{}) error {
	var val = reflect.Indirect(reflect.ValueOf(holder))
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	for rows.Next() {
		err = scanRowsToTableValue(rows, columns, val)
		if err != nil {
			return err
		}
	}

	return nil
}

func scanRowsToTableValue(rows *sql.Rows, columns []string, table reflect.Value) error {
	var receivers = tableReceivers(table)
	var holders []interface{}
	for _, column := range columns {
		holders = append(holders, receivers[column])
	}
	return rows.Scan(holders...)
}

func scanRowsToMap(rows *sql.Rows, dst map[string]interface{}) error {
	types, err := rows.ColumnTypes()
	if err != nil {
		return err
	}

	for rows.Next() {
		var scanValues []interface{}
		for _, t := range types {
			var holder = reflect.New(t.ScanType())
			dst[t.Name()] = holder.Elem()
			scanValues = append(scanValues, holder.Interface())
		}

		var err = rows.Scan(scanValues...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) scan(any interface{}) *DB {
	var val = reflect.ValueOf(any)
	if val.Kind() != reflect.Ptr {
		db.err = errScanHolderMustBeValidPointer
		return db
	}

	if val.IsNil() {
		db.err = errScanHolderMustBeValidPointer
		return db
	}

	var query = db.sentence.String()
	log.Infofln("query: %s", query)
	rows, err := db.DB.Query(query)
	db.err = err
	if err != nil {
		return db
	}
	defer rows.Close()

	if val.Type().Kind() == reflect.Ptr {
		switch val.Type().Elem().Kind() {
		// scan to slice
		case reflect.Slice:
			scanRowsToTableSlice(rows, any)
			// scan to table struct
		case reflect.Struct:
			scanRowsToTable(rows, any)
			// scan to map
		case reflect.Map:
			initMap(any)
			m, ok := (any).(*(map[string]interface{}))
			if ok {
				db.err = scanRowsToMap(rows, *m)
			}
		default:
		}
	}

	return db
}