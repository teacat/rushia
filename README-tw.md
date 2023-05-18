# Rushia [![GoDoc](https://godoc.org/github.com/teacat/rushia/v3?status.svg)](https://godoc.org/github.com/teacat/rushia/v3) [![Coverage Status](https://coveralls.io/repos/github/teacat/rushia/badge.svg?branch=master)](https://coveralls.io/github/teacat/rushia?branch=master) [![Build Status](https://app.travis-ci.com/teacat/rushia.svg?branch=master)](https://app.travis-ci.com/github/teacat/rushia) [![Go Report Card](https://goreportcard.com/badge/github.com/teacat/rushia)](https://goreportcard.com/report/github.com/teacat/rushia)

一個由 [Golang](https://golang.org/) 撰寫且比起部分 [ORM](https://zh.wikipedia.org/wiki/%E5%AF%B9%E8%B1%A1%E5%85%B3%E7%B3%BB%E6%98%A0%E5%B0%84) 還要讚的 [MySQL](https://www.mysql.com/) 指令建置函式庫。彈性高、不需要建構體標籤。原生想法基於 [PHP-MySQLi-Database-Class](https://github.com/joshcam/PHP-MySQLi-Database-Class) 和 [Laravel 查詢建構器](https://laravel.com/docs/8.x/queries) 但多了些功能。

這是一個 SQL 指令建構庫，本身不帶有任何 SQL 連線，適合用於某些套件的基底。

## 特色

-   幾乎全功能的函式庫。
-   容易理解與記住、且使用方式十分簡單。
-   SQL 指令建構函式。
-   資料庫表格建構協助函式。
-   彈性的建構體映射。
-   可串連的使用方式。
-   支援子指令（Sub Query）。
-   透過預置聲明（[Prepared Statement](https://en.wikipedia.org/wiki/Prepared_statement)），99.9% 避免 SQL 注入攻擊。

## 為什麼？

[Gorm](https://github.com/jinzhu/gorm) 已經是 [Golang](https://golang.org/) 裡的 [ORM](https://zh.wikipedia.org/wiki/%E5%AF%B9%E8%B1%A1%E5%85%B3%E7%B3%BB%E6%98%A0%E5%B0%84) 典範，但實際上要操作複雜與關聯性高的 SQL 指令時並不是很合適，而 Rushia 解決了這個問題。Rushia 也試圖不要和建構體扯上關係，不希望使用者需要手動指定任何標籤在建構體中。

## 安裝方式

打開終端機並且透過 `go get` 安裝此套件即可。

```bash
$ go get github.com/teacat/rushia/v3
```

## 使用方式

Rushia 的使用方式十分直覺與簡易，類似基本的 SQL 指令集但是更加地簡化了。

### 初始化語法

最基本的資料庫執行語法起始於 `NewQuery(...)`，其中可以帶入資料表的名稱，又或者是子指令（Sub Query）。而更為複雜的使用方式請參閱之後的章節。

```go
q := rushia.NewQuery("Users")
```

### 複製語法

預設情況下你的所有修改總是會改到同一個 Rushia 語法，如果你的環境可能有多執行緒或是希望複製一份規則額外更改，請使用 `Copy`。

```go
a := rushia.NewQuery("Users")
a.Where("Type = ?", "VIP")

b := a.Copy()
b.Where("Name = ?", "YamiOdymel")

Build(a.Select())
// 等效於：SELECT * FROM Users WHERE Type = ?
Build(b.Select())
// 等效於：SELECT * FROM Users WHERE Type = ? AND Name = ?
```

### 建置語法

當完成撰寫一個查詢語法後，必須透過 `Build` 將其建置便能得到建置的語句與其參數。一個語句必須要有 `Select`、`Exists`、`Replace`、`Update`、`Delete`…等作為結尾，否則會無法建置。

```go
query, params := rushia.Build(rushia.NewQuery("Users").Select())
// 等效於：SELECT * FROM Users
```

### 與其他資料庫套件搭配

由於 Rushia 是一個語法建置套件，這讓你可以得心應手地與自己喜好的資料庫連線函式庫進行搭配。舉例來說你可以使用 [jmoiron/sqlx](https://github.com/jmoiron/sqlx)：

```go
// 初始化 SQLX 的連線。
db, err := sqlx.Open("mysql", "root:password@tcp(localhost:3306)/db")

// 透過 Rushia 建置語法。
q := rushia.NewQuery("Users").Where("Usernam = ?", "YamiOdymel").Select()
query, params := rushia.Build(q)

// 將相關語法與參數傳入給 SQLX 的函式並執行。
rows, err := db.Query(query, params...)
// 等效於：SELECT * FROM Users WHERE Username = ?
```

又或者是 [go-gorm/gorm](https://github.com/go-gorm/gorm)：

```go
// 初始化 Gorm 的連線。
db, err := gorm.Open(mysql.Open("root:password@tcp(localhost:3306)/db"), &gorm.Config{})

// 透過 Rushia 建置語法。
q := rushia.NewQuery("Users").Where("Username = ?", "YamiOdymel").Select()
query, params := rushia.Build(q)

// 將相關語法與參數傳入給 Gorm 的函式並執行。
db.Raw(query, params...).Scan(&myUser)
// 等效於：SELECT * FROM Users WHERE Username = ?
```

### 結構體映射

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
rushia.NewQuery("Users").Insert(u)
// 等效於：INSERT INTO Users (Username, Password) VALUES (?, ?)
```

#### 結構體標籤

透過指定 `rushia` 結構體標籤，你可以省略一個欄位或是重新命名其欄位在 MySQL 查詢語法裡的名稱。

```go
type User struct {
	Username string `rushia:"-"`
	RealName string `rushia:"real_name"`
	Password string
}
u := User{
	Username: "YamiOdymel",
	RealName: "洨洨安",
	Password: "test",
}
rushia.NewQuery("Users").Insert(u)
// 等效於：INSERT INTO Users (real_name, Password) VALUES (?, ?)
```

### 省略

透過 `Omit`，你可以省略建構體中的某些欄位。

```go
type User struct {
	Username string
	Password string
	Age      int    `rushia:"my_age"`
}
u := User{
	Username: "YamiOdymel",
	Password: "test",
	Age     : "32"
}
rushia.NewQuery("Users").Omit("Username", "my_age").Insert(u)
// 等效於：INSERT INTO Users (Password) VALUES (?)
```

### 插入

Rushia 提供一個簡短的 `H`（`map[string]interface{}` 別名），這趨近於 [`gin.H`](https://pkg.go.dev/github.com/gin-gonic/gin#H)。建立一個插入語法的時候可以傳入 `H`、`H` 或是結構體。

```go
rushia.NewQuery("Users").Insert(rushia.H{
	"Username": "YamiOdymel",
	"Password": "test",
})
// 等效於：INSERT INTO Users (Username, Password) VALUES (?, ?)

rushia.NewQuery("Users").Insert(map[string]interface{
	"Username": "YamiOdymel",
	"Password": "test",
})
// 等效於：INSERT INTO Users (Username, Password) VALUES (?, ?)
```

### 多筆資料

Rushia 允許你透過 `[]H` 或 `[]map[string]interface{}` 一次插入多筆資料。

```go
data := []H{
	{
		"Username": "YamiOdymel",
		"Password": "test",
	}, {
		"Username": "Karisu",
		"Password": "12345",
	},
}
rushia.NewQuery("Users").Insert(data)
// 等效於：INSERT INTO Users (Username, Password) VALUES (?, ?), (?, ?)
```

### 覆蓋

覆蓋的用法與插入相同。當有同筆資料時會先進行刪除，然後再插入一筆新的，這對有外鍵的表格來說十分危險。若需要更為安全的方式請使用 `OnDuplicate`（`ON DUPLICATE KEY UPDATE`）函式。

```go
rushia.NewQuery("Users").Replace(rushia.H{
	"Username": "YamiOdymel",
	"Password": "test",
})
// 等效於：REPLACE INTO Users (Username, Password) VALUES (?, ?)
```

### 當重複時

Rushia 支援了插入資料若重複時可以更新該筆資料的指定欄位，這類似「覆蓋」，但這並不會先刪除原先的資料，這種方式僅會在插入時檢查是否重複，若重複則更新該筆資料。

```go
rushia.NewQuery("Users").As("New").OnDuplicate(rushia.H{
	"UpdatedAt": rushia.NewExpr("New.UpdatedAt"),
}).Insert(rushia.H{
	"Username":  "YamiOdymel",
	"UpdatedAt": rushia.NewExpr("NOW()"),
})
// 等效於：INSERT INTO Users (Username, UpdatedAt) VALUES (?, NOW()) AS New ON DUPLICATE KEY UPDATE UpdatedAt = New.UpdatedAt

rushia.NewQuery("Users").OnDuplicate(rushia.H{
	"UpdatedAt": rushia.NewExpr("VALUES(UpdatedAt)"),
}).Insert(rushia.H{
	"Username":  "YamiOdymel",
	"UpdatedAt": rushia.NewExpr("NOW()"),
})
// 注意！`VALUES` 這個用法已經在 MySQL 8.0.20 被棄用！請使用上面的方法！
// 等效於：INSERT INTO Users (Username, UpdatedAt) VALUES (?, NOW()) ON DUPLICATE KEY UPDATE UpdatedAt = VALUES(UpdatedAt)
```

### 表達式

插入較為複雜的值時，可以使用 `NewExpr` 建立一個新的表達式，便能傳入生指令與相關參數執行像是 `SHA1()` 或者取得目前時間的 `NOW()`，甚至將目前時間加上一年 ⋯ 等。

```go
rushia.NewQuery("Users").Insert(rushia.H{
	"Username":  "YamiOdymel",
	"Password":  rushia.NewExpr("SHA1(?)", "secretpassword+salt"),
	"Expires":   rushia.NewExpr("NOW() + INTERVAL 1 YEAR"),
	"CreatedAt": rushia.NewExpr("NOW()"),
})
// 等效於：INSERT INTO Users (Username, Password, Expires, CreatedAt) VALUES (?, SHA1(?), NOW() + INTERVAL 1 YEAR, NOW())
```

### 筆數限制

`Limit` 能夠限制 SQL 執行的筆數，如果指定 `10`，那就表示只處理最前面 10 筆資料而非全部（例如：選擇、更新、移除）。如果指定 `10, 20`，那就是忽略前面 10 筆，並處理之後的 20 筆資料（`11, 12... 30`）。

```go
rushia.NewQuery("Users").Limit(10).Update(data)
// 等效於：UPDATE Users SET ... LIMIT 10

rushia.NewQuery("Users").Limit(10, 20).Select(data)
// 等效於：SELECT * from Users LIMIT 10, 20
```

### 筆數偏移

透過 `Offset` 能夠以偏移的方式取得資料，這類似 `Limit` 但參數是反過來的。例如：`10, 20` 則會從 `21` 開始取得 10 筆資料（`21, 22... 30`）。

```go
rushia.NewQuery("Users").Offset(10, 20).Select()
// 等效於：SELECT * from Users LIMIT 10 OFFSET 20
```

### 分頁

`Paginate` 是一個較親近於人類的友善好函式，其參數為：`頁數, 單頁筆數`。例如：`1, 20` 會取得首 20 筆資料，而 `2, 20` 則會取得第二頁的 20 筆資料（基本上為 21 至 40）。

```go
rushia.NewQuery("Users").Paginate(1, 20).Select()
// 等效於：SELECT * from Users LIMIT 0, 20

rushia.NewQuery("Users").Paginate(2, 20).Select()
// 等效於：SELECT * from Users LIMIT 20, 20
```

### 更新

更新一筆資料在 Rushia 中極為簡單，你只需要指定表格名稱還有資料即可。

```go
rushia.NewQuery("Users").Where("Username = ?", "YamiOdymel").Update(rushia.H{
	"Username": "Karisu",
	"Password": "123456",
})
// 等效於：UPDATE Users SET Username = ?, Password = ? WHERE Username = ?
```

### 片段更新

當你希望某些欄位在零值的時候不要進行更新，那麼你就可以使用 `Patch` 來做片段更新（也叫小修補）。

```go
rushia.NewQuery("Users").Where("Username = ?", "YamiOdymel").Patch(rushia.H{
	"Age": 0,
	"Username": "",
	"Password": "123456",
})
// 等效於：UPDATE Users SET Password = ? WHERE Username = ?
```

如果你希望有些欄位雖然是零值（如：`false`、`0`）但仍該在 `Patch` 時照樣更新，那麼就可以使用 `Exclude`。傳入資料型態（如：`reflect.Bool`、`reflect.String`）來以型態排除特定欄位、而字串則表示欲忽略的欄位名稱。

排除的資料型態或欄位會在零值時一樣被更新到資料庫中。

```go
rushia.NewQuery("Users").Where("Username = ?", "YamiOdymel").Exclude("Username", reflect.Int).Patch(rushia.H{
	"Age":      0,
	"Username": "",
	"Password": "123456",
})
// 等效於：UPDATE Users SET Age = ?, Password = ?, Username = ? WHERE Username = ?
```

### 刪除

刪除一筆資料再簡單不過了。

```go
rushia.NewQuery("Users").Where("ID = ?", 1).Delete()
// 等效於：DELETE FROM Users WHERE ID = ?
```

### 選擇與取得

最基本的資料取得在 Rushia 中透過 `Select` 使用。

```go
rushia.NewQuery("Users").Select()
// 等效於：SELECT * FROM Users
```

#### 指定欄位

在 `Select` 中傳遞欄位名稱作為參數，多個欄位由逗點區分，亦能是函式。

```go
rushia.NewQuery("Users").Select("Username", "Nickname")
// 等效於：SELECT Username, Nickname FROM Users

rushia.NewQuery("Users").Select(rushia.NewExpr("COUNT(*) AS Count"))
// 等效於：SELECT COUNT(*) AS Count FROM Users
```

#### 單行資料

如果只想要取得單筆資料，那麼就可以用上 `SelectOne`，這簡單來說就是 `.Limit(1).Select(...)` 的縮寫。

```go
rushia.NewQuery("Users").SelectOne("Username")
// 等效於：SELECT Username FROM Users LIMIT 1
```

#### 排除重複

取得資料的時候可以指定 `Distinct` 過濾重複內容。

```go
rushia.NewQuery("Products").Distinct().Select()
// 等效於：SELECT DISTINCT * FROM Products
```

#### 聯集

可以透過 `Union` 或 `UnionAll` 整合多個表格選取之間的資料。

```go
locationQuery := rushia.NewQuery("Locations").Select()

rushia.NewQuery("Users").Union(locationQuery).Select()
// 等效於：SELECT * FROM Users UNION SELECT * FROM Locations

rushia.NewQuery("Users").UnionAll(locationQuery).Select()
// 等效於：SELECT * FROM Users UNION ALL SELECT * FROM Locations
```

### 選擇是否存在

透過 `Exists` 來執行一個 `SELECT EXISTS`。

```go
rushia.NewQuery("Users").Where("Username = ?", "YamiOdymel").Exists()
// 等效於：SELECT EXISTS(SELECT * FROM Users WHERE Username = ?)
```

### 表格別名

`As` 能夠替目前的查詢語句賦予表格別名，通常會應用在子查詢。若是在表格加入（JOIN）或是一般場景，則可以使用 `NewAlias`。

```go
rushia.NewQuery(NewQuery("Users").Select()).As("Result").Where("Username = ?", "YamiOdymel").Select())
// 等效於：SELECT * FROM (SELECT * FROM Users) AS Result WHERE Username = ?

rushia.NewQuery(rushia.NewAlias("UserFriendRelationships", "relations")).Where("relations.ID = ?", 5).Select()
// 等效於： SELECT * FROM UserFriendRelationships AS relations relations.WHERE ID = ?
```

### 執行生指令

Rushia 已經提供了近乎日常中 80% 會用到的方式，但如果好死不死你想使用的功能在那 20% 之中，我們還提供了原生的方法能讓你直接輸入 SQL 指令執行自己想要的鳥東西。一個最基本的生指令（Raw Query）就像這樣。

其中亦能帶有預置聲明（Prepared Statement），也就是指令中的問號符號替代了原本的值。這能避免你的 SQL 指令遭受注入攻擊。

正如標準的 `NewQuery` 一樣，`NewRawQuery` 也需要透過 `Build` 才能建置。而且 `NewRawQuery` 不能使用所有輔助功能，如：`Limit`、`OrderBy`…。

```go
q := rushia.NewRawQuery("SELECT * FROM Users WHERE ID >= ?", 10)
```

## 條件宣告

透過 Rushia 宣告 `WHERE` 或 `HAVING` 條件也能夠很輕鬆。這裡是實際應用最常派上用場的條件函式：

| SQL 語法                                                 | 使用方式                                                                                   |
| -------------------------------------------------------- | ------------------------------------------------------------------------------------------ |
| `Column = ?`<br>`Column > ?`                             | `.Where("Column = ?", "Value")`<br>`.Where("Column > ?", "Value")`                         |
| `Column = Column`                                        | `.Where("Column = Column")`                                                                |
| `Column IN (?, ?)`<br>`Column NOT IN (?, ?)`             | `.Where("Column IN (?, ?)", "A", "B")`<br>`.Where("Column NOT IN (?, ?)", "A", "B")`       |
| `Column IN (?, ?)`                                       | `.Where("Column IN ?", []interface{}{"A", "B"})`                                           |
| `Column BETWEEN ? AND ?`<br>`Column NOT BETWEEN ? AND ?` | `.Where("Column BETWEEN ? AND ?", 1, 20)`<br>`.Where("Column NOT BETWEEN ? AND ?", 1, 20)` |
| `Column IS NULL`<br>`Column IS NOT NULL`                 | `.Where("Column IS NULL")`<br>`.Where("Column IS NOT NULL")`                               |
| `Column EXISTS Query`<br>`Column NOT EXISTS Query`       | `.Where("Column EXISTS ?", subQuery)`<br>`.Where("Column NOT EXISTS ?", subQuery)`         |
| `Column LIKE ?`<br>`Column NOT LIKE ?`                   | `.Where("Column LIKE ?", "Value")`<br>`.Where("Column NOT LIKE ?", "Value")`               |
| `(Column = Column OR Column = ?)`                        | `.Where("(Column = Column OR Column = ?)", "Value")`                                       |

這些函式總共有幾種變形，分別適用於 `Where`、`OrWhere`、`Having`、`OrHaving`、`JoinWhere`、`OrJoinWhere`。

```go
rushia.NewQuery("Users").Where("ID = ?", 1).Where("Username = ?", "admin").Select()
// 等效於：SELECT * FROM Users WHERE ID = ? AND Username = ?

rushia.NewQuery("Users").Having("ID = ?", 1).Having("Username = ?", "admin").Select()
// 等效於：SELECT * FROM Users HAVING ID = ? AND Username = ?

rushia.NewQuery("Users").Where("ID != CompanyID").Where("DATE(CreatedAt) = DATE(LastLogin)").Select()
// 等效於：SELECT * FROM Users WHERE ID != CompanyID AND DATE(CreatedAt) = DATE(LastLogin)
```

### 預置聲明展開

透過預置聲明（[Prepared Statement](https://en.wikipedia.org/wiki/Prepared_statement)）可以避免 SQL 指令遭受注入攻擊。

在 Rushia 中有個額外的功能，傳入一個 Slice（無論是：`[]interface{}` 或 `[]int`…等）給其中的單個 `?` 會自動展開成為預置聲明。

```go
rushia.NewQuery("Users").Where("ID IN ?", []interface{}{"A", "B", "C"}).Select()
// 等效於：SELECT * FROM Users WHERE ID IN (?, ?, ?)
```

### 脫逸值

與 [mysqljs/mysql](https://github.com/mysqljs/mysql) 套件中的 `??` 雙問號用法相同，你可以透過 `??` 產生出 (\`) 來脫逸字元。這對於欄位名稱很有用。

```go
var ColumnUserID = "ID"
rushia.NewQuery("Users").Where("?? = ?", ColumnUserID, 3).Select()
// 等效於：SELECT * FROM Users WHERE `ID` = ?
```

### 排序

Rushia 亦支援排序功能，如遞增或遞減，亦能擺放函式。

```go
rushia.NewQuery("Users").OrderBy("ID ASC").OrderBy("Login DESC").OrderBy("RAND()").Select()
// 等效於：SELECT * FROM Users ORDER BY ID ASC, Login DESC, RAND()
```

#### 從值排序

也能夠從值進行排序，只需要傳入一個切片即可。

```go
rushia.NewQuery("Users").OrderByField("UserGroup", "SuperUser", "Admin", "Users").Select()
// 等效於：SELECT * FROM Users ORDER BY FIELD (UserGroup, ?, ?, ?) ASC
```

### 分組

簡單的透過 `GroupBy` 就能夠將資料由指定欄位分組。

```go
rushia.NewQuery("Users").GroupBy("Name").Select()
// 等效於：SELECT * FROM Users GROUP BY Name
```

### 加入表格

Rushia 支援多種表格加入方式，如：`InnerJoin`、`LeftJoin`、`RightJoin`、`NaturalJoin`、`CrossJoin`。在 Join 時，最後一個參數預設可以擺入條件式。

```go
rushia.
	NewQuery("Products").
	LeftJoin("Users", "Products.TenantID = Users.TenantID").
	Select("Users.Name", "Products.ProductName")
// 等效於：SELECT Users.Name, Products.ProductName FROM Products AS Products LEFT JOIN Users AS Users ON (Products.TenantID = Users.TenantID)
```

但你也可以將加入條件式拆開放到後面定義。

```go
rushia.
	NewQuery("Products").
	LeftJoin("Users").
	JoinWhere("Products.TenantID = Users.TenantID")
	Select("Users.Name", "Products.ProductName")
// 等效於：SELECT Users.Name, Products.ProductName FROM Products AS Products LEFT JOIN Users AS Users ON (Products.TenantID = Users.TenantID)
```

#### 條件限制

你亦能透過 `JoinWhere` 或 `OrJoinWhere` 擴展表格加入的限制條件，使用時這個條件總是會加到最後一個 `Join` 的表格。

```go
rushia.
	NewQuery("Products").
	LeftJoin("Users", "Products.TenantID = Users.TenantID").
	OrJoinWhere("Users.TenantID = ?", 5).
	Select("Users.Name", "Products.ProductName")
// 等效於：SELECT Users.Name, Products.ProductName FROM Products AS Products LEFT JOIN Users AS Users ON (Products.TenantID = Users.TenantID OR Users.TenantID = ?)
```

### 子指令

Rushia 支援複雜的子指令，將一個指令語法帶入當成值使用就能夠將其當作子指令。

```go
subQuery := rushia.NewQuery("VIPUsers").Select("UserID")

rushia.NewQuery("Users").Where("ID IN ?", subQuery).Select()
// 等效於：SELECT * FROM Users WHERE ID IN (SELECT UserID FROM VIPUsers)
```

#### 插入

插入新資料時也可以使用子指令，但必須確保子指令只會回傳一個欄位與單行資料。

```go
subQuery := rushia.NewQuery("Users").Where("ID = ?", 6).SelectOne("Name")

rushia.NewQuery("Products").Insert(rushia.H{
	"ProductName": "測試商品",
	"UserID":      subQuery,
	"LastUpdated": rushia.NewExpr("NOW()")
})
// 等效於：INSERT INTO Products (ProductName, UserID, LastUpdated) VALUES (?, (SELECT Name FROM Users WHERE ID = 6 LIMIT 1), NOW())
```

#### 加入

就算是加入表格的時候也可以用上子指令，但你需要使用 `As` 為子指令建立別名。

```go
subQuery := rushia.NewQuery("Users").As("Users").Where("Active = ?", 1).Select()

rushia.
	NewQuery("Products").
	LeftJoin(subQuery, "Products.UserID = Users.ID").
	Select("Users.Username", "Products.ProductName")
// 等效於：SELECT Users.Username, Products.ProductName FROM Products AS Products LEFT JOIN (SELECT * FROM Users WHERE Active = ?) AS Users ON Products.UserID = Users.ID
```

### 子指令置換

在使用表達式或生指令的時候可能會希望用上子指令，這個時候可以傳入一個子指令則會替換相對應的 `?` 預置變數。

```go
subQuery := rushia.NewQuery("Locations").Select()
rawQuery := rushia.NewRawQuery("SELECT UserID FROM Users WHERE EXISTS (?)", subQuery)

NewQuery("Products").WhereExists(rawQuery).Select()
// 等效於：SELECT * FROM Products WHERE EXISTS (SELECT UserID FROM Users WHERE EXISTS (SELECT * FROM Locations))
```

### 指令關鍵字

Rushia 也支援設置指令關鍵字。

```go
rushia.NewQuery("Users").SetQueryOption("FOR UPDATE").Select()
// 等效於：SELECT * FROM Users FOR UPDATE

rushia.NewQuery("Users").SetQueryOption("SQL_NO_CACHE").Select()
// 等效於：SELECT SQL_NO_CACHE * FROM Users

rushia.NewQuery("Users").SetQueryOption("LOW_PRIORITY", "IGNORE").Insert(data)
// Gives: INSERT LOW_PRIORITY IGNORE INTO Users ...
```

## 複雜場景範例

```go
jobHistories := rushia.NewQuery("JobHistories").
	Where("DepartmentID BETWEEN ? AND ?", 50, 100).
	Select("JobID")
jobs := rushia.NewQuery("Jobs").
	Where("JobID IN ?", jobHistories).
	GroupBy("JobID").
	Select("JobID", "AVG(MinSalary) AS MyAVG")
maxAverage := rushia.NewQuery(jobs).
	As("SS").
	Select("MAX(MyAVG)")
employees := rushia.NewQuery("Employees").
	GroupBy("JobID").
	Having("AVG(Salary) < ?", maxAverage).
	Select("JobID", "AVG(Salary)")

// 等效於：
// SELECT JobID,
//        AVG(Salary)
// FROM   Employees
// HAVING AVG(Salary) < (SELECT MAX(MyAVG)
//                       FROM   (SELECT JobID,
//                                      AVG(MinSalary) AS MyAVG
//                               FROM   Jobs
//                               WHERE  JobID IN (SELECT JobID
//                                                FROM   JobHistories
//                                                WHERE  DepartmentID BETWEEN 50
//                                                       AND 100
//                                               )
//                               GROUP  BY JobID) AS SS)
// GROUP  BY job_id;

agents := rushia.NewQuery("Agents").
	Where("Commission < ?", 0.12).
	Select()
customers := rushia.NewQuery("Customers").
	Where("Grade = ?", 3).
	Where("CustomerCountry <> ?", "India").
	Where("OpeningAmount < ?", 7000).
	Where("EXISTS ?", agents).
	Select("OutstandingAmount")
orders := rushia.NewQuery("Orders").
	Where("OrderAmount > ?", 2000).
	Where("OrderDate < ?", "01-SEP-08").
	Where("AdvanceAmount < ?", rushia.NewExpr("ANY (?)", customers)).
	Select("OrderNum", "OrderDate", "OrderAmount", "AdvanceAmount")

// 等效於：
// SELECT OrderNum,
//        OrderDate,
//        OrderAmount,
//        AdvanceAmount
// FROM   Orders
// WHERE  OrderAmount > 2000
//        AND OrderDate < '01-SEP-08'
//        AND AdvanceAmount < ANY (SELECT OutstandingAmount
//                                 FROM   Customers
//                                 WHERE  Grade = 3
//                                        AND CustomerCountry <> 'India'
//                                        AND OpeningAmount < 7000
//                                        AND EXISTS (SELECT *
//                                                    FROM   Agents
//                                                    WHERE  Commission < 0.12));
```

## 相關連結

這裡是 Rushia 受啟發，或是和資料庫有所關聯的連結。

-   [kisielk/sqlstruct](http://godoc.org/github.com/kisielk/sqlstruct)
-   [jmoiron/sqlx](https://github.com/jmoiron/sqlx)
-   [russross/meddler](https://github.com/russross/meddler)
-   [jinzhu/gorm](https://github.com/jinzhu/gorm)
-   [doug-martin/goqu](https://github.com/doug-martin/goqu)
-   [gocraft/dbr](https://github.com/gocraft/dbr)
-   [go-ozzo/ozzo-dbx](https://github.com/go-ozzo/ozzo-dbx)
-   [kyleconroy/sqlc](https://github.com/kyleconroy/sqlc)
