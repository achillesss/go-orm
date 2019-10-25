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

// and xx and xx and xx
func (j *joinSquel) addAndRaw(raw string) *joinSquel {
	j.and = append(j.and, raw)
	return j
}

func (j *joinSquel) addAnd(column string, value interface{}) *joinSquel {
	j.and = append(j.and, joinColumnValue(column, value))
	return j
}

func (j *joinSquel) setSubAnd(s *joinSquel) *joinSquel {
	j.subAnd = s
	return j
}

func (s *joinSquel) andMap(src map[string]interface{}) *joinSquel {
	var j = newJoinSquelFromMap(src)
	s.addAndRaw(j.String())
	return s
}

func (s *joinSquel) andTable(src reflect.Value) *joinSquel {
	var j = newJoinSquelFromTable(src)
	s.addAndRaw(j.String())
	return s
}

func (db *DB) Where(where interface{}, args ...interface{}) *DB {
	if db.sentence.where == nil {
		db.sentence.where = new(joinSquel)
	}
	var val = reflect.Indirect(reflect.ValueOf(where))
	var typ = val.Type()
	switch typ.Kind() {
	case reflect.String:
		db.sentence.where.addAndRaw(fmt.Sprintf(val.String(), args...))
	case reflect.Struct:
		db.sentence.where.andTable(val)
		db.mod = where

	case reflect.Map:
		m, ok := where.(map[string]interface{})
		if ok {
			db.sentence.where.andMap(m)
		}
	}

	return db
}
