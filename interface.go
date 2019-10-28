package orm

import "database/sql"

type SqlReadDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type SqlWriteDB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type SqlDB interface {
	SqlReadDB
	SqlWriteDB
	Begin() (*sql.Tx, error)
}

type SqlTx interface {
	Commit() error
	Rollback() error
}

type SqlTxDB interface {
	SqlReadDB
	SqlWriteDB
	SqlTx
}
