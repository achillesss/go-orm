package orm

import (
	"fmt"

	"github.com/google/uuid"
)

// begin transaction
func (db *DB) begin() *DB {
	var d = db.copy()
	d.SqlTxDB, d.err = d.SqlDB.Begin()
	d.isTxOn = db.err == nil
	if dbConfig.beginTxMonitor != nil {
		go func() {
			d.txUUID = uuid.New().String()
			var now = GetNowTime()
			beginTxChan <- &BeginTx{
				ID:      d.txUUID,
				BeginAt: now,
			}
		}()
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
