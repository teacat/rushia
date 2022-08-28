package rushia

import (
	"fmt"
	"reflect"
	"regexp"
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
		return q.buildNothing()
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

func (q *Query) buildNothing() string {
	tableQuery := q.bindParam(q.table, &bindOptions{
		keepStringValue: true,
	})
	return tableQuery
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
			unionQuery += fmt.Sprintf("UNION (%s)", query)
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

func removeIndex(s []interface{}, index int) []interface{} {
	return append(s[:index], s[index+1:]...)
}

func (q *Query) buildConditions(conditions []condition) string {
	var qu string
	for i, condition := range conditions {
		// Don't apply the AND/OR connector to the first item.
		if i != 0 {
			qu += fmt.Sprintf("%s ", condition.connector.toQuery())
		}
		if len(condition.args) == 0 {
			qu += fmt.Sprintf("%s ", condition.query)
			continue
		}
		// ?
		if !strings.Contains(condition.query, "?") {
			panic("rushia: incorrect where condition usage")
		}
		if strings.Contains(condition.query, "??") {
			r := regexp.MustCompile(`(?m)(\?\?|\?)`)
			found := r.FindAllString(condition.query, -1)
			count := strings.Count(condition.query, "??")
			for i := len(found) - 1; i >= 0; i-- {
				if found[i] != "??" {
					continue
				}
				condition.query = replaceNth(condition.query, "??", fmt.Sprintf("`%s`", condition.args[i].(string)), count)
				count--
				condition.args = removeIndex(condition.args, i)
			}
		}

		//
		for argIndex, arg := range condition.args {
			//
			if v, ok := arg.(*Query); ok {
				query, params := Build(v)
				condition.query = replaceNth(condition.query, "?", query, argIndex+1)
				q.bindParams(params, nil)
				continue
			}
			//
			if reflect.TypeOf(arg).Kind() == reflect.Slice {
				var params []interface{}
				s := reflect.ValueOf(arg)
				if s.Len() == 0 {
					panic("rushia: no len slice was passed as arg, stop it before sending to rushia")
				}
				for i := 0; i < s.Len(); i++ {
					params = append(params, s.Index(i).Interface())
				}
				condition.query = replaceNth(condition.query, "?", fmt.Sprintf("(%s)", q.bindParams(params, nil)), argIndex+1)
				continue
			}
			q.bindParam(arg, nil)
		}
		qu += fmt.Sprintf("%s ", condition.query)
	}
	return q.trim(qu)
}

func (q *Query) processEscaped(qu string, args ...interface{}) (string, []interface{}) {
	if !strings.Contains(qu, "??") {
		return qu, args
	}
	r := regexp.MustCompile(`(?m)(\?\?|\?)`)
	found := r.FindAllString(qu, -1)
	count := strings.Count(qu, "??")
	for i := len(found) - 1; i >= 0; i-- {
		if found[i] != "??" {
			continue
		}
		qu = replaceNth(qu, "??", fmt.Sprintf("`%s`", args[i].(string)), count)
		count--
		args = removeIndex(args, i)
	}
	return qu, args
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
	} else if q.limit.count != 0 {
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
				query: conditions[0].(string),
				args:  conditions[1:],
				// It's fine to be `And` or `Or`
				// since the build doesn't build the first connector.
				connector: connectorTypeAnd,
			},
		}
	}
	q.joins = append(q.joins, j)
	return q
}
