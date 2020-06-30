package orm

import (
	"database/sql"
	"time"

	"github.com/achillesss/log"
	_ "github.com/go-sql-driver/mysql"
)

func (c *connConfig) Open() (*DB, error) {
	db, err := sql.Open(c.driverName, c.loginString())
	if err != nil {
		return nil, err
	}
	dbConfig = *c
	if dbConfig.dbStatsInterval > 0 {
		var ticker = time.NewTicker(dbConfig.dbStatsInterval)
		for range ticker.C {
			log.Infofln("DB STATS: %+#v", db.Stats())
		}
	}
	return &DB{SqlDB: db, OriginDB: db}, nil
}
