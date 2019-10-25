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
	s = strings.Join(j.and, " and ")

	var or = j.or
	or = append([]string{s}, or...)
	s = strings.Join(or, " or ")

	if j.subAnd != nil {
		var sa = j.subAnd.String()
		if sa != "" {
			s = strings.Join([]string{s, bracket(sa)}, " and ")
		}
	}

	if j.subOr != nil {
		var so = j.subOr.String()
		if so != "" {
			s = strings.Join([]string{s, bracket(so)}, " or ")
		}
	}

	return s
}

func joinColumnValue(column string, value interface{}) string {
	return strings.Join([]string{
		convertToSqlColumn(column),
		fmt.Sprintf("%v", convertToSqlValue(value)),
	}, "=")
}

func newJoinSquel(k string, v interface{}) *joinSquel {
	var j joinSquel
	j.and = append(j.and, joinColumnValue(k, v))
	return &j
}

func newJoinSquelFromMap(src map[string]interface{}) *joinSquel {
	var j joinSquel
	for k, v := range src {
		j.addAnd(k, v)
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
		j.addAnd(k, v)
	}
	return &j
}
