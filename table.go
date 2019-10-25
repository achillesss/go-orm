package orm

import (
	"reflect"
	"strings"
)

func getValueTableName(table reflect.Value, getTableNameMethod string) (tableName string, ok bool) {
	var typ = table.Type()
	if typ.Kind() != reflect.Struct {
		return
	}

	tableName = camalToSnake(typ.Name())
	ok = true

	m, methodOK := typ.MethodByName(getTableNameMethod)
	if !methodOK {
		return
	}

	var ft = m.Func.Type()
	if ft.NumIn() != 1 {
		return
	}

	if ft.NumOut() != 1 {
		return
	}

	if ft.Out(0).Kind() != reflect.String {
		return
	}

	var names = m.Func.Call([]reflect.Value{table})
	tableName = names[0].String()
	return
}

func getTableName(t interface{}, getTableNameMethod string) (tableName string, ok bool) {
	var val = reflect.Indirect(reflect.ValueOf(t))
	return getValueTableName(val, getTableNameMethod)
}

// parse tag to map[string]string
// tag: `gorm:"k1:v1;k2:v2"`
func parseTag(t reflect.StructTag) map[string]string {
	var valueString = t.Get("gorm")
	if valueString == "" {
		return nil
	}

	var valueSlice = strings.Split(valueString, ";")
	var tagKV = make(map[string]string)
	for _, v := range valueSlice {
		v = strings.Trim(v, " ")
		var kv = strings.Split(v, ":")
		if len(kv) != 2 {
			continue
		}
		var key = strings.Trim(kv[0], " ")
		var val = strings.Trim(kv[1], " ")
		tagKV[key] = val
	}

	if len(tagKV) == 0 {
		return nil
	}

	return tagKV
}

func getColumnName(valueField reflect.Value, typeField reflect.StructField) (name string, isColumn bool, isStruct bool) {
	if !valueField.CanInterface() {
		return
	}

	var typ = valueField.Type()
	k := typ.Kind()
	if k == reflect.Ptr {
		typ = typ.Elem()
		k = typ.Kind()
	}

	_, ok := validColumn[k]
	if !ok {
		return
	}

	if k == reflect.Struct && typ.Name() != "Time" {
		isStruct = true
		return
	}

	var tags = parseTag(typeField.Tag)
	if tags == nil {
		return
	}

	isColumn = true
	name = camalToSnake(typeField.Name)

	n, ok := tags["column"]
	if !ok {
		return
	}

	name = n
	return
}

func readTable(val reflect.Value, readDefaultValue bool) (columns []string, values []interface{}) {
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		var columeName, isColumn, isStruct = getColumnName(valueField, typeField)

		if isStruct {
			cs, vs := readTable(valueField, readDefaultValue)
			columns = append(columns, cs...)
			values = append(values, vs...)
			continue
		}

		if !isColumn {
			continue
		}

		if !readDefaultValue {
			var defaultValue = reflect.Zero(valueField.Type())
			if reflect.DeepEqual(defaultValue.Interface(), valueField.Interface()) {
				continue
			}
		}

		columns = append(columns, columeName)
		values = append(values, valueField.Interface())
	}

	return
}

func tableReceivers(val reflect.Value) map[string]interface{} {
	var tr = make(map[string]interface{})
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		var columnName, isColumn, isStruct = getColumnName(valueField, typeField)

		if isStruct {
			rs := tableReceivers(valueField)
			for k, v := range rs {
				tr[k] = v
			}
			continue
		}

		if !isColumn {
			continue
		}

		tr[columnName] = valueField.Addr().Interface()
	}
	return tr
}

func (q *sqlSentence) table(t interface{}) *sqlSentence {
	var val = reflect.Indirect(reflect.ValueOf(t))
	return q.valueTable(val)
}

func (q *sqlSentence) valueTable(table reflect.Value) *sqlSentence {
	tName, ok := getValueTableName(table, dbConfig.getTableNameMethod)
	if !ok {
		return q
	}

	q.tableName = tName
	return q
}

func (db *DB) Table(t interface{}) *DB {
	db.mod = t
	return db
}
