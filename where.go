package orm

import (
	"fmt"
	"math"
	"reflect"
	"strings"
)

type joinSquel struct {
	and    []string
	or     []string
	subAnd *joinSquel
	subOr  *joinSquel
}

func (j *joinSquel) String() string {
	if j.and == nil {
		return ""
	}

	var s string
	s = strings.Join(j.and, " AND ")

	var or = j.or
	or = append([]string{s}, or...)
	s = strings.Join(or, " OR ")

	if j.subAnd != nil {
		var sa = j.subAnd.String()
		if sa != "" {
			s = strings.Join([]string{s, bracket(sa)}, " AND ")
		}
	}

	if j.subOr != nil {
		var so = j.subOr.String()
		if so != "" {
			s = strings.Join([]string{s, bracket(so)}, " OR ")
		}
	}

	return s
}

func joinColumnValue(column string, symbol string, values ...interface{}) string {
	if len(values) == 0 {
		return ""
	}

	if len(values) == 1 {
		var val = fmt.Sprintf("%v", values[0])
		if symbol == "" {
			symbol = "="
			if strings.EqualFold(val, "NULL") {
				symbol = " IS "
			}
		}

		return strings.Join([]string{
			convertToSqlColumn(column),
			val,
		}, symbol)
	}

	if symbol == "" {
		symbol = " IN "
	}

	var strValues []string
	for _, v := range values {
		strValues = append(strValues, fmt.Sprintf("%v", v))
	}

	return strings.Join([]string{
		convertToSqlColumn(column),
		bracket(strings.Join(strValues, ",")),
	}, symbol)
}

func newJoinSquel(k string, symbol string, v interface{}) *joinSquel {
	var j joinSquel
	j.and = append(j.and, joinColumnValue(k, symbol, v))
	return &j
}

func newJoinSquelFromMap(src map[string][]interface{}, symbol string) *joinSquel {
	var j joinSquel
	for k, v := range src {
		j.addAnd(k, symbol, v...)
	}
	return &j
}

func newJoinSquelFromTable(table reflect.Value) *joinSquel {
	var columns, values = readTable(table, false)
	var maxLength = int(math.Max(float64(len(columns)), float64(len(values))))
	var j joinSquel
	for i := 0; i < maxLength; i++ {
		var k = columns[i]
		var v = values[i]
		j.addAnd(k, "=", v)
	}
	return &j
}
