package orm

import "fmt"

func (db *DB) raw(format string, args ...interface{}) *DB {
	var d = db.copy()
	d.sentence.raw_(format, args...)
	return d
}

func (h *sqlHead) raw() {
	h.option = optionRaw
}

func (s *sqlSentence) raw_(format string, args ...interface{}) {
	s.head.raw()
	s.raw = fmt.Sprintf(format, args...)
}
