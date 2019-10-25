package orm

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type connConfig struct {
	driverName         string
	user               string
	password           string
	address            string
	db                 string
	charSet            string
	parseTime          bool
	loc                string
	timeout            time.Duration
	readTimeout        time.Duration
	writeTimeout       time.Duration
	getTableNameMethod string
}

type ConnOption interface {
	update(*connConfig)
}

type optionHolder struct {
	newly func(*connConfig)
}

func (o *optionHolder) update(option *connConfig) {
	o.newly(option)
}

func newOptionHolder(f func(*connConfig)) *optionHolder {
	return &optionHolder{
		newly: f,
	}
}

// "mysql" by default
func WithDriverName(driverName string) ConnOption {
	return newOptionHolder(func(o *connConfig) {
		o.driverName = driverName
	})
}

// "" by default
func WithUser(user string) ConnOption {
	return newOptionHolder(func(o *connConfig) {
		o.user = user
	})
}

// "" by default
func WithPassword(password string) ConnOption {
	return newOptionHolder(func(o *connConfig) {
		o.password = password
	})
}

// "" by default
func WithAddress(address string) ConnOption {
	return newOptionHolder(func(o *connConfig) {
		o.address = address
	})
}

// "" by default
func WithDB(db string) ConnOption {
	return newOptionHolder(func(o *connConfig) {
		o.db = db
	})
}

// utf8mb4 by default
func WithCharSet(charSet string) ConnOption {
	return newOptionHolder(func(o *connConfig) {
		o.charSet = charSet
	})
}

// true by default
func WithParseTime(parseTime bool) ConnOption {
	return newOptionHolder(func(o *connConfig) {
		o.parseTime = parseTime
	})
}

// UTC by default
func WithLoc(loc string) ConnOption {
	return newOptionHolder(func(o *connConfig) {
		o.loc = loc
	})
}

// time.Minute by default
func WithTimeout(timeout time.Duration) ConnOption {
	return newOptionHolder(func(o *connConfig) {
		o.timeout = timeout
	})
}

// time.Minute by default
func WithReadTimeout(timeout time.Duration) ConnOption {
	return newOptionHolder(func(o *connConfig) {
		o.readTimeout = timeout
	})
}

// time.Minute by default
func WithWriteTimeout(timeout time.Duration) ConnOption {
	return newOptionHolder(func(o *connConfig) {
		o.readTimeout = timeout
	})
}

// GetTableNameMethod should be like
// func (t *AnyTable) XXX() string { return "xxx" }
// which XXX is the method name below
// method is "Name" by default
func WithGetTableNameMethod(method string) ConnOption {
	return newOptionHolder(func(o *connConfig) {
		o.getTableNameMethod = method
	})
}

// default connection configuration
func defaultConnConfig() *connConfig {
	return &connConfig{
		driverName:         "mysql",
		charSet:            "utf8mb4",
		parseTime:          true,
		loc:                "UTC",
		timeout:            time.Minute,
		readTimeout:        time.Minute,
		writeTimeout:       time.Minute,
		getTableNameMethod: "Name",
	}
}

func NewConnConfig(options ...ConnOption) *connConfig {
	var conf = defaultConnConfig()
	for _, option := range options {
		option.update(conf)
	}
	return conf
}

func (c *connConfig) loginString() string {
	var format = "%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s&timeout=%s&readTimeout=%s&writeTimeout=%s"
	return fmt.Sprintf(
		format,
		c.user,
		c.password,
		c.address,
		c.db,
		c.charSet,
		c.parseTime,
		c.loc,
		c.timeout,
		c.readTimeout,
		c.writeTimeout,
	)
}

func (c *connConfig) Open() (*DB, error) {
	db, err := sql.Open(c.driverName, c.loginString())
	if err != nil {
		return nil, err
	}
	dbConfig = *c
	return &DB{DB: db, sentence: &sqlSentence{}}, nil
}
