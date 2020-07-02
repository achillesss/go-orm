package orm

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var startQueryChan chan *StartQuery
var endQueryChan chan *EndQuery
var beginTxChan chan *BeginTx
var endTxChan chan *EndTx

func (c *connConfig) Open() (*DB, error) {
	db, err := sql.Open(c.driverName, c.loginString())
	if err != nil {
		return nil, err
	}
	dbConfig = *c

	var ormDB DB
	ormDB.OriginDB = db
	ormDB.SqlDB = db
	ormDB.StartCount = new(int64)
	ormDB.EndCount = new(int64)

	if dbConfig.dbStatsMonitor != nil {
		var getDBStatsFunc = func() DBStats {
			var stats DBStats
			stats.DBStats = db.Stats()
			stats.StartCount = *ormDB.StartCount
			stats.EndCount = *ormDB.EndCount
			return stats
		}

		go dbConfig.dbStatsMonitor(getDBStatsFunc)
	}

	if dbConfig.startQueryMonitor != nil {
		startQueryChan = make(chan *StartQuery)
		go dbConfig.startQueryMonitor(startQueryChan)
	}

	if dbConfig.endQueryMonitor != nil {
		endQueryChan = make(chan *EndQuery)
		go dbConfig.endQueryMonitor(endQueryChan)
	}

	if dbConfig.beginTxMonitor != nil {
		beginTxChan = make(chan *BeginTx)
		go dbConfig.beginTxMonitor(beginTxChan)
	}

	if dbConfig.endTxMonitor != nil {
		endTxChan = make(chan *EndTx)
		go dbConfig.endTxMonitor(endTxChan)
	}

	return &ormDB, nil
}
