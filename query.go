package rushia

import (
	"reflect"
)

// Copy
func (q *Query) Copy() *Query {
	a := *q
	b := a
	return &b
}

// Insert
func (q *Query) Insert(v interface{}) *Query {
	q.typ = queryTypeInsert
	q.data = v
	return q
}

// Replace
func (q *Query) Replace(v interface{}) *Query {
	q.typ = queryTypeReplace
	q.data = v
	return q
}

// Update
func (q *Query) Update(v interface{}) *Query {
	q.typ = queryTypeUpdate
	q.data = v
	return q
}

// Select
func (q *Query) Select(columns ...interface{}) *Query {
	q.typ = queryTypeSelect
	q.selects = columns
	return q
}

// SelectOne
func (q *Query) SelectOne(columns ...interface{}) *Query {
	q.Limit(1)
	q.typ = queryTypeSelect
	q.selects = columns
	return q
}

// Patch
func (q *Query) Patch(v interface{}) *Query {
	q.typ = queryTypePatch
	q.data = v
	return q
}

// Exists
func (q *Query) Exists() *Query {
	q.typ = queryTypeExists
	return q
}

// InsertSelect
func (q *Query) InsertSelect(qu *Query, columns ...interface{}) *Query {
	q.typ = queryTypeInsertSelect
	q.subQuery = qu
	q.selects = columns
	return q
}

// Delete
func (q *Query) Delete() *Query {
	q.typ = queryTypeDelete
	return q
}

// Omit
func (q *Query) Omit(fields ...string) *Query {
	q.omits = append(q.omits, fields...)
	return q
}

// OnDuplicate
func (q *Query) OnDuplicate(v H) *Query {
	if q.duplicate == nil {
		q.duplicate = make(H)
	}
	for k, j := range v {
		q.duplicate[k] = j
	}
	return q
}

// Exclude
func (q *Query) Exclude(fields ...interface{}) *Query {
	for _, v := range fields {
		switch j := v.(type) {
		case reflect.Kind:
			q.exclude.kinds = append(q.exclude.kinds, j)
		case string:
			q.exclude.fields = append(q.exclude.fields, j)
		}
	}
	return q
}

// Limit
func (q *Query) Limit(from int, count ...int) *Query {
	q.limit.from = from
	if len(count) > 0 {
		q.limit.count = count[0]
	}
	return q
}

// As
func (q *Query) As(alias string) *Query {
	q.alias = alias
	return q
}

// Offset
func (q *Query) Offset(count int, offset int) *Query {
	q.offset.count = count
	q.offset.offset = offset
	return q
}

// Having
func (q *Query) Having(args ...interface{}) *Query {
	q.havings = append(q.havings, condition{
		args:      args,
		connector: connectorTypeAnd,
	})
	return q
}

// OrHaving
func (q *Query) OrHaving(args ...interface{}) *Query {
	q.havings = append(q.havings, condition{
		args:      args,
		connector: connectorTypeOr,
	})
	return q
}

// Where
func (q *Query) Where(args ...interface{}) *Query {
	q.wheres = append(q.wheres, condition{
		args:      args,
		connector: connectorTypeAnd,
	})
	return q
}

// OrWhere
func (q *Query) OrWhere(args ...interface{}) *Query {
	q.wheres = append(q.wheres, condition{
		args:      args,
		connector: connectorTypeOr,
	})
	return q
}

// JoinWhere
func (q *Query) JoinWhere(args ...interface{}) *Query {
	q.joins[len(q.joins)-1].conditions = append(q.joins[len(q.joins)-1].conditions, condition{
		args:      args,
		connector: connectorTypeAnd,
	})
	return q
}

// OrJoinWhere
func (q *Query) OrJoinWhere(args ...interface{}) *Query {
	q.joins[len(q.joins)-1].conditions = append(q.joins[len(q.joins)-1].conditions, condition{
		args:      args,
		connector: connectorTypeOr,
	})
	return q
}

// Distinct
func (q *Query) Distinct() *Query {
	q.SetQueryOption("DISTINCT")
	return q
}

// Union
func (q *Query) Union(qu *Query) *Query {
	q.unions = append(q.unions, union{
		query: qu,
	})
	return q
}

// UnionAll
func (q *Query) UnionAll(qu *Query) *Query {
	q.unions = append(q.unions, union{
		query: qu,
		all:   true,
	})
	return q
}

// OrderBy
func (q *Query) OrderBy(columns ...string) *Query {
	for _, v := range columns {
		q.orders = append(q.orders, order{
			column: v,
		})
	}
	return q
}

// OrderByField
func (q *Query) OrderByField(field string, values ...interface{}) *Query {
	q.orders = append(q.orders, order{
		field:  field,
		values: values,
	})
	return q
}

// GroupBy
func (q *Query) GroupBy(columns ...string) *Query {
	q.groups = append(q.groups, columns...)
	return q
}

func (q *Query) putJoin(typ joinType, t interface{}, conditions ...interface{}) *Query {
	j := join{
		typ: typ,
	}
	switch v := t.(type) {
	case *Query:
		j.subQuery = v
	case string:
		j.table = v
	}
	if len(conditions) != 0 {
		j.conditions = []condition{
			{
				args: conditions,
				// It's fine to be `And` or `Or`
				// since the build doesn't build the first connector.
				connector: connectorTypeAnd,
			},
		}
	}
	q.joins = append(q.joins, j)
	return q
}

// CrossJoin
func (q *Query) CrossJoin(table interface{}, conditions ...interface{}) *Query {
	return q.putJoin(joinTypeCross, table, conditions...)
}

// LeftJoin
func (q *Query) LeftJoin(table interface{}, conditions ...interface{}) *Query {
	return q.putJoin(joinTypeLeft, table, conditions...)
}

// RightJoin
func (q *Query) RightJoin(table interface{}, conditions ...interface{}) *Query {
	return q.putJoin(joinTypeRight, table, conditions...)
}

// InnerJoin
func (q *Query) InnerJoin(table interface{}, conditions ...interface{}) *Query {
	return q.putJoin(joinTypeInner, table, conditions...)
}

// NaturalJoin
func (q *Query) NaturalJoin(table interface{}, conditions ...interface{}) *Query {
	return q.putJoin(joinTypeNatural, table, conditions...)
}

// SetQueryOption
func (q *Query) SetQueryOption(option string) *Query {
	q.queryOptions = append(q.queryOptions, option)
	return q
}
