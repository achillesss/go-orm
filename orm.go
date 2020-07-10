package orm

import (
	"database/sql"

	"github.com/wizhodl/go-utils/log"
	"github.com/wizhodl/go-utils/stack"
)

type DBStats struct {
	sql.DBStats
	StartCount int64
	EndCount   int64
}

type DB struct {
	SqlDB
	SqlTxDB

	// sentence gen sql sentence
	sentence *sqlSentence
	// err returns any err
	err      error
	isTxOn   bool
	txCaller string
	txUUID   string

	OriginDB   *sql.DB
	StartCount *int64
	EndCount   *int64
}

func (db *DB) copy() *DB {
	var d DB
	d.SqlDB = db.SqlDB
	d.SqlTxDB = db.SqlTxDB
	d.isTxOn = db.isTxOn
	d.StartCount = db.StartCount
	d.EndCount = db.EndCount

	if db.sentence == nil {
		d.sentence = newSentence()
	} else {
		d.sentence = db.sentence.copy()
	}

	return &d
}

func (db *DB) Table(t interface{}) *DB {
	return db.table(t)
}

// Select *: none input
// Select columns: []string
func (db *DB) Select(columns ...string) *DB {
	return db.select_(columns...)
}

func (db *DB) Where(where interface{}, args ...interface{}) *DB {
	return db.and(where, args...)
}

func (db *DB) And(where interface{}, args ...interface{}) *DB {
	return db.and(where, args...)
}

func (db *DB) Or(where interface{}, args ...interface{}) *DB {
	return db.or(where, args...)
}

func (db *DB) GroupBy(columns ...string) *DB { return db.groupBy(columns...) }

// OrderBy(column string, isAsc bool)
// OrderBy(columns []string, isAsc bool)
// OrderBy(isAscColumns map[string]bool, columns []string)
// OrderBy(isAscColumns map[string]bool, columns ...string)
func (db *DB) OrderBy(columns interface{}, args ...interface{}) *DB {
	return db.orderBy(columns, args...)
}

func (db *DB) Limit(limit int) *DB   { return db.limit(limit) }
func (db *DB) Offset(offset int) *DB { return db.offset(offset) }

// Insert(&struct{})
// Insert([]*struct{})
// Insert(format string, args ...interface{})
func (db *DB) Insert(set interface{}, args ...interface{}) *DB {
	return db.insert(set, args...)
}

// Update(&struct{}, string, string...)
// Update([]*struct{})
// Update(map[string]interface{}, string, string ...)
// Update(format string, args ...interface{})
func (db *DB) Update(set interface{}, args ...interface{}) *DB {
	return db.update(set, args...)
}

func (db *DB) Delete() *DB {
	return db.delete_()
}

func (db *DB) Raw(format string, args ...interface{}) *DB {
	return db.raw(format, args...)
}

// Do a sql query
func (db *DB) Do(any ...interface{}) error {
	return db.do(any...).err
}

// begin transaction
func (db *DB) Begin() *DB {
	return db.begin()
}

// Commit is Commit
func (db *DB) Commit() error {
	if db.isTxOn {
		db.isTxOn = false
		return db.SqlTxDB.Commit()
	}
	return nil
}

// Rollback is Rollback
func (db *DB) Rollback() error {
	if db.isTxOn {
		db.isTxOn = false
		return db.SqlTxDB.Rollback()
	}
	return nil
}

// end transaction
func (db *DB) End(ok bool) error {
	defer func() { db.isTxOn = false }()
	var err = end(db.SqlTxDB, ok)
	var caller = log.CallerLine(1)
	if dbConfig.endTxMonitor != nil {
		go func() {
			var stackKey = stack.GetStackHash()
			endTxChan <- &EndTx{
				ID:       db.txUUID,
				EndAt:    GetNowTime(),
				Error:    err,
				IsCommit: ok,
				Caller:   caller,
				StackKey: stackKey,
			}
		}()
	}
	return err
}

func (db *DB) Close() error {
	return db.OriginDB.Close()
}
