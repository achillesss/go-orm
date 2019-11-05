package orm

import (
	"fmt"
	"reflect"
)

// or xx or xx or xx
func (j *joinSquel) addOrRaw(raw string) *joinSquel {
	j.or = append(j.or, raw)
	return j
}

func (j *joinSquel) addOr(column string, symbol string, values ...interface{}) *joinSquel {
	j.or = append(j.or, joinColumnValue(column, symbol, values...))
	return j
}

func (j *joinSquel) setSubOr(s *joinSquel) *joinSquel {
	j.subOr = s
	return j
}

func (s *joinSquel) orMap(src map[string][]interface{}, symbol string) *joinSquel {
	for k, vs := range src {
		for i := range vs {
			vs[i] = convertToSqlValue(k, vs[i])
		}
	}
	var j = newJoinSquelFromMap(src, symbol)
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

func (s *joinSquel) orDefaultColumns(columns ...string) *joinSquel {
	for _, column := range columns {
		if column == "" {
			continue
		}
		s.addOr(column, "", defaultSQLValue(column))
	}
	return s
}

func (s *sqlSentence) or(where interface{}, args ...interface{}) {
	if s.where == nil {
		s.where = new(joinSquel)
	}

	var val = reflect.Indirect(reflect.ValueOf(where))
	switch val.Kind() {
	case reflect.String:
		s.where.addOrRaw(fmt.Sprintf(val.String(), args...))

	case reflect.Map:
		var symbol string
		if len(args) != 0 {
			s, ok := args[0].(string)
			if ok {
				symbol = s
			}
		}

		switch m := val.Interface().(type) {
		case map[string]interface{}:
			var srcMap = make(map[string][]interface{})
			for k, v := range m {
				srcMap[k] = append(srcMap[k], v)
			}
			s.where.orMap(srcMap, symbol)

		case map[string][]interface{}:
			s.where.orMap(m, symbol)
		}
	case reflect.Struct:
		s.where.orTable(val)
		s.mod = where
	}
}

func (db *DB) or(where interface{}, args ...interface{}) *DB {
	var d = db.copy()
	d.sentence.or(where, args...)
	return d
}
