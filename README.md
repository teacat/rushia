# Rushia [台灣正體](./README-tw.md) [![GoDoc](https://godoc.org/github.com/teacat/rushia/v3?status.svg)](https://godoc.org/github.com/teacat/rushia/v3) [![Coverage Status](https://coveralls.io/repos/github/teacat/rushia/badge.svg?branch=master)](https://coveralls.io/github/teacat/rushia?branch=master) [![Build Status](https://travis-ci.org/teacat/rushia.svg?branch=master)](https://travis-ci.org/teacat/rushia) [![Go Report Card](https://goreportcard.com/badge/github.com/teacat/rushia)](https://goreportcard.com/report/github.com/teacat/rushia)

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

## Installation

Install the package via `go get` command.

```bash
$ go get github.com/teacat/rushia/v2
```

## Usage

Rushia is easy to use, it's kinda like a SQL query but simplized.

### Create query

A basic query starts from `NewQuery(...)` with a table name or a sub query. A complex example with sub query will be mentioned in the later chapters.

```go
q := rushia.NewQuery("Users")
```

### Copy query

By default, Rushia creates a pointer query where you will always modify to the same query. To copy the query with existing rules simply use `Copy`.

```go
a := rushia.NewQuery("Users")
a.Where("Type = ?", "VIP")

b := a.Copy()
b.Where("Name = ?", "YamiOdymel")

Build(a.Select())
// Equals: SELECT * FROM Users WHERE Type = ?
Build(b.Select())
// Equals: SELECT * FROM Users WHERE Type = ? AND Name = ?
```

### Build query

Execute the `Build` function when you completed a query with `Select`, `Exists`, `Replace`, `Update`, `Delete`... etc. To get the generated query and the params.

```go
query, params := rushia.Build(rushia.NewQuery("Users").Select())
// Equals: SELECT * FROM Users
```

### Use with the other libraries

Since Rushia is just a SQL Builder, you are able to use it with any other database execution libraries. For example with [jmoiron/sqlx](https://github.com/jmoiron/sqlx):

```go
// Initialize a SQLX connection.
db, err := sqlx.Open("mysql", "root:password@tcp(localhost:3306)/db")

// Build the query via Rushia.
q := rushia.NewQuery("Users").Where("Username = ?", "YamiOdymel").Select()
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
q := rushia.NewQuery("Users").Where("Username = ?", "YamiOdymel").Select()
query, params := rushia.Build(q)

// Pass the query and the parameters to Gorm to execute.
db.Raw(query, params...).Scan(&myUser)
// Equals: SELECT * FROM Users WHERE Username = ?
```

### Struct mapping

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

#### Struct tag

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
	Age      int    `rushia:"my_age"`
}
u := User{
	Username: "YamiOdymel",
	Password: "test",
	Age     : "32"
}
rushia.NewQuery("Users").Omit("Username", "my_age").Insert(u)
// Equals: INSERT INTO Users (Password) VALUES (?)
```

### Insert

Rushia provides a shorthand `H` alias, stands for `map[string]interface{}`. It's the same as [`gin.H`](https://pkg.go.dev/github.com/gin-gonic/gin#H). You can pass a struct or a `H`, `H` into a Insert query.

```go
rushia.NewQuery("Users").Insert(rushia.H{
	"Username": "YamiOdymel",
	"Password": "test",
})
// Equals: INSERT INTO Users (Username, Password) VALUES (?, ?)

rushia.NewQuery("Users").Insert(map[string]interface{
	"Username": "YamiOdymel",
	"Password": "test",
})
// Equals: INSERT INTO Users (Username, Password) VALUES (?, ?)
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

### Replace

The usage Replace is the same as Insert but it deletes the duplicated data and insert a new one. It's dangerous for any data that contains foregin keys. To be safe, use `OnDuplicate` (`ON DUPLICATE KEY UPDATE`) instead.

```go
rushia.NewQuery("Users").Replace(rushia.H{
	"Username": "YamiOdymel",
	"Password": "test",
})
// Equals: REPLACE INTO Users (Username, Password) VALUES (?, ?)
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
// Equals: INSERT INTO Users (Username, UpdatedAt) VALUES (?, NOW()) AS New ON DUPLICATE KEY UPDATE UpdatedAt = New.UpdatedAt

rushia.NewQuery("Users").OnDuplicate(rushia.H{
	"UpdatedAt": rushia.NewExpr("VALUES(UpdatedAt)"),
}).Insert(rushia.H{
	"Username":  "YamiOdymel",
	"UpdatedAt": rushia.NewExpr("NOW()"),
})
// CAUTION! `VALUES` has been deprecated since MySQL 8.0.20! Use the above example instead!
// Equals: INSERT INTO Users (Username, UpdatedAt) VALUES (?, NOW()) ON DUPLICATE KEY UPDATE UpdatedAt = VALUES(UpdatedAt)
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

### Limit

`Limit` limits the rows to process (Select, Update, Delete). Only the first `10` rows will be affected if it was set to `10`. If `10, 20` was specified, it will skip the first 10 results and process the next 20 results.

```go
rushia.NewQuery("Users").Limit(10).Update(data)
// Equals: UPDATE Users SET ... LIMIT 10

rushia.NewQuery("Users").Limit(10, 20).Select(data)
// Equals: SELECT * from Users LIMIT 10, 20
```

### Offset

The usage of `Offset` works a bit like `Limit` but opposite arguments. If `10, 20` was specified, it skips the first 20 results and deal with the rest 10 results.

```go
rushia.NewQuery("Users").Offset(10, 20).Select()
// Equals: SELECT * from Users LIMIT 10 OFFSET 20
```

### Paginate

`Paginate` is human-friendly, the argument works as `page, count`. With `1, 20` it fetches the first 20 results, with `2, 20` it fetches the other 20 results from page 2 (basically from 21 to 40).

```go
rushia.NewQuery("Users").Paginate(1, 20).Select()
// Equals: SELECT * from Users LIMIT 0, 20

rushia.NewQuery("Users").Paginate(2, 20).Select()
// Equals: SELECT * from Users LIMIT 20, 20
```

### Update

To update a data in Rushia is easy as a rocket launch (wat? (todo: update this description later)).

```go
rushia.NewQuery("Users").Where("Username = ?", "YamiOdymel").Update(rushia.H{
	"Username": "Karisu",
	"Password": "123456",
})
// Equals: UPDATE Users SET Username = ?, Password = ? WHERE Username = ?
```

### Patch

By using `Patch`, it's possible to ignore the zero value fields while updating.

```go
rushia.NewQuery("Users").Where("Username = ?", "YamiOdymel").Patch(rushia.H{
	"Age": 0,
	"Username": "",
	"Password": "123456",
})
// Equals: UPDATE Users SET Password = ? WHERE Username = ?
```

With `Exclude`, you can also exclude the fields to force it update even if it's a zero value (e.g. `false`, `0`). Passing strings as column names to exclude, and `reflect.Kind` to exclude by data types.

Any fields that was excluded will still be updated even if it's a zero value.

```go
rushia.NewQuery("Users").Where("Username = ?", "YamiOdymel").Exclude("Username", reflect.Int).Patch(rushia.H{
	"Age":      0,
	"Username": "",
	"Password": "123456",
})
// Equals: UPDATE Users SET Age = ?, Password = ?, Username = ? WHERE Username = ?
```

### Delete

Deletes everything! Remember to add a condition to prevent it really deletes everything.

```go
rushia.NewQuery("Users").Where("ID = ", 1).Delete()
// Equals: DELETE FROM Users WHERE ID = ?
```

### Select

Use `Select` to get the data.

```go
rushia.NewQuery("Users").Select()
// Equals: SELECT * FROM Users
```

#### Specify columns

Specify the columns to select in the `Select` arguments, It colud also be a expression.

```go
rushia.NewQuery("Users").Select("Username", "Nickname")
// Equals: SELECT Username, Nickname FROM Users

rushia.NewQuery("Users").Select(rushia.NewExpr("COUNT(*) AS Count"))
// Equals: SELECT COUNT(*) AS Count FROM Users
```

#### Select One

To get a single row data, use `SelectOne`. It's a shorthand for `.Limit(1).Select(...)`.

```go
rushia.NewQuery("Users").SelectOne("Username")
// Equals: SELECT Username FROM Users LIMIT 1
```

#### Distinct

Specifing `Distinct` to eliminate the duplicate rows while fetching the data.

```go
rushia.NewQuery("Products").Distinct().Select()
// Equals: SELECT DISTINCT * FROM Products
```

#### Union

`Union` or `UnionAll` allows you to merge the data between different table selections.

```go
locationQuery := rushia.NewQuery("Locations").Select()

rushia.NewQuery("Users").Union(locationQuery).Select()
// Equals: SELECT * FROM Users UNION SELECT * FROM Locations

rushia.NewQuery("Users").UnionAll(locationQuery).Select()
// Equals: SELECT * FROM Users UNION ALL SELECT * FROM Locations
```

### Select exists

To execute `SELECT EXISTS` by calling `Exists`.

```go
rushia.NewQuery("Users").Where("Username = ?", "YamiOdymel").Exists()
// Equals: SELECT EXISTS(SELECT * FROM Users WHERE Username = ?)
```

### Table alias

`As` assign an alias to the query, it's useful if you are creating a sub query. In a joining or common scenario, use `NewAlias` instead.

```go
rushia.NewQuery(NewQuery("Users").Select()).As("Result").Where("Username = ?", "YamiOdymel").Select())
// Equals: SELECT * FROM (SELECT * FROM Users) AS Result WHERE Username = ?

rushia.NewQuery(rushia.NewAlias("UserFriendRelationships", "relations")).Where("relations.ID = ?", 5).Select()
// Equals: SELECT * FROM UserFriendRelationships AS relations WHERE relations.ID = ?
```

### Raw Query

Rushia provides you the most 80% things you will use, but if you are in the bad luck to request for the rest 20%, the only hope is to use Raw Query.

A raw query does also support the prepared statement, to replace the value as `?` to prevent the SQL injection.

`NewRawQuery` is the same as `NewQuery` that required to be `Build`, and the helper functions such as: `Limit`, `OrderBy`...etc, are not able to be used.

```go
q := rushia.NewRawQuery("SELECT * FROM Users WHERE ID >= ?", 10)
```

### Conditions

To define a `WHERE` or `HAVING` condition in Rushia is a piece of cake!

| SQL Query                                                | Usage                                                                                      |
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

The condition functions has it's own transform for `Where`, `OrWhere`, `Having`, `OrHaving`, `JoinWhere`, `OrJoinWhere`.

```go
rushia.NewQuery("Users").Where("ID = ?", 1).Where("Username = ?", "admin").Select()
// Equals: SELECT * FROM Users WHERE ID = ? AND Username = ?

rushia.NewQuery("Users").Having("ID = ?", 1).Having("Username = ?", "admin").Select()
// Equals: SELECT * FROM Users HAVING ID = ? AND Username = ?

rushia.NewQuery("Users").WhereC("ID != CompanyID").Where("DATE(CreatedAt) = DATE(LastLogin)").Select()
// Equals: SELECT * FROM Users WHERE ID != CompanyID AND DATE(CreatedAt) = DATE(LastLogin)
```

### Escaped Values

The same usage as `??` double question marks in [mysqljs/mysql](https://github.com/mysqljs/mysql) package, it's possible to escape the values with backticks (\`) by using `??`. It's useful for column names.

```go
var ColumnUserID = "ID"
rushia.NewQuery("Users").Where("?? = ?", ColumnUserID, 3).Select()
// Equals: SELECT * FROM Users WHERE `ID` = ?
```

### Order

Ordering is also supported in Rushia and can be used with functions.

```go
rushia.NewQuery("Users").OrderBy("ID ASC").OrderBy("Login DESC").OrderBy("RAND()").Select()
// Equals: SELECT * FROM Users ORDER BY ID ASC, Login DESC, RAND()
```

#### Order by field

Or ordering by custom field values:

```go
rushia.NewQuery("Users").OrderByField("UserGroup", "SuperUser", "Admin", "Users").Select()
// Equals: SELECT * FROM Users ORDER BY FIELD (UserGroup, ?, ?, ?)
```

### Group by

The result can also be grouped with `GroupBy`.

```go
rushia.NewQuery("Users").GroupBy("Name").Select()
// Equals: SELECT * FROM Users GROUP BY Name
```

### Table joins

Rushia supports multiple ways to join the tables, such as: `InerrJoin`, `LeftJoin`, `RightJoin`, `NaturalJoin`, `CrossJoin`. While joining, the last argument is always a raw condition and colud be useful.

```go
rushia.
	NewQuery("Products").
	LeftJoin("Users", "Products.TenantID = Users.TenantID").
	Select("Users.Name", "Products.ProductName")
// Equals: SELECT Users.Name, Products.ProductName FROM Products AS Products LEFT JOIN Users AS Users ON (Products.TenantID = Users.TenantID)

rushia.
	NewQuery("Products").
	LeftJoin("Users", "Products.TenantID = ?", 3).
	Select("Users.Name", "Products.ProductName")
// Equals: SELECT Users.Name, Products.ProductName FROM Products AS Products LEFT JOIN Users AS Users ON (Products.TenantID = ?)
```

Or just omit the condition and define it later in the function chaining.

```go
rushia.
	NewQuery("Products").
	LeftJoin("Users").
	JoinWhere("Products.TenantID = Users.TenantID")
	Select("Users.Name", "Products.ProductName")
// Equals: SELECT Users.Name, Products.ProductName FROM Products AS Products LEFT JOIN Users AS Users ON (Products.TenantID = Users.TenantID)
```

#### Join condition

With `JoinWhere` or `OrJoinWhere` to expand the conditions for the table joins. The condition will always to be added into the latest joined table.

```go
rushia.
	NewQuery("Products").
	LeftJoin("Users", "Products.TenantID = Users.TenantID").
	OrJoinWhere("Users.TenantID = ?", 5).
	Select("Users.Name", "Products.ProductName")
// Equals: SELECT Users.Name, Products.ProductName FROM Products AS Products LEFT JOIN Users AS Users ON (Products.TenantID = Users.TenantID OR Users.TenantID = ?)
```

### Sub query

Rushia supports nested query which is called Sub Query. Use a query as a value to make it sub query.

```go
subQuery := rushia.NewQuery("VIPUsers").Select("UserID")

rushia.NewQuery("Users").WhereIn("ID", subQuery).Select()
// Equals: SELECT * FROM Users WHERE ID IN (SELECT UserID FROM VIPUsers)
```

#### Sub query insertion

To insert a value from a sub query, simply use the query as a value and make sure the sub query only returns one column and one row as result.

```go
subQuery := rushia.NewQuery("Users").Where("ID = ?", 6).SelectOne("Name")

rushia.NewQuery("Products").Insert(rushia.H{
	"ProductName": "測試商品",
	"UserID":      subQuery,
	"LastUpdated": rushia.NewExpr("NOW()")
})
// Equals: INSERT INTO Products (ProductName, UserID, LastUpdated) VALUES (?, (SELECT Name FROM Users WHERE ID = 6 LIMIT 1), NOW())
```

#### Sub query joining

Join a table from a sub query is possible, but requires to assign an alias to the sub query by using `As`.

```go
subQuery := rushia.NewQuery("Users").As("Users").Where("Active = ?", 1).Select()

rushia.
	NewQuery("Products").
	LeftJoin(subQuery, "Products.UserID = Users.ID").
	Select("Users.Username", "Products.ProductName")
// Equals: SELECT Users.Username, Products.ProductName FROM Products AS Products LEFT JOIN (SELECT * FROM Users WHERE Active = ?) AS Users ON Products.UserID = Users.ID
```

#### Sub query swapping

Passing a sub query to a raw query or an expression will automatically looking for the prepared statement `?` to replace as a built sub query.

```go
subQuery := rushia.NewQuery("Locations").Select()
rawQuery := rushia.NewRawQuery("SELECT UserID FROM Users WHERE EXISTS (?)", subQuery)

NewQuery("Products").Where("EXISTS ?", rawQuery).Select()
// Equals: SELECT * FROM Products WHERE EXISTS (SELECT UserID FROM Users WHERE EXISTS (SELECT * FROM Locations))
```

### Set query options

You can set the query options with Rushia.

```go
rushia.NewQuery("Users").SetQueryOption("FOR UPDATE").Select()
// Equals: SELECT * FROM Users FOR UPDATE

rushia.NewQuery("Users").SetQueryOption("SQL_NO_CACHE").Select()
// Equals: SELECT SQL_NO_CACHE * FROM Users

rushia.NewQuery("Users").SetQueryOption("LOW_PRIORITY", "IGNORE").Insert(data)
// Gives: INSERT LOW_PRIORITY IGNORE INTO Users ...
```

## Complex query example

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

// Equals:
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
	Where("AdvanceAmount < ANY (?)", customers).
	Select("OrderNum", "OrderDate", "OrderAmount", "AdvanceAmount")

// Equals:
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

## References

Let's see what inspired Rushia.

-   [kisielk/sqlstruct](http://godoc.org/github.com/kisielk/sqlstruct)
-   [jmoiron/sqlx](https://github.com/jmoiron/sqlx)
-   [russross/meddler](https://github.com/russross/meddler)
-   [jinzhu/gorm](https://github.com/jinzhu/gorm)
-   [doug-martin/goqu](https://github.com/doug-martin/goqu)
-   [gocraft/dbr](https://github.com/gocraft/dbr)
-   [go-ozzo/ozzo-dbx](https://github.com/go-ozzo/ozzo-dbx)
-   [kyleconroy/sqlc](https://github.com/kyleconroy/sqlc)
