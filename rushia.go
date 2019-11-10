package rushia

// NewQuery 會建立一個新的 SQL 建置工具。
func NewQuery() Query {
	return Query{}
}

// NewSubQuery 會建立一個新的子指令（Sub Query），這讓你可以將子指令傳入其他的條件式（例如：`WHERE`），
// 若欲將子指令傳入插入（Join）條件中，必須在參數指定此子指令的別名。
func NewSubQuery(alias ...string) SubQuery {
	subQuery := SubQuery{
		query: NewQuery(),
	}
	if len(alias) > 0 {
		subQuery.query.alias = alias[0]
	}
	return subQuery
}

// NewTimestamp 會建立一個新的 SQL 時間戳記輔助工具。
func NewTimestamp() Timestamp {
	return Timestamp{}
}

// NewMigration 會建立表格 SQL 輔助工具。
func NewMigration() Migration {
	return Migration{}
}
