package rushia

import (
	"fmt"
	"strings"
)

// Function 重現了一個像 `SHA(?)` 或 `NOW()` 的資料庫函式。
type Function struct {
	query  string
	values []interface{}
}

// condition 是一個 `WHERE` 或 `HAVING` 的條件式。
type condition struct {
	args      []interface{}
	connector string
}

// order 是個基於 `ORDER` 的排序資訊。
type order struct {
	column string
	args   []interface{}
}

// join 帶有資料表格的加入資訊。
type join struct {
	typ        string
	query      string
	table      interface{}
	condition  string
	conditions []condition
}

// Builder 是個資料庫的 SQL 指令建置系統，同時也帶有資料庫的連線資料。
type Builder struct {
	// alias 是作為子指令時所帶有的別名，這會用在子指令資料表格的加入上。
	alias string
	// destination 呈現了資料的映射目的地指針。
	destination        interface{}
	tableName          []string
	conditions         []condition
	havingConditions   []condition
	queryOptions       []string
	joins              []join
	onDuplicateColumns []string
	lastInsertIDColumn string
	limit              []int
	orders             []order
	groupBy            []string
	lockMethod         string
	query              string
	params             []interface{}

	// Timestamp 是時間戳的相關輔助函式。
	Timestamp *Timestamp
	// PageLimit 限制了一頁僅能有幾筆資料。
	PageLimit int
}

// newBuilder 會基於傳入的資料庫連線來建立一個新的 SQL 指令建置系統。
func newBuilder() Builder {
	return Builder{Timestamp: &Timestamp{}, PageLimit: 20}
}

//=======================================================
// 保存函式
//=======================================================

// saveJoin 會保存資料表格的加入資訊。
func (b Builder) saveJoin(table interface{}, typ string, condition string) Builder {
	var query string
	switch v := table.(type) {
	// 子指令加入。
	case SubQuery:
		query = v.builder.query
	// 普通的表格加入。
	case string:
		query = v
	}

	b.joins = append(b.joins, join{
		query:     query,
		typ:       typ,
		table:     table,
		condition: condition,
	})
	return b
}

// saveJoinCondition 會將資料表格的加入條件式資訊保存到指定的資料表格加入資訊中。
func (b Builder) saveJoinCondition(connector string, table interface{}, args ...interface{}) Builder {
	var query string
	switch v := table.(type) {
	// 子指令條件式。
	case SubQuery:
		query = v.builder.query
	// 普通條件式。
	case string:
		query = v
	}

	var joins []join
	for _, v := range b.joins {
		if v.query == query {
			v.conditions = append(v.conditions, condition{
				args:      args,
				connector: connector,
			})
		}
		joins = append(joins, v)
	}
	b.joins = joins
	return b
}

// saveCondition 會保存欄位的查詢條件。
func (b Builder) saveCondition(typ, connector string, args ...interface{}) Builder {
	var c condition
	c.connector = connector
	c.args = args
	if typ == "HAVING" {
		b.havingConditions = append(b.havingConditions, c)
	} else {
		b.conditions = append(b.conditions, c)
	}
	return b
}

//=======================================================
// 參數函式
//=======================================================

// bindParams 會將接收到的多個變數綁定到本次的建置工作中，並且產生、回傳相對應的 SQL 指令片段。
func (b Builder) bindParams(data interface{}) (query string) {
	switch d := data.(type) {
	case []interface{}:
		for _, v := range d {
			query += fmt.Sprintf("%s, ", b.bindParam(v))
		}
	case []int:
		for _, v := range d {
			query += fmt.Sprintf("%s, ", b.bindParam(v))
		}
	case []string:
		for _, v := range d {
			query += fmt.Sprintf("%s, ", b.bindParam(v))
		}
	}
	query = trim(query)
	return
}

// bindParam 會將單個傳入的變數綁定到本次的建置工作中，並且依照變數型態來產生並回傳相對應的 SQL 指令片段與決定是否要以括號包覆。
func (b Builder) bindParam(data interface{}, parentheses ...bool) (param string) {
	switch v := data.(type) {
	case SubQuery:
		if len(v.builder.params) > 0 {
			b.params = append(b.params, v.builder.params...)
		}
	case Function:
		if len(v.values) > 0 {
			b.params = append(b.params, v.values...)
		}
	case nil:
	case Timestamp:
		b.params = append(b.params, v.value)
	default:
		b.params = append(b.params, data)
	}
	param = b.paramToQuery(data, parentheses...)
	return
}

// paramToQuery 會將參數的變數資料型態轉換成 SQL 指令片段，並決定是否要加上括號。
func (b Builder) paramToQuery(data interface{}, parentheses ...bool) (param string) {
	switch v := data.(type) {
	case SubQuery:
		if len(parentheses) > 0 {
			if parentheses[0] == false {
				param = fmt.Sprintf("%s", v.builder.query)
			}
		} else {
			param = fmt.Sprintf("(%s)", v.builder.query)
		}
	case Function:
		param = v.query
	case nil:
		param = "NULL"
	default:
		param = "?"
	}
	return
}

//=======================================================
// 建置函式
//=======================================================

// buildWhere 會基於目前所擁有的條件式來建置一串 `WHERE` 和 `HAVING` 的 SQL 指令。
func (b Builder) buildWhere(typ string) (query string) {
	var conditions []condition
	if typ == "HAVING" {
		conditions = b.havingConditions
		query = "HAVING "
	} else {
		conditions = b.conditions
		query = "WHERE "
	}
	if len(conditions) == 0 {
		query = ""
		return
	}
	query += b.buildConditions(conditions)
	return
}

// buildUpdate 會建置 `UPDATE` 的 SQL 指令。
func (b Builder) buildUpdate(data interface{}) (query string) {
	var set string
	beforeOptions, _ := b.buildQueryOptions()
	query = fmt.Sprintf("UPDATE %s%s SET ", beforeOptions, b.tableName[0])

	switch realData := data.(type) {
	case map[string]interface{}:
		for column, value := range realData {
			set += fmt.Sprintf("%s = %s, ", column, b.bindParam(value))
		}
	}
	query += fmt.Sprintf("%s ", trim(set))
	return
}

// buildLimit 會建置 `LIMIT` 的 SQL 指令。
func (b Builder) buildLimit() (query string) {
	switch len(b.limit) {
	case 0:
		return
	case 1:
		query = fmt.Sprintf("LIMIT %d ", b.limit[0])
	case 2:
		query = fmt.Sprintf("LIMIT %d, %d ", b.limit[0], b.limit[1])
	}
	return
}

// buildSelect 會建置 `SELECT` 的 SQL 指令。
func (b Builder) buildSelect(columns ...string) (query string) {
	beforeOptions, _ := b.buildQueryOptions()

	if len(columns) == 0 {
		query = fmt.Sprintf("SELECT %s* FROM %s ", beforeOptions, b.tableName[0])
	} else {
		query = fmt.Sprintf("SELECT %s%s FROM %s ", beforeOptions, strings.Join(columns, ", "), b.tableName[0])
	}
	return
}

// buildConditions 會將傳入的條件式轉換成指定的 `WHERE` 或 `HAVING` SQL 指令。
func (b Builder) buildConditions(conditions []condition) (query string) {
	for i, v := range conditions {
		// 如果不是第一個條件式的話，那麼就增加連結語句。
		if i != 0 {
			query += fmt.Sprintf("%s ", v.connector)
		}

		// 取得欄位名稱的種類，有可能是個 SQL 指令或普通的欄位名稱、甚至是子指令。
		var typ string
		switch q := v.args[0].(type) {
		case string:
			if strings.Contains(q, "?") || strings.Contains(q, "(") || len(v.args) == 1 {
				typ = "Query"
			} else {
				typ = "Column"
			}
		case SubQuery:
			typ = "SubQuery"
		}

		// 基於種類來建置相對應的條件式。
		switch len(v.args) {
		// .Where("Column = Column")
		case 1:
			query += fmt.Sprintf("%s ", v.args[0].(string))
		// .Where("Column = ?", "Value")
		// .Where("Column", "Value")
		// .Where(subQuery, "EXISTS")
		case 2:
			switch typ {
			case "Query":
				query += fmt.Sprintf("%s ", v.args[0].(string))
				b.bindParam(v.args[1])
			case "Column":
				switch d := v.args[1].(type) {
				case Timestamp:
					query += fmt.Sprintf(d.query, v.args[0].(string), b.bindParam(d))
				default:
					query += fmt.Sprintf("%s = %s ", v.args[0].(string), b.bindParam(d))
				}
			case "SubQuery":
				query += fmt.Sprintf("%s %s ", v.args[1].(string), b.bindParam(v.args[0]))
			}
		// .Where("Column", ">", "Value")
		// .Where("Column", "IN", subQuery)
		// .Where("Column", "IS", nil)
		case 3:
			if typ == "Query" {
				query += fmt.Sprintf("%s ", v.args[0].(string))
				b.bindParams(v.args[1:])
			} else {
				if v.args[1].(string) == "IN" || v.args[1].(string) == "NOT IN" {
					query += fmt.Sprintf("%s %s (%s) ", v.args[0].(string), v.args[1].(string), b.bindParam(v.args[2], false))
				} else {
					query += fmt.Sprintf("%s %s %s ", v.args[0].(string), v.args[1].(string), b.bindParam(v.args[2]))
				}
			}
		// .Where("(Column = ? OR Column = SHA(?))", "Value", "Value")
		// .Where("Column", "BETWEEN", 1, 20)
		default:
			if typ == "Query" {
				query += fmt.Sprintf("%s ", v.args[0].(string))
				b.bindParams(v.args[1:])
			} else {
				switch v.args[1].(string) {
				case "BETWEEN", "NOT BETWEEN":
					query += fmt.Sprintf("%s %s %s AND %s ", v.args[0].(string), v.args[1].(string), b.bindParam(v.args[2]), b.bindParam(v.args[3]))
				case "IN", "NOT IN":
					query += fmt.Sprintf("%s %s (%s) ", v.args[0].(string), v.args[1].(string), b.bindParams(v.args[2:]))
				}
			}
		}
	}
	return
}

// buildDelete 會建置 `DELETE` 的 SQL 指令。
func (b Builder) buildDelete(tableNames ...string) (query string) {
	beforeOptions, _ := b.buildQueryOptions()
	query += fmt.Sprintf("DELETE %sFROM %s ", beforeOptions, strings.Join(tableNames, ", "))
	return
}

// buildQueryOptions 依照以保存的語句選項來建置執行選項的 SQL 指令片段。
// 這會回傳兩個 SQL 指令片段，分別是放在整體 SQL 指令的前面與後面。
func (b Builder) buildQueryOptions() (before string, after string) {
	for _, v := range b.queryOptions {
		switch v {
		case "ALL", "DISTINCT", "SQL_CACHE", "SQL_NO_CACHE", "DISTINCTROW", "HIGH_PRIORITY", "STRAIGHT_JOIN", "SQL_SMALL_RESULT", "SQL_BIG_RESULT", "SQL_BUFFER_RESULT", "SQL_CALC_FOUND_ROWS", "LOW_PRIORITY", "QUICK", "IGNORE", "DELAYED":
			before += fmt.Sprintf("%s, ", v)
		case "FOR UPDATE", "LOCK IN SHARE MODE":
			after += fmt.Sprintf("%s, ", v)
		}
	}
	if before != "" {
		before = fmt.Sprintf("%s ", trim(before))
	}
	if after != "" {
		after = fmt.Sprintf("%s ", trim(after))
	}
	return
}

// buildQuery 會將所有建置工作串連起來並且依序執行來建置整個可用的 SQL 指令。
func (b Builder) buildQuery() Builder {
	b.query += b.buildDuplicate()
	b.query += b.buildJoin()
	b.query += b.buildWhere("WHERE")
	b.query += b.buildWhere("HAVING")
	b.query += b.buildOrderBy()
	b.query += b.buildGroupBy()
	b.query += b.buildLimit()

	_, afterOptions := b.buildQueryOptions()
	b.query += afterOptions
	b.query = strings.TrimSpace(b.query)
	return b
}

// buildOrderBy 會基於現有的排序資料來建置 `ORDERY BY` 的 SQL 指令。
func (b Builder) buildOrderBy() (query string) {
	if len(b.orders) == 0 {
		return
	}
	query += "ORDER BY "
	for _, v := range b.orders {
		switch len(v.args) {
		// .OrderBy("RAND()")
		case 0:
			query += fmt.Sprintf("%s, ", v.column)
		// .OrderBy("ID", "ASC")
		case 1:
			query += fmt.Sprintf("%s %s, ", v.column, v.args[0])
		// .OrderBy("UserGroup", "ASC", "SuperUser", "Admin")
		default:
			query += fmt.Sprintf("FIELD (%s, %s) %s, ", v.column, b.bindParams(v.args[1:]), v.args[0])
		}
	}
	query = trim(query) + " "
	return
}

// buildGroupBy 會建置 `GROUP BY` 的 SQL 指令。
func (b Builder) buildGroupBy() (query string) {
	if len(b.groupBy) == 0 {
		return
	}
	query += "GROUP BY "
	for _, v := range b.groupBy {
		query += fmt.Sprintf("%s, ", v)
	}
	query = trim(query) + " "
	return
}

// buildDuplicate 會建置 `ON DUPLICATE KEY UPDATE` 的 SQL 指令。
func (b Builder) buildDuplicate() (query string) {
	if len(b.onDuplicateColumns) == 0 {
		return
	}
	query += "ON DUPLICATE KEY UPDATE "
	if b.lastInsertIDColumn != "" {
		query += fmt.Sprintf("%s=LAST_INSERT_ID(%s), ", b.lastInsertIDColumn, b.lastInsertIDColumn)
	}
	for _, v := range b.onDuplicateColumns {
		query += fmt.Sprintf("%s = VALUES(%s), ", v, v)
	}
	query = trim(query)
	return
}

// buildInsert 會建置 `INSERT INTO` 的 SQL 指令。
func (b Builder) buildInsert(operator string, data interface{}) (query string) {
	var columns, values string
	beforeOptions, _ := b.buildQueryOptions()

	// 會基於資料型態建置不同的指令。
	switch realData := data.(type) {
	case map[string]interface{}:
		for column, value := range realData {
			columns += fmt.Sprintf("%s, ", column)
			values += fmt.Sprintf("%s, ", b.bindParam(value))
		}
		values = fmt.Sprintf("(%s)", trim(values))

	case []map[string]interface{}:
		var columnNames []string
		// 先取得欄位的名稱，這樣才能照順序遍歷整個 `map`。
		for name := range realData[0] {
			columnNames = append(columnNames, name)
			// 先建置欄位名稱的 SQL 指令片段。
			columns += fmt.Sprintf("%s, ", name)
		}
		for _, single := range realData {
			var currentValues string
			for _, name := range columnNames {
				currentValues += fmt.Sprintf("%s, ", b.bindParam(single[name]))
			}
			values += fmt.Sprintf("(%s), ", trim(currentValues))
		}
		values = trim(values)
	}
	columns = trim(columns)
	query = fmt.Sprintf("%s %sINTO %s (%s) VALUES %s ", operator, beforeOptions, b.tableName[0], columns, values)
	return
}

// buildJoin 會建置資料表的插入 SQL 指令。
func (b Builder) buildJoin() (query string) {
	for _, v := range b.joins {
		// 插入的種類（例如：`LEFT JOIN`、`RIGHT JOIN`、`INNER JOIN`）。
		query += fmt.Sprintf("%s ", v.typ)
		switch d := v.table.(type) {
		// 子指令。
		case SubQuery:
			query += fmt.Sprintf("%s AS %s ON ", b.bindParam(d), d.builder.alias)
		// 資料表格名稱。
		case string:
			query += fmt.Sprintf("%s ON ", d)
		}

		if len(v.conditions) == 0 {
			query += fmt.Sprintf("(%s) ", v.condition)
		} else {
			conditionsQuery := strings.TrimSpace(b.buildConditions(v.conditions))
			query += fmt.Sprintf("(%s %s %s) ", v.condition, v.conditions[0].connector, conditionsQuery)
		}
	}
	return
}

//=======================================================
// 執行函式
//=======================================================

// runQuery 會以 `Query` 的方式執行建置出來的 SQL 指令。
func (b Builder) runQuery() (query string, params []interface{}) {
	b = b.buildQuery()
	query, params = b.query, b.params
	return
}

//=======================================================
// 輸出函式
//=======================================================

// Table 能夠指定資料表格的名稱。
func (b Builder) Table(tableName ...string) Builder {
	b.tableName = tableName
	return b
}

//=======================================================
// 選擇函式
//=======================================================

// Get 會取得多列的資料結果，傳入的參數為欲取得的欄位名稱，不傳入參數表示取得所有欄位。
func (b Builder) Get(columns ...string) (query string, params []interface{}) {
	b.query = b.buildSelect(columns...)
	query, params = b.runQuery()
	return
}

// GetOne 會取得僅單列的資料作為結果，傳入的參數為欲取得的欄位名稱，不傳入參數表示取得所有欄位。
// 簡單說，這就是 `.Limit(1).Get()` 的縮寫用法。
func (b Builder) GetOne(columns ...string) (query string, params []interface{}) {
	query, params = b.Limit(1).Get(columns...)
	return
}

// WithTotalCount 會在 SQL 執行指令中安插 `SQL_CALC_FOUND_ROWS` 選項，
// 如此一來就能夠在執行完 SQL 指令後取得查詢的總計行數。在不同情況下，這可能會拖低執行效能。
func (b Builder) WithTotalCount() Builder {
	return b.SetQueryOption("SQL_CALC_FOUND_ROWS")
}

//=======================================================
// 插入函式
//=======================================================

// Insert 會插入一筆新的資料。
func (b Builder) Insert(data interface{}) (query string, params []interface{}) {
	b.query = b.buildInsert("INSERT", data)
	query, params = b.runQuery()
	return
}

// InsertMulti 會一次插入多筆資料。
func (b Builder) InsertMulti(data interface{}) (query string, params []interface{}) {
	b.query = b.buildInsert("INSERT", data)
	query, params = b.runQuery()
	return
}

// Delete 會移除相符的資料列，記得用上 `Where` 條件式來避免整個資料表格被清空。
// 這很重要好嗎，因為⋯你懂的⋯。喔，不。
func (b Builder) Delete() (query string, params []interface{}) {
	b.query = b.buildDelete(b.tableName...)
	query, params = b.runQuery()
	return
}

//=======================================================
// 更新函式
//=======================================================

// Replace 基本上和 `Insert` 無異，這會在有重複資料時移除該筆資料並重新插入。
// 若無該筆資料則插入新的資料。
func (b Builder) Replace(data interface{}) (query string, params []interface{}) {
	b.query = b.buildInsert("REPLACE", data)
	query, params = b.runQuery()
	return
}

// Update 會以指定的資料來更新相對應的資料列。
func (b Builder) Update(data interface{}) (query string, params []interface{}) {
	b.query = b.buildUpdate(data)
	query, params = b.runQuery()
	return
}

// OnDuplicate 能夠指定欲更新的欄位名稱，這會在插入的資料重複時自動更新相對應的欄位。
func (b Builder) OnDuplicate(columns []string, lastInsertID ...string) Builder {
	b.onDuplicateColumns = columns
	if len(lastInsertID) != 0 {
		b.lastInsertIDColumn = lastInsertID[0]
	}
	return b
}

//=======================================================
// 限制函式
//=======================================================

// Limit 能夠在 SQL 查詢指令中建立限制筆數的條件。
func (b Builder) Limit(from int, count ...int) Builder {
	if len(count) == 0 {
		b.limit = []int{from}
	} else {
		b.limit = []int{from, count[0]}
	}
	return b
}

// OrderBy 會依照指定的欄位來替結果做出排序（例如：`DESC`、`ASC`）。
func (b Builder) OrderBy(column string, args ...interface{}) Builder {
	b.orders = append(b.orders, order{
		column: column,
		args:   args,
	})
	return b
}

// GroupBy 會在執行 SQL 指令時依照特定的欄位來做執行區分。
func (b Builder) GroupBy(columns ...string) Builder {
	b.groupBy = columns
	return b
}

//=======================================================
// 指令函式
//=======================================================

// RawQuery 會接收傳入的變數來執行傳入的 SQL 執行語句，變數可以在語句中以 `?`（Prepared Statements）使用來避免 SQL 注入攻擊。
// 這會將多筆資料映射到本地的建構體切片、陣列上。
func (b Builder) RawQuery(q string, values ...interface{}) (query string, params []interface{}) {
	b.query = q
	b.params = values
	query, params = b.runQuery()
	return
}

//=======================================================
// 條件函式
//=======================================================

// Where 會增加一個 `WHERE AND` 條件式。
func (b Builder) Where(args ...interface{}) Builder {
	b = b.saveCondition("WHERE", "AND", args...)
	return b
}

// OrWhere 會增加一個 `WHERE OR` 條件式。
func (b Builder) OrWhere(args ...interface{}) Builder {
	b = b.saveCondition("WHERE", "OR", args...)
	return b
}

// Having 會增加一個 `HAVING AND` 條件式。
func (b Builder) Having(args ...interface{}) Builder {
	b = b.saveCondition("HAVING", "AND", args...)
	return b
}

// OrHaving 會增加一個 `HAVING OR` 條件式。
func (b Builder) OrHaving(args ...interface{}) Builder {
	b = b.saveCondition("HAVING", "OR", args...)
	return b
}

//=======================================================
// 加入函式
//=======================================================

// LeftJoin 會向左插入一個資料表格。
func (b Builder) LeftJoin(table interface{}, condition string) Builder {
	b = b.saveJoin(table, "LEFT JOIN", condition)
	return b
}

// RightJoin 會向右插入一個資料表格。
func (b Builder) RightJoin(table interface{}, condition string) Builder {
	b = b.saveJoin(table, "RIGHT JOIN", condition)
	return b
}

// InnerJoin 會內部插入一個資料表格。
func (b Builder) InnerJoin(table interface{}, condition string) Builder {
	b = b.saveJoin(table, "INNER JOIN", condition)
	return b
}

// NaturalJoin 會自然插入一個資料表格。
func (b Builder) NaturalJoin(table interface{}, condition string) Builder {
	b = b.saveJoin(table, "NATURAL JOIN", condition)
	return b
}

// JoinWhere 能夠建立一個基於 `WHERE AND` 的條件式給某個指定的插入資料表格。
func (b Builder) JoinWhere(table interface{}, args ...interface{}) Builder {
	b = b.saveJoinCondition("AND", table, args...)
	return b
}

// JoinOrWhere 能夠建立一個基於 `WHERE OR` 的條件式給某個指定的插入資料表格。
func (b Builder) JoinOrWhere(table interface{}, args ...interface{}) Builder {
	b = b.saveJoinCondition("OR", table, args...)
	return b
}

// SubQuery 能夠將目前的 SQL 指令轉換為子指令（Sub Query）來防止建置後直接被執行，這讓你可以將子指令傳入其他的條件式（例如：`WHERE`），
// 若欲將子指令傳入插入（Join）條件中，必須在參數指定此子指令的別名。
func (b Builder) SubQuery(alias ...string) (subQuery SubQuery) {
	subQuery = SubQuery{
		PageLimit: b.PageLimit,
		builder:   newBuilder(),
	}
	if len(alias) > 0 {
		subQuery.builder.alias = alias[0]
	}
	return
}

//=======================================================
// 輔助函式
//=======================================================

// Func 會基於參數來返回一個新的 SQL 資料庫函式，
// 這能夠當作函式放置於查詢指令中，而不會被當作普通的資料執行。
func (b Builder) Func(query string, data ...interface{}) Function {
	return Function{
		query:  query,
		values: data,
	}
}

// Now 會回傳一個基於 `INTERVAL` 的 SQL 資料庫函式，
// 傳入的參數格式可以是 `+1Y`、`-2M`，同時也可以像 `Now("+1Y", "-2M")` 一樣地串連使用。
// 支援的格式為：`Y`(年)、`M`(月)、`D`(日)、`W`(星期)、`h`(小時)、`m`(分鐘)、`s`(秒數)。
func (b Builder) Now(formats ...string) Function {
	query := "NOW() "
	unitMap := map[string]string{
		"Y": "YEAR",
		"M": "MONTH",
		"D": "DAY",
		"W": "WEEK",
		"h": "HOUR",
		"m": "MINUTE",
		"s": "SECOND",
	}
	for _, v := range formats {
		operator := string(v[0])
		interval := v[1 : len(v)-1]
		unit := string(v[len(v)-1])
		query += fmt.Sprintf("%s INTERVAL %s %s ", operator, interval, unitMap[unit])
	}
	return b.Func(strings.TrimSpace(query))
}

// SetLockMethod 會設置鎖定資料表格的方式（例如：`WRITE`、`READ`）。
func (b Builder) SetLockMethod(method string) Builder {
	b.lockMethod = strings.ToUpper(method)
	return b
}

// Lock 會以指定的上鎖方式來鎖定某個指定的資料表格，這能用以避免資料競爭問題。
func (b Builder) Lock(tableNames ...string) (query string, params []interface{}) {
	var tables string
	for _, v := range tableNames {
		tables += fmt.Sprintf("%s %s, ", v, b.lockMethod)
	}
	tables = trim(tables)

	query, params = b.RawQuery(fmt.Sprintf("LOCK TABLES %s", tables))
	return
}

// Unlock 能解鎖已鎖上的資料表格。
func (b Builder) Unlock(tableNames ...string) (query string, params []interface{}) {
	query, params = b.RawQuery("UNLOCK TABLES")
	return
}

// SetQueryOption 會設置 SQL 指令的額外選項（例如：`SQL_NO_CACHE`）。
func (b Builder) SetQueryOption(options ...string) Builder {
	b.queryOptions = append(b.queryOptions, options...)
	return b
}

//=======================================================
// 輔助函式
//=======================================================

// trim 會清理接收到的字串，移除最後無謂的逗點與空白。
func trim(input string) (result string) {
	if len(input) == 0 {
		result = strings.TrimSpace(input)
	} else {
		result = strings.TrimSpace(input[0 : len(input)-2])
	}
	return
}
