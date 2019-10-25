package orm

import (
	"fmt"
	"reflect"
)

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

func (db *DB) and(where interface{}, args ...interface{}) *DB {
	if db.sentence.where == nil {
		db.sentence.where = new(joinSquel)
	}
	var val = reflect.Indirect(reflect.ValueOf(where))
	switch val.Kind() {
	case reflect.String:
		db.sentence.where.addAndRaw(fmt.Sprintf(val.String(), args...))
	case reflect.Map:
		m, ok := val.Interface().(map[string]interface{})
		if ok {
			db.sentence.where.andMap(m)
		}
	case reflect.Struct:
		db.sentence.where.andTable(val)
		db.mod = where
	}
	return db
}
