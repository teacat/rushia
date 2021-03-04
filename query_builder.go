package rushia

import (
	"fmt"
	"reflect"
	"strings"
)

// isOmitted searchs for the field in the query omit option.
func (q *Query) isOmitted(field string) bool {
	for _, v := range q.omits {
		if v == field {
			return true
		}
	}
	return false
}

// bindOptions is the option for different binding situations.
type bindOptions struct {
	// noParentheses decides to wrap the sub query in the parentheses or not.
	noParentheses bool
	// keepStringValue returns the original string value instead of treating it like a prepared statement.
	// usually used for column names, so it won't be convert to `?` symbol.
	keepStringValue bool
}

// bindParams loops the bindParam function for each value in the slice, and the values will be bind into the Query.
func (q *Query) bindParams(data []interface{}, options *bindOptions) string {
	var qu string
	for _, v := range data {
		qu += fmt.Sprintf("%s, ", q.bindParam(v, options))
	}
	return q.trim(qu)
}

// bindParam binds the value to the Query and returns how it should look in SQL based on it's type.
// If the value was a sub query, `bindParam` builds it and push the params from the sub query to the current query, and returns the sub query SQL.
func (q *Query) bindParam(data interface{}, options *bindOptions) string {
	switch v := data.(type) {
	case *Query:
		qu, p := Build(v)
		q.params = append(q.params, p...)
		if options != nil && options.noParentheses {
			return qu
		}
		return fmt.Sprintf("(%s)", qu)
	case *Expr:
		exprQ, exprP := buildExpr(v)
		q.params = append(q.params, exprP...)
		return exprQ
	case nil:
		return "NULL"
	case string:
		if options != nil && options.keepStringValue {
			return v
		}
		q.params = append(q.params, data)
		return "?"
	default:
		q.params = append(q.params, data)
		return "?"
	}
}

// separateStrings separates the strings with commas.
func (q *Query) separateStrings(v []string) string {
	return strings.Join(v, ", ")
}

// separateParams binds the values while separating them.
func (q *Query) separateParams(j []interface{}) string {
	var qu string
	for _, v := range j {
		qu += fmt.Sprintf("%s, ", q.bindParam(v, nil))
	}
	return q.trim(qu)
}

// separatePairs binds the values and making the key value as a pair.
func (q *Query) separatePairs(h H) string {
	var qu string
	for k, v := range h {
		qu += fmt.Sprintf("%s = %s, ", k, q.bindParam(v, nil))
	}
	return q.trim(qu)
}

// separateGroups binds the value group and wraps each group in parentheses.
func (q *Query) separateGroups(j [][]interface{}) string {
	var qu string
	for _, v := range j {
		qu += fmt.Sprintf("(%s), ", q.separateParams(v))
	}
	return q.trim(qu)
}

// padSpace adds the space in the end of the string if it was not empty.
func (q *Query) padSpace(s string) string {
	if s != "" {
		return fmt.Sprintf("%s ", s)
	}
	return s
}

//=======================================================
// Build
//=======================================================

func (q *Query) buildQuery() string {
	switch q.typ {
	case queryTypeInsert:
		return q.buildInsert(insertTypeInsert)
	case queryTypeReplace:
		return q.buildReplace()
	case queryTypeUpdate:
		return q.buildUpdate(false)
	case queryTypeSelect:
		return q.buildSelect()
	case queryTypePatch:
		return q.buildPatch()
	case queryTypeExists:
		return q.buildExists()
	case queryTypeInsertSelect:
		return q.buildInsertSelect()
	case queryTypeRawQuery:
		return q.buildRawQuery()
	case queryTypeDelete:
		return q.buildDelete()
	default:
		panic(ErrQueryTypeUnspecified)
	}
}

func buildExpr(expr *Expr) (query string, params []interface{}) {
	for i, j := range expr.params {
		switch v := j.(type) {
		case *Query:
			q, p := Build(v)
			expr.rawQuery = replaceNth(expr.rawQuery, "?", q, i+1)
			params = append(params, p...)
		default:
			params = append(params, j)
		}
	}
	query = expr.rawQuery
	return
}

func (q *Query) buildInsert(typ insertType) string {
	columns, values, _ := q.flattenData(q.data)
	insertQuery := typ.toQuery()
	beforeQuery := q.padSpace(q.trim(q.buildBeforeQueryOptions()))
	tableQuery := q.bindParam(q.table, &bindOptions{
		keepStringValue: true,
	})
	columnsQuery := q.separateStrings(columns)
	valuesQuery := q.separateGroups(values)

	return fmt.Sprintf("%s %sINTO %s (%s) VALUES %s",
		insertQuery,
		beforeQuery,
		tableQuery,
		columnsQuery,
		valuesQuery,
	)
}

func (q *Query) buildReplace() string {
	return q.buildInsert(insertTypeReplace)
}

func (q *Query) buildUpdate(isPatch bool) string {
	_, _, h := q.flattenData(q.data)
	data := h[0]
	if isPatch {
		data = q.patchH(data)
	}
	beforeQuery := q.padSpace(q.trim(q.buildBeforeQueryOptions()))
	tableQuery := q.bindParam(q.table, &bindOptions{
		keepStringValue: true,
	})
	pairsQuery := q.separatePairs(data)

	return fmt.Sprintf("UPDATE %s%s SET %s",
		beforeQuery,
		tableQuery,
		pairsQuery,
	)
}

func (q *Query) buildDelete() string {
	tableQuery := q.bindParam(q.table, &bindOptions{
		keepStringValue: true,
	})
	return fmt.Sprintf("DELETE FROM %s", tableQuery)
}

func (q *Query) buildSelect() string {
	beforeQuery := q.padSpace(q.trim(q.buildBeforeQueryOptions()))
	selectQuery := "*"
	if len(q.selects) != 0 {
		selectQuery = q.bindParams(q.selects, &bindOptions{keepStringValue: true})
	}
	tableQuery := q.bindParam(q.table, &bindOptions{keepStringValue: true})

	return fmt.Sprintf("SELECT %s%s FROM %s", beforeQuery, selectQuery, tableQuery)
}

func (q *Query) buildPatch() string {
	return q.buildUpdate(true)
}

func (q *Query) buildExists() string {
	query, params := Build(NewRawQuery("SELECT EXISTS(?)", q.Copy().Select()))
	q.params = params
	return query
}

func (q *Query) buildInsertSelect() string {
	beforeQuery := q.padSpace(q.trim(q.buildBeforeQueryOptions()))
	tableQuery := q.bindParam(q.table, &bindOptions{
		keepStringValue: true,
	})
	fieldQuery := q.bindParams(q.selects, &bindOptions{
		keepStringValue: true,
	})
	selectQuery, selectParams := Build(q.subQuery)
	q.bindParams(selectParams, nil)

	return fmt.Sprintf("INSERT %sINTO %s (%s) %s",
		beforeQuery,
		tableQuery,
		fieldQuery,
		selectQuery,
	)
}

func (q *Query) buildRawQuery() string {
	query, params := buildExpr(NewExpr(q.rawQuery, q.params...))
	q.params = params
	return query
}

func (q *Query) buildUnion() string {
	if len(q.unions) == 0 {
		return ""
	}
	var unionQuery string
	for _, v := range q.unions {
		query, params := Build(v.query)
		q.bindParams(params, nil)
		if v.all {
			unionQuery += fmt.Sprintf("UNION ALL %s", query)
		} else {
			unionQuery += fmt.Sprintf("UNION %s", query)
		}
	}
	return unionQuery
}

func (q *Query) buildAs() string {
	if q.alias == "" {
		return ""
	}
	return fmt.Sprintf("AS %s", q.alias)
}

func (q *Query) buildDuplicate() string {
	if q.duplicate == nil {
		return ""
	}
	duplicateQuery := q.separatePairs(q.duplicate)
	return fmt.Sprintf("ON DUPLICATE KEY UPDATE %s", duplicateQuery)
}

func (q *Query) buildJoin() string {
	var jqu string
	for _, v := range q.joins {
		var table string
		switch {
		// .Join(subQuery, "Column = Column")
		case v.subQuery != nil:
			table = q.bindParam(v.subQuery, nil)

		// .Join("Table", "Column = Column")
		case v.table != "":
			table = v.table
		}
		jqu += fmt.Sprintf("%s %s ON (%s) ", v.typ.toQuery(), table, q.buildConditions(v.conditions))
	}
	return q.trim(jqu)
}

func (q *Query) buildConditions(c []condition) string {
	var qu string
	for i, j := range c {
		// Don't apply the AND/OR connector to the first item.
		if i != 0 {
			qu += fmt.Sprintf("%s ", j.connector.toQuery())
		}
		// Judge the type of the first argument.
		typ := "String"
		switch v := j.args[0].(type) {
		case string:
			if strings.Contains(v, "?") || strings.Contains(v, "(") || len(j.args) == 1 {
				typ = "Query"

				// NOTE: Workaround for `.Having("Avg(Column)", "<", subQuery)`
				if len(j.args) == 3 {
					if _, ok := j.args[2].(*Query); ok {
						typ = "String"
					}
				}
			}
		case *Expr:
			typ = "Expr"
		}

		// Build the query by the amount of the arguments.
		switch len(j.args) {
		// ※ Query, String
		// .Where("Column = Column")
		// ※ Expr
		// .Where(NewExpr("Column = ?", "Value"))
		// .Where(NewExpr("EXISTS ?", subQuery))
		case 1:
			switch typ {
			case "Query", "String":
				qu += fmt.Sprintf("%s ", j.args[0].(string))
			case "Expr":
				qu += fmt.Sprintf("%s ", q.bindParam(j.args[0].(*Expr), nil))
			}

		// ※ Query
		// .Where("Column = ?", "Value")
		// ※ String
		// .Where("Column", "Value")
		// .Where("Column", NewExpr("= Column"))
		// .Where("EXISTS", subQuery)
		case 2:
			switch typ {
			case "Query":
				q.bindParam(j.args[1], nil)
				qu += fmt.Sprintf("%s ", j.args[0].(string))
			case "String":
				p := q.bindParam(j.args[1], nil)
				switch j.args[0].(string) {
				case "NOT EXISTS", "EXISTS":
					qu += fmt.Sprintf("%s %s ", j.args[0].(string), p)
				default:
					qu += fmt.Sprintf("%s = %s ", j.args[0].(string), p)
				}
			}

		// ※ Query
		// .Where("(Column = ? OR Column = ?)", "A", "B")
		// ※ String
		// .Where("Column", ">", "Value")
		// .Where("Column", ">", NewExpr("ANY (Query)"))
		// .Where("Column", "IN", subQuery)
		// .Where("Column", "IS", nil)
		// .Having("Avg(Column)", "<", subQuery)
		case 3:
			switch typ {
			case "Query":
				q.bindParams(j.args[1:], nil)
				qu += fmt.Sprintf("%s ", j.args[0].(string))
			case "String":
				switch j.args[1].(string) {
				case "IN", "NOT IN":
					qu += fmt.Sprintf("%s %s (%s) ", j.args[0].(string), j.args[1].(string), q.bindParam(j.args[2], &bindOptions{
						noParentheses: true,
					}))
				default:
					qu += fmt.Sprintf("%s %s %s ", j.args[0].(string), j.args[1].(string), q.bindParam(j.args[2], nil))
				}
			}

		// ※ Query
		// .Where("(Column = ? OR Column = ? OR Column = ?)", "Value", "Value", "Value")
		// ※ String
		// .Where("Column", "BETWEEN", 1, 20)
		// .Where("Column", "IN", 1, "foo", 20)
		default:
			switch typ {
			case "Query":
				q.bindParams(j.args[1:], nil)
				qu += fmt.Sprintf("%s ", j.args[0].(string))
			case "String":
				switch j.args[1].(string) {
				case "BETWEEN", "NOT BETWEEN":
					qu += fmt.Sprintf("%s %s %s AND %s ", j.args[0].(string), j.args[1].(string), q.bindParam(j.args[2], nil), q.bindParam(j.args[3], nil))
				case "IN", "NOT IN":
					qu += fmt.Sprintf("%s %s (%s) ", j.args[0].(string), j.args[1].(string), q.bindParams(j.args[2:], nil))
				}
			}
		}
	}
	return q.trim(qu)
}

func (q *Query) buildWhere() string {
	if len(q.wheres) == 0 {
		return ""
	}
	return fmt.Sprintf("WHERE %s", q.buildConditions(q.wheres))
}

func (q *Query) buildHaving() string {
	if len(q.havings) == 0 {
		return ""
	}
	return fmt.Sprintf("HAVING %s", q.buildConditions(q.havings))
}

func (q *Query) buildOrderBy() string {
	if len(q.orders) == 0 {
		return ""
	}
	var qu string
	for _, v := range q.orders {
		switch {
		// .OrderBy("RAND()")
		// .OrderBy("ID ASC")
		case v.column != "":
			qu += fmt.Sprintf("%s, ", v.column)

		// .OrderByField("UserGroup ASC", "SuperUser", "Admin")
		case v.field != "":
			qu += fmt.Sprintf("FIELD (%s, %s), ", v.field, q.bindParams(v.values, nil))
		}
	}
	return fmt.Sprintf("ORDER BY %s", q.trim(qu))
}

func (q *Query) buildGroupBy() string {
	if len(q.groups) == 0 {
		return ""
	}
	return fmt.Sprintf("GROUP BY %s", strings.Join(q.groups, ", "))
}

func (q *Query) buildLimit() string {
	if q.limit.from != 0 && q.limit.count == 0 {
		return fmt.Sprintf("LIMIT %d", q.limit.from)
	} else if q.limit.from != 0 && q.limit.count != 0 {
		return fmt.Sprintf("LIMIT %d, %d", q.limit.from, q.limit.count)
	} else {
		return ""
	}
}

func (q *Query) buildOffset() string {
	if q.offset.count == 0 && q.offset.offset == 0 {
		return ""
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", q.offset.count, q.offset.offset)
}

func (q *Query) buildBeforeQueryOptions() string {
	var qu string
	for _, v := range q.queryOptions {
		switch v {
		case "ALL", "DISTINCT", "SQL_CACHE", "SQL_NO_CACHE", "DISTINCTROW", "HIGH_PRIORITY", "STRAIGHT_JOIN", "SQL_SMALL_RESULT", "SQL_BIG_RESULT", "SQL_BUFFER_RESULT", "SQL_CALC_FOUND_ROWS", "LOW_PRIORITY", "QUICK", "IGNORE", "DELAYED":
			qu += fmt.Sprintf("%s, ", v)
		}
	}
	return qu
}

func (q *Query) buildAfterQueryOptions() string {
	var qu string
	for _, v := range q.queryOptions {
		switch v {
		case "FOR UPDATE", "LOCK IN SHARE MODE":
			qu += fmt.Sprintf("%s, ", v)
		}
	}
	return qu
}

//=======================================================
// Helpers
//=======================================================

// flattenData parses the `interface{}` to a flatten columns, values group pair.
// It also returns a H slice which converted from interface{}.
func (q *Query) flattenData(data interface{}) (columns []string, values [][]interface{}, h []H) {
	switch v := data.(type) {
	case H:
		var k []interface{}
		v = q.omitH(v)
		columns, k = q.flattenH(v)
		values = [][]interface{}{k}
		h = []H{v}
	case []H:
		columns, values, h = q.flattenHs(v)
	case map[string]interface{}:
		columns, values, h = q.flattenData(H(v))
	case []map[string]interface{}:
		columns, values, h = q.flattenData(q.mapsToHs(v))
	default:
		columns, values, h = q.flattenData(q.structToH(v))
	}
	return
}

// patchH eliminates the zero values of a H data,
// and it also refers to the Query exclude option.
func (q *Query) patchH(data H) H {
	for k, v := range data {
		if q.shouldEliminate(k, v) {
			delete(data, k)
		}
	}
	return data
}

// omitH omits the fields of a H data based on the Query omit option.
func (q *Query) omitH(data H) H {
	for k := range data {
		if q.isOmitted(k) {
			delete(data, k)
		}
	}
	return data
}

// flattenHs flatten a slice of H by passing it back to the `flattenData` and collects the result.
func (q *Query) flattenHs(data []H) (columns []string, values [][]interface{}, hs []H) {
	for k, j := range data {
		cols, vals, h := q.flattenData(j)
		if k == 0 {
			columns = cols
		}
		values = append(values, vals[0])
		hs = append(hs, h[0])
	}
	return
}

// flattenH flatten a H to column names, values.
func (q *Query) flattenH(data H) (columns []string, values []interface{}) {
	for k, v := range data {
		columns = append(columns, k)
		values = append(values, v)
	}
	return
}

// structToH converts a struct to H data and rename/omit it by the rushia struct tag.
func (q *Query) structToH(data interface{}) H {
	h := make(H)

	var t reflect.Type
	var v reflect.Value

	// Get the real data type underneath the pointer if the data is a pointer.
	if reflect.TypeOf(data).Kind() == reflect.Ptr {
		t = reflect.Indirect(reflect.ValueOf(data)).Type()
		v = reflect.Indirect(reflect.ValueOf(data))
	} else {
		t = reflect.TypeOf(data)
		v = reflect.ValueOf(data)
	}

	for i := 0; i < t.NumField(); i++ {
		k := t.Field(i).Name
		if name, ok := t.Field(i).Tag.Lookup("rushia"); ok {
			if name == "" || name == "-" {
				continue
			}
			k = name
		}
		h[k] = v.Field(i).Interface()
	}
	return h
}

// mapsToHs converts map slice to H slice.
func (q *Query) mapsToHs(data []map[string]interface{}) []H {
	var hs []H
	for _, j := range data {
		hs = append(hs, H(j))
	}
	return hs
}

// shouldEliminate is designed for Patch.
// Returns true if the value was a zero value to indicates the value should be skipped,
// returns false if the value was not a zero value, either the type/field name was in the exclude list.
func (q *Query) shouldEliminate(k string, v interface{}) bool {
	var isExcludedColumn bool
	for _, j := range q.exclude.fields {
		if k == j {
			isExcludedColumn = true
			break
		}
	}
	valueOf := reflect.ValueOf(v)
	var isExcludedKind bool
	kind := valueOf.Kind()
	for _, j := range q.exclude.kinds {
		if kind == j {
			isExcludedKind = true
			break
		}
	}
	return (!isExcludedColumn && !isExcludedKind) && valueOf.IsZero()
}

// trim trims the unnecessary commas in the end of the string.
func (q *Query) trim(s string) string {
	return strings.TrimRight(strings.TrimSpace(s), ",")
}

// replaceNth removes the nth repeated occurrence of the specified string,
// usually used for sub query prepared statment `?` symbol replacement.
func replaceNth(s, old, new string, n int) string {
	i := 0
	for m := 1; m <= n; m++ {
		x := strings.Index(s[i:], old)
		if x < 0 {
			break
		}
		i += x
		if m == n {
			return s[:i] + new + s[i+len(old):]
		}
		i += len(old)
	}
	return s
}

// putJoin
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
