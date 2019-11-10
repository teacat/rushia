package rushia

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var builder Query

// assertEqual 會將期望的 SQL 指令與實際的 SQL 指令拆分，因為 Reiner 裡有 Map 會導致產生出來的結果每次都不如預期地按照順序排。
// 拆分後便會比對是否有相同的「字詞」，若短缺則是執行結果不符合預期即報錯。
func assertEqual(a *assert.Assertions, expected string, actual string) {
	originalExpected := expected
	originalActual := actual
	expected = strings.Replace(expected, "(", "", -1)
	expected = strings.Replace(expected, ")", "", -1)
	expected = strings.Replace(expected, ",", "", -1)
	expectedParts := strings.Split(expected, " ")
	actual = strings.Replace(actual, "(", "", -1)
	actual = strings.Replace(actual, ")", "", -1)
	actual = strings.Replace(actual, ",", "", -1)
	actualParts := strings.Split(actual, " ")
	passed := []bool{}
	for _, v := range expectedParts {
		for _, vv := range actualParts {
			if v == vv {
				passed = append(passed, true)
				break
			}

		}
	}
	if len(passed) != len(actualParts) {
		a.Fail(`Not equal:`, "expected: \"%s\"\nreceived: \"%s\"", originalExpected, originalActual)
	}
	return
}

func TestMain(t *testing.T) {
	builder = NewQuery()
}

func TestInsert(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Insert(map[string]interface{}{
		"Username": "YamiOdymel",
		"Password": "test",
	})
	assertEqual(assert, "INSERT INTO Users (Username, Password) VALUES (?, ?)", query)
}

func TestInsertMulti(t *testing.T) {
	assert := assert.New(t)
	data := []map[string]interface{}{
		{
			"Username": "YamiOdymel",
			"Password": "test",
		}, {
			"Username": "Karisu",
			"Password": "12345",
		},
	}
	query, _ := builder.Table("Users").InsertMulti(data)
	assertEqual(assert, "INSERT INTO Users (Password, Username) VALUES (?, ?), (?, ?)", query)
}

func TestReplace(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Replace(map[string]interface{}{
		"Username": "YamiOdymel",
		"Password": "test",
	})
	assertEqual(assert, "REPLACE INTO Users (Password, Username) VALUES (?, ?)", query)
}

func TestInsertFunc(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Insert(map[string]interface{}{
		"Username":  "YamiOdymel",
		"Password":  builder.Func("SHA1(?)", "secretpassword+salt"),
		"Expires":   builder.Now("+1Y"),
		"CreatedAt": builder.Now(),
	})
	assertEqual(assert, "INSERT INTO Users (CreatedAt, Expires, Password, Username) VALUES (NOW(), NOW() + INTERVAL 1 YEAR, SHA1(?), ?)", query)
}

func TestOnDuplicateInsert(t *testing.T) {
	assert := assert.New(t)
	lastInsertID := "ID"
	query, _ := builder.Table("Users").OnDuplicate([]string{"UpdatedAt"}, lastInsertID).Insert(map[string]interface{}{
		"Username":  "YamiOdymel",
		"Password":  "test",
		"UpdatedAt": builder.Now(),
	})
	assertEqual(assert, "INSERT INTO Users (Password, UpdatedAt, Username) VALUES (?, NOW(), ?) ON DUPLICATE KEY UPDATE ID=LAST_INSERT_ID(ID), UpdatedAt = VALUES(UpdatedAt)", query)
}

func TestUpdate(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Where("Username", "YamiOdymel").Update(map[string]interface{}{
		"Username": "Karisu",
		"Password": "123456",
	})
	assertEqual(assert, "UPDATE Users SET Password = ?, Username = ? WHERE Username = ?", query)
}

func TestLimitUpdate(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Limit(10).Update(map[string]interface{}{
		"Username": "Karisu",
		"Password": "123456",
	})
	assertEqual(assert, "UPDATE Users SET Password = ?, Username = ? LIMIT 10", query)
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Get()
	assertEqual(assert, "SELECT * FROM Users", query)
}

func TestLimitGet(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Limit(10).Get()
	assertEqual(assert, "SELECT * FROM Users LIMIT 10", query)
}

func TestLimitMultiGet(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Limit(10, 20).Get()
	assertEqual(assert, "SELECT * FROM Users LIMIT 10, 20", query)
}

func TestGetColumns(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Get("Username", "Nickname")
	assertEqual(assert, "SELECT Username, Nickname FROM Users", query)

	query, _ = builder.Table("Users").Get("COUNT(*) AS Count")
	assertEqual(assert, "SELECT COUNT(*) AS Count FROM Users", query)
}

func TestGetOne(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Where("ID", 1).Get()
	assertEqual(assert, "SELECT * FROM Users WHERE ID = ?", query)

	query, _ = builder.Table("Users").GetOne()
	assertEqual(assert, "SELECT * FROM Users LIMIT 1", query)

	query, _ = builder.Table("Users").Get("SUM(ID)", "COUNT(*) AS Count")
	assertEqual(assert, "SELECT SUM(ID), COUNT(*) AS Count FROM Users", query)
}

func TestRawQuery(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.RawQuery("SELECT * FROM Users WHERE ID >= ?", 10)
	assertEqual(assert, "SELECT * FROM Users WHERE ID >= ?", query)
}

func TestWhere(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Where("ID", 1).Where("Username", "admin").Get()
	assertEqual(assert, "SELECT * FROM Users WHERE ID = ? AND Username = ?", query)
}

func TestWhereQuery(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Where("(ID = ?)", 1).Where("Username", "admin").Get()
	assertEqual(assert, "SELECT * FROM Users WHERE (ID = ?) AND Username = ?", query)
	query, _ = builder.Table("Users").Where("(ID = ? OR Password = SHA(?))", 1, "password").Where("Username", "admin").Get()
	assertEqual(assert, "SELECT * FROM Users WHERE (ID = ? OR Password = SHA(?)) AND Username = ?", query)
}

func TestWhereHaving(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Where("ID", 1).Having("Username", "admin").Get()
	assertEqual(assert, "SELECT * FROM Users WHERE ID = ? HAVING Username = ?", query)
	query, _ = builder.Table("Users").Where("ID", 1).Having("Username", "admin").OrHaving("Password", "test").Get()
	assertEqual(assert, "SELECT * FROM Users WHERE ID = ? HAVING Username = ? OR Password = ?", query)
}

func TestWhereColumns(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Where("LastLogin = CreatedAt").Get()
	assertEqual(assert, "SELECT * FROM Users WHERE LastLogin = CreatedAt", query)
}

func TestWhereOperator(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Where("ID", ">=", 50).Get()
	assertEqual(assert, "SELECT * FROM Users WHERE ID >= ?", query)
}

func TestWhereBetween(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Where("ID", "BETWEEN", 0, 20).Get()
	assertEqual(assert, "SELECT * FROM Users WHERE ID BETWEEN ? AND ?", query)

	query, _ = builder.Table("Users").Where("ID", "NOT BETWEEN", 0, 20).Get()
	assertEqual(assert, "SELECT * FROM Users WHERE ID NOT BETWEEN ? AND ?", query)
}

func TestWhereIn(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Where("ID", "IN", 1, 5, 27, -1, "d").Get()
	assertEqual(assert, "SELECT * FROM Users WHERE ID IN (?, ?, ?, ?, ?)", query)

	query, _ = builder.Table("Users").Where("ID", "NOT IN", 1, 5, 27, -1, "d").Get()
	assertEqual(assert, "SELECT * FROM Users WHERE ID NOT IN (?, ?, ?, ?, ?)", query)
}

func TestOrWhere(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Where("FirstName", "John").OrWhere("FirstName", "Peter").Get()
	assertEqual(assert, "SELECT * FROM Users WHERE FirstName = ? OR FirstName = ?", query)

	query, _ = builder.Table("Users").Where("A = B").OrWhere("(A = C OR A = D)").Get()
	assertEqual(assert, "SELECT * FROM Users WHERE A = B OR (A = C OR A = D)", query)
}

func TestWhereNull(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Where("LastName", "IS", nil).Get()
	assertEqual(assert, "SELECT * FROM Users WHERE LastName IS NULL", query)

	query, _ = builder.Table("Users").Where("LastName", "IS NOT", nil).Get()
	assertEqual(assert, "SELECT * FROM Users WHERE LastName IS NOT NULL", query)
}

func TestTimestampDate(t *testing.T) {
	assert := assert.New(t)
	ts := NewTimestamp()
	query, _ := builder.Table("Users").Where("CreatedAt", ts.IsDate("2017-07-13")).Get()
	assertEqual(assert, "SELECT * FROM Users WHERE DATE(FROM_UNIXTIME(CreatedAt)) = ?", query)

	query, _ = builder.Table("Users").Where("CreatedAt", ts.IsYear(2017)).Get()
	assertEqual(assert, "SELECT * FROM Users WHERE YEAR(FROM_UNIXTIME(CreatedAt)) = ?", query)

	query, _ = builder.Table("Users").Where("CreatedAt", ts.IsMonth(1)).Get()
	assertEqual(assert, "SELECT * FROM Users WHERE MONTH(FROM_UNIXTIME(CreatedAt)) = ?", query)
	query, _ = builder.Table("Users").Where("CreatedAt", ts.IsMonth("January")).Get()
	assertEqual(assert, "SELECT * FROM Users WHERE MONTH(FROM_UNIXTIME(CreatedAt)) = ?", query)

	query, _ = builder.Table("Users").Where("CreatedAt", ts.IsDay(16)).Get()
	assertEqual(assert, "SELECT * FROM Users WHERE DAY(FROM_UNIXTIME(CreatedAt)) = ?", query)

	query, _ = builder.Table("Users").Where("CreatedAt", ts.IsWeekday(5)).Get()
	assertEqual(assert, "SELECT * FROM Users WHERE WEEKDAY(FROM_UNIXTIME(CreatedAt)) = ?", query)
	query, _ = builder.Table("Users").Where("CreatedAt", ts.IsWeekday("Friday")).Get()
	assertEqual(assert, "SELECT * FROM Users WHERE WEEKDAY(FROM_UNIXTIME(CreatedAt)) = ?", query)
}

func TestTimestampTime(t *testing.T) {
	assert := assert.New(t)
	ts := NewTimestamp()
	query, _ := builder.Table("Users").Where("CreatedAt", ts.IsHour(18)).Get()
	assertEqual(assert, "SELECT * FROM Users WHERE HOUR(FROM_UNIXTIME(CreatedAt)) = ?", query)

	query, _ = builder.Table("Users").Where("CreatedAt", ts.IsMinute(25)).Get()
	assertEqual(assert, "SELECT * FROM Users WHERE MINUTE(FROM_UNIXTIME(CreatedAt)) = ?", query)

	query, _ = builder.Table("Users").Where("CreatedAt", ts.IsSecond(16)).Get()
	assertEqual(assert, "SELECT * FROM Users WHERE SECOND(FROM_UNIXTIME(CreatedAt)) = ?", query)

	query, _ = builder.Table("Users").Where("CreatedAt", ts.IsWeekday(5)).Get()
	assertEqual(assert, "SELECT * FROM Users WHERE WEEKDAY(FROM_UNIXTIME(CreatedAt)) = ?", query)
}

func TestRawWhere(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Where("ID != CompanyID").Where("DATE(CreatedAt) = DATE(LastLogin)").Get()
	assertEqual(assert, "SELECT * FROM Users WHERE ID != CompanyID AND DATE(CreatedAt) = DATE(LastLogin)", query)
}

func TestRawWhereParams(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Where("(ID = ? OR ID = ?)", 6, 2).Where("Login", "Mike").Get()
	assertEqual(assert, "SELECT * FROM Users WHERE (ID = ? OR ID = ?) AND Login = ?", query)
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").Where("ID", 1).Delete()
	assertEqual(assert, "DELETE FROM Users WHERE ID = ?", query)
}

func TestOrderBy(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").OrderBy("ID", "ASC").OrderBy("Login", "DESC").OrderBy("RAND()").Get()
	assertEqual(assert, "SELECT * FROM Users ORDER BY ID ASC, Login DESC, RAND()", query)
}

func TestOrderByField(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").OrderBy("UserGroup", "ASC", "SuperUser", "Admin", "Users").Get()
	assertEqual(assert, "SELECT * FROM Users ORDER BY FIELD (UserGroup, ?, ?, ?) ASC", query)
}

func TestGroupBy(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").GroupBy("Name").Get()
	assertEqual(assert, "SELECT * FROM Users GROUP BY Name", query)
	query, _ = builder.Table("Users").GroupBy("Name", "ID").Get()
	assertEqual(assert, "SELECT * FROM Users GROUP BY Name, ID", query)
}

func TestJoin(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.
		Table("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		Where("Users.ID", 6).
		Get("Users.Name", "Products.ProductName")
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products LEFT JOIN Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?", query)

	query, _ = builder.
		Table("Products").
		RightJoin("Users", "Products.TenantID = Users.TenantID").
		Where("Users.ID", 6).
		Get("Users.Name", "Products.ProductName")
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products RIGHT JOIN Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?", query)

	query, _ = builder.
		Table("Products").
		InnerJoin("Users", "Products.TenantID = Users.TenantID").
		Where("Users.ID", 6).
		Get("Users.Name", "Products.ProductName")
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products INNER JOIN Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?", query)

	query, _ = builder.
		Table("Products").
		NaturalJoin("Users", "Products.TenantID = Users.TenantID").
		Where("Users.ID", 6).
		Get("Users.Name", "Products.ProductName")
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products NATURAL JOIN Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?", query)

	query, _ = builder.
		Table("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		RightJoin("Posts", "Products.TenantID = Posts.TenantID").
		Where("Users.ID", 6).
		Get("Users.Name", "Products.ProductName")
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products LEFT JOIN Users ON (Products.TenantID = Users.TenantID) RIGHT JOIN Posts ON (Products.TenantID = Posts.TenantID) WHERE Users.ID = ?", query)
}

func TestJoinWhere(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.
		Table("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		JoinOrWhere("Users", "Users.TenantID", 5).
		Get("Users.Name", "Products.ProductName")
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products LEFT JOIN Users ON (Products.TenantID = Users.TenantID OR Users.TenantID = ?)", query)
	query, _ = builder.
		Table("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		JoinWhere("Users", "Users.Username", "Wow").
		Get("Users.Name", "Products.ProductName")
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products LEFT JOIN Users ON (Products.TenantID = Users.TenantID AND Users.Username = ?)", query)
	query, _ = builder.
		Table("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		RightJoin("Posts", "Products.TenantID = Posts.TenantID").
		JoinWhere("Posts", "Posts.Username", "Wow").
		JoinWhere("Users", "Users.Username", "Wow").
		Get("Users.Name", "Products.ProductName")
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products LEFT JOIN Users ON (Products.TenantID = Users.TenantID AND Users.Username = ?) RIGHT JOIN Posts ON (Products.TenantID = Posts.TenantID AND Posts.Username = ?)", query)
}

func TestWithTotalCount(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").WithTotalCount().Get("Username")
	assertEqual(assert, "SELECT SQL_CALC_FOUND_ROWS Username FROM Users", query)
}

func TestSubQueryGet(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewSubQuery().Table("Products").Where("Quantity", ">", 2).Get("UserID")
	query, _ := builder.Table("Users").Where("ID", "IN", subQuery).Get()
	assertEqual(assert, "SELECT * FROM Users WHERE ID IN (SELECT UserID FROM Products WHERE Quantity > ?)", query)
}

func TestSubQueryInsert(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewSubQuery().Table("Users").Where("ID", 6).Get("Name")
	query, _ := builder.Table("Products").Insert(map[string]interface{}{
		"ProductName": "測試商品",
		"UserID":      subQuery,
		"LastUpdated": builder.Now(),
	})
	assertEqual(assert, "INSERT INTO Products (LastUpdated, ProductName, UserID) VALUES (NOW(), ?, (SELECT Name FROM Users WHERE ID = ?))", query)
}

func TestSubQueryJoin(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewSubQuery("Users").Table("Users").Where("Active", 1).Get()
	query, _ := builder.
		Table("Products").
		LeftJoin(subQuery, "Products.UserID = Users.ID").
		Get("Users.Username", "Products.ProductName")
	assertEqual(assert, "SELECT Users.Username, Products.ProductName FROM Products LEFT JOIN (SELECT * FROM Users WHERE Active = ?) AS Users ON (Products.UserID = Users.ID)", query)
}

func TestSubQueryExist(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewSubQuery("Users").Table("Users").Where("Company", "測試公司").Get("UserID")
	query, _ := builder.Table("Products").Where(subQuery, "EXISTS").Get()
	assertEqual(assert, "SELECT * FROM Products WHERE EXISTS (SELECT UserID FROM Users WHERE Company = ?)", query)
}

func TestSubQueryRawQuery(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewSubQuery("Users").RawQuery("SELECT UserID FROM Users WHERE Company = ?", "測試公司")
	query, _ := builder.Table("Products").Where(subQuery, "EXISTS").Get()
	assertEqual(assert, "SELECT * FROM Products WHERE EXISTS (SELECT UserID FROM Users WHERE Company = ?)", query)
}

func TestLock(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Lock("Users")
	assertEqual(assert, "LOCK TABLES Users", query)
}

func TestSetQueryOption(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Table("Users").SetQueryOption("FOR UPDATE").Get("Username")
	assertEqual(assert, "SELECT Username FROM Users FOR UPDATE", query)
}

func TestSetLockMethod(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.SetLockMethod("WRITE").Lock("Users")
	assertEqual(assert, "LOCK TABLES Users WRITE", query)
}

func TestUnlock(t *testing.T) {
	assert := assert.New(t)
	query, _ := builder.Unlock("Users")
	assertEqual(assert, "UNLOCK TABLES Users", query)
}
