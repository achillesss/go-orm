package orm

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/wizhodl/go-utils/log"
)

var space = ""

var upper = regexp.MustCompile("[A-Z]")

func camelToSnake(s string) string {
	var a []string
	var lastIndex int
	var lastIsUpper bool
	for i, r := range s {
		var b = []byte{byte(r)}
		var isUpper = upper.Match(b)
		if i == 0 {
			lastIsUpper = isUpper
			continue
		}

		if isUpper && !lastIsUpper {
			a = append(a, s[lastIndex:i])
			lastIndex = i
		}

		lastIsUpper = isUpper

		if i == len(s)-1 {
			a = append(a, s[lastIndex:])
		}

	}
	return strings.ToLower(strings.Join(a, "_"))
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

func defaultSQLValue(column string) string {
	return fmt.Sprintf("DEFAULT(%s)", convertToSqlColumn(column))
}

func convertValueToSqlValue(column string, val reflect.Value, keepOriginValue bool) interface{} {
	// is pointer
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return NULL
		}
		val = val.Elem()
	}
	var typ = val.Type()

	if !keepOriginValue {
		var zeroValue = reflect.Zero(typ)
		if reflect.DeepEqual(zeroValue.Interface(), val.Interface()) {
			return defaultSQLValue(column)
		}
	}

	if typ.Kind() == reflect.String {
		return strconv.Quote(val.String())
	}

	if typ.Kind() == reflect.Struct && typ.Name() == "Time" {
		return dateTimeStr(val.Interface().(time.Time))
	}

	return val.Interface()
}

func convertToSqlValue(column string, src interface{}, keepOriginValue ...bool) interface{} {
	var val = reflect.ValueOf(src)
	return convertValueToSqlValue(column, val, len(keepOriginValue) > 0 && keepOriginValue[0])
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

func debugLog(format string, args ...interface{}) {
	if dbConfig.debugOn {
		log.Infofln(format, args...)
	}
}

func GetNowTime() time.Time { return time.Now() }
