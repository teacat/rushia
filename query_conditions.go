package rushia

import "fmt"

//=======================================================
// Where
//=======================================================

// WhereValue
func (q *Query) WhereValue(column string, operator string, v interface{}) *Query {
	q.Where(column, operator, v)
	return q
}

// WhereColumn
func (q *Query) WhereColumn(columnA string, operator string, columnB string) *Query {
	q.Where(NewExpr(fmt.Sprintf("%s %s %s", columnA, operator, columnB)))
	return q
}

// WhereRaw
func (q *Query) WhereRaw(query string, args ...interface{}) *Query {
	q.Where(NewExpr(query, args...))
	return q
}

// WhereNotLike
func (q *Query) WhereNotLike(column string, v interface{}) *Query {
	q.Where(column, "NOT LIKE", v)
	return q
}

// WhereLike
func (q *Query) WhereLike(column string, v interface{}) *Query {
	q.Where(column, "LIKE", v)
	return q
}

// WhereBetween
func (q *Query) WhereBetween(column string, start, end interface{}) *Query {
	q.Where(column, "BETWEEN", start, end)
	return q
}

// WhereNotBetween
func (q *Query) WhereNotBetween(column string, start, end interface{}) *Query {
	q.Where(column, "NOT BETWEEN", start, end)
	return q
}

// WhereIn
func (q *Query) WhereIn(column string, v ...interface{}) *Query {
	q.Where(append([]interface{}{column, "IN"}, v...)...)
	return q
}

// WhereNotIn
func (q *Query) WhereNotIn(column string, v ...interface{}) *Query {
	q.Where(append([]interface{}{column, "NOT IN"}, v...)...)
	return q
}

// WhereIsNull
func (q *Query) WhereIsNull(column string) *Query {
	q.Where(column, "IS", nil)
	return q
}

// WhereIsNotNull
func (q *Query) WhereIsNotNull(column string) *Query {
	q.Where(column, "IS NOT", nil)
	return q
}

// WhereExists
func (q *Query) WhereExists(qu *Query) *Query {
	q.Where("EXISTS", qu)
	return q
}

// WhereNotExists
func (q *Query) WhereNotExists(qu *Query) *Query {
	q.Where("NOT EXISTS", qu)
	return q
}

//=======================================================
// Where Or
//=======================================================

// OrWhereValue
func (q *Query) OrWhereValue(column string, operator string, v interface{}) *Query {
	q.OrWhere(column, operator, v)
	return q
}

// OrWhereColumn
func (q *Query) OrWhereColumn(columnA string, operator string, columnB string) *Query {
	q.OrWhere(NewExpr(fmt.Sprintf("%s %s %s", columnA, operator, columnB)))
	return q
}

// OrWhereRaw
func (q *Query) OrWhereRaw(query string, args ...interface{}) *Query {
	q.OrWhere(NewExpr(query, args...))
	return q
}

// OrWhereNotLike
func (q *Query) OrWhereNotLike(column string, v interface{}) *Query {
	q.OrWhere(column, "NOT LIKE", v)
	return q
}

// OrWhereLike
func (q *Query) OrWhereLike(column string, v interface{}) *Query {
	q.OrWhere(column, "LIKE", v)
	return q
}

// OrWhereBetween
func (q *Query) OrWhereBetween(column string, start, end interface{}) *Query {
	q.OrWhere(column, "BETWEEN", start, end)
	return q
}

// OrWhereNotBetween
func (q *Query) OrWhereNotBetween(column string, start, end interface{}) *Query {
	q.OrWhere(column, "NOT BETWEEN", start, end)
	return q
}

// OrWhereIn
func (q *Query) OrWhereIn(column string, v ...interface{}) *Query {
	q.OrWhere(append([]interface{}{column, "IN"}, v...)...)
	return q
}

// OrWhereNotIn
func (q *Query) OrWhereNotIn(column string, v ...interface{}) *Query {
	q.OrWhere(append([]interface{}{column, "NOT IN"}, v...)...)
	return q
}

// OrWhereIsNull
func (q *Query) OrWhereIsNull(column string) *Query {
	q.OrWhere(column, "IS", nil)
	return q
}

// OrWhereIsNotNull
func (q *Query) OrWhereIsNotNull(column string) *Query {
	q.OrWhere(column, "IS NOT", nil)
	return q
}

// OrWhereExists
func (q *Query) OrWhereExists(qu *Query) *Query {
	q.OrWhere("EXISTS", qu)
	return q
}

// OrWhereNotExists
func (q *Query) OrWhereNotExists(qu *Query) *Query {
	q.OrWhere("NOT EXISTS", qu)
	return q
}

//=======================================================
// Join
//=======================================================

// JoinValue
func (q *Query) JoinValue(column string, operator string, v interface{}) *Query {
	q.JoinWhere(column, operator, v)
	return q
}

// JoinColumn
func (q *Query) JoinColumn(columnA string, operator string, columnB string) *Query {
	q.JoinWhere(NewExpr(fmt.Sprintf("%s %s %s", columnA, operator, columnB)))
	return q
}

// JoinRaw
func (q *Query) JoinRaw(query string, args ...interface{}) *Query {
	q.JoinWhere(NewExpr(query, args...))
	return q
}

// JoinNotLike
func (q *Query) JoinNotLike(column string, v interface{}) *Query {
	q.JoinWhere(column, "NOT LIKE", v)
	return q
}

// JoinLike
func (q *Query) JoinLike(column string, v interface{}) *Query {
	q.JoinWhere(column, "LIKE", v)
	return q
}

// JoinBetween
func (q *Query) JoinBetween(column string, start, end interface{}) *Query {
	q.JoinWhere(column, "BETWEEN", start, end)
	return q
}

// JoinNotBetween
func (q *Query) JoinNotBetween(column string, start, end interface{}) *Query {
	q.JoinWhere(column, "NOT BETWEEN", start, end)
	return q
}

// JoinIn
func (q *Query) JoinIn(column string, v ...interface{}) *Query {
	q.JoinWhere(append([]interface{}{column, "IN"}, v...)...)
	return q
}

// JoinNotIn
func (q *Query) JoinNotIn(column string, v ...interface{}) *Query {
	q.JoinWhere(append([]interface{}{column, "NOT IN"}, v...)...)
	return q
}

// JoinIsNull
func (q *Query) JoinIsNull(column string) *Query {
	q.JoinWhere(column, "IS", nil)
	return q
}

// JoinIsNotNull
func (q *Query) JoinIsNotNull(column string) *Query {
	q.JoinWhere(column, "IS NOT", nil)
	return q
}

// JoinExists
func (q *Query) JoinExists(qu *Query) *Query {
	q.JoinWhere("EXISTS", qu)
	return q
}

// JoinNotExists
func (q *Query) JoinNotExists(qu *Query) *Query {
	q.JoinWhere("NOT EXISTS", qu)
	return q
}

//=======================================================
// Join Or
//=======================================================

// OrJoinValue
func (q *Query) OrJoinValue(column string, operator string, v interface{}) *Query {
	q.OrJoinWhere(column, operator, v)
	return q
}

// OrJoinColumn
func (q *Query) OrJoinColumn(columnA string, operator string, columnB string) *Query {
	q.OrJoinWhere(NewExpr(fmt.Sprintf("%s %s %s", columnA, operator, columnB)))
	return q
}

// OrJoinRaw
func (q *Query) OrJoinRaw(query string, args ...interface{}) *Query {
	q.OrJoinWhere(NewExpr(query, args...))
	return q
}

// OrJoinNotLike
func (q *Query) OrJoinNotLike(column string, v interface{}) *Query {
	q.OrJoinWhere(column, "NOT LIKE", v)
	return q
}

// OrJoinLike
func (q *Query) OrJoinLike(column string, v interface{}) *Query {
	q.OrJoinWhere(column, "LIKE", v)
	return q
}

// OrJoinBetween
func (q *Query) OrJoinBetween(column string, start, end interface{}) *Query {
	q.OrJoinWhere(column, "BETWEEN", start, end)
	return q
}

// OrJoinNotBetween
func (q *Query) OrJoinNotBetween(column string, start, end interface{}) *Query {
	q.OrJoinWhere(column, "NOT BETWEEN", start, end)
	return q
}

// OrJoinIn
func (q *Query) OrJoinIn(column string, v ...interface{}) *Query {
	q.OrJoinWhere(append([]interface{}{column, "IN"}, v...)...)
	return q
}

// OrJoinNotIn
func (q *Query) OrJoinNotIn(column string, v ...interface{}) *Query {
	q.OrJoinWhere(append([]interface{}{column, "NOT IN"}, v...)...)
	return q
}

// OrJoinIsNull
func (q *Query) OrJoinIsNull(column string) *Query {
	q.OrJoinWhere(column, "IS", nil)
	return q
}

// OrJoinIsNotNull
func (q *Query) OrJoinIsNotNull(column string) *Query {
	q.OrJoinWhere(column, "IS NOT", nil)
	return q
}

// OrJoinExists
func (q *Query) OrJoinExists(qu *Query) *Query {
	q.OrJoinWhere("EXISTS", qu)
	return q
}

// OrJoinNotExists
func (q *Query) OrJoinNotExists(qu *Query) *Query {
	q.OrJoinWhere("NOT EXISTS", qu)
	return q
}

//=======================================================
// Having
//=======================================================

// HavingValue
func (q *Query) HavingValue(column string, operator string, v interface{}) *Query {
	q.Having(column, operator, v)
	return q
}

// HavingColumn
func (q *Query) HavingColumn(columnA string, operator string, columnB string) *Query {
	q.Having(NewExpr(fmt.Sprintf("%s %s %s", columnA, operator, columnB)))
	return q
}

// HavingRaw
func (q *Query) HavingRaw(query string, args ...interface{}) *Query {
	q.Having(NewExpr(query, args...))
	return q
}

// HavingNotLike
func (q *Query) HavingNotLike(column string, v interface{}) *Query {
	q.Having(column, "NOT LIKE", v)
	return q
}

// HavingLike
func (q *Query) HavingLike(column string, v interface{}) *Query {
	q.Having(column, "LIKE", v)
	return q
}

// HavingBetween
func (q *Query) HavingBetween(column string, start, end interface{}) *Query {
	q.Having(column, "BETWEEN", start, end)
	return q
}

// HavingNotBetween
func (q *Query) HavingNotBetween(column string, start, end interface{}) *Query {
	q.Having(column, "NOT BETWEEN", start, end)
	return q
}

// HavingIn
func (q *Query) HavingIn(column string, v ...interface{}) *Query {
	q.Having(append([]interface{}{column, "IN"}, v...)...)
	return q
}

// HavingNotIn
func (q *Query) HavingNotIn(column string, v ...interface{}) *Query {
	q.Having(append([]interface{}{column, "NOT IN"}, v...)...)
	return q
}

// HavingIsNull
func (q *Query) HavingIsNull(column string) *Query {
	q.Having(column, "IS", nil)
	return q
}

// HavingIsNotNull
func (q *Query) HavingIsNotNull(column string) *Query {
	q.Having(column, "IS NOT", nil)
	return q
}

// HavingExists
func (q *Query) HavingExists(qu *Query) *Query {
	q.Having("EXISTS", qu)
	return q
}

// HavingNotExists
func (q *Query) HavingNotExists(qu *Query) *Query {
	q.Having("NOT EXISTS", qu)
	return q
}

//=======================================================
// Having Or
//=======================================================

// OrHavingValue
func (q *Query) OrHavingValue(column string, operator string, v interface{}) *Query {
	q.OrHaving(column, operator, v)
	return q
}

// OrHavingColumn
func (q *Query) OrHavingColumn(columnA string, operator string, columnB string) *Query {
	q.OrHaving(NewExpr(fmt.Sprintf("%s %s %s", columnA, operator, columnB)))
	return q
}

// OrHavingRaw
func (q *Query) OrHavingRaw(query string, args ...interface{}) *Query {
	q.OrHaving(NewExpr(query, args...))
	return q
}

// OrHavingNotLike
func (q *Query) OrHavingNotLike(column string, v interface{}) *Query {
	q.OrHaving(column, "NOT LIKE", v)
	return q
}

// OrHavingLike
func (q *Query) OrHavingLike(column string, v interface{}) *Query {
	q.OrHaving(column, "LIKE", v)
	return q
}

// OrHavingBetween
func (q *Query) OrHavingBetween(column string, start, end interface{}) *Query {
	q.OrHaving(column, "BETWEEN", start, end)
	return q
}

// OrHavingNotBetween
func (q *Query) OrHavingNotBetween(column string, start, end interface{}) *Query {
	q.OrHaving(column, "NOT BETWEEN", start, end)
	return q
}

// OrHavingIn
func (q *Query) OrHavingIn(column string, v ...interface{}) *Query {
	q.OrHaving(append([]interface{}{column, "IN"}, v...)...)
	return q
}

// OrHavingNotIn
func (q *Query) OrHavingNotIn(column string, v ...interface{}) *Query {
	q.OrHaving(append([]interface{}{column, "NOT IN"}, v...)...)
	return q
}

// OrHavingIsNull
func (q *Query) OrHavingIsNull(column string) *Query {
	q.OrHaving(column, "IS", nil)
	return q
}

// OrHavingIsNotNull
func (q *Query) OrHavingIsNotNull(column string) *Query {
	q.OrHaving(column, "IS NOT", nil)
	return q
}

// OrHavingExists
func (q *Query) OrHavingExists(qu *Query) *Query {
	q.OrHaving("EXISTS", qu)
	return q
}

// OrHavingNotExists
func (q *Query) OrHavingNotExists(qu *Query) *Query {
	q.OrHaving("NOT EXISTS", qu)
	return q
}
