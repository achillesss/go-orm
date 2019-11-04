package orm

func (s *sqlSentence) group(columns ...string) {
	if len(columns) == 0 {
		return
	}

	if s.groupBy == nil {
		s.groupBy = new(sqlGroup)
	}

	s.groupBy.columns = append(s.groupBy.columns, columns...)
}

func (db *DB) groupBy(columns ...string) *DB {
	db.sentence.group(columns...)
	return db
}
