package orm

import (
	"fmt"
	"reflect"
	"strings"
)

type sqlValue struct {
	columns []string
	values  []interface{}
}

func (v sqlValue) InsertString() string {
	if len(v.columns) != len(v.values) {
		panic(fmt.Sprintf("columns %d with %d values", len(v.columns), len(v.values)))
		return ""
	}

	var columns []string
	var values []string
	for i := range v.columns {
		columns = append(columns, convertToSqlColumn(v.columns[i]))
		values = append(values, fmt.Sprintf("%v", v.values[i]))
	}

	return strings.Join([]string{
		bracket(strings.Join(columns, ",")),
		bracket(strings.Join(values, ",")),
	}, " VALUES ")
}

func (v sqlValue) UpdateString() string {
	if len(v.columns) != len(v.values) {
		panic(fmt.Sprintf("columns %d with %d values", len(v.columns), len(v.values)))
		return ""
	}

	var pairs []string
	for i := range v.columns {
		pairs = append(pairs, strings.Join([]string{
			convertToSqlColumn(v.columns[i]),
			fmt.Sprintf("%v", v.values[i]),
		}, "="))
	}

	return strings.Join(pairs, ",")
}

type sqlValues struct {
	columns []string
	values  []map[string]interface{}
}

func (v sqlValues) String() string {
	var columns []string
	var values []map[string]string
	for i := range v.columns {
		columns = append(columns, convertToSqlColumn(v.columns[i]))
		for j, row := range v.values {
			if len(values)-1 < j {
				values = append(values, make(map[string]string))
			}
			values[j][v.columns[i]] = fmt.Sprintf("%v", row[v.columns[i]])
		}
	}

	var rows []string
	for _, row := range values {
		var r []string
		for _, column := range v.columns {
			r = append(r, row[column])
		}
		rows = append(rows, bracket(strings.Join(r, ",")))
	}

	return strings.Join([]string{
		bracket(strings.Join(columns, ",")),
		strings.Join(rows, ","),
	}, " VALUES ")
}

// batch insert
func (h *sqlHead) insert() {
	h.option = optionInsert
}

func (s *sqlValue) insertSingle(table interface{}) {
	var val = reflect.Indirect(reflect.ValueOf(table))
	s.columns, s.values = readTable(val, false)
}

func (s *sqlValues) insertSlice(tables ...interface{}) {
	for _, table := range tables {
		var val = reflect.Indirect(reflect.ValueOf(table))
		var columnValues = make(map[string]interface{})
		var columns, values = readTable(val, true)
		for i := range columns {
			columnValues[columns[i]] = values[i]
		}
		s.values = append(s.values, columnValues)
		if s.columns != nil {
			continue
		}
		s.columns = columns
	}
}

func (s *sqlSentence) insert(set interface{}, args ...interface{}) {
	s.head.insert()
	var val = reflect.Indirect(reflect.ValueOf(set))
	switch val.Kind() {
	case reflect.Slice:
		var values sqlValues
		s.values = &values

		var tables []interface{}
		for i := 0; i < val.Len(); i++ {
			tables = append(tables, val.Index(i).Interface())
		}
		values.insertSlice(tables...)
		s.mod = tables[0]

	case reflect.Struct:
		s.mod = set
		var value sqlValue
		s.value = &value
		value.insertSingle(set)

	}
}

func (db *DB) insert(set interface{}, args ...interface{}) *DB {
	var d = db.copy()
	d.sentence.insert(set, args...)
	return d
}
