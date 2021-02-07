package rushia

import (
	"reflect"
)

// Copy creates a copy of the current query,
// so you are able to make changes and it won't modify the original query.
func (q *Query) Copy() *Query {
	a := *q
	b := a
	return &b
}

// Insert creates a `INSERT INTO` query with specified data.
// It inserts a data to the database.
func (q *Query) Insert(v interface{}) *Query {
	q.typ = queryTypeInsert
	q.data = v
	return q
}

// Replace creates a `REPLACE INTO` query with specified data.
// It deletes the original data and creates a new one instead, prettry dangerous if the data contains a foreign key.
func (q *Query) Replace(v interface{}) *Query {
	q.typ = queryTypeReplace
	q.data = v
	return q
}

// Update creates a `UPDATE` query with specified data.
// It updates the data with new data, normally use with `WHERE` condition.
func (q *Query) Update(v interface{}) *Query {
	q.typ = queryTypeUpdate
	q.data = v
	return q
}

// Select creates a `SELECT` query with specified columns, can be empty for select everything (`*`).
// It fetches the data from database.
func (q *Query) Select(columns ...interface{}) *Query {
	q.typ = queryTypeSelect
	q.selects = columns
	return q
}

// SelectOne works the same as `Select` but returns only one row as result.
// It's the combination of `.Limit(1).Select()`.
func (q *Query) SelectOne(columns ...interface{}) *Query {
	q.Limit(1)
	q.typ = queryTypeSelect
	q.selects = columns
	return q
}

// Patch works the same as `Update` but ignores the zero value.
// The zero value fields won't be updated unless it's in exclude list, to define the list, call `Exclude`.
func (q *Query) Patch(v interface{}) *Query {
	q.typ = queryTypePatch
	q.data = v
	return q
}

// Exists creates a `SELECT EXISTS` query, returns a result if the query does match a row.
func (q *Query) Exists() *Query {
	q.typ = queryTypeExists
	return q
}

// InsertSelect creates a `INSERT SELECT` query, it works a bit like table copy.
// The insert data is from another selection, pass a `SELECT` query to the first argument.
func (q *Query) InsertSelect(qu *Query, columns ...interface{}) *Query {
	q.typ = queryTypeInsertSelect
	q.subQuery = qu
	q.selects = columns
	return q
}

// Delete creates a `DELETE` query to delete the data.
// Make sure you are using it with `WHERE` condition to not delete all the data.
func (q *Query) Delete() *Query {
	q.typ = queryTypeDelete
	return q
}

// Omit omits specified fields in the data so it won't be insert/update into the database.
func (q *Query) Omit(fields ...string) *Query {
	q.omits = append(q.omits, fields...)
	return q
}

// OnDuplicate creates `ON DUPLICATE KEY UPDATE` query, works when inserting a duplicated data,
// the data will be automatically updated to the new value.
func (q *Query) OnDuplicate(v H) *Query {
	if q.duplicate == nil {
		q.duplicate = make(H)
	}
	for k, j := range v {
		q.duplicate[k] = j
	}
	return q
}

// Exclude excludes the specified fields, data types while patching with `Patch` method.
// Pass string values as field names, and `reflect.Kind` as data types to exclude.
// While patching, all the zero values will be ignored unless it's in the exclude list.
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

// Limit creates `LIMIT` option to the query.
func (q *Query) Limit(from int, count ...int) *Query {
	q.limit.from = from
	if len(count) > 0 {
		q.limit.count = count[0]
	}
	return q
}

// As creates an alias for current query.
func (q *Query) As(alias string) *Query {
	q.alias = alias
	return q
}

// Offset creates `LIMIT OFFSET` option to the query.
func (q *Query) Offset(count int, offset int) *Query {
	q.offset.count = count
	q.offset.offset = offset
	return q
}

// Having creates a `HAVING` condition.
func (q *Query) Having(args ...interface{}) *Query {
	q.havings = append(q.havings, condition{
		args:      args,
		connector: connectorTypeAnd,
	})
	return q
}

// OrHaving creates a `HAVING OR` condition.
func (q *Query) OrHaving(args ...interface{}) *Query {
	q.havings = append(q.havings, condition{
		args:      args,
		connector: connectorTypeOr,
	})
	return q
}

// Where creates a `WHERE` condition.
func (q *Query) Where(args ...interface{}) *Query {
	q.wheres = append(q.wheres, condition{
		args:      args,
		connector: connectorTypeAnd,
	})
	return q
}

// OrWhere creates a `WHERE OR` condition.
func (q *Query) OrWhere(args ...interface{}) *Query {
	q.wheres = append(q.wheres, condition{
		args:      args,
		connector: connectorTypeOr,
	})
	return q
}

// JoinWhere creates the `AND` joining condition for latest table join.
func (q *Query) JoinWhere(args ...interface{}) *Query {
	q.joins[len(q.joins)-1].conditions = append(q.joins[len(q.joins)-1].conditions, condition{
		args:      args,
		connector: connectorTypeAnd,
	})
	return q
}

// OrJoinWhere creates the `OR` joining condition for latest table join.
func (q *Query) OrJoinWhere(args ...interface{}) *Query {
	q.joins[len(q.joins)-1].conditions = append(q.joins[len(q.joins)-1].conditions, condition{
		args:      args,
		connector: connectorTypeOr,
	})
	return q
}

// Distinct adds the `DISTINCT` option to the query.
func (q *Query) Distinct() *Query {
	q.SetQueryOption("DISTINCT")
	return q
}

// Union creates a `UNION` query that connects two tables.
// It groups the result from multiple tables but eliminates the duplicates.
func (q *Query) Union(qu *Query) *Query {
	q.unions = append(q.unions, union{
		query: qu,
	})
	return q
}

// UnionAll creates a `UNION ALL` query that connects two tables.
// It groups the result from multiple tables and it keeps the duplicates.
func (q *Query) UnionAll(qu *Query) *Query {
	q.unions = append(q.unions, union{
		query: qu,
		all:   true,
	})
	return q
}

// OrderBy creates a `ORDER BY` option to the query.
func (q *Query) OrderBy(columns ...string) *Query {
	for _, v := range columns {
		q.orders = append(q.orders, order{
			column: v,
		})
	}
	return q
}

// OrderByField creates a `ORDER BY FIELD` option to the query.
func (q *Query) OrderByField(field string, values ...interface{}) *Query {
	q.orders = append(q.orders, order{
		field:  field,
		values: values,
	})
	return q
}

// GroupBy creates a `GROUP BY` option to the query.
func (q *Query) GroupBy(columns ...string) *Query {
	q.groups = append(q.groups, columns...)
	return q
}

// CrossJoin creates a `CROSS JOIN` to join a table.
func (q *Query) CrossJoin(table interface{}, conditions ...interface{}) *Query {
	return q.putJoin(joinTypeCross, table, conditions...)
}

// LeftJoin creates a `LEFT JOIN` to join a table.
func (q *Query) LeftJoin(table interface{}, conditions ...interface{}) *Query {
	return q.putJoin(joinTypeLeft, table, conditions...)
}

// RightJoin creates a `RIGHT JOIN` to join a table.
func (q *Query) RightJoin(table interface{}, conditions ...interface{}) *Query {
	return q.putJoin(joinTypeRight, table, conditions...)
}

// InnerJoin creates a `INNER JOIN` to join a table.
func (q *Query) InnerJoin(table interface{}, conditions ...interface{}) *Query {
	return q.putJoin(joinTypeInner, table, conditions...)
}

// NaturalJoin creates a `NATURAL JOIN` to join a table.
func (q *Query) NaturalJoin(table interface{}, conditions ...interface{}) *Query {
	return q.putJoin(joinTypeNatural, table, conditions...)
}

// SetQueryOption sets the query options, and it will be automatically be appended after or before the query.
func (q *Query) SetQueryOption(option string) *Query {
	q.queryOptions = append(q.queryOptions, option)
	return q
}
