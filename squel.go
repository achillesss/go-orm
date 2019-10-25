package orm

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	optionNone = iota
	optionSelect
	optionInsert
	optionUpdate
)

type sqlHead struct {
	option    int
	columns   []string
	receivers []reflect.Type
}

func (h *sqlHead) String() string {
	switch h.option {
	case optionSelect:
		var columns = "*"
		if h.columns != nil {
			columns = strings.Join(h.columns, ",")
		}
		return fmt.Sprintf("SELECT %s FROM", columns)
	case optionInsert:
		return "INSERT INTO"
	case optionUpdate:
		return "UPDATE"
	default:
		return ""
	}
}

func (h *sqlHead) newReceivers() []interface{} {
	var receivers []interface{}
	for _, typ := range h.receivers {
		receivers = append(receivers, reflect.New(typ).Interface())
	}
	return receivers
}

type sqlOrder struct {
	column string
	isAsc  bool
}

func (s *sqlOrder) String() string {
	var order = "ASC"
	if !s.isAsc {
		order = "DESC"
	}
	return strings.Join([]string{convertToSqlColumn(s.column), order}, " ")
}

type sqlOrders []*sqlOrder

func (s sqlOrders) String() string {
	var os []string
	for _, o := range s {
		os = append(os, o.String())
	}
	return strings.Join(os, ",")
}

type sqlGroup struct {
	column []string
}

func (s *sqlGroup) String() string {
	return fmt.Sprintf("GROUP BY %s", strings.Join(convertToSqlColumns(s.column), ","))
}

type sqlSentence struct {
	head      sqlHead
	tableName string
	where     *joinSquel
	groupBy   *sqlGroup
	orderBy   sqlOrders
	offset    int
	limit     int
}

func (q *sqlSentence) copy() *sqlSentence {
	var s sqlSentence
	s.head = q.head
	s.tableName = q.tableName
	return &s
}

func limitSquel(limit int) string {
	return fmt.Sprintf("LIMIT %d", limit)
}

func offsetSquel(offset int) string {
	return fmt.Sprintf("OFFSET %d", offset)
}

func (q *sqlSentence) String() string {
	var sentenceSlice []string
	sentenceSlice = append(sentenceSlice, q.head.String())
	sentenceSlice = append(sentenceSlice, q.tableName)

	if q.where != nil {
		sentenceSlice = append(sentenceSlice, "WHERE", q.where.String())
	}

	if q.groupBy != nil {
		sentenceSlice = append(sentenceSlice, q.groupBy.String())
	}

	if q.orderBy != nil {
		sentenceSlice = append(sentenceSlice, q.orderBy.String())
	}

	if q.limit != 0 {
		sentenceSlice = append(sentenceSlice, limitSquel(q.limit))
	}

	if q.offset != 0 {
		sentenceSlice = append(sentenceSlice, offsetSquel(q.offset))
	}

	return strings.Join(sentenceSlice, " ") + ";"
}
