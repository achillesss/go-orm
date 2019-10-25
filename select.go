package orm

func (h *sqlHead) select_(columns ...string) {
	h.option = optionSelect
	if columns == nil {
		return
	}

	for _, column := range columns {
		h.columns = append(h.columns, column)
	}
}

func (q *sqlSentence) select_(columns ...string) *sqlSentence {
	q.head.select_(columns...)
	return q

}

func (db *DB) select_(columns ...string) *DB {
	db.sentence.select_(columns...)
	return db
}

// Select *: none input
// Select table: struct{}, &struct{}
// Select columns: []string
func (db *DB) Select(columns ...string) *DB {
	return db.select_(columns...)
}
