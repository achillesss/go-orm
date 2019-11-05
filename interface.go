package orm

import "database/sql"

type SqlReadDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type SqlWriteDB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type SqlBaseDB interface {
	SqlReadDB
	SqlWriteDB
}

type SqlDB interface {
	SqlBaseDB
	Begin() (*sql.Tx, error)
	Close() error
}

type SqlTx interface {
	Commit() error
	Rollback() error
}

type SqlTxDB interface {
	SqlBaseDB
	SqlTx
}
