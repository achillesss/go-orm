package orm

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wizhodl/go-utils/log"
)

func (c *connConfig) Open() (*DB, error) {
	db, err := sql.Open(c.driverName, c.loginString())
	if err != nil {
		return nil, err
	}
	dbConfig = *c
	if dbConfig.dbStatsInterval > 0 {
		var ticker = time.NewTicker(dbConfig.dbStatsInterval)
		go func() {
			for range ticker.C {
				if dbConfig.logLevel <= dbConfig.infoLevel {
					log.Infofln("DB STATS: %+#v", db.Stats())
				}
			}
		}()
	}
	return &DB{SqlDB: db, OriginDB: db}, nil
}
