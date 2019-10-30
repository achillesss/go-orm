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

func (s *sqlSentence) and(where interface{}, args ...interface{}) {
	if s.where == nil {
		s.where = new(joinSquel)
	}

	var val = reflect.Indirect(reflect.ValueOf(where))
	switch val.Kind() {
	case reflect.String:
		s.where.addAndRaw(fmt.Sprintf(val.String(), args...))

	case reflect.Map:
		m, ok := val.Interface().(map[string]interface{})
		if ok {
			s.where.andMap(m)
		}

	case reflect.Struct:
		s.where.andTable(val)
		s.mod = where
	}
}

func (db *DB) and(where interface{}, args ...interface{}) *DB {
	var d = db.copy()
	d.sentence.and(where, args...)
	return d
}
