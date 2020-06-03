package orm

func (db *DB) delete_() *DB {
	var d = db.copy()
	d.sentence.delete_()
	return d
}

func (q *sqlSentence) delete_() *sqlSentence {
	q.head.delete_()
	return q
}

func (h *sqlHead) delete_() {
	h.option = optionDelete
}
