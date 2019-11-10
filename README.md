# Reiner [![GoDoc](https://godoc.org/github.com/teacat/rushia?status.svg)](https://godoc.org/github.com/teacat/rushia) [![Coverage Status](https://coveralls.io/repos/github/teacat/rushia/badge.svg?branch=master)](https://coveralls.io/github/teacat/rushia?branch=master) [![Build Status](https://travis-ci.org/teacat/rushia.svg?branch=master)](https://travis-ci.org/teacat/rushia) [![Go Report Card](https://goreportcard.com/badge/github.com/teacat/rushia)](https://goreportcard.com/report/github.com/teacat/rushia)

一個由 [Golang](https://golang.org/) 撰寫且比起部分 [ORM](https://zh.wikipedia.org/wiki/%E5%AF%B9%E8%B1%A1%E5%85%B3%E7%B3%BB%E6%98%A0%E5%B0%84) 還要讚的 [MySQL](https://www.mysql.com/) 指令建置函式庫。彈性高、不需要建構體標籤。實際上，這就只是 [PHP-MySQLi-Database-Class](https://github.com/joshcam/PHP-MySQLi-Database-Class) 不過是用在 [Golang](https://golang.org/) 而已（但還是多了些功能）。

# 這是什麼？

萊納是一個由 [Golang](https://golang.org/) 撰寫的 [MySQL](https://www.mysql.com/) 的指令建置函式庫（不是 [ORM](https://zh.wikipedia.org/wiki/%E5%AF%B9%E8%B1%A1%E5%85%B3%E7%B3%BB%E6%98%A0%E5%B0%84)，永遠也不會是），幾乎所有東西都能操控於你手中。類似自己撰寫資料庫指令但是更簡單，JOIN 表格也變得比以前更方便了。

* 幾乎全功能的函式庫。
* 自動避免於 Goroutine 發生資料競爭的設計。
* 支援 MySQL 複寫橫向擴展機制（區分讀／寫連線）。
* 容易理解與記住、且使用方式十分簡單。
* SQL 指令建構函式。
* 資料庫表格建構協助函式。
* 可串連的使用方式。
* 支援子指令（Sub Query）。
* 支援多樣的結果物件綁定（如：建構體切片）。
* 可手動操作的交易機制（Transaction）和回溯（Rollback）功能。
* 透過預置聲明（[Prepared Statement](https://en.wikipedia.org/wiki/Prepared_statement)），99.9% 避免 SQL 注入攻擊。

# 為什麼？

[Gorm](https://github.com/jinzhu/gorm) 已經是 [Golang](https://golang.org/) 裡的 [ORM](https://zh.wikipedia.org/wiki/%E5%AF%B9%E8%B1%A1%E5%85%B3%E7%B3%BB%E6%98%A0%E5%B0%84) 典範，但實際上要操作複雜與關聯性高的 SQL 指令時並不是很合適，而 Reiner 解決了這個問題。Reiner 也試圖不要和建構體扯上關係，不希望使用者需要手動指定任何標籤在建構體中。

# 執行緒與併發安全性？

我們都知道 [Golang](https://golang.org/) 的目標就是併發程式，為了避免 Goroutine 導致資料競爭問題，Reiner 會在每有變更的時候自動複製 SQL 指令建置函式庫來避免所有併發程式共用同個 SQL 指令建置函式庫（此方式並不會使資料庫連線遞增而造成效能問題）。

在原先的舊版本中則需要手動透過 `Copy` 或 `Clone` 複製建置函式庫，這繁雜的手續正是重新設計的原因。但也因為如此，現在 SQL 指令建置函式若需要分散串接則需要重新地不斷賦值，簡單來說就是像這樣。

```go
package main

func main() {
	db, _ := reiner.New("...")

	// 在舊有的版本中原本能夠這樣直覺地分散串接一段 SQL 指令。
	// 注意！這是舊版本的做法，目前已經被廢除。
	db.Table("Users")
	if ... {
		db.Where("Username", "YamiOdymel")
	}
	if ... {
		db.Limit(1, 10)
	}
	db.Get()

	// 新的版本中因為 Reiner 會不斷地回傳一個新的建置資料函式，
	// 因此必須不斷地重新賦值。
	myDB := db.Table("Users")
	if ... {
		myDB = myDB.Where("Username", "YamiOdymel")
	}
	if ... {
		myDB = myDB.Limit(1, 10)
	}
	myDB, _ = myDB.Get()
}
```

# 效能如何？

這裡有份簡略化的[效能測試報表](https://github.com/teacat/reiner-benchmark)。目前仍會持續優化並且增加快取以避免重複建置相同指令而費時。

```
測試規格：
1.7 GHz Intel Core i7 (4650U)
8 GB 1600 MHz DDR3

插入：Dbr > SQL > SQLx > Xorm > Reiner > Gorm
BenchmarkReinerInsert-4             3000            571298 ns/op            1719 B/op         49 allocs/op
BenchmarkSQLInsert-4                3000            429340 ns/op             901 B/op         17 allocs/op
BenchmarkDbrInsert-4                5000            413442 ns/op            2210 B/op         37 allocs/op
BenchmarkSQLxInsert-4               3000            444055 ns/op             902 B/op         17 allocs/op
BenchmarkGormInsert-4               2000            776838 ns/op            5319 B/op        101 allocs/op
BenchmarkXormInsert-4               3000            562341 ns/op            2921 B/op         64 allocs/op

選擇 100 筆資料：SQL > SQLx > Dbr > Reiner > Gorm > Xorm
BenchmarkReinerSelect100-4          2000            659189 ns/op           42907 B/op       1155 allocs/op
BenchmarkSQLSelect100-4             5000            336121 ns/op           28864 B/op        723 allocs/op
BenchmarkDbrSelect100-4             3000            529430 ns/op           87496 B/op       1638 allocs/op
BenchmarkSQLxSelect100-4            3000            376810 ns/op           32368 B/op        829 allocs/op
BenchmarkGormSelect100-4            2000            726107 ns/op          209236 B/op       3870 allocs/op
BenchmarkXormSelect100-4            2000            868688 ns/op          103358 B/op       4583 allocs/op
```

# 索引

* [安裝方式](#安裝方式)
* [命名建議](#命名建議)
* [NULL 值](#null-值)
* [使用方式](#使用方式)
    * [資料庫連線](#資料庫連線)
    	* [水平擴展（讀／寫分離）](#水平擴展讀寫分離)
		* [SQL 建構模式](#sql-建構模式)
	* [資料綁定與處理](#資料綁定與處理)
		* [逐行掃描](#逐行掃描)
	* [插入](#插入)
		* [覆蓋](#覆蓋)
		* [函式](#函式)
		* [當重複時](#當重複時)
		* [多筆資料](#多筆資料)
			* [省略重複鍵名](#省略重複鍵名)
	* [筆數限制](#筆數限制)
	* [更新](#更新)
	* [選擇與取得](#選擇與取得)
		* [筆數限制](#筆數限制-1)
		* [指定欄位](#指定欄位)
		* [單行資料](#單行資料)
		* [單欄位值](#單欄位值)
		* [分頁功能](#分頁功能)
	* [執行生指令](#執行生指令)
		* [單行資料](#單行資料-1)
		* [單欄位值](#單欄位值-1)
		* [進階方式](#進階方式)
	* [條件宣告](#條件宣告)
		* [擁有](#擁有)
		* [欄位比較](#欄位比較)
		* [自訂運算子](#自訂運算子)
		* [介於／不介於](#介於不介於)
		* [於清單／不於清單內](#於清單不於清單內)
		* [或／還有或](#或還有或)
		* [空值](#空值)
		* [時間戳](#時間戳)
			* [相對](#相對)
			* [日期](#日期)
			* [時間](#時間)
		* [生條件](#生條件)
			* [條件變數](#條件變數)
	* [刪除](#刪除)
	* [排序](#排序)
		* [從值排序](#從值排序)
	* [群組](#群組)
	* [加入](#加入)
		* [條件限制](#條件限制)
	* [子指令](#子指令)
		* [選擇／取得](#選擇取得)
		* [插入](#插入-1)
		* [加入](#加入-1)
		* [存在／不存在](#存在不存在)
	* [是否擁有該筆資料](#是否擁有該筆資料)
	* [輔助函式](#輔助函式)
		* [資料庫連線](#資料庫連線)
		* [最後執行的 SQL 指令](#最後執行的-sql-指令)
		* [結果／影響的行數](#結果影響的行數)
		* [最後插入的編號](#最後插入的編號)
		* [總筆數](#總筆數)
	* [交易函式](#交易函式)
	* [鎖定表格](#鎖定表格)
	* [指令關鍵字](#指令關鍵字)
		* [多個選項](#多個選項)
* [表格建構函式](#表格建構函式)

# 安裝方式

打開終端機並且透過 `go get` 安裝此套件即可。

```bash
$ go get gopkg.in/teacat/reiner.v2
```

# 命名建議

在 Reiner 中為了配合 [Golang](https://golang.org/) 程式命名規範，我們建議你將所有事情以[駝峰式大小寫](https://zh.wikipedia.org/zh-tw/%E9%A7%9D%E5%B3%B0%E5%BC%8F%E5%A4%A7%E5%B0%8F%E5%AF%AB)命名，因為這能夠確保兩邊的風格相同。事實上，甚至連資料庫內的表格名稱、欄位名稱都該這麼做。當遇上 `ip`、`id`、`url` 時，請遵循 Golang 的命名方式皆以大寫使用，如 `AddrIP`、`UserID`、`PhotoURL`，而不是 `AddrIp`、`UserId`、`PhotoUrl`。

# NULL 值

在 Golang 裏處理資料庫的 NULL 值向來都不是很方便，因此不建議允許資料庫中可有 NULL 欄位。基於 Reiner 底層的 [`go-sql-driver/mysql`](https://github.com/go-sql-driver/mysql) 因素，Reiner 並不會將接收到的 NULL 值轉換成指定型態的零值（Zero Value），這意味著當你從資料庫中取得一個可能為 NULL 值的字串，你必須透過 `*string` 或者 `sql.NullString` 而非普通的 `string` 型態（會發生 Scan 錯誤），實際用法像這樣。

```go
type User struct {
	Username string
	Nickname sql.NullString
	Age 	 sql.NullInt64
}

// 然後綁定這個建構體到資料庫結果。
var u User
db.Table("Users").Bind(&u).GetOne()

// 輸出取得的結果。
if !u.Nickname.Valid() {
	panic("Nickname 的內容是 NULL！不可饒恕！")
}
fmt.Println(u.Nickname.Value)
```

# 使用方式

Reiner 的使用方式十分直覺與簡易，類似基本的 SQL 指令集但是更加地簡化了。

## 資料庫連線

首先你需要透過函式來將 Reiner 連上資料庫，如此一來才能夠初始化建置函式庫與相關的資料庫表格建構函式。一個最基本的單資料庫連線，讀寫都將透過此連線，連線字串共用於其它套件是基於 DSN（[Data Source Name](https://en.wikipedia.org/wiki/Data_source_name)）。

```go
import "github.com/teacat/reiner"

db, err := reiner.New("root:root@/test?charset=utf8")
if err != nil {
    panic(err)
}
```

### 水平擴展（讀／寫分離）

這種方式可以有好幾個主要資料庫、副從資料庫，這意味著寫入時都會流向到主要資料庫，而讀取時都會向副從資料庫請求。這很適合用在大型結構還有水平擴展上。當你有多個資料庫來源時，Reiner 會逐一遞詢每個資料庫來源，英文稱其為 [Round Robin](https://zh.wikipedia.org/zh-tw/%E5%BE%AA%E7%92%B0%E5%88%B6)，也就是每個資料庫都會輪流呼叫而避免單個資料庫負荷過重，也不會有隨機呼叫的事情發生。

```go
import "github.com/teacat/reiner"

db, err := reiner.New("root:root@/master?charset=utf8", []string{
	"root:root@/slave?charset=utf8",
	"root:root@/slave2?charset=utf8",
	"root:root@/slave3?charset=utf8",
})
if err != nil {
    panic(err)
}
```

### SQL 建構模式

如果你已經有喜好的 SQL 資料庫處理套件，那麼你就可以在建立 Reiner 時不要傳入任何資料，這會使 Reiner 避免與資料庫互動，透過這個設計你可以將 Reiner 作為你的 SQL 指令建構函式。

```go
// 當沒有傳入 MySQL 連線資料時，Reiner 僅會建置 SQL 執行指令而非與資料庫有實際互動。
builder, _ := reiner.New()
// 然後像這樣透過 Reiner 建立執行指令。
myQuery, _ := builder.Table("Users").Where("Username", "YamiOdymel").Get()

// 透過 `Query` 取得 Reiner 所建立的 Query 當作欲執行的資料庫指令。
sql.Prepare(myQuery.Query())
// 接著展開 `Params` 即是我們在 Reiner 中存放的值。
sql.Exec(myQuery.Params()...)
// 等效於：SELECT * FROM Users WHERE Username = ?
```

## 資料綁定與處理

Reiner 允許你將查詢結果映射到結構體切片或結構體。

```go
var user []*User
db.Bind(&user).Table("Users").Get()
```

## 插入

透過 Reiner 你可以很輕鬆地透過建構體或是 map 來插入一筆資料。這是最傳統的插入方式，若該表格有自動遞增的編號欄位，插入後你就能透過 `LastInsertID` 獲得最後一次插入的編號。

```go
db.Table("Users").Insert(map[string]interface{}{
	"Username": "YamiOdymel",
	"Password": "test",
})
// 等效於：INSERT INTO Users (Username, Password) VALUES (?, ?)
```

### 覆蓋

覆蓋的用法與插入相同，當有同筆資料時會先進行刪除，然後再插入一筆新的，這對有外鍵的表格來說十分危險。

```go
db.Table("Users").Replace(map[string]interface{}{
	"Username": "YamiOdymel",
	"Password": "test",
})
// 等效於：REPLACE INTO Users (Username, Password) VALUES (?, ?)
```

### 函式

插入時你可以透過 Reiner 提供的函式來執行像是 `SHA1()` 或者取得目前時間的 `NOW()`，甚至將目前時間加上一年⋯等。

```go
db.Table("Users").Insert(map[string]interface{}{
	"Username":  "YamiOdymel",
	"Password":  db.Func("SHA1(?)", "secretpassword+salt"),
	"Expires":   db.Now("+1Y"),
	"CreatedAt": db.Now(),
})
// 等效於：INSERT INTO Users (Username, Password, Expires, CreatedAt) VALUES (?, SHA1(?), NOW() + INTERVAL 1 YEAR, NOW())
```

### 當重複時

Reiner 支援了插入資料若重複時可以更新該筆資料的指定欄位>這類似「覆蓋」，但這並不會先刪除原先的資料，這種方式僅會在插入時檢查是否重複，若重複則更新該筆資料。

```go
lastInsertID := "ID"
db.Table("Users").OnDuplicate([]string{"UpdatedAt"}, lastInsertID).Insert(map[string]interface{}{
	"Username":  "YamiOdymel",
	"Password":  "test",
	"UpdatedAt": db.Now(),
})
// 等效於：INSERT INTO Users (Username, Password, UpdatedAt) VALUES (?, ?, NOW()) ON DUPLICATE KEY UPDATE UpdatedAt = VALUES(UpdatedAt)
```

### 多筆資料

Reiner 允許你透過 `InsertMulti` 同時間插入多筆資料（單指令插入多筆資料），這省去了透過迴圈不斷執行單筆插入的困擾，這種方式亦大幅度提升了效能。

```go
data := []map[string]interface{}{
	{
		"Username": "YamiOdymel",
		"Password": "test",
	}, {
		"Username": "Karisu",
		"Password": "12345",
	},
}
db.Table("Users").InsertMulti(data)
// 等效於：INSERT INTO Users (Username, Password) VALUES (?, ?), (?, ?)
```

## 筆數限制

`Limit` 能夠限制 SQL 執行的筆數，如果是 10，那就表示只處理最前面 10 筆資料而非全部（例如：選擇、更新、移除）。

```go
db.Table("Users").Limit(10).Update(data)
// 等效於：UPDATE Users SET ... LIMIT 10
```

## 更新

更新一筆資料在 Reiner 中極為簡單，你只需要指定表格名稱還有資料即可。

```go
db.Table("Users").Where("Username", "YamiOdymel").Update(map[string]interface{}{
	"Username": "Karisu",
	"Password": "123456",
})
// 等效於：UPDATE Users SET Username = ?, Password = ? WHERE Username = ?
```

## 選擇與取得

最基本的選擇在 Reiner 中稱之為 `Get` 而不是 `Select`。

```go
db.Table("Users").Get()
// 等效於：SELECT * FROM Users
```

### 指定欄位

在 `Get` 中傳遞欄位名稱作為參數，多個欄位由逗點區分，亦能是函式。

```go
db.Table("Users").Get("Username", "Nickname")
// 等效於：SELECT Username, Nickname FROM Users

db.Table("Users").Get("COUNT(*) AS Count")
// 等效於：SELECT COUNT(*) AS Count FROM Users
```

### 單行資料

通常多筆結果會映射到一個切片或是陣列，而 `GetOne` 可以取得單筆資料並將其結果映射到單個建構體或 `map`，令使用上更加方便。

當透過 `map[string]interface{}` 當作映射對象的時候，請注意資料庫並不會自動辨別 `int`、`string` 等資料型態，反倒有可能會是 `int64`、`[]uint8{[]byte}`，因此使用 `map` 時請多加注意在型態轉換上的部分。

```go
var u User
db.Bind(&u).Table("Users").Where("ID", 1).GetOne()
// 等效於：SELECT * FROM Users WHERE ID = ? LIMIT 1

var d map[string]interface{}
db.Bind(&d).Table("Users").GetOne("SUM(ID) AS Sum", "COUNT(*) AS Count")
// 等效於：SELECT SUM(ID), COUNT(*) AS Count FROM Users LIMIT 1

fmt.Println(d["Sum"])
fmt.Println(d["Count"])
```

### 單欄位值

透過 `GetValue` 和 `GetValues` 來取得單個欄位的內容。例如說：你想要單個使用者的暱稱，甚至是多個使用者的暱稱陣列就很適用。

```go
// 取得多筆資料的 `Username` 欄位資料。
var us []string
db.Bind(&u).Table("Users").GetValues("Username")
// 等效於：SELECT Username FROM Users

// 取得單筆資料的某個欄位值。
var u string
db.Bind(&u).Table("Users").GetValue("Username")
// 等效於：SELECT Username FROM Users LIMIT 1

// 或者是函式。
var i int
db.Bind(&i).Table("Users").GetValue("COUNT(*)")
// 等效於：SELECT COUNT(*) FROM Users LIMIT 1
```

### 分頁功能

分頁就像是取得資料ㄧ樣，但更擅長用於多筆資料、不會一次顯示完畢的內容。Reiner 能夠幫你自動處理換頁功能，讓你不需要自行計算換頁時的筆數應該從何開始。為此，你需要定義兩個變數，一個是目前的頁數，另一個是單頁能有幾筆資料。

```go
// 目前的頁數。
page := 1
// 設置一頁最多能有幾筆資料。
db.PageLimit = 10
db = db.Table("Users").Paginate(page)
// 等效於：SELECT SQL_CALC_FOUND_ROWS * FROM Users LIMIT 0, 10

fmt.Println("目前頁數為 %d，共有 %d 頁", page, db.TotalPages)
```

## 執行生指令

Reiner 已經提供了近乎日常中 80% 會用到的方式，但如果好死不死你想使用的功能在那 20% 之中，我們還提供了原生的方法能讓你直接輸入 SQL 指令執行自己想要的鳥東西。一個最基本的生指令（Raw Query）就像這樣。

其中亦能帶有預置聲明（Prepared Statement），也就是指令中的問號符號替代了原本的值。這能避免你的 SQL 指令遭受注入攻擊。

```go
var us []User
db.Bind(&us).RawQuery("SELECT * FROM Users WHERE ID >= ?", 10)
```

### 單行資料

`RawQueryOne` 是個僅選擇單筆資料的生指令函式，這意味著你能夠將取得的資料映射到建構體或是 `map` 上。

```go
var u User
db.Bind(&u).RawQueryOne("SELECT * FROM Users WHERE ID = ?", 10)
// 等效於：SELECT * FROM Users WHERE ID = ? LIMIT 1

var d map[string]interface{}
db.Bind(&d).RawQueryOne("SELECT SUM(ID), COUNT(*) AS Count FROM Users")
// 等效於：SELECT SUM(ID), COUNT(*) AS Count FROM Users LIMIT 1

fmt.Println(d["Sum"])
fmt.Println(d["Count"])
```

### 單欄位值

透過 `RawQueryValue` 與 `RawQueryValues` 可以取得單個欄位的內容。例如說：你想要單個使用者的暱稱，甚至是多個使用者的暱稱陣列就很適用。

```go
// 取得多筆資料的 `Username` 欄位資料。
var us []string
db.Bind(&us).RawQueryValues("SELECT Username FROM Users")

// 取得單筆資料的某個欄位值。
var pwd string
db.Bind(&pwd).RawQueryValue("SELECT Password FROM Users WHERE ID = ?", 10)
// 等效於：SELECT Password FROM Users WHERE ID = ? LIMIT 1

// 或者是函式。
var i int
db.Bind(&i).RawQueryValue("SELECT COUNT(*) FROM Users")
// 等效於：SELECT COUNT(*) FROM Users LIMIT 1
```

### 進階方式

如果你對 SQL 指令夠熟悉，你也可以使用更進階且複雜的用法。

```go
db.RawQuery("SELECT ID, FirstName, LastName FROM Users WHERE ID = ? AND Username = ?", 1, "admin")

params := []int{10, 1, 10, 11, 2, 10}
query := (`
	(SELECT A FROM t1 WHERE A = ? AND B = ?)
	UNION ALL
	(SELECT A FROM t2 WHERE A = ? AND B = ?)
	UNION ALL
	(SELECT A FROM t3 WHERE A = ? AND B = ?)
`)
db.RawQuery(query, params...)
```

## 條件宣告

透過 Reiner 宣告 `WHERE` 條件也能夠很輕鬆。一個最基本的 `WHERE AND` 像這樣使用。

```go
db.Table("Users").Where("ID", 1).Where("Username", "admin").Get()
// 等效於：SELECT * FROM Users WHERE ID = ? AND Username = ?
```

### 擁有

`HAVING` 能夠與 `WHERE` 一同使用。

```go
db.Table("Users").Where("ID", 1).Having("Username", "admin").Get()
// 等效於：SELECT * FROM Users WHERE ID = ? HAVING Username = ?
```

### 欄位比較

如果你想要在條件中宣告某個欄位是否等於某個欄位⋯你能夠像這樣。

```go
// 別這樣。
db.Table("Users").Where("LastLogin", "CreatedAt").Get()
// 這樣才對。
db.Table("Users").Where("LastLogin = CreatedAt").Get()
// 等效於：SELECT * FROM Users WHERE LastLogin = CreatedAt
```

### 自訂運算子

在 `Where` 或 `Having` 中，你可以自訂條件的運算子，如 >=、<=、<>⋯等。

```go
db.Table("Users").Where("ID", ">=", 50).Get()
// 等效於：SELECT * FROM Users WHERE ID >= ?
```

### 介於／不介於

條件也可以用來限制數值內容是否在某數之間（相反之，也能夠限制是否不在某範圍內）。

```go
db.Table("Users").Where("ID", "BETWEEN", 0, 20).Get()
// 等效於：SELECT * FROM Users WHERE ID BETWEEN ? AND ?
```

### 於清單／不於清單內

條件能夠限制並確保取得的內容不在（或者在）指定清單內。

```go
db.Table("Users").Where("ID", "IN", 1, 5, 27, -1, "d").Get()
// 等效於：SELECT * FROM Users WHERE ID IN (?, ?, ?, ?, ?)

list := []interface{}{1, 5, 27, -1, "d"}
db.Table("Users").Where("ID", "IN", list...).Get()
// 等效於：SELECT * FROM Users WHERE ID IN (?, ?, ?, ?, ?)
```

### 或／還有或

通常來說多個 `Where` 會產生 `AND` 條件，這意味著所有條件都必須符合，有些時候你只希望符合部分條件即可，就能夠用上 `OrWhere`。

```go
db.Table("Users").Where("FirstNamte", "John").OrWhere("FirstNamte", "Peter").Get()
// 等效於：SELECT * FROM Users WHERE FirstName = ? OR FirstName = ?
```

如果你的要求比較多，希望達到「A = B 或者 (A = C 或 A = D)」的話，你可以嘗試這樣。

```go
db.Table("Users").Where("A = B").OrWhere("(A = C OR A = D)").Get()
// 等效於：SELECT * FROM Users WHERE A = B OR (A = C OR A = D)
```

### 空值

確定某個欄位是否為空值。

```go
// 別這樣。
db.Table("Users").Where("LastName", "NULL").Get()
// 這樣才對。
db.Table("Users").Where("LastName", "IS", nil).Get()
// 等效於：SELECT * FROM Users WHERE LastName IS NULL
```

### 時間戳

[Unix Timestamp](https://en.wikipedia.org/wiki/Unix_time) 是一項將日期與時間秒數換算成數字的格式（範例：`1498001308`），這令你能夠輕易地換算其秒數，但當你要判斷時間是否為某一年、月、日，甚至範圍的時候就會有些許困難，而 Reiner 也替你想到了這一點。

需要注意的是 Reiner 中的 `Timestamp` 工具無法串聯使用，這意味著當你想要確認時間戳是否為某年某月時，你需要有兩個 `Where` 條件，而不行使用 `IsYear().IsMonth()`。更多的用法可以在原生文件中找到，這裡僅列出不完全的範例供大略參考。

#### 日期

判斷是否為特定年、月、日、星期或完整日期。

```go
t := db.Timestamp

db.Table("Users").Where("CreatedAt", t.IsDate("2017-07-13")).Get()
// 等效於：SELECT * FROM Users WHERE DATE(FROM_UNIXTIME(CreatedAt)) = ?

db.Table("Users").Where("CreatedAt", t.IsYear(2017)).Get()
// 等效於：SELECT * FROM Users WHERE YEAR(FROM_UNIXTIME(CreatedAt)) = ?

db.Table("Users").Where("CreatedAt", t.IsMonth(1)).Get()
db.Table("Users").Where("CreatedAt", t.IsMonth("January")).Get()
// 等效於：SELECT * FROM Users WHERE MONTH(FROM_UNIXTIME(CreatedAt)) = ?

db.Table("Users").Where("CreatedAt", t.IsDay(16)).Get()
// 等效於：SELECT * FROM Users WHERE DAY(FROM_UNIXTIME(CreatedAt)) = ?

db.Table("Users").Where("CreatedAt", t.IsWeekday(5)).Get()
db.Table("Users").Where("CreatedAt", t.IsWeekday("Friday")).Get()
// 等效於：SELECT * FROM Users WHERE WEEKDAY(FROM_UNIXTIME(CreatedAt)) = ?
```

#### 時間

確定是否為特定時間。

```go
t := db.Timestamp

db.Table("Users").Where("CreatedAt", t.IsHour(18)).Get()
// 等效於：SELECT * FROM Users WHERE HOUR(FROM_UNIXTIME(CreatedAt)) = ?

db.Table("Users").Where("CreatedAt", t.IsMinute(25)).Get()
// 等效於：SELECT * FROM Users WHERE MINUTE(FROM_UNIXTIME(CreatedAt)) = ?

db.Table("Users").Where("CreatedAt", t.IsSecond(16)).Get()
// 等效於：SELECT * FROM Users WHERE SECOND(FROM_UNIXTIME(CreatedAt)) = ?

db.Table("Users").Where("CreatedAt", t.IsWeekday(5)).Get()
// 等效於：SELECT * FROM Users WHERE WEEKDAY(FROM_UNIXTIME(CreatedAt)) = ?
```

### 生條件

你也能夠直接在條件中輸入指令。

```go
db.Table("Users").Where("ID != CompanyID").Where("DATE(CreatedAt) = DATE(LastLogin)").Get()
// 等效於：SELECT * FROM Users WHERE ID != CompanyID AND DATE(CreatedAt) = DATE(LastLogin)
```

#### 條件變數

生條件中可以透過 `?` 符號，並且在後面傳入自訂變數。

```go
db.Table("Users").Where("(ID = ? OR ID = ?)", 6, 2).Where("Login", "Mike").Get()
// 等效於：SELECT * FROM Users WHERE (ID = ? OR ID = ?) AND Login = ?
```

## 刪除

刪除一筆資料再簡單不過了，透過 `Count` 計數能夠清楚知道你的 SQL 指令影響了幾行資料，如果是零的話即是無刪除任何資料。

```go
var err error
db, err = db.Table("Users").Where("ID", 1).Delete()
if count := db.Count(); err == nil && count != 0 {
    fmt.Printf("成功地刪除了 %d 筆資料！", count)
}
// 等效於：DELETE FROM Users WHERE ID = ?
```

## 排序

Reiner 亦支援排序功能，如遞增或遞減，亦能擺放函式。

```go
db.Table("Users").OrderBy("ID", "ASC").OrderBy("Login", "DESC").OrderBy("RAND()").Get()
// 等效於：SELECT * FROM Users ORDER BY ID ASC, Login DESC, RAND()
```

### 從值排序

也能夠從值進行排序，只需要傳入一個切片即可。

```go
db.Table("Users").OrderBy("UserGroup", "ASC", "SuperUser", "Admin", "Users").Get()
// 等效於：SELECT * FROM Users ORDER BY FIELD (UserGroup, ?, ?, ?) ASC
```

## 群組

簡單的透過 `GroupBy` 就能夠將資料由指定欄位群組排序。

```go
db.Table("Users").GroupBy("Name").Get()
// 等效於：SELECT * FROM Users GROUP BY Name
```

## 加入

Reiner 支援多種表格加入方式，如：`InnerJoin`、`LeftJoin`、`RightJoin`、`NaturalJoin`、`CrossJoin`。

```go
db.
	Table("Products").
	LeftJoin("Users", "Products.TenantID = Users.TenantID").
	Where("Users.ID", 6).
	Get("Users.Name", "Products.ProductName")
// 等效於：SELECT Users.Name, Products.ProductName FROM Products AS Products LEFT JOIN Users AS Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?
```

### 條件限制

你亦能透過 `JoinWhere` 或 `JoinOrWhere` 擴展表格加入的限制條件。

```go
db.
	Table("Products").
	LeftJoin("Users", "Products.TenantID = Users.TenantID").
	JoinOrWhere("Users", "Users.TenantID", 5).
	Get("Users.Name", "Products.ProductName")
// 等效於：SELECT Users.Name, Products.ProductName FROM Products AS Products LEFT JOIN Users AS Users ON (Products.TenantID = Users.TenantID OR Users.TenantID = ?)
```

## 子指令

Reiner 支援複雜的子指令，欲要建立一個子指令請透過 `SubQuery` 函式，這將會建立一個不能被執行的資料庫建置函式庫，令你可以透過 `Get`、`Update` 等建立相關 SQL 指令，但不會被資料庫執行。將其帶入到一個正常的資料庫函式中即可成為子指令。

```go
subQuery := db.SubQuery().Table("Users").Get()
// 等效於不會被執行的：SELECT * FROM Users
```

### 選擇／取得

你能夠輕易地將子指令放置在選擇／取得指令中。

```go
subQuery := db.SubQuery().Table("Products").Where("Quantity", ">", 2).Get("UserID")

db.Table("Users").Where("ID", "IN", subQuery).Get()
// 等效於：SELECT * FROM Users WHERE ID IN (SELECT UserID FROM Products WHERE Quantity > ?)
```

### 插入

插入新資料時也可以使用子指令。

```go
subQuery := db.SubQuery().Table("Users").Where("ID", 6).Get("Name")

db.Table("Products").Insert(map[string]interface{}{
	"ProductName": "測試商品",
	"UserID":      subQuery,
	"LastUpdated": db.Now(),
})
// 等效於：INSERT INTO Products (ProductName, UserID, LastUpdated) VALUES (?, (SELECT Name FROM Users WHERE ID = 6), NOW())
```

### 加入

就算是加入表格的時候也可以用上子指令，但你需要為子指令建立別名。

```go
subQuery := db.SubQuery("Users").Table("Users").Where("Active", 1).Get()

db.
	Table("Products").
	LeftJoin(subQuery, "Products.UserID = U.ID").
	Get("Users.Username", "Products.ProductName")
// 等效於：SELECT Users.Username, Products.ProductName FROM Products AS Products LEFT JOIN (SELECT * FROM Users WHERE Active = ?) AS Users ON Products.UserID = Users.ID
```

### 存在／不存在

你同時也能夠透過子指令來確定某筆資料是否存在。

```go
subQuery := db.SubQuery().Table("Users").Where("Company", "測試公司").Get("UserID")

db.Table("Products").Where(subQuery, "EXISTS").Get()
// 等效於：SELECT * FROM Products WHERE EXISTS (SELECT UserID FROM Users WHERE Company = ?)
```

## 是否擁有該筆資料

有些時候我們只想知道資料庫是否有符合的資料，但並不是要取得其資料，舉例來說就像是登入是僅是要確認帳號密碼是否吻合，此時就可以透過 `Has` 用來確定資料庫是否有這筆資料。

```go
has, err := db.Table("Users").Where("Username", "yamiodymel").Where("Password", "123456").Has()
if has {
	fmt.Println("登入成功！")
} else {
	fmt.Println("帳號或密碼錯誤。")
}
```

## 輔助函式

Reiner 有提供一些輔助用的函式協助你除錯、紀錄，或者更加地得心應手。

### 資料庫連線

透過 Disconnect 結束一段連線。

```go
if err := db.Disconnect(); err != nil {
	panic(err)
}
```

你也能在資料庫發生錯誤、連線遺失時透過 `Connect` 來重新手動連線。

```go
if err := db.Ping(); err != nil {
	db.Connect()
}
```

### 最後執行的 SQL 指令

取得最後一次所執行的 SQL 指令，這能夠用來記錄你所執行的所有動作。

```go
db = db.Table("Users").Get()
fmt.Println("最後一次執行的 SQL 指令是：%s", db.LastQuery)
// 輸出：SELECT * FROM Users
```

### 結果／影響的行數

行數很常用於檢查是否有資料、作出變更。資料庫不會因為沒有變更任何資料而回傳一個錯誤（資料庫僅會在真正發生錯誤時回傳錯誤資料），所以這是很好的檢查方法。

```go
db, _ = db.Table("Users").Get()
fmt.Println("總共獲取 %s 筆資料", db.Count())
db, _ = db.Table("Users").Delete()
fmt.Println("總共刪除 %s 筆資料", db.Count())
db, _ = db.Table("Users").Update(data)
fmt.Println("總共更新 %s 筆資料", db.Count())
```

### 最後插入的編號

當插入一筆新的資料，而該表格帶有自動遞增的欄位時，就能透過 `LastInsertID` 取得最新一筆資料的編號。

```go
var id int

db, _ = db.Table("Users").Insert(data)
id = db.LastInsertID
```

如果你是經由 `InsertMulti` 同時間插入多筆資料，基於 MySQL 底層的設定，你並沒有辦法透過 `LastInsertID` 取得剛才插入的所有資料編號。如果你仍希望取得插入編號，請透過迴圈不斷地執行 `Insert` 並保存其 `LastInsertID` 資料。

```go
var ids []int

for ... {
	var err error
	db, err = db.Table("Users").Insert(data)
	if err != nil {
		ids = append(ids, db.LastInsertID)
	}
}
```

### 總筆數

如果你想取得這個指令總共能夠取得多少筆資料，透過 `WithTotalCount` 就能夠啟用總筆數查詢，這可能會稍微降低一點資料庫效能。

```go
db, _ = db.Table("Users").WithTotalCount().Get()
fmt.Println(db.TotalCount)
```

## 交易函式

交易函式僅限於 [InnoDB](https://zh.wikipedia.org/zh-tw/InnoDB) 型態的資料表格，這能令你的資料寫入更加安全。你可以透過 `Begin` 開始記錄並繼續你的資料庫寫入行為，如果途中發生錯誤，你能透過 `Rollback` 回到紀錄之前的狀態，即為回溯（或滾回、退回），如果這筆交易已經沒有問題了，透過 `Commit` 將這次的變更永久地儲存到資料庫中。

```go
// 當交易開始時請使用回傳的 `tx` 而不是原先的 `db`，這樣才能確保交易繼續。
tx, err := db.Begin()
if err != nil {
	panic(err)
}

// 如果插入資料時發生錯誤，則呼叫 `Rollback()` 回到交易剛開始的時候。
if _, err = tx.Table("Wallets").Insert(data); err != nil {
	tx.Rollback()
	panic(err)
}
if _, err = tx.Table("Users").Insert(data); err != nil {
	tx.Rollback()
	panic(err)
}

// 透過 `Commit()` 確保上列變更都已經永久地儲存到資料庫。
if err := tx.Commit(); err != nil {
	panic(err)
}
```

## 鎖定表格

你能夠手動鎖定資料表格，避免同時間寫入相同資料而發生錯誤。

```go
db.Table("Users").SetLockMethod("WRITE").Lock()

// 呼叫其他的 Lock() 函式也會自動將前一個上鎖解鎖，當然你也可以手動呼叫 Unlock() 解鎖。
db.Unlock()
// 等效於：UNLOCK TABLES

// 同時間要鎖上兩個表格也很簡單。
db.Table("Users", "Logs").SetLockMethod("READ").Lock()
// 等效於：LOCK TABLES Users READ, Logs READ
```

## 指令關鍵字

Reiner 也支援設置指令關鍵字。

```go
db.Table("Users").SetQueryOption("LOW_PRIORITY").Insert(data)
// 等效於：INSERT LOW_PRIORITY INTO Users ...

db.Table("Users").SetQueryOption("FOR UPDATE").Get()
// 等效於：SELECT * FROM Users FOR UPDATE

db.Table("Users").SetQueryOption("SQL_NO_CACHE").Get()
// 等效於：SELECT SQL_NO_CACHE * FROM Users
```

### 多個選項

你亦能同時設置多個關鍵字給同個指令。

```go
db.Table("Users").SetQueryOption("LOW_PRIORITY", "IGNORE").Insert(data)
// Gives: INSERT LOW_PRIORITY IGNORE INTO Users ...
```

## 效能追蹤

這會降低執行效能，但透過追蹤功能能夠有效地得知每個指令所花費的執行時間和建置指令，並且取得相關執行檔案路徑與行號。

```go
db = db.SetTrace(true).Table("Users").Get()
fmt.Printf("%+v", db.Traces[0])

//[{Query:SELECT * FROM Users Duration:808.698µs Stacks:[map
//[File:/Users/YamiOdymel/go/src/github.com/teacat/reiner/builder.go Line:559 Skip:0 PC:19399228] map[Line:666 Skip:1 PC:19405153 //File:/Users/YamiOdymel/go/src/github.com/teacat/reiner/builder.go] map[Skip:2 PC:19407043 //File:/Users/YamiOdymel/go/src/github.com/teacat/reiner/builder.go Line:705] map[Line:74 Skip:3 PC:19548011 //File:/Users/YamiOdymel/go/src/github.com/teacat/reiner/builder.go] map[PC:17610310 //File:/usr/local/Cellar/go/1.8/libexec/src/testing/testing.go Line:657 Skip:4] map
//[File:/usr/local/Cellar/go/1.8/libexec/src/runtime/asm_amd64.s Line:2197 Skip:5 PC:17143345]] Error:<nil>}]
```

# 表格建構函式

Reiner 除了基本的資料庫函式可供使用外，還能夠建立一個表格並且規劃其索引、外鍵、型態。

```go
migration := db.Migration()

migration.Table("Users").Column("Username").Varchar(32).Primary().Create()
// 等效於：CREATE TABLE Users (Username VARCHAR(32) NOT NULL PRIMARY KEY) ENGINE=INNODB
```


| 數值       | 字串       | 二進制     | 檔案資料     | 時間      | 浮點數     | 固組   |
|-----------|------------|-----------|------------|-----------|-----------|-------|
| TinyInt   | Char       | Binary    | Blob       | Date      | Double    | Enum  |
| SmallInt  | Varchar    | VarBinary | MediumBlob | DateTime  | Decimal   | Set   |
| MediumInt | TinyText   | Bit       | LongBlob   | Time      | Float     |       |
| Int       | Text       |           |            | Timestamp |           |       |
| BigInt    | MediumText |           |            | Year      |           |       |
|           | LongText   |           |            |           |           |       |

# 相關連結

這裡是 Reiner 受啟發，或是和資料庫有所關聯的連結。

* [kisielk/sqlstruct](http://godoc.org/github.com/kisielk/sqlstruct)
* [jmoiron/sqlx](https://github.com/jmoiron/sqlx)
* [russross/meddler](https://github.com/russross/meddler)
* [jinzhu/gorm](https://github.com/jinzhu/gorm)
