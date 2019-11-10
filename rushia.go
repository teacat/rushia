package rushia

import (
	"fmt"
	"strings"
)

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

// NewFunc 會基於參數來返回一個新的 SQL 資料庫函式，
// 這能夠當作函式放置於查詢指令中，而不會被當作普通的資料執行。
func NewFunc(query string, data ...interface{}) Function {
	return Function{
		query:  query,
		values: data,
	}
}

// NewNow 會回傳一個基於 `INTERVAL` 的 SQL 資料庫函式，
// 傳入的參數格式可以是 `+1Y`、`-2M`，同時也可以像 `Now("+1Y", "-2M")` 一樣地串連使用。
// 支援的格式為：`Y`(年)、`M`(月)、`D`(日)、`W`(星期)、`h`(小時)、`m`(分鐘)、`s`(秒數)。
func NewNow(formats ...string) Function {
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
	return NewFunc(strings.TrimSpace(query))
}
