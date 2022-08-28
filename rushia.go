package rushia

import (
	"fmt"
	"reflect"
	"strings"
)

// Expr
type Expr struct {
	rawQuery string
	params   []interface{}
}

// H
type H map[string]interface{}

const (
	queryTypeUnknown queryType = iota
	queryTypeInsert
	queryTypeReplace
	queryTypeUpdate
	queryTypeSelect
	queryTypePatch
	queryTypeExists
	queryTypeInsertSelect
	queryTypeRawQuery
	queryTypeSubQuery
	queryTypeDelete
)

type queryType int

const (
	connectorTypeAnd connectorType = iota
	connectorTypeOr
)

type connectorType int

func (t connectorType) toQuery() string {
	switch t {
	case connectorTypeOr:
		return "OR"
	default:
		return "AND"
	}
}

const (
	joinTypeLeft joinType = iota
	joinTypeRight
	joinTypeInner
	joinTypeNatural
	joinTypeCross
)

type joinType int

func (t joinType) toQuery() string {
	switch t {
	case joinTypeRight:
		return "RIGHT JOIN"
	case joinTypeInner:
		return "INNER JOIN"
	case joinTypeNatural:
		return "NATURAL JOIN"
	case joinTypeCross:
		return "CROSS JOIN"
	default:
		return "LEFT JOIN"
	}
}

const (
	insertTypeInsert insertType = iota
	insertTypeReplace
)

type insertType int

func (t insertType) toQuery() string {
	switch t {
	case insertTypeReplace:
		return "REPLACE"
	default:
		return "INSERT"
	}
}

type condition struct {
	query     string
	args      []interface{}
	connector connectorType
}

type join struct {
	table    string
	subQuery *Query
	typ      joinType

	conditions []condition
}

type limit struct {
	from  int
	count int
}

type offset struct {
	count  int
	offset int
}

type order struct {
	field  string
	values []interface{}

	column string
	// sort orderSortType
}

type exclude struct {
	kinds  []reflect.Kind
	fields []string
}

type union struct {
	all   bool
	query *Query
}

// Query
type Query struct {
	alias string

	typ      queryType
	subQuery *Query

	table        interface{}
	wheres       []condition
	havings      []condition
	queryOptions []string

	unions []union

	data interface{}

	selects []interface{}

	joins     []join
	duplicate H

	limit  limit
	offset offset

	orders []order

	groups []string

	rawQuery string
	params   []interface{}

	omits   []string
	exclude exclude
}

// NewQuery creates a Query based on a table name or a sub query.
func NewQuery(table interface{}) *Query {
	q := &Query{
		table: table,
	}
	return q
}

// NewRawQuery creates a Query based on the passed in raw query and the parameters.
func NewRawQuery(q string, params ...interface{}) *Query {
	if strings.Contains(q, "??") {
		panic("rushia: raw query doesn't support escape ?? sign yet")
	}
	return &Query{
		typ:      queryTypeRawQuery,
		rawQuery: q,
		params:   params,
	}
}

// NewExpr creates an Expression that accepts raw query and the parameters. Could be useful as the value if you are representing a complex query.
func NewExpr(query string, params ...interface{}) *Expr {
	return &Expr{
		rawQuery: query,
		params:   params,
	}
}

// NewAlias creates an alias for a table.
func NewAlias(table string, alias string) string {
	return fmt.Sprintf("%s AS %s", table, alias)
}

// Build builds the Query.
func Build(q *Query) (query string, params []interface{}) {
	query += q.padSpace(q.buildQuery())
	if q.typ == queryTypeRawQuery || q.typ == queryTypeExists {
		return q.trim(query), q.params
	}
	query += q.padSpace(q.buildAs())
	query += q.padSpace(q.buildDuplicate())
	query += q.padSpace(q.buildUnion())
	query += q.padSpace(q.buildJoin())
	query += q.padSpace(q.buildWhere())
	query += q.padSpace(q.buildHaving())
	query += q.padSpace(q.buildOrderBy())
	query += q.padSpace(q.buildGroupBy())
	query += q.padSpace(q.buildLimit())
	query += q.padSpace(q.buildOffset())
	query += q.padSpace(q.buildAfterQueryOptions())
	return q.trim(query), q.params
}
