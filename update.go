package orm

import (
	"reflect"
)

func (h *sqlHead) update() {
	h.option = optionUpdate
}

func (s *sqlValue) updateMap(m map[string]interface{}) {
	for k, v := range m {
		s.columns = append(s.columns, k)
		s.values = append(s.values, convertToSqlValue(k, v, true))
	}
}

func (s *sqlSentence) updateMap(m map[string]interface{}) {
	s.value.updateMap(m)
}

// defaultColumns
// update table set column = DEFAULT, ...

func (s *sqlValue) updateTable(table interface{}, defaultColumns ...interface{}) {
	var val = reflect.Indirect(reflect.ValueOf(table))
	s.columns, s.values = readTable(val, false)
	for _, c := range defaultColumns {
		column, ok := c.(string)
		if !ok {
			continue
		}

		s.columns = append(s.columns, column)
		s.values = append(s.values, defaultSQLValue(column))
	}
}

func (s *sqlSentence) update(set interface{}, args ...interface{}) {
	s.head.update()

	var val = reflect.Indirect(reflect.ValueOf(set))
	switch val.Kind() {
	// TODO:
	//  case reflect.Slice:
	// batch update

	case reflect.Struct:
		s.mod = set

		if s.value == nil {
			var value sqlValue
			s.value = &value
		}

		s.value.updateTable(set)

	case reflect.Map:
		if s.value == nil {
			var value sqlValue
			s.value = &value
		}

		m, ok := set.(map[string]interface{})
		if ok {
			s.value.updateMap(m)
		}
	}
}

func (db *DB) update(set interface{}, args ...interface{}) *DB {
	var d = db.copy()
	d.sentence.update(set, args...)
	return d
}

func (db *DB) doUpdate() *DB {
	var query = db.sentence.String()
	if db.isTxOn {
		_, db.err = db.SqlTxDB.Exec(query)
	} else {
		_, db.err = db.SqlDB.Exec(query)
	}
	return db
}
