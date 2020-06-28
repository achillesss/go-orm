package orm

import "github.com/wizhodl/go-utils/log"

// begin transaction
func (db *DB) begin() *DB {
	var d = db.copy()
	d.SqlTxDB, d.err = d.SqlDB.Begin()
	d.isTxOn = db.err == nil
	return d
}

func End(tx SqlTx, ok bool) error {
	if ok {
		if err := tx.Commit(); err != nil {
			if dbConfig.sentryCaptureMessageFunc != nil {
				dbConfig.sentryCaptureMessageFunc(fmt.Sprintf("CommitFailed: %v", err))
			}
			return tx.Rollback()
		}
	}

	return tx.Rollback()
}
