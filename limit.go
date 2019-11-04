package orm

func (s *sqlSentence) limitBy(l int) {
	s.limit = l
}

func (s *sqlSentence) offsetBy(o int) {
	s.offset = o
}

func (db *DB) limit(limit int) *DB {
	db.sentence.limitBy(limit)
	return db
}

func (db *DB) offset(offset int) *DB {
	db.sentence.offsetBy(offset)
	return db
}
