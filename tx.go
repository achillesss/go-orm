package orm

// begin transaction
func (db *DB) begin() *DB {
	var d = db.copy()
	d.SqlTxDB, d.err = d.SqlDB.Begin()
	d.isTxOn = db.err == nil
	return d
}

func End(tx SqlTx, ok bool) error {
	if ok {
		return tx.Commit()
	}

	return tx.Rollback()
}
