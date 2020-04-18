# Rushia [![GoDoc](https://godoc.org/github.com/teacat/rushia?status.svg)](https://godoc.org/github.com/teacat/rushia) [![Coverage Status](https://coveralls.io/repos/github/teacat/rushia/badge.svg?branch=master)](https://coveralls.io/github/teacat/rushia?branch=master) [![Build Status](https://travis-ci.org/teacat/rushia.svg?branch=master)](https://travis-ci.org/teacat/rushia) [![Go Report Card](https://goreportcard.com/badge/github.com/teacat/rushia)](https://goreportcard.com/report/github.com/teacat/rushia)

一個由 [Golang](https://golang.org/) 撰寫且比起部分 [ORM](https://zh.wikipedia.org/wiki/%E5%AF%B9%E8%B1%A1%E5%85%B3%E7%B3%BB%E6%98%A0%E5%B0%84) 還要讚的 [MySQL](https://www.mysql.com/) 指令建置函式庫。彈性高、不需要建構體標籤。實際上，這就只是 [PHP-MySQLi-Database-Class](https://github.com/joshcam/PHP-MySQLi-Database-Class) 不過是用在 [Golang](https://golang.org/) 而已（但還是多了些功能）。

這是一個 SQL 指令建構庫，本身不帶有任何 SQL 連線，適合用於某些套件的基底。

# 這是什麼？

露西婭是一個由 [Golang](https://golang.org/) 撰寫的 [MySQL](https://www.mysql.com/) 的指令建置函式庫（不是 [ORM](https://zh.wikipedia.org/wiki/%E5%AF%B9%E8%B1%A1%E5%85%B3%E7%B3%BB%E6%98%A0%E5%B0%84)，永遠也不會是），幾乎所有東西都能操控於你手中。類似自己撰寫資料庫指令但是更簡單，JOIN 表格也變得比以前更方便了。

* 幾乎全功能的函式庫。
* 自動避免於 Goroutine 發生資料競爭的設計。
* 容易理解與記住、且使用方式十分簡單。
* SQL 指令建構函式。
* 資料庫表格建構協助函式。
* 彈性的建構體映射。
* 可串連的使用方式。
* 支援子指令（Sub Query）。
* 透過預置聲明（[Prepared Statement](https://en.wikipedia.org/wiki/Prepared_statement)），99.9% 避免 SQL 注入攻擊。

# 為什麼？

[Gorm](https://github.com/jinzhu/gorm) 已經是 [Golang](https://golang.org/) 裡的 [ORM](https://zh.wikipedia.org/wiki/%E5%AF%B9%E8%B1%A1%E5%85%B3%E7%B3%BB%E6%98%A0%E5%B0%84) 典範，但實際上要操作複雜與關聯性高的 SQL 指令時並不是很合適，而 Rushia 解決了這個問題。Rushia 也試圖不要和建構體扯上關係，不希望使用者需要手動指定任何標籤在建構體中。

# 索引

* [安裝方式](#安裝方式)
* [命名建議](#命名建議)
* [NULL 值](#null-值)
* [使用方式](#使用方式)
    * [映射](#映射)
    	* [省略](#省略)
	* [插入](#插入)
		* [覆蓋](#覆蓋)
		* [函式](#函式)
		* [當重複時](#當重複時)
		* [多筆資料](#多筆資料)
	* [筆數限制](#筆數限制)
	* [筆數偏移](#筆數偏移)
	* [更新](#更新)
		* [片段更新](#片段更新)
	* [選擇與取得](#選擇與取得)
		* [指定欄位](#指定欄位)
		* [單行資料](#單行資料)
	* [執行生指令](#執行生指令)
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
	* [輔助函式](#輔助函式)
		* [總筆數](#總筆數)
	* [鎖定表格](#鎖定表格)
	* [指令關鍵字](#指令關鍵字)
		* [多個選項](#多個選項)
* [表格建構函式](#表格建構函式)

# 安裝方式

打開終端機並且透過 `go get` 安裝此套件即可。

```bash
$ go get gopkg.in/teacat/rushia.v1
```

# 命名建議

在 Rushia 中為了配合 [Golang](https://golang.org/) 程式命名規範，我們建議你將所有事情以[駝峰式大小寫](https://zh.wikipedia.org/zh-tw/%E9%A7%9D%E5%B3%B0%E5%BC%8F%E5%A4%A7%E5%B0%8F%E5%AF%AB)命名，因為這能夠確保兩邊的風格相同。事實上，甚至連資料庫內的表格名稱、欄位名稱都該這麼做。當遇上 `ip`、`id`、`url` 時，請遵循 Golang 的命名方式皆以大寫使用，如 `AddrIP`、`UserID`、`PhotoURL`，而不是 `AddrIp`、`UserId`、`PhotoUrl`。

# NULL 值

在 Golang 裏處理資料庫的 NULL 值向來都不是很方便，因此不建議允許資料庫中可有 NULL 欄位。

# 使用方式

Rushia 的使用方式十分直覺與簡易，類似基本的 SQL 指令集但是更加地簡化了。

## 映射

你能夠直接將一個建構體傳入 `Insert` 或是 `Update` 之中，其欄位名稱與值都會被自動轉換 (注意！這並不會轉換成 MySQL 最常用的 `snake_case`！)。

```go
type User struct {
	Username string
	Password string
}
u := User{
	Username: "YamiOdymel",
	Password: "test",
}
db.Table("Users").Insert(u)
// 等效於：INSERT INTO Users (Username, Password) VALUES (?, ?)
```

### 省略

透過 `Omit`，你可以省略建構體中的某些欄位。

```go
type User struct {
	Username string
	Password string
}
u := User{
	Username: "YamiOdymel",
	Password: "test",
}
db.Table("Users").Omit("Username").Insert(u)
// 等效於：INSERT INTO Users (Password) VALUES (?)
```

## 插入

透過 Rushia 你可以很輕鬆地透過建構體或是 map 來插入一筆資料。這是最傳統的插入方式。

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

插入時你可以透過 Rushia 提供的函式來執行像是 `SHA1()` 或者取得目前時間的 `NOW()`，甚至將目前時間加上一年⋯等。

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

Rushia 支援了插入資料若重複時可以更新該筆資料的指定欄位，這類似「覆蓋」，但這並不會先刪除原先的資料，這種方式僅會在插入時檢查是否重複，若重複則更新該筆資料。

```go
lastInsertID := "ID"
db.Table("Users").OnDuplicate([]string{"UpdatedAt"}, lastInsertID).Insert(map[string]interface{}{
	"Username":  "YamiOdymel",
	"Password":  "test",
	"UpdatedAt": db.Now(),
})
// 等效於：INSERT INTO Users (Username, Password, UpdatedAt) VALUES (?, ?, NOW()) ON DUPLICATE KEY UPDATE ID=LAST_INSERT_ID(ID), UpdatedAt = VALUES(UpdatedAt)

db.Table("Users").OnDuplicate([]string{"UpdatedAt"}).Insert(map[string]interface{}{
	"Username":  "YamiOdymel",
	"Password":  "test",
	"UpdatedAt": db.Now(),
})
// 等效於：INSERT INTO Users (Username, Password, UpdatedAt) VALUES (?, ?, NOW()) ON DUPLICATE KEY UPDATE UpdatedAt = VALUES(UpdatedAt)
```

### 多筆資料

Rushia 允許你透過 `InsertMulti` 同時間插入多筆資料（單指令插入多筆資料），這省去了透過迴圈不斷執行單筆插入的困擾，這種方式亦大幅度提升了效能。

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

db.Table("Users").Limit(10, 20).Select(data)
// 等效於：SELECT * from Users LIMIT 10, 20
```

## 筆數偏移

透過 `Offset` 能夠以 `筆數, 上次索引編號` 的方式取得資料，例如：`10, 20` 則會從 `21` 開始取得 10 筆資料（`21, 22, 23...`）。

```go
db.Table("Users").Offset(10, 20).Select()
// 等效於：SELECT * from Users LIMIT 10 OFFSET 20
```

## 更新

更新一筆資料在 Rushia 中極為簡單，你只需要指定表格名稱還有資料即可。

```go
db.Table("Users").Where("Username", "YamiOdymel").Update(map[string]interface{}{
	"Username": "Karisu",
	"Password": "123456",
})
// 等效於：UPDATE Users SET Username = ?, Password = ? WHERE Username = ?
```

### 片段更新

當你希望某些欄位在零值的時候不要進行更新，那麼你就可以使用 `Patch` 來做片段更新（也叫小修補）。

```go
db.Table("Users").Where("Username", "YamiOdymel").Patch(map[string]interface{}{
	"Age": 0,
	"Username": "",
	"Password": "123456",
})
// 等效於：UPDATE Users SET Password = ? WHERE Username = ?
```

如果你希望有些欄位雖然是零值（如：`false`、`0`）但仍該在 `Patch` 時照樣更新，那麼就可以傳入一個 `PatchOptions` 選項。`ExcludedTypes` 表示欲排除的資料型態（如：`reflect.Bool`、`reflect.String`）、`ExcludedColumns` 表示欲忽略的欄位名稱。

排除的資料型態或欄位會在零值時一樣被更新到資料庫中。

```go
db.Table("Users").Where("Username", "YamiOdymel").Patch(map[string]interface{}{
	"Age": 0,
	"Username": "",
	"Password": "123456",
}, PatchOptions{
	ExcludedTypes: []reflect.Kind{reflect.Int},
	ExcludedColumns: []string{"Username"},
})
// 等效於：UPDATE Users SET Age = ?, Password = ?, Username = ? WHERE Username = ?
```

## 選擇與取得

最基本的資料取得在 Rushia 中透過 `Select` 使用。

```go
db.Table("Users").Select()
// 等效於：SELECT * FROM Users
```

### 指定欄位

在 `Select` 中傳遞欄位名稱作為參數，多個欄位由逗點區分，亦能是函式。

```go
db.Table("Users").Select("Username", "Nickname")
// 等效於：SELECT Username, Nickname FROM Users

db.Table("Users").Select("COUNT(*) AS Count")
// 等效於：SELECT COUNT(*) AS Count FROM Users
```

### 單行資料

如果只想要取得單筆資料，那麼就可以用上 `SelectOne`，這簡單來說就是 `.Limit(1)` 的縮寫。

```go
db.Table("Users").SelectOne("Username")
// 等效於：SELECT Username FROM Users LIMIT 1
```

## 執行生指令

Rushia 已經提供了近乎日常中 80% 會用到的方式，但如果好死不死你想使用的功能在那 20% 之中，我們還提供了原生的方法能讓你直接輸入 SQL 指令執行自己想要的鳥東西。一個最基本的生指令（Raw Query）就像這樣。

其中亦能帶有預置聲明（Prepared Statement），也就是指令中的問號符號替代了原本的值。這能避免你的 SQL 指令遭受注入攻擊。

```go
var us []User
db.Bind(&us).RawQuery("SELECT * FROM Users WHERE ID >= ?", 10)
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

透過 Rushia 宣告 `WHERE` 條件也能夠很輕鬆。一個最基本的 `WHERE AND` 像這樣使用。

```go
db.Table("Users").Where("ID", 1).Where("Username", "admin").Select()
// 等效於：SELECT * FROM Users WHERE ID = ? AND Username = ?
```

### 擁有

`HAVING` 能夠與 `WHERE` 一同使用。

```go
db.Table("Users").Where("ID", 1).Having("Username", "admin").Select()
// 等效於：SELECT * FROM Users WHERE ID = ? HAVING Username = ?
```

### 欄位比較

如果你想要在條件中宣告某個欄位是否等於某個欄位⋯你能夠像這樣。

```go
// 別這樣。
db.Table("Users").Where("LastLogin", "CreatedAt").Select()
// 這樣才對。
db.Table("Users").Where("LastLogin = CreatedAt").Select()
// 等效於：SELECT * FROM Users WHERE LastLogin = CreatedAt
```

### 自訂運算子

在 `Where` 或 `Having` 中，你可以自訂條件的運算子，如 >=、<=、<>⋯等。

```go
db.Table("Users").Where("ID", ">=", 50).Select()
// 等效於：SELECT * FROM Users WHERE ID >= ?
```

### 介於／不介於

條件也可以用來限制數值內容是否在某數之間（相反之，也能夠限制是否不在某範圍內）。

```go
db.Table("Users").Where("ID", "BETWEEN", 0, 20).Select()
// 等效於：SELECT * FROM Users WHERE ID BETWEEN ? AND ?
```

### 於清單／不於清單內

條件能夠限制並確保取得的內容不在（或者在）指定清單內。

```go
db.Table("Users").Where("ID", "IN", 1, 5, 27, -1, "d").Select()
// 等效於：SELECT * FROM Users WHERE ID IN (?, ?, ?, ?, ?)

list := []interface{}{1, 5, 27, -1, "d"}
db.Table("Users").Where("ID", "IN", list...).Select()
// 等效於：SELECT * FROM Users WHERE ID IN (?, ?, ?, ?, ?)
```

### 或／還有或

通常來說多個 `Where` 會產生 `AND` 條件，這意味著所有條件都必須符合，有些時候你只希望符合部分條件即可，就能夠用上 `OrWhere`。

```go
db.Table("Users").Where("FirstNamte", "John").OrWhere("FirstNamte", "Peter").Select()
// 等效於：SELECT * FROM Users WHERE FirstName = ? OR FirstName = ?
```

如果你的要求比較多，希望達到「A = B 或者 (A = C 或 A = D)」的話，你可以嘗試這樣。

```go
db.Table("Users").Where("A = B").OrWhere("(A = C OR A = D)").Select()
// 等效於：SELECT * FROM Users WHERE A = B OR (A = C OR A = D)
```

### 空值

確定某個欄位是否為空值。

```go
// 別這樣。
db.Table("Users").Where("LastName", "NULL").Select()
// 這樣才對。
db.Table("Users").Where("LastName", "IS", nil).Select()
// 等效於：SELECT * FROM Users WHERE LastName IS NULL
```

### 時間戳

[Unix Timestamp](https://en.wikipedia.org/wiki/Unix_time) 是一項將日期與時間秒數換算成數字的格式（範例：`1498001308`），這令你能夠輕易地換算其秒數，但當你要判斷時間是否為某一年、月、日，甚至範圍的時候就會有些許困難，而 Rushia 也替你想到了這一點。

需要注意的是 Rushia 中的 `Timestamp` 工具無法串聯使用，這意味著當你想要確認時間戳是否為某年某月時，你需要有兩個 `Where` 條件，而不行使用 `IsYear().IsMonth()`。更多的用法可以在原生文件中找到，這裡僅列出不完全的範例供大略參考。

#### 日期

判斷是否為特定年、月、日、星期或完整日期。

```go
t := rushia.NewTimestamp()

db.Table("Users").Where("CreatedAt", t.IsDate("2017-07-13")).Select()
// 等效於：SELECT * FROM Users WHERE DATE(FROM_UNIXTIME(CreatedAt)) = ?

db.Table("Users").Where("CreatedAt", t.IsYear(2017)).Select()
// 等效於：SELECT * FROM Users WHERE YEAR(FROM_UNIXTIME(CreatedAt)) = ?

db.Table("Users").Where("CreatedAt", t.IsMonth(1)).Select()
db.Table("Users").Where("CreatedAt", t.IsMonth("January")).Select()
// 等效於：SELECT * FROM Users WHERE MONTH(FROM_UNIXTIME(CreatedAt)) = ?

db.Table("Users").Where("CreatedAt", t.IsDay(16)).Select()
// 等效於：SELECT * FROM Users WHERE DAY(FROM_UNIXTIME(CreatedAt)) = ?

db.Table("Users").Where("CreatedAt", t.IsWeekday(5)).Select()
db.Table("Users").Where("CreatedAt", t.IsWeekday("Friday")).Select()
// 等效於：SELECT * FROM Users WHERE WEEKDAY(FROM_UNIXTIME(CreatedAt)) = ?
```

#### 時間

確定是否為特定時間。

```go
t := rushia.NewTimestamp()

db.Table("Users").Where("CreatedAt", t.IsHour(18)).Select()
// 等效於：SELECT * FROM Users WHERE HOUR(FROM_UNIXTIME(CreatedAt)) = ?

db.Table("Users").Where("CreatedAt", t.IsMinute(25)).Select()
// 等效於：SELECT * FROM Users WHERE MINUTE(FROM_UNIXTIME(CreatedAt)) = ?

db.Table("Users").Where("CreatedAt", t.IsSecond(16)).Select()
// 等效於：SELECT * FROM Users WHERE SECOND(FROM_UNIXTIME(CreatedAt)) = ?

db.Table("Users").Where("CreatedAt", t.IsWeekday(5)).Select()
// 等效於：SELECT * FROM Users WHERE WEEKDAY(FROM_UNIXTIME(CreatedAt)) = ?
```

### 生條件

你也能夠直接在條件中輸入指令。

```go
db.Table("Users").Where("ID != CompanyID").Where("DATE(CreatedAt) = DATE(LastLogin)").Select()
// 等效於：SELECT * FROM Users WHERE ID != CompanyID AND DATE(CreatedAt) = DATE(LastLogin)
```

#### 條件變數

生條件中可以透過 `?` 符號，並且在後面傳入自訂變數。

```go
db.Table("Users").Where("(ID = ? OR ID = ?)", 6, 2).Where("Login", "Mike").Select()
// 等效於：SELECT * FROM Users WHERE (ID = ? OR ID = ?) AND Login = ?
```

## 刪除

刪除一筆資料再簡單不過了。

```go
db.Table("Users").Where("ID", 1).Delete()
// 等效於：DELETE FROM Users WHERE ID = ?
```

## 排序

Rushia 亦支援排序功能，如遞增或遞減，亦能擺放函式。

```go
db.Table("Users").OrderBy("ID", "ASC").OrderBy("Login", "DESC").OrderBy("RAND()").Select()
// 等效於：SELECT * FROM Users ORDER BY ID ASC, Login DESC, RAND()
```

### 從值排序

也能夠從值進行排序，只需要傳入一個切片即可。

```go
db.Table("Users").OrderBy("UserGroup", "ASC", "SuperUser", "Admin", "Users").Select()
// 等效於：SELECT * FROM Users ORDER BY FIELD (UserGroup, ?, ?, ?) ASC
```

## 群組

簡單的透過 `GroupBy` 就能夠將資料由指定欄位群組排序。

```go
db.Table("Users").GroupBy("Name").Select()
// 等效於：SELECT * FROM Users GROUP BY Name
```

## 加入

Rushia 支援多種表格加入方式，如：`InnerJoin`、`LeftJoin`、`RightJoin`、`NaturalJoin`、`CrossJoin`。

```go
db.
	Table("Products").
	LeftJoin("Users", "Products.TenantID = Users.TenantID").
	Where("Users.ID", 6).
	Select("Users.Name", "Products.ProductName")
// 等效於：SELECT Users.Name, Products.ProductName FROM Products AS Products LEFT JOIN Users AS Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?
```

### 條件限制

你亦能透過 `JoinWhere` 或 `JoinOrWhere` 擴展表格加入的限制條件。

```go
db.
	Table("Products").
	LeftJoin("Users", "Products.TenantID = Users.TenantID").
	JoinOrWhere("Users", "Users.TenantID", 5).
	Select("Users.Name", "Products.ProductName")
// 等效於：SELECT Users.Name, Products.ProductName FROM Products AS Products LEFT JOIN Users AS Users ON (Products.TenantID = Users.TenantID OR Users.TenantID = ?)
```

## 子指令

Rushia 支援複雜的子指令，欲要建立一個子指令請透過 `SubQuery` 函式。將其帶入到一個正常的建置函式中即可成為子指令。

```go
subQuery := db.SubQuery().Table("Users").Select()
// 等效於：SELECT * FROM Users
```

### 選擇／取得

你能夠輕易地將子指令放置在選擇／取得指令中。

```go
subQuery := db.SubQuery().Table("Products").Where("Quantity", ">", 2).Select("UserID")

db.Table("Users").Where("ID", "IN", subQuery).Select()
// 等效於：SELECT * FROM Users WHERE ID IN (SELECT UserID FROM Products WHERE Quantity > ?)
```

### 插入

插入新資料時也可以使用子指令。

```go
subQuery := db.SubQuery().Table("Users").Where("ID", 6).Select("Name")

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
subQuery := db.SubQuery("Users").Table("Users").Where("Active", 1).Select()

db.
	Table("Products").
	LeftJoin(subQuery, "Products.UserID = U.ID").
	Select("Users.Username", "Products.ProductName")
// 等效於：SELECT Users.Username, Products.ProductName FROM Products AS Products LEFT JOIN (SELECT * FROM Users WHERE Active = ?) AS Users ON Products.UserID = Users.ID
```

### 存在／不存在

你同時也能夠透過子指令來確定某筆資料是否存在。

```go
subQuery := db.SubQuery().Table("Users").Where("Company", "測試公司").Select("UserID")

db.Table("Products").Where(subQuery, "EXISTS").Select()
// 等效於：SELECT * FROM Products WHERE EXISTS (SELECT UserID FROM Users WHERE Company = ?)
```

## 輔助函式

Rushia 有提供一些輔助用的函式協助你除錯、紀錄，或者更加地得心應手。

### 總筆數

如果你想取得這個指令總共能夠取得多少筆資料，透過 `WithTotalCount` 就能夠啟用總筆數查詢，這可能會稍微降低一點資料庫效能。

```go
db.Table("Users").WithTotalCount().Select()
// 等效於：SELECT SQL_CALC_FOUND_ROWS * FROM Users
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

Rushia 也支援設置指令關鍵字。

```go
db.Table("Users").SetQueryOption("LOW_PRIORITY").Insert(data)
// 等效於：INSERT LOW_PRIORITY INTO Users ...

db.Table("Users").SetQueryOption("FOR UPDATE").Select()
// 等效於：SELECT * FROM Users FOR UPDATE

db.Table("Users").SetQueryOption("SQL_NO_CACHE").Select()
// 等效於：SELECT SQL_NO_CACHE * FROM Users
```

### 多個選項

你亦能同時設置多個關鍵字給同個指令。

```go
db.Table("Users").SetQueryOption("LOW_PRIORITY", "IGNORE").Insert(data)
// Gives: INSERT LOW_PRIORITY IGNORE INTO Users ...
```

# 表格建構函式

Rushia 除了基本的資料庫函式可供使用外，還能夠建立一個表格並且規劃其索引、外鍵、型態。

```go
migration := rushia.NewMigration()

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

這裡是 Rushia 受啟發，或是和資料庫有所關聯的連結。

* [kisielk/sqlstruct](http://godoc.org/github.com/kisielk/sqlstruct)
* [jmoiron/sqlx](https://github.com/jmoiron/sqlx)
* [russross/meddler](https://github.com/russross/meddler)
* [jinzhu/gorm](https://github.com/jinzhu/gorm)
