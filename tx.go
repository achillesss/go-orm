package orm

import (
	"fmt"

	"github.com/wizhodl/go-utils/log"
)

// begin transaction
func (db *DB) begin() *DB {
	var d = db.copy()
	d.SqlTxDB, d.err = d.SqlDB.Begin()
	d.isTxOn = db.err == nil
	if dbConfig.txUUIDFunc != nil {
		d.txUUID = dbConfig.txUUIDFunc()
		if dbConfig.logLevel <= dbConfig.infoLevel {
			log.InfoflnN(2, "SQL TX START|%s", d.txUUID)
		}
	}
	return d
}

func end(tx SqlTx, ok bool) error {
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
