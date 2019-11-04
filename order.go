package orm

func newOrderBy(column string, isAsc bool) *sqlOrder {
	var s sqlOrder
	s.column = column
	s.isAsc = isAsc
	return &s
}

func (s *sqlSentence) orderByColumn(column string, isAsc bool) {
	var order = newOrderBy(column, isAsc)
	s.orderBy = append(s.orderBy, order)
}

func (s *sqlSentence) orderByMap(orders map[string]bool, columns ...string) {
	for _, column := range columns {
		s.orderByColumn(column, orders[column])
	}
}

func (s *sqlSentence) orderBySlice(columns []string, isAsc bool) {
	for _, column := range columns {
		s.orderByColumn(column, isAsc)
	}
}

func (db *DB) orderBy(columns interface{}, args ...interface{}) *DB {
	switch c := columns.(type) {
	case string:
		var isAsc bool
		if len(args) > 0 {
			isAsc, _ = args[0].(bool)
		}
		db.sentence.orderByColumn(c, isAsc)

	case []string:
		var isAsc bool
		if len(args) > 0 {
			isAsc, _ = args[0].(bool)
		}
		db.sentence.orderBySlice(c, isAsc)

	case map[string]bool:
		var columns []string
		for i, arg := range args {
			switch cs := arg.(type) {
			case string:
				columns = append(columns, cs)

			case []string:
				if i == 0 {
					db.sentence.orderByMap(c, cs...)
					return db
				}

			default:
			}
		}

		db.sentence.orderByMap(c, columns...)
	}

	return db
}
