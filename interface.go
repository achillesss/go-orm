package orm

import (
	"database/sql"
	"reflect"
)

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

// convert slice to []interface{}
func ConvertToInterfaceSlice(src interface{}) []interface{} {
	var val = reflect.ValueOf(src)
	if val.Kind() != reflect.Slice {
		return nil
	}

	var res []interface{}
	for n := 0; n < val.Len(); n++ {
		res = append(res, val.Index(n).Interface())
	}

	return res
}
