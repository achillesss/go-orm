package orm

import "fmt"

// begin transaction
func (db *DB) begin() *DB {
	var d = db.copy()
	d.SqlTxDB, d.err = d.SqlDB.Begin()
	d.isTxOn = db.err == nil
	return d
}

func End(tx SqlTx, ok bool) error {
	if !ok {
		return tx.Rollback()
	}

	var err = tx.Commit()
	if err == nil {
		return nil
	}

	if dbConfig.handleCommitError != nil {
		dbConfig.handleCommitError(fmt.Errorf("CommitFailed: %v", err))
	}

	return tx.Rollback()
}
