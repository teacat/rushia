package rushia

// SubQuery 是單個子指令，任何的變更都會回傳一份複製子指令來避免多個 Goroutine 編輯同個子指令指標建構體。
type SubQuery struct {
	query Query
}

//=======================================================
// 輸出函式
//=======================================================

// Table 能夠指定資料表格的名稱。
func (s SubQuery) Table(tableName ...string) SubQuery {
	s.query = s.query.Table(tableName...)
	return s
}

//=======================================================
// 選擇函式
//=======================================================

// Get 會取得多列的資料結果，傳入的參數為欲取得的欄位名稱，不傳入參數表示取得所有欄位。
func (s SubQuery) Get(columns ...string) SubQuery {
	s.query.query, s.query.params = s.query.Get(columns...)
	return s
}

//=======================================================
// 限制函式
//=======================================================

// Limit 能夠在 SQL 查詢指令中建立限制筆數的條件。
func (s SubQuery) Limit(from int, count ...int) SubQuery {
	s.query = s.query.Limit(from, count...)
	return s
}

// OrderBy 會依照指定的欄位來替結果做出排序（例如：`DESC`、`ASC`）。
func (s SubQuery) OrderBy(column string, args ...interface{}) SubQuery {
	s.query = s.query.OrderBy(column, args...)
	return s
}

// GroupBy 會在執行 SQL 指令時依照特定的欄位來做執行區分。
func (s SubQuery) GroupBy(columns ...string) SubQuery {
	s.query = s.query.GroupBy(columns...)
	return s
}

//=======================================================
// 指令函式
//=======================================================

// RawQuery 會接收傳入的變數來執行傳入的 SQL 執行語句，變數可以在語句中以 `?`（Prepared Statements）使用來避免 SQL 注入攻擊。
// 這會將多筆資料映射到本地的建構體切片、陣列上。
func (s SubQuery) RawQuery(query string, values ...interface{}) SubQuery {
	s.query.query, s.query.params = s.query.RawQuery(query, values...)
	return s
}

//=======================================================
// 條件函式
//=======================================================

// Where 會增加一個 `WHERE AND` 條件式。
func (s SubQuery) Where(args ...interface{}) SubQuery {
	s.query = s.query.Where(args...)
	return s
}

// OrWhere 會增加一個 `WHERE OR` 條件式。
func (s SubQuery) OrWhere(args ...interface{}) SubQuery {
	s.query = s.query.OrWhere(args...)
	return s
}

// Having 會增加一個 `HAVING AND` 條件式。
func (s SubQuery) Having(args ...interface{}) SubQuery {
	s.query = s.query.Having(args...)
	return s
}

// OrHaving 會增加一個 `HAVING OR` 條件式。
func (s SubQuery) OrHaving(args ...interface{}) SubQuery {
	s.query = s.query.OrHaving(args...)
	return s
}

//=======================================================
// 加入函式
//=======================================================

// LeftJoin 會向左插入一個資料表格。
func (s SubQuery) LeftJoin(table interface{}, condition string) SubQuery {
	s.query = s.query.LeftJoin(table, condition)
	return s
}

// RightJoin 會向右插入一個資料表格。
func (s SubQuery) RightJoin(table interface{}, condition string) SubQuery {
	s.query = s.query.RightJoin(table, condition)
	return s
}

// InnerJoin 會內部插入一個資料表格。
func (s SubQuery) InnerJoin(table interface{}, condition string) SubQuery {
	s.query = s.query.InnerJoin(table, condition)
	return s
}

// NaturalJoin 會自然插入一個資料表格。
func (s SubQuery) NaturalJoin(table interface{}, condition string) SubQuery {
	s.query = s.query.NaturalJoin(table, condition)
	return s
}

// JoinWhere 能夠建立一個基於 `WHERE AND` 的條件式給某個指定的插入資料表格。
func (s SubQuery) JoinWhere(table interface{}, args ...interface{}) SubQuery {
	s.query = s.query.JoinWhere(table, args...)
	return s
}

// JoinOrWhere 能夠建立一個基於 `WHERE OR` 的條件式給某個指定的插入資料表格。
func (s SubQuery) JoinOrWhere(table interface{}, args ...interface{}) SubQuery {
	s.query = s.query.JoinOrWhere(table, args...)
	return s
}
