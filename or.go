package orm

import (
	"fmt"
	"reflect"
)

func (j *joinSquel) addOrRaw(raw string) *joinSquel {
	j.or = append(j.or, raw)
	return j
}

// or xx or xx or xx
func (j *joinSquel) addOr(column string, value interface{}) *joinSquel {
	j.or = append(j.or, joinColumnValue(column, value))
	return j
}

func (j *joinSquel) setSubOr(s *joinSquel) *joinSquel {
	j.subOr = s
	return j
}

func (s *joinSquel) orMap(src map[string]interface{}) *joinSquel {
	var j = newJoinSquelFromMap(src)
	if j.subOr == nil {
		s.subOr = j
		return s
	}
	s.subOr.addAndRaw(j.String())
	return s
}

func (s *joinSquel) orTable(src reflect.Value) *joinSquel {
	var j = newJoinSquelFromTable(src)
	s.addOrRaw(j.String())
	return s
}

func (db *DB) or(where interface{}, args ...interface{}) *DB {
	var val = reflect.Indirect(reflect.ValueOf(where))
	switch w := val.Kind(); w {
	case reflect.String:
		db.sentence.where.addOrRaw(fmt.Sprintf(val.String(), args...))
	case reflect.Map:
		m, ok := val.Interface().(map[string]interface{})
		if ok {
			db.sentence.where.orMap(m)
		}
	case reflect.Struct:
		db.sentence.where.orTable(val)
	}
	return db
}

func (db *DB) Or(where interface{}, args ...interface{}) *DB {
	return db.or(where, args...)
}
