package orm

func (h *sqlHead) select_(columns ...string) {
	h.option = optionSelect
	if columns == nil {
		return
	}

	h.fields = columns
}

func (q *sqlSentence) select_(columns ...string) *sqlSentence {
	q.head.select_(columns...)
	return q

}

func (db *DB) select_(columns ...string) *DB {
	db.sentence.select_(columns...)
	return db
}
