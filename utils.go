package orm

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var space = ""

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func camalToSnake(camal string) string {
	var snake = matchFirstCap.ReplaceAllString(camal, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

var invalidColumn = map[reflect.Kind]interface{}{
	reflect.Invalid:       nil,
	reflect.Uintptr:       nil,
	reflect.Complex64:     nil,
	reflect.Complex128:    nil,
	reflect.Array:         nil,
	reflect.Chan:          nil,
	reflect.Func:          nil,
	reflect.Interface:     nil,
	reflect.Map:           nil,
	reflect.Ptr:           nil,
	reflect.Slice:         nil,
	reflect.UnsafePointer: nil,
}

var validColumn = map[reflect.Kind]interface{}{
	reflect.Bool:    nil,
	reflect.Int:     nil,
	reflect.Int8:    nil,
	reflect.Int16:   nil,
	reflect.Int32:   nil,
	reflect.Int64:   nil,
	reflect.Uint:    nil,
	reflect.Uint8:   nil,
	reflect.Uint16:  nil,
	reflect.Uint32:  nil,
	reflect.Uint64:  nil,
	reflect.Float32: nil,
	reflect.Float64: nil,
	reflect.String:  nil,
	reflect.Struct:  nil,
}

const (
	NULL       = "NULL"
	timeFormat = "2006-01-02 15:04:05"
	dateFormat = "2006-01-02"
)

func dateTimeStr(t time.Time) string {
	return strconv.Quote(t.Format(timeFormat))
}

func dateStr(t time.Time) string {
	return fmt.Sprintf("%q", t.Format(dateFormat))
}

func convertValueToSqlValue(val reflect.Value) interface{} {
	var typ = val.Type()
	// is pointer
	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			return NULL
		}
		val = val.Elem()
		typ = val.Type()
	}

	if typ.Kind() == reflect.String {
		return strconv.Quote(val.String())
	}

	if typ.Kind() == reflect.Struct && typ.Name() == "Time" {
		return dateTimeStr(val.Interface().(time.Time))
	}

	return val.Interface()
}

func convertToSqlValue(src interface{}) interface{} {
	var val = reflect.ValueOf(src)
	return convertValueToSqlValue(val)
}

func convertToSqlColumns(columns []string) []string {
	var cs []string
	for _, c := range columns {
		cs = append(cs, convertToSqlColumn(c))
	}
	return cs
}

func convertToSqlColumn(column string) string {
	column = strings.Trim(column, "`")
	return "`" + column + "`"
}

func bracket(str string) string { return "(" + str + ")" }

func initMap(src interface{}) bool {
	var val = reflect.ValueOf(src)
	if val.Kind() != reflect.Ptr {
		if val.Kind() == reflect.Map {
			return true
		}
		return false
	}

	if val.Kind() == reflect.Ptr {
		var elm = val.Elem()
		if elm.Kind() != reflect.Map {
			return false
		}
		elm.Set(reflect.MakeMap(elm.Type()))
	}
	return true
}
