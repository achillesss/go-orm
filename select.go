package orm

import (
	"database/sql"
	"reflect"
)

// TODO: able to register any holder type

func (h *sqlHead) select_(columns ...string) {
	h.option = optionSelect
	if columns == nil {
		return
	}

	h.fields = columns
}

func (q *sqlSentence) select_(columns ...string) *sqlSentence {
	q.head.select_(columns...)
	return q

}

func (db *DB) select_(columns ...string) *DB {
	var d = db.copy()
	d.sentence.select_(columns...)
	return d
}

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

	err = ErrNotFound

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

	return err
}

func scanRowsToTable(rows *sql.Rows, holder interface{}) error {
	var val = reflect.Indirect(reflect.ValueOf(holder))
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	err = ErrNotFound

	for rows.Next() {
		err = scanRowsToTableValue(rows, columns, val)
		if err != nil {
			return err
		}
	}

	return err
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

	err = ErrNotFound

	for rows.Next() {
		var scanValues []interface{}
		for _, t := range types {
			var holder = reflect.New(t.ScanType())
			dst[t.Name()] = holder.Elem().Interface()
			scanValues = append(scanValues, holder.Interface())
		}

		err = rows.Scan(scanValues...)
		if err != nil {
			return err
		}
	}

	return err
}

func scanRowsToAny(rows *sql.Rows, any ...interface{}) error {
	var err error = ErrNotFound

	for rows.Next() {
		err = rows.Scan(any...)
	}

	return err
}
