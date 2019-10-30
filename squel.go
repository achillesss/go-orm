package orm

import (
	"fmt"
	"strings"
)

const (
	optionNone = iota
	optionSelect
	optionInsert
	optionUpdate
)

type sqlHead struct {
	option int
	fields []string
}

func (h *sqlHead) String() string {
	switch h.option {
	case optionSelect:
		var fields = "*"
		if h.fields != nil {
			fields = strings.Join(h.fields, ",")
		}
		return fmt.Sprintf("SELECT %s FROM", fields)

	case optionInsert:
		return "INSERT INTO"

	case optionUpdate:
		return "UPDATE"

	default:
		return ""
	}
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
	// table
	mod interface{}

	head      sqlHead
	tableName string

	// insert
	values *sqlValues
	value  *sqlValue

	where   *joinSquel
	groupBy *sqlGroup
	orderBy sqlOrders
	offset  int
	limit   int

	raw string
}

func (s *sqlSentence) copy() *sqlSentence {
	var sen = *s
	return &sen
}

func newSentence() *sqlSentence {
	var s sqlSentence
	return &s
}

func limitSquel(limit int) string {
	return fmt.Sprintf("LIMIT %d", limit)
}

func offsetSquel(offset int) string {
	return fmt.Sprintf("OFFSET %d", offset)
}

func (q *sqlSentence) String() string {
	if q.raw != "" {
		return q.raw
	}

	var sentenceSlice []string
	sentenceSlice = append(sentenceSlice, q.head.String())
	sentenceSlice = append(sentenceSlice, convertToSqlColumn(q.tableName))

	switch q.head.option {
	case optionUpdate:
		sentenceSlice = append(sentenceSlice, "SET")

		if q.value != nil {
			sentenceSlice = append(sentenceSlice, q.value.UpdateString())
		}

		var w = q.where.String()
		if q.where != nil && w != "" {
			sentenceSlice = append(sentenceSlice, "WHERE", w)
		}

	case optionInsert:
		if q.values != nil {
			sentenceSlice = append(sentenceSlice, q.values.String())
		}

		if q.value != nil {
			sentenceSlice = append(sentenceSlice, q.value.InsertString())
		}

	case optionSelect:
		var w = q.where.String()
		if q.where != nil && w != "" {
			sentenceSlice = append(sentenceSlice, "WHERE", w)
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
	}

	q.raw = strings.Join(sentenceSlice, " ") + ";"
	return q.raw
}
