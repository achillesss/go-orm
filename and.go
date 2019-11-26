package orm

import (
	"fmt"
	"reflect"
)

// and xx and xx and xx
func (j *joinSquel) addAndRaw(raw string) *joinSquel {
	if raw == "" {
		return j
	}
	j.and = append(j.and, raw)
	return j
}

func (j *joinSquel) addAnd(column string, symbol string, values ...interface{}) *joinSquel {
	j.and = append(j.and, joinColumnValue(column, symbol, values...))
	return j
}

func (j *joinSquel) setSubAnd(s *joinSquel) *joinSquel {
	j.subAnd = s
	return j
}

func (s *joinSquel) andMap(src map[string][]interface{}, symbol string) *joinSquel {
	var valMap = make(map[string][]interface{})
	for k, vs := range src {
		for i := range vs {
			valMap[k] = append(valMap[k], convertToSqlValue(k, vs[i], true))
		}
	}
	var j = newJoinSquelFromMap(valMap, symbol)
	s.addAndRaw(j.String())
	return s
}

func (s *joinSquel) andTable(src reflect.Value) *joinSquel {
	var j = newJoinSquelFromTable(src)
	s.addAndRaw(j.String())
	return s
}

func (s *joinSquel) andDefaultColumns(columns ...string) *joinSquel {
	for _, column := range columns {
		if column == "" {
			continue
		}
		s.addAnd(column, "", defaultSQLValue(column))
	}
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
			s.where.andMap(srcMap, symbol)

		case map[string][]interface{}:
			s.where.andMap(m, symbol)
		}

	case reflect.Struct:
		s.where.andTable(val)
		s.mod = where
		for i, arg := range args {
			switch c := arg.(type) {
			case string:
				s.where.andDefaultColumns(c)

			case []string:
				if i == 0 {
					s.where.andDefaultColumns(c...)
					break
				}

			default:
			}
		}
	}
}

func (db *DB) and(where interface{}, args ...interface{}) *DB {
	var d = db.copy()
	d.sentence.and(where, args...)
	return d
}
