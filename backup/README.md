# Rushia [![GoDoc](https://godoc.org/github.com/teacat/rushia?status.svg)](https://godoc.org/github.com/teacat/rushia) [![Coverage Status](https://coveralls.io/repos/github/teacat/rushia/badge.svg?branch=master)](https://coveralls.io/github/teacat/rushia?branch=master) [![Build Status](https://travis-ci.org/teacat/rushia.svg?branch=master)](https://travis-ci.org/teacat/rushia) [![Go Report Card](https://goreportcard.com/badge/github.com/teacat/rushia)](https://goreportcard.com/report/github.com/teacat/rushia)

A MySQL query builder that's way- better than most the [ORM](https://zh.wikipedia.org/wiki/%E5%AF%B9%E8%B1%A1%E5%85%B3%E7%B3%BB%E6%98%A0%E5%B0%84) that written in [Golang](https://golang.org/). Flexible and no struct tags needed. The original idea was from [PHP-MySQLi-Database-Class](https://github.com/joshcam/PHP-MySQLi-Database-Class) and [Laravel Query Builder](https://laravel.com/docs/8.x/queries) with extra functions.

This is a query builder without any database connection implmentation, fits for any library as base.

## Features

-   Fully functional.
-   Easy to use.
-   SQL query builder.
-   Table migration.
-   Struct mapping.
-   Method chaining.
-   Sub query supported.
-   [Prepared Statement](https://en.wikipedia.org/wiki/Prepared_statement) supported to prevent 99.9% of SQL injection.

## Why?

[Gorm](https://github.com/jinzhu/gorm) is a famous [ORM](https://zh.wikipedia.org/wiki/%E5%AF%B9%E8%B1%A1%E5%85%B3%E7%B3%BB%E6%98%A0%E5%B0%84) in [Golang](https://golang.org/) community, it's really good to use until you meet the JOINs with complex quries. Rushia solved the problem by making a better query builder and omits the dependency with structs.

## Indexes

-   [Installation](#安裝方式)
-   [Naming convention](#命名建議)
-   [NULL values](#null-值)
-   [Usages](#使用方式)
    -   [Mapping](#映射)
        -   [Omit](#省略)
    -   [Insert](#插入)
        -   [Replace](#覆蓋)
        -   [函式](#函式)
        -   [當重複時](#當重複時)
        -   [多筆資料](#多筆資料)
    -   [筆數限制](#筆數限制)
    -   [筆數偏移](#筆數偏移)
    -   [更新](#更新)
        -   [片段更新](#片段更新)
    -   [選擇與取得](#選擇與取得)
        -   [指定欄位](#指定欄位)
        -   [單行資料](#單行資料)
    -   [執行生指令](#執行生指令)
        -   [進階方式](#進階方式)
    -   [條件宣告](#條件宣告)
        -   [擁有](#擁有)
        -   [欄位比較](#欄位比較)
        -   [自訂運算子](#自訂運算子)
        -   [介於／不介於](#介於不介於)
        -   [於清單／不於清單內](#於清單不於清單內)
        -   [或／還有或](#或還有或)
        -   [空值](#空值)
        -   [時間戳](#時間戳)
            -   [相對](#相對)
            -   [日期](#日期)
            -   [時間](#時間)
        -   [生條件](#生條件)
            -   [條件變數](#條件變數)
    -   [刪除](#刪除)
    -   [排序](#排序)
        -   [從值排序](#從值排序)
    -   [群組](#群組)
    -   [加入](#加入)
        -   [條件限制](#條件限制)
    -   [子指令](#子指令)
        -   [選擇／取得](#選擇取得)
        -   [插入](#插入-1)
        -   [加入](#加入-1)
        -   [存在／不存在](#存在不存在)
    -   [輔助函式](#輔助函式)
        -   [總筆數](#總筆數)
    -   [鎖定表格](#鎖定表格)
    -   [指令關鍵字](#指令關鍵字)
        -   [多個選項](#多個選項)
-   [表格建構函式](#表格建構函式)

## Installation

Install the package via `go get` command.

```bash
$ go get github.com/teacat/rushia
```

## NULL values

We suggest you to make all the columns in the database as non-nullable since Golang sucks at supporting NULL fields.

# Usage

Rushia is easy to use, it's kinda like a SQL query but simplized.

# Create query

A basic query starts from `NewQuery(...)` with a table name or a sub query. A complex example with sub query will be mentioned in the later chapters.

```
q := rushia.NewQuery("Users")
```

# Copy

adasdsadasdasdasd

# Build query

Execute the `Build` function when you completed a query with `Select`, `Exists`, `Replace`, `Update`, `Delete`... etc. To get the generated query and the params.

```go
query, params := rushia.Build(rushia.NewQuery("Users").Select())
// Equals: SELECT * FROM Users
```

# Use with the other libraries

Since Rushia is just a SQL Builder, you are able to use it with any other database execution libraries. For example with [jmoiron/sqlx](https://github.com/jmoiron/sqlx):

```go
// Initialize a SQLX connection.
db := sqlx.Open("mysql", "root:password@tcp(localhost:3306)/db")

// Build the query via Rushia.
q := rushia.NewQuery("Users").Where("Username", "YamiOdymel").Select()
query, params := rushia.Build(q)

// Pass the query and the parameters to SQLX to execute.
rows, err := db.Query(query, params...)
// Equals: SELECT * FROM Users WHERE Username = ?
```

Or [go-gorm/gorm](https://github.com/go-gorm/gorm) if you like:

```go
// Initialize a Gorm connection.
db, err := gorm.Open(mysql.Open("root:password@tcp(localhost:3306)/db"), &gorm.Config{})

// Build the query via Rushia.
q := rushia.NewQuery("Users").Where("Username", "YamiOdymel").Select()
query, params := rushia.Build(q)

// Pass the query and the parameters to Gorm to execute.
db.Raw(query, params...).Scan(&myUser)
// Equals: SELECT * FROM Users WHERE Username = ?
```

## Mapping

You are able to pass a struct to `Insert` or `Update` functions and it will be automatically applies the field names and the values into the query.

But be careful! It won't converts the `CamelCase` field names into `snake_cases`.

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
// Equals: INSERT INTO Users (Username, Password) VALUES (?, ?)
```

### Struct tag

You could omit or rename a field by specify the `rushia` struct tag.

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
// Equals：INSERT INTO Users (real_name, Password) VALUES (?, ?)
```

### Omit

Ignore the fields in the SQL query by using `Omit`.

```go
type User struct {
	Username string
	Password string
}
u := User{
	Username: "YamiOdymel",
	Password: "test",
}
rushia.NewQuery("Users").Omit("Username").Insert(u)
// Equals: INSERT INTO Users (Password) VALUES (?)
```

## Insert

Rushia provides a shorthand `H` alias for `H`, it the same as [`gin.H`](https://pkg.go.dev/github.com/gin-gonic/gin#H). You can pass a struct or a `H`, `H` into a Insert query.

```go
rushia.NewQuery("Users").Insert(rushia.H{
	"Username": "YamiOdymel",
	"Password": "test",
})
// Equals: INSERT INTO Users (Username, Password) VALUES (?, ?)

rushia.NewQuery("Users").Insert(rushia.H{
	"Username": "YamiOdymel",
	"Password": "test",
})
// Equals: INSERT INTO Users (Username, Password) VALUES (?, ?)
```

### Replace

The usage Replace is the same as Insert but it deletes the duplicated data and insert a new one. It's dangerous for any data that contains foregin keys. To be safe, use `OnDuplicate` (`ON DUPLICATE KEY UPDATE`) instead.

```go
rushia.NewQuery("Users").Replace(rushia.H{
	"Username": "YamiOdymel",
	"Password": "test",
})
// Equals: REPLACE INTO Users (Username, Password) VALUES (?, ?)
```

### Expression

By using `NewExpr` to create an Expression, you can represent a complex value that accepts a raw query, and the parameters to create functions such as: `SHA1()` or `NOW()` and intervals.

```go
rushia.NewQuery("Users").Insert(rushia.H{
	"Username":  "YamiOdymel",
	"Password":  rushia.NewExpr("SHA1(?)", "secretpassword+salt"),
	"Expires":   rushia.NewExpr("NOW() + INTERVAL 1 YEAR"),
	"CreatedAt": rushia.NewExpr("NOW()"),
})
// Equals: INSERT INTO Users (Username, Password, Expires, CreatedAt) VALUES (?, SHA1(?), NOW() + INTERVAL 1 YEAR, NOW())
```

### On duplicate

Rushia supports `ON DUPLICATE KEY UPDATE` to update the specified data when it's duplicated on insertion. It's like `Replace` but it won't delete the duplicated data but update it instead.

```go
rushia.NewQuery("Users").As("New").OnDuplicate(rushia.H{
	"UpdatedAt": rushia.NewExpr("New.UpdatedAt"),
}).Insert(rushia.H{
	"Username":  "YamiOdymel",
	"UpdatedAt": rushia.NewExpr("NOW()"),
})
// Equals: INSERT INTO Users (Username, UpdatedAt) VALUES (?, NOW()) ON DUPLICATE KEY UPDATE UpdatedAt = New.UpdatedAt

rushia.NewQuery("Users").OnDuplicate(rushia.H{
	"UpdatedAt": rushia.NewExpr("VALUES(UpdatedAt)"),
}).Insert(rushia.H{
	"Username":  "YamiOdymel",
	"UpdatedAt": rushia.NewExpr("NOW()"),
})
// CAUTION! `VALUES` has been deprecated since MySQL 8.0.20! Use the above example instead!
// Equals: INSERT INTO Users (Username, UpdatedAt) VALUES (?, NOW()) ON DUPLICATE KEY UPDATE UpdatedAt = VALUES(UpdatedAt)
```

### Insert multiple

By passing a `[]H` or `[]map[string]interface{}` to insert multiple values at once.

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
// Equals: INSERT INTO Users (Username, Password) VALUES (?, ?), (?, ?)
```

## Limit

`Limit` limits the rows to process (Select, Update, Delete). Only the first `10` rows will be affected if it was set to `10`.

```go
rushia.NewQuery("Users").Limit(10).Update(data)
// Equals: UPDATE Users SET ... LIMIT 10

rushia.NewQuery("Users").Limit(10, 20).Select(data)
// Equals: SELECT * from Users LIMIT 10, 20
```

## Offset

The usage of `Offset` is a bit like pagination, the arguments work as `count, last_index`. If `Offset(10, 20)` was called, the result `21, 22... 30` will be fetched.

```go
rushia.NewQuery("Users").Offset(10, 20).Select()
// Equals: SELECT * from Users LIMIT 10 OFFSET 20
```

## Update

To update a data in Rushia is easy as a rocket launch (wat? (todo: update this description later)).

```go
rushia.NewQuery("Users").Where("Username", "YamiOdymel").Update(rushia.H{
	"Username": "Karisu",
	"Password": "123456",
})
// Equals: UPDATE Users SET Username = ?, Password = ? WHERE Username = ?
```

### Patch

By using `Patch`, it's possible to ignore the zero value fields while updating.

```go
rushia.NewQuery("Users").Where("Username", "YamiOdymel").Patch(rushia.H{
	"Age": 0,
	"Username": "",
	"Password": "123456",
})
// Equals: UPDATE Users SET Password = ? WHERE Username = ?
```

With `Exclude`, you can also exclude the fields to force it update even if it's a zero value (e.g. `false`, `0`). Passing strings as column names to exclude, and `reflect.Kind` to exclude by data types.

Any fields that was excluded will still be updated even if it's a zero value.

```go
rushia.NewQuery("Users").Where("Username", "YamiOdymel").Exclude("Username", reflect.Int).Patch(rushia.H{
	"Age":      0,
	"Username": "",
	"Password": "123456",
})
// Equals: UPDATE Users SET Age = ?, Password = ?, Username = ? WHERE Username = ?
```

## Select

To simpliy get a data by using `Select`.

```go
rushia.NewQuery("Users").Select()
// Equals: SELECT * FROM Users
```

### Specify columns

Specify the columns to select in the `Select` arguments, It colud also be a expression.

```go
rushia.NewQuery("Users").Select("Username", "Nickname")
// Equals: SELECT Username, Nickname FROM Users

rushia.NewQuery("Users").Select(NewExpr("COUNT(*) AS Count"))
// Equals: SELECT COUNT(*) AS Count FROM Users
```

### Select One

To get a single row data, use `SelectOne`. It's a shorthand for `.Limit(1).Select(...)`.

```go
rushia.NewQuery("Users").SelectOne("Username")
// Equals: SELECT Username FROM Users LIMIT 1
```

## Union

`Union` or `UnionAll` allows you to merge the data between different table selections.

```go
locationQuery := rushia.NewQuery("Locations").Select()

rushia.NewQuery("Users").Union(locationQuery).Select()
// Equals: SELECT * FROM Users UNION SELECT * FROM Locations

rushia.NewQuery("Users").UnionAll(locationQuery).Select()
// Equals: SELECT * FROM Users UNION ALL SELECT * FROM Locations
```

## Table alias

When creating a sub query or table joins, you might need `As` to assign an alias to a table.

```go
rushia.NewQuery("Users").As("U").Select()
// Equals: SELECT * FROM Users AS U
```

## Raw Query

Rushia provides you the most 80% things you will use, but if you are in the bad luck to request for the rest 20%, the only hope is to use Raw Query.

A raw query does also support the prepared statement, to replace the value as `?` to prevent the SQL injection.

`NewRawQuery` is the same as `NewQuery` that required to be `Build`, and the helper functions such as: `Limit`, `OrderBy`...etc, are not able to be used.

```go
q := rushia.NewRawQuery("SELECT * FROM Users WHERE ID >= ?", 10)
```

## Where

To define a `WHERE` condition in Rushia is a piece of cake! The basic `WHERE AND` works like:

```go
rushia.NewQuery("Users").Where("ID", 1).Where("Username", "admin").Select()
// Equals: SELECT * FROM Users WHERE ID = ? AND Username = ?
```

### Having

It's possible to use `HAVING` with `WHERE` conditions.

```go
rushia.NewQuery("Users").Where("ID", 1).Having("Username", "admin").Select()
// Equals: SELECT * FROM Users WHERE ID = ? HAVING Username = ?
```

### Column comparison

To judge between two columns:

```go
// ✓ DO.
rushia.NewQuery("Users").Where("LastLogin = CreatedAt").Select()
// ✖ DON'T!
rushia.NewQuery("Users").Where("LastLogin", "CreatedAt").Select()

// Equals: SELECT * FROM Users WHERE LastLogin = CreatedAt
```

### Operators

You are able to change the operators (e.g. >=, <=, <>) in `Where` and `Having`:

```go
rushia.NewQuery("Users").Where("ID", ">=", 50).Select()
// Equals: SELECT * FROM Users WHERE ID >= ?
```

### Between

Use `BETWEEN` make sure a value was in (or not) a specified range.

```go
rushia.NewQuery("Users").Where("ID", "BETWEEN", 0, 20).Select()
// Equals: SELECT * FROM Users WHERE ID BETWEEN ? AND ?

rushia.NewQuery("Users").Where("ID", "NOT BETWEEN", 0, 20).Select()
// Equals: SELECT * FROM Users WHERE ID NOT BETWEEN ? AND ?
```

### In

Use `IN` to make sure the value was in (or not) the list.

```go
rushia.NewQuery("Users").Where("ID", "IN", 1, 5, 27, -1, "d").Select()
// Equals: SELECT * FROM Users WHERE ID IN (?, ?, ?, ?, ?)

rushia.NewQuery("Users").Where("ID", "NOT IN", 1, 5, 27, -1, "d").Select()
// Equals: SELECT * FROM Users WHERE ID NOT IN (?, ?, ?, ?, ?)
```

### Or

With `Where` and `Having`, it creates `AND` conditions. If you want to create a `OR`, simply using `OrWhere` or `OrHaving`.

```go
rushia.NewQuery("Users").Where("FirstNamte", "John").OrWhere("FirstNamte", "Peter").Select()
// Equals: SELECT * FROM Users WHERE FirstName = ? OR FirstName = ?
```

You'll need to manually write a query if you are triying to create a condition group such as `A = B OR (A = C OR A = D)`.

```go
rushia.NewQuery("Users").Where("A = B").OrWhere("(A = C OR A = D)").Select()
// Equals: SELECT * FROM Users WHERE A = B OR (A = C OR A = D)
```

### NULL

To verify if a value is a NULL value or not:

```go
// ✓ DO.
rushia.NewQuery("Users").Where("LastName", "IS", nil).Select()
// ✖ DON'T!
rushia.NewQuery("Users").Where("LastName", "NULL").Select()

// Equals: SELECT * FROM Users WHERE LastName IS NULL
```

### Raw condition

You are able to pass a raw query into a condition.

```go
rushia.NewQuery("Users").Where("ID != CompanyID").Where("DATE(CreatedAt) = DATE(LastLogin)").Select()
// Equals: SELECT * FROM Users WHERE ID != CompanyID AND DATE(CreatedAt) = DATE(LastLogin)
```

#### Condition parameters

A raw query condition can also be used with the `?` parameters.

```go
rushia.NewQuery("Users").Where("(ID = ? OR ID = ?)", 6, 2).Where("Login", "Mike").Select()
// Equals: SELECT * FROM Users WHERE (ID = ? OR ID = ?) AND Login = ?
```

## Distinct

Specifing `Distinct` to eliminate the duplicate rows while fetching the data.

```go
rushia.NewQuery("Products").Distinct().Select()
// Equals: SELECT DISTINCT * FROM Products
```

## Delete

Deletes everything! Remember to add a condition to prevent it really deletes everything.

```go
rushia.NewQuery("Users").Where("ID", 1).Delete()
// Equals: DELETE FROM Users WHERE ID = ?
```

## Order

Ordering is also supported in Rushia and can be used with functions.

```go
rushia.NewQuery("Users").OrderBy("ID", "ASC").OrderBy("Login", "DESC").OrderBy("RAND()").Select()
// Equals: SELECT * FROM Users ORDER BY ID ASC, Login DESC, RAND()
```

### Order by field

Or ordering by custom field values:

```go
rushia.NewQuery("Users").OrderByField("UserGroup", "SuperUser", "Admin", "Users").Select()
// Equals: SELECT * FROM Users ORDER BY FIELD (UserGroup, ?, ?, ?)
```

## Group by

The result can also be grouped with `GroupBy`.

```go
rushia.NewQuery("Users").GroupBy("Name").Select()
// Equals: SELECT * FROM Users GROUP BY Name
```

## Join

Rushia supports multiple ways to join the tables, such as: `InerrJoin`, `LeftJoin`, `RightJoin`, `NaturalJoin`, `CrossJoin`.

```go
rushia.
	NewQuery("Products").
	LeftJoin("Users", "Products.TenantID = Users.TenantID").
	Where("Users.ID", 6).
	Select("Users.Name", "Products.ProductName")
// Equals: SELECT Users.Name, Products.ProductName FROM Products AS Products LEFT JOIN Users AS Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?
```

### Join condition

With `JoinWhere` or `OrJoinWhere` to expand the conditions for the table joins.

```go
rushia.
	NewQuery("Products").
	LeftJoin("Users", "Products.TenantID = Users.TenantID").
	OrJoinWhere("Users", "Users.TenantID", 5).
	Select("Users.Name", "Products.ProductName")
// Equals: SELECT Users.Name, Products.ProductName FROM Products AS Products LEFT JOIN Users AS Users ON (Products.TenantID = Users.TenantID OR Users.TenantID = ?)
```

## Sub query

Rushia supports nested query which is called Sub Query. Use a query as a value to make it sub query.

```go
subQuery := rushia.NewQuery("VIPUsers").Select("UserID")

rushia.NewQuery("Users").Where("ID", "IN", subQuery).Select()
// Equals: SELECT * FROM Users WHERE ID IN (SELECT UserID FROM VIPUsers)
```

### insert

To insert a value from a sub query, simply use the query as a value.

```go
subQuery := rushia.NewQuery("Users").Where("ID", 6).Select("Name")

rushia.NewQuery("Products").Insert(rushia.H{
	"ProductName": "測試商品",
	"UserID":      subQuery,
	"LastUpdated": rushia.NewExpr("NOW()")
})
// Equals: INSERT INTO Products (ProductName, UserID, LastUpdated) VALUES (?, (SELECT Name FROM Users WHERE ID = 6), NOW())
```

### join

Join a table from a sub query is possible, but requires to assign an alias to the sub query by using `As`.

```go
subQuery := rushia.NewQuery("Users").As("Users").Where("Active", 1).Select()

rushia.
	NewQuery("Products").
	LeftJoin(subQuery, "Products.UserID = Users.ID").
	Select("Users.Username", "Products.ProductName")
// Equals: SELECT Users.Username, Products.ProductName FROM Products AS Products LEFT JOIN (SELECT * FROM Users WHERE Active = ?) AS Users ON Products.UserID = Users.ID
```

### exists

To see if a value does exist or not by using the sub query as a `WHERE` condition statement.

```go
subQuery := rushia.NewQuery("Users").Where("Company", "測試公司").Select("UserID")

rushia.NewQuery("Products").Where("EXISTS", subQuery).Select()
// Equals: SELECT * FROM Products WHERE EXISTS (SELECT UserID FROM Users WHERE Company = ?)
```

## 輔助函式

Rushia 有提供一些輔助用的函式協助你除錯、紀錄，或者更加地得心應手。

## Set query options

You can set the query options with Rushia.

```go
rushia.NewQuery("Users").SetQueryOption("FOR UPDATE").Select()
// Equals: SELECT * FROM Users FOR UPDATE

rushia.NewQuery("Users").SetQueryOption("SQL_NO_CACHE").Select()
// Equals: SELECT SQL_NO_CACHE * FROM Users

rushia.NewQuery("Users").SetQueryOption("LOW_PRIORITY", "IGNORE").Insert(data)
// Gives: INSERT LOW_PRIORITY IGNORE INTO Users ...
```

# 表格建構函式

Rushia 除了基本的資料庫函式可供使用外，還能夠建立一個表格並且規劃其索引、外鍵、型態。

```go
migration := rushia.NewMigration()

migration.Table("Users").Column("Username").Varchar(32).Primary().Create()
// Equals: CREATE TABLE Users (Username VARCHAR(32) NOT NULL PRIMARY KEY) ENGINE=INNODB
```

| 數值      | 字串       | 二進制    | 檔案資料   | 時間      | 浮點數  | 固組 |
| --------- | ---------- | --------- | ---------- | --------- | ------- | ---- |
| TinyInt   | Char       | Binary    | Blob       | Date      | Double  | Enum |
| SmallInt  | Varchar    | VarBinary | MediumBlob | DateTime  | Decimal | Set  |
| MediumInt | TinyText   | Bit       | LongBlob   | Time      | Float   |      |
| Int       | Text       |           |            | Timestamp |         |      |
| BigInt    | MediumText |           |            | Year      |         |      |
|           | LongText   |           |            |           |         |      |

# References

Let's see what inspired Rushia.

-   [kisielk/sqlstruct](http://godoc.org/github.com/kisielk/sqlstruct)
-   [jmoiron/sqlx](https://github.com/jmoiron/sqlx)
-   [russross/meddler](https://github.com/russross/meddler)
-   [jinzhu/gorm](https://github.com/jinzhu/gorm)