package orm

import (
	"fmt"
	"time"
)

var dbConfig connConfig

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
	debugOn            bool
}

type ConnOption interface {
	runUpdate(*connConfig)
}

type optionHolder struct {
	update func(*connConfig)
}

func (o *optionHolder) runUpdate(option *connConfig) {
	o.update(option)
}

func newOptionHolder(f func(*connConfig)) *optionHolder {
	return &optionHolder{
		update: f,
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

// time.Minute by default
func WithDebug(debugOn bool) ConnOption {
	return newOptionHolder(func(o *connConfig) {
		o.debugOn = debugOn
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
		option.runUpdate(conf)
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
