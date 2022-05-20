package rushia

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

	expectedLen := len(strings.Split(originalExpected, " "))
	actualLen := len(strings.Split(originalActual, " "))

	if len(passed) != len(actualParts) {
		a.Fail(`Not equal:`, "expected: \"%s\"\nreceived: \"%s\"", originalExpected, originalActual)
	}
	if expectedLen != actualLen {
		a.Fail(`Not same length:`, "expected: \"%s\"\nreceived: \"%s\"", originalExpected, originalActual)
	}
}

func assertParams(a *assert.Assertions, expected []interface{}, actual []interface{}) {
	if len(expected) != len(actual) {
		a.Fail("Not same params length", "expected: \"%d\" %+v\nreceived: \"%d\" %+v", len(expected), expected, len(actual), actual)
	}
	for _, v := range expected {
		var yes bool
		for _, j := range actual {
			if v == j {
				yes = true
			}
		}
		if !yes {
			a.Fail(`Not in params`)
		}
	}
}

func assertParamOrders(a *assert.Assertions, expected []interface{}, actual []interface{}) {
	if len(expected) != len(actual) {
		a.Fail("Not same params length", "expected: \"%d\" %+v\nreceived: \"%d\" %+v", len(expected), expected, len(actual), actual)
	}
	for k, v := range expected {
		if actual[k] != v {
			a.Fail(`Not same params order`)
		}
	}
}

func BenchmarkInsertStruct(b *testing.B) {
	u := struct {
		Username string
		Password string
	}{
		Username: "YamiOdymel",
		Password: "test",
	}
	for i := 0; i < b.N; i++ {
		NewQuery("Users").Insert(u)
	}
}

func BenchmarkInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewQuery("Users").Insert(H{
			"Username": "YamiOdymel",
			"Password": "test",
		})
	}
}

//=======================================================
// Insert
//=======================================================

func TestInsertStruct(t *testing.T) {
	u := struct {
		Username string
		Password string
	}{
		Username: "YamiOdymel",
		Password: "test",
	}
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Insert(u))
	assertEqual(assert, "INSERT INTO Users (Username, Password) VALUES (?, ?)", query)
	assertParams(assert, []interface{}{"YamiOdymel", "test"}, params)
}

func TestInsertStructPointer(t *testing.T) {
	type user struct {
		Username string
		Password string
	}
	u := &user{
		Username: "YamiOdymel",
		Password: "test",
	}
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Insert(u))
	assertEqual(assert, "INSERT INTO Users (Username, Password) VALUES (?, ?)", query)
	assertParams(assert, []interface{}{"YamiOdymel", "test"}, params)
}

func TestInsert(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Insert(H{
		"Username": "YamiOdymel",
		"Password": "test",
	}))
	assertEqual(assert, "INSERT INTO Users (Username, Password) VALUES (?, ?)", query)
	assertParams(assert, []interface{}{"YamiOdymel", "test"}, params)
}

func TestInsertMap(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Insert(map[string]interface{}{
		"Username": "YamiOdymel",
		"Password": "test",
	}))
	assertEqual(assert, "INSERT INTO Users (Username, Password) VALUES (?, ?)", query)
	assertParams(assert, []interface{}{"YamiOdymel", "test"}, params)
}

func TestInsertMulti(t *testing.T) {
	assert := assert.New(t)
	data := []H{
		{
			"Username": "YamiOdymel",
			"Password": "test",
		}, {
			"Username": "Karisu",
			"Password": "12345",
		},
	}
	query, params := Build(NewQuery("Users").Insert(data))
	assertEqual(assert, "INSERT INTO Users (Password, Username) VALUES (?, ?), (?, ?)", query)
	assertParams(assert, []interface{}{"YamiOdymel", "test", "Karisu", "12345"}, params)
}

func TestInsertMultiMaps(t *testing.T) {
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
	query, params := Build(NewQuery("Users").Insert(data))
	assertEqual(assert, "INSERT INTO Users (Password, Username) VALUES (?, ?), (?, ?)", query)
	assertParams(assert, []interface{}{"YamiOdymel", "test", "Karisu", "12345"}, params)
}

func TestInsertStructOmit(t *testing.T) {
	u := struct {
		Username string
		Password string
	}{
		Username: "YamiOdymel",
		Password: "test",
	}
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Omit("Username").Insert(u))
	assertEqual(assert, "INSERT INTO Users (Password) VALUES (?)", query)
	assertParams(assert, []interface{}{"test"}, params)
}

func TestInsertStructOmitByTag(t *testing.T) {
	u := struct {
		Username string `rushia:"-"`
		Password string
	}{
		Username: "YamiOdymel",
		Password: "test",
	}
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Insert(u))
	assertEqual(assert, "INSERT INTO Users (Password) VALUES (?)", query)
	assertParams(assert, []interface{}{"test"}, params)
}

func TestInsertStructTagOmit(t *testing.T) {
	assert := assert.New(t)
	u := struct {
		Username string `rushia:"user_name"`
		Password string
	}{
		Username: "YamiOdymel",
		Password: "test",
	}
	query, params := Build(NewQuery("Users").Omit("user_name").Insert(u))
	assertEqual(assert, "INSERT INTO Users (Password) VALUES (?)", query)
	assertParams(assert, []interface{}{"test"}, params)
}

func TestInsertStructTagRename(t *testing.T) {
	u := struct {
		Username string `rushia:"user_name"`
		Password string
	}{
		Username: "YamiOdymel",
		Password: "test",
	}
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Insert(u))
	assertEqual(assert, "INSERT INTO Users (user_name, Password) VALUES (?, ?)", query)
	assertParams(assert, []interface{}{"YamiOdymel", "test"}, params)
}

func TestInsertOmit(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Omit("Username").Insert(H{
		"Username": "YamiOdymel",
		"Password": "test",
	}))
	assertEqual(assert, "INSERT INTO Users (Password) VALUES (?)", query)
	assertParams(assert, []interface{}{"test"}, params)
}

func TestInsertMultiOmit(t *testing.T) {
	assert := assert.New(t)
	data := []H{
		{
			"Username": "YamiOdymel",
			"Password": "test",
		}, {
			"Username": "Karisu",
			"Password": "12345",
		},
	}
	query, params := Build(NewQuery("Users").Omit("Username").Insert(data))
	assertEqual(assert, "INSERT INTO Users (Password) VALUES (?), (?)", query)
	assertParams(assert, []interface{}{"test", "12345"}, params)
}

func TestInsertExpr(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Insert(H{
		"Username":  "YamiOdymel",
		"Password":  NewExpr("SHA1(?)", "secretpassword+salt"),
		"Expires":   NewExpr("NOW() + INTERVAL 1 YEAR"),
		"CreatedAt": NewExpr("NOW()"),
	}))
	assertEqual(assert, "INSERT INTO Users (CreatedAt, Expires, Password, Username) VALUES (NOW(), NOW() + INTERVAL 1 YEAR, SHA1(?), ?)", query)
	assertParams(assert, []interface{}{"YamiOdymel", "secretpassword+salt"}, params)
}

func TestInsertSubQueryExpr(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewQuery("Salaries").Where("Username = ?", "YamiOdymel").Select("Salary")
	query, params := Build(NewQuery("Users").Insert(H{
		"Username":  "YamiOdymel",
		"AvgSalary": NewExpr("SUM((?))", subQuery),
	}))
	assertEqual(assert, "INSERT INTO Users (Username, AvgSalary) VALUES (?, SUM((SELECT Salary FROM Salaries WHERE Username = ?)))", query)
	assertParams(assert, []interface{}{"YamiOdymel", "YamiOdymel"}, params)

	subQuery = NewQuery("Salaries").Where("Username = ?", "YamiOdymel").Select("Salary")
	query, params = Build(NewQuery("Users").Insert(H{
		"Username": "YamiOdymel",
		"Salary":   subQuery,
	}))
	assertEqual(assert, "INSERT INTO Users (Username, Salary) VALUES (?, (SELECT Salary FROM Salaries WHERE Username = ?))", query)
	assertParams(assert, []interface{}{"YamiOdymel", "YamiOdymel"}, params)
}

func TestInsertSelect(t *testing.T) {
	assert := assert.New(t)
	from := NewQuery("AdditionalUsers").Where("Name LIKE ?", "ABC%").Select("ID", "Username", "Nickname")
	query, params := Build(NewQuery("Users").InsertSelect(from, "ID", "Username", "Nickname"))
	assertEqual(assert, "INSERT INTO Users (ID, Username, Nickname) SELECT ID, Username, Nickname FROM AdditionalUsers WHERE Username LIKE ?", query)
	assertParams(assert, []interface{}{"ABC%"}, params)
}

func TestOnDuplicateInsert(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").OnDuplicate(H{
		"UpdatedAt": NewExpr("VALUES(UpdatedAt)"), // Deprecated in MySQL 8.0.20
		"ID":        NewExpr("LAST_INSERT_ID(ID)"),
	}).Insert(H{
		"Username":  "YamiOdymel",
		"Password":  "test",
		"UpdatedAt": NewExpr("NOW()"),
	}))
	assertEqual(assert, "INSERT INTO Users (Password, UpdatedAt, Username) VALUES (?, NOW(), ?) ON DUPLICATE KEY UPDATE ID = LAST_INSERT_ID(ID), UpdatedAt = VALUES(UpdatedAt)", query)
	assertParams(assert, []interface{}{"YamiOdymel", "test"}, params)

	query, params = Build(NewQuery("Users").As("New").OnDuplicate(H{
		"UpdatedAt": NewExpr("New.UpdatedAt"),
		"ID":        NewExpr("LAST_INSERT_ID(ID)"),
	}).Insert(H{
		"Username":  "YamiOdymel",
		"Password":  "test",
		"UpdatedAt": NewExpr("NOW()"),
	}))
	assertEqual(assert, "INSERT INTO Users (Password, UpdatedAt, Username) VALUES (?, NOW(), ?) AS New ON DUPLICATE KEY UPDATE ID = LAST_INSERT_ID(ID), UpdatedAt = New.UpdatedAt", query)
	assertParams(assert, []interface{}{"YamiOdymel", "test"}, params)
}

//=======================================================
// Build
//=======================================================

func TestBuildNoType(t *testing.T) {
	assert := assert.New(t)
	assert.Panics(func() {
		Build(NewQuery("Users"))
	})
}

//=======================================================
// Replace
//=======================================================

func TestReplace(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Replace(H{
		"Username": "YamiOdymel",
		"Password": "test",
	}))
	assertEqual(assert, "REPLACE INTO Users (Password, Username) VALUES (?, ?)", query)
	assertParams(assert, []interface{}{"YamiOdymel", "test"}, params)
}

//=======================================================
// Update
//=======================================================

func TestUpdateOmit(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("Username = ?", "YamiOdymel").Omit("Username").Update(H{
		"Username": "Karisu",
		"Password": "123456",
	}))
	assertEqual(assert, "UPDATE Users SET Password = ? WHERE Username = ?", query)
	assertParams(assert, []interface{}{"YamiOdymel", "123456"}, params)
}

func TestUpdateOmitStruct(t *testing.T) {
	u := struct {
		Username string
		Password string
	}{
		Username: "YamiOdymel",
		Password: "test",
	}
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("Username = ?", "YamiOdymel").Omit("Username").Update(u))
	assertEqual(assert, "UPDATE Users SET Password = ? WHERE Username = ?", query)
	assertParams(assert, []interface{}{"test", "YamiOdymel"}, params)
}

func TestUpdate(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("Username = ?", "YamiOdymel").Update(H{
		"Username": "",
		"Password": "123456",
	}))
	assertEqual(assert, "UPDATE Users SET Password = ?, Username = ? WHERE Username = ?", query)
	assertParams(assert, []interface{}{"YamiOdymel", "", "123456"}, params)
}

func TestUpdateCase(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("Username = ?", "YamiOdymel").Update(H{
		"Username": "",
		"Password": NewExpr(`CASE Username WHEN "YamiOdymel" THEN 123 WHEN "Foobar" THEN 456 ELSE 789 END`),
	}))
	assertEqual(assert, `UPDATE Users SET Password = CASE Username WHEN "YamiOdymel" THEN 123 WHEN "Foobar" THEN 456 ELSE 789 END, Username = ? WHERE Username = ?`, query)
	assertParams(assert, []interface{}{"YamiOdymel", ""}, params)
}

func TestUpdateStruct(t *testing.T) {
	u := struct {
		Username string
		Password string
	}{
		Username: "YamiOdymel",
		Password: "test",
	}
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("Username = ?", "YamiOdymel").Update(u))
	assertEqual(assert, "UPDATE Users SET Password = ?, Username = ? WHERE Username = ?", query)
	assertParams(assert, []interface{}{"YamiOdymel", "YamiOdymel", "test"}, params)
}

func TestLimitUpdate(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Limit(10).Update(H{
		"Username": "Karisu",
		"Password": "123456",
	}))
	assertEqual(assert, "UPDATE Users SET Password = ?, Username = ? LIMIT 10", query)
	assertParams(assert, []interface{}{"Karisu", "123456"}, params)
}

//=======================================================
// Patch
//=======================================================

func TestPatch(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("Username = ?", "YamiOdymel").Patch(H{
		"Username": "",
		"Age":      0,
		"Height":   183,
		"Password": "123456",
	}))
	assertEqual(assert, "UPDATE Users SET Password = ?, Height = ? WHERE Username = ?", query)
	assertParams(assert, []interface{}{"YamiOdymel", "123456", 183}, params)
}

func TestPatchExcludeTypes(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("Username = ?", "YamiOdymel").Exclude(reflect.String).Patch(H{
		"Username": "",
		"Age":      0,
		"Height":   183,
		"Password": "123456",
	}))
	assertEqual(assert, "UPDATE Users SET Password = ?, Username = ?, Height = ? WHERE Username = ?", query)
	assertParams(assert, []interface{}{"YamiOdymel", "", 183, "YamiOdymel"}, params)
}

func TestPatchExcludeColumns(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("Username = ?", "YamiOdymel").Exclude("Username", "Age").Patch(H{
		"Username": "",
		"Age":      0,
		"Height":   0,
		"Password": "123456",
	}))
	assertEqual(assert, "UPDATE Users SET Password = ?, Username = ?, Age = ? WHERE Username = ?", query)
	assertParams(assert, []interface{}{"YamiOdymel", "123456", 0, ""}, params)
}

func TestPatchStruct(t *testing.T) {
	u := struct {
		Username string
		Password string
	}{
		Username: "",
		Password: "test",
	}
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("Username = ?", "YamiOdymel").Patch(u))
	assertEqual(assert, "UPDATE Users SET Password = ? WHERE Username = ?", query)
	assertParams(assert, []interface{}{"YamiOdymel", "test"}, params)
}

//=======================================================
// Exists
//=======================================================

func TestExists(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("Username = ?", "YamiOdymel").Exists())
	assertEqual(assert, "SELECT EXISTS(SELECT * FROM Users WHERE Username = ?)", query)
	assertParams(assert, []interface{}{"YamiOdymel"}, params)
}

//=======================================================
// Select
//=======================================================

func TestSelect(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Select())
	assertEqual(assert, "SELECT * FROM Users", query)
}

func TestOffsetSelect(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Offset(10, 20).Select())
	assertEqual(assert, "SELECT * FROM Users LIMIT 10 OFFSET 20", query)
}

func TestLimitSelect(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Limit(10).Select())
	assertEqual(assert, "SELECT * FROM Users LIMIT 10", query)
}

func TestClearLimit(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Limit(10).Select().ClearLimit())
	assertEqual(assert, "SELECT * FROM Users", query)

	query, _ = Build(NewQuery("Users").Offset(10, 20).Select().ClearLimit())
	assertEqual(assert, "SELECT * FROM Users", query)

	query, _ = Build(NewQuery("Users").Limit(10, 100).Select().ClearLimit())
	assertEqual(assert, "SELECT * FROM Users", query)
}

func TestLimitMultiSelect(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Limit(10, 20).Select())
	assertEqual(assert, "SELECT * FROM Users LIMIT 10, 20", query)
}

func TestPaginateSelect(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Paginate(1, 100).Select())
	assertEqual(assert, "SELECT * FROM Users LIMIT 0, 100", query)

	query, _ = Build(NewQuery("Users").Paginate(2, 100).Select())
	assertEqual(assert, "SELECT * FROM Users LIMIT 100, 100", query)

	query, _ = Build(NewQuery("Users").Paginate(3, 100).Select())
	assertEqual(assert, "SELECT * FROM Users LIMIT 200, 100", query)
}

func TestGetColumns(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Select("Username", "Nickname"))
	assertEqual(assert, "SELECT Username, Nickname FROM Users", query)

	query, _ = Build(NewQuery("Users").Select("COUNT(*) AS Count"))
	assertEqual(assert, "SELECT COUNT(*) AS Count FROM Users", query)

	query, _ = Build(NewQuery("Users").Select("SUM(ID)", "COUNT(*) AS Count"))
	assertEqual(assert, "SELECT SUM(ID), COUNT(*) AS Count FROM Users", query)
}

func TestSelectOne(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("ID = ?", 1).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID = ?", query)
	assertParams(assert, []interface{}{1}, params)

	query, _ = Build(NewQuery("Users").Limit(1).Select())
	assertEqual(assert, "SELECT * FROM Users LIMIT 1", query)

	query, _ = Build(NewQuery("Users").SelectOne())
	assertEqual(assert, "SELECT * FROM Users LIMIT 1", query)

	query, _ = Build(NewQuery("Users").SelectOne("Username", "Nickname"))
	assertEqual(assert, "SELECT Username, Nickname FROM Users LIMIT 1", query)
}

//=======================================================
// Raw Query
//=======================================================

func TestRawQuery(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewRawQuery("SELECT * FROM Users WHERE ID >= ?", 10))
	assertEqual(assert, "SELECT * FROM Users WHERE ID >= ?", query)
	assertParams(assert, []interface{}{10}, params)
}

//=======================================================
// Where
//=======================================================

func TestWhere(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("ID = ?", 1).Where("Username = ?", "admin").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID = ? AND Username = ?", query)
	assertParams(assert, []interface{}{1, "admin"}, params)

	query, params = Build(NewQuery("Users").Where("ID = ?", 1).OrWhere("Username = ?", "admin").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID = ? OR Username = ?", query)
	assertParams(assert, []interface{}{1, "admin"}, params)
}

func TestWhereEscape(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("?? = ?", "ID", 1).Where("?? = ?", "Username", "admin").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE `ID` = ? AND `Username` = ?", query)
	assertParams(assert, []interface{}{1, "admin"}, params)

	query, params = Build(NewQuery("Users").Where("?? = ?", "ID", 1).OrWhere("?? = ?", "Username", "admin").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE `ID` = ? OR `Username` = ?", query)
	assertParams(assert, []interface{}{1, "admin"}, params)
}

func TestWhereQuery(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("(ID = ?)", 1).Where("Username = ?", "admin").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE (ID = ?) AND Username = ?", query)
	assertParams(assert, []interface{}{1, "admin"}, params)

	query, params = Build(NewQuery("Users").Where("(ID = ? OR Password = SHA(?))", 1, "password").Where("Username = ?", "admin").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE (ID = ? OR Password = SHA(?)) AND Username = ?", query)
	assertParams(assert, []interface{}{1, "password", "admin"}, params)

	query, params = Build(NewQuery("Users").Where("(ID = ? OR Password = SHA(?))", 1, "password").OrWhere("Username = ?", "admin").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE (ID = ? OR Password = SHA(?)) OR Username = ?", query)
	assertParams(assert, []interface{}{1, "password", "admin"}, params)

	query, params = Build(NewQuery("Users").Where("(ID = ? OR Password = SHA(?) OR Username = ?)", 1, "password", "Hello").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE (ID = ? OR Password = SHA(?) OR Username = ?)", query)
	assertParams(assert, []interface{}{1, "password", "Hello"}, params)
}

func TestWhereExpr(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("(ID = ?)", 1).Where("Username = ?", "admin").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE (ID = ?) AND Username = ?", query)
	assertParams(assert, []interface{}{1, "admin"}, params)

	query, params = Build(NewQuery("Users").Where("(ID = ? OR Password = SHA(?))", 1, "password").Where("Username = ?", "admin").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE (ID = ? OR Password = SHA(?)) AND Username = ?", query)
	assertParams(assert, []interface{}{1, "password", "admin"}, params)

	query, params = Build(NewQuery("Users").Where("(ID = ? OR Password = SHA(?))", 1, "password").Where("Username = ?", "admin").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE (ID = ? OR Password = SHA(?)) AND Username = ?", query)
	assertParams(assert, []interface{}{1, "password", "admin"}, params)
}

func TestWhereHaving(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("ID = ?", 1).Having("Username = ?", "admin").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID = ? HAVING Username = ?", query)
	assertParams(assert, []interface{}{1, "admin"}, params)

	query, params = Build(NewQuery("Users").Where("ID = ?", 1).Having("Username = ?", "admin").OrHaving("Password = ?", "test").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID = ? HAVING Username = ? OR Password = ?", query)
	assertParams(assert, []interface{}{1, "admin", "test"}, params)
}

func TestWhereColumns(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("LastLogin = CreatedAt").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE LastLogin = CreatedAt", query)

	query, _ = Build(NewQuery("Users").Where("LastLogin = CreatedAt").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE LastLogin = CreatedAt", query)
}

func TestWhereOperator(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("ID >= ?", 50).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID >= ?", query)
	assertParams(assert, []interface{}{50}, params)
}

func TestWhereLike(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("ID LIKE ?", 50).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID LIKE ?", query)
	assertParams(assert, []interface{}{50}, params)

	query, params = Build(NewQuery("Users").Where("ID NOT LIKE ?", 50).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID NOT LIKE ?", query)
	assertParams(assert, []interface{}{50}, params)
}

func TestWhereBetween(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("ID BETWEEN ? AND ?", 0, 20).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID BETWEEN ? AND ?", query)
	assertParams(assert, []interface{}{0, 20}, params)

	now := time.Now()
	nowAdd := time.Now().Add(time.Second * 60)
	query, params = Build(NewQuery("Users").Where("ID BETWEEN ? AND ?", now, nowAdd).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID BETWEEN ? AND ?", query)
	assertParams(assert, []interface{}{now, nowAdd}, params)

	query, params = Build(NewQuery("Users").Where("ID NOT BETWEEN ? AND ?", 0, 20).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID NOT BETWEEN ? AND ?", query)
	assertParams(assert, []interface{}{0, 20}, params)

	query, params = Build(NewQuery("Users").Where("ID NOT BETWEEN ? AND ?", now, nowAdd).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID NOT BETWEEN ? AND ?", query)
	assertParams(assert, []interface{}{now, nowAdd}, params)
}

func TestWhereIn(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("ID IN (?, ?, ?, ?, ?)", 1, 5, 27, -1, "d").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID IN (?, ?, ?, ?, ?)", query)
	assertParams(assert, []interface{}{1, 5, 27, -1, "d"}, params)

	query, params = Build(NewQuery("Users").Where("ID IN ?", []interface{}{1, 5, 27, -1, "d"}).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID IN (?, ?, ?, ?, ?)", query)
	assertParams(assert, []interface{}{1, 5, 27, -1, "d"}, params)

	query, params = Build(NewQuery("Users").Where("ID IN ?", []int{1, 5, 27, -1}).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID IN (?, ?, ?, ?)", query)
	assertParams(assert, []interface{}{1, 5, 27, -1}, params)

	query, params = Build(NewQuery("Users").Where("ID IN (?)", 1).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID IN (?)", query)
	assertParams(assert, []interface{}{1}, params)

	query, params = Build(NewQuery("Users").Where("ID NOT IN (?, ?, ?, ?, ?)", 1, 5, 27, -1, "d").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID NOT IN (?, ?, ?, ?, ?)", query)
	assertParams(assert, []interface{}{1, 5, 27, -1, "d"}, params)

	query, params = Build(NewQuery("Users").Where("ID NOT IN ?", []interface{}{1, 5, 27, -1, "d"}).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID NOT IN (?, ?, ?, ?, ?)", query)
	assertParams(assert, []interface{}{1, 5, 27, -1, "d"}, params)
}

func TestOrWhere(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("FirstName = ?", "John").OrWhere("FirstName = ?", "Peter").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE FirstName = ? OR FirstName = ?", query)
	assertParams(assert, []interface{}{"John", "Peter"}, params)

	query, _ = Build(NewQuery("Users").Where("A = B").OrWhere("(A = C OR A = D)").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE A = B OR (A = C OR A = D)", query)
}

func TestWhereNull(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("LastName IS NULL").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE LastName IS NULL", query)

	query, _ = Build(NewQuery("Users").Where("LastName IS NOT NULL").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE LastName IS NOT NULL", query)
}

func TestWhereExists(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewQuery("Products").Select()

	query, _ := Build(NewQuery("Users").Where("EXISTS ?", subQuery).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE EXISTS (SELECT * FROM Products)", query)

	query, _ = Build(NewQuery("Users").Where("NOT EXISTS ?", subQuery).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE NOT EXISTS (SELECT * FROM Products)", query)
}

func TestRawWhere(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("ID != CompanyID").Where("DATE(CreatedAt) = DATE(LastLogin)").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID != CompanyID AND DATE(CreatedAt) = DATE(LastLogin)", query)

	query, params := Build(NewQuery("Users").Where("ID != CompanyID").Where("DATE(CreatedAt) = DATE(LastLogin)").Where("Username = ?", "YamiOdymel").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID != CompanyID AND DATE(CreatedAt) = DATE(LastLogin) AND Username = ?", query)
	assertParams(assert, []interface{}{"YamiOdymel"}, params)
}

//=======================================================
// As
//=======================================================

func TestAs(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery(NewAlias("Products", "p")).Select())
	assertEqual(assert, "SELECT * FROM Products AS p", query)
}

//=======================================================
// Distinct
//=======================================================

func TestDistinct(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Distinct().Select("Username"))
	assertEqual(assert, "SELECT DISTINCT Username FROM Users", query)
}

//=======================================================
// Union
//=======================================================

func TestUnion(t *testing.T) {
	assert := assert.New(t)
	tableQuery := NewQuery("Locations").Select()

	query, _ := Build(NewQuery("Users").Union(tableQuery).Select())
	assertEqual(assert, "SELECT * FROM Users UNION SELECT * FROM Locations", query)

	query, params := Build(NewQuery(NewQuery("Users").Union(tableQuery).Select()).Where("Username = ?", "YamiOdymel").Select())
	assertEqual(assert, "SELECT * FROM (SELECT * FROM Users UNION SELECT * FROM Locations) WHERE Username = ?)", query)
	assertParams(assert, []interface{}{"YamiOdymel"}, params)
}

func TestUnionAll(t *testing.T) {
	assert := assert.New(t)
	tableQuery := NewQuery("Locations").Select()

	query, _ := Build(NewQuery("Users").UnionAll(tableQuery).Select())
	assertEqual(assert, "SELECT * FROM Users UNION ALL SELECT * FROM Locations", query)

	query, params := Build(NewQuery(
		NewQuery("Users").UnionAll(tableQuery).Select(),
	).As("Result").Where("Username = ?", "YamiOdymel").Select())
	assertEqual(assert, "SELECT * FROM (SELECT * FROM Users UNION ALL SELECT * FROM Locations) AS Result WHERE Username = ?", query)
	assertParams(assert, []interface{}{"YamiOdymel"}, params)
}

//=======================================================
// Delete
//=======================================================

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Where("ID = ?", 1).Delete())
	assertEqual(assert, "DELETE FROM Users WHERE ID = ?", query)
	assertParams(assert, []interface{}{1}, params)
}

//=======================================================
// OrderBy
//=======================================================

func TestOrderBy(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").OrderBy("ID ASC").OrderBy("Login DESC").OrderBy("RAND()").Select())
	assertEqual(assert, "SELECT * FROM Users ORDER BY ID ASC, Login DESC, RAND()", query)
}

func TestOrderByField(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").OrderByField("UserGroup", "SuperUser", "Admin", "Users").Select())
	assertEqual(assert, "SELECT * FROM Users ORDER BY FIELD (UserGroup, ?, ?, ?)", query)
	assertParams(assert, []interface{}{"SuperUser", "Admin", "Users"}, params)
}

//=======================================================
// GroupBy
//=======================================================

func TestGroupBy(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").GroupBy("Name").Select())
	assertEqual(assert, "SELECT * FROM Users GROUP BY Name", query)
	query, _ = Build(NewQuery("Users").GroupBy("Name", "ID").Select())
	assertEqual(assert, "SELECT * FROM Users GROUP BY Name, ID", query)
}

//=======================================================
// Join
//=======================================================

func TestJoin(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Products").
		CrossJoin("Users", "Products.TenantID = Users.TenantID").
		Where("Users.ID = ?", 6).
		Select("Users.Name", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products CROSS JOIN Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?", query)
	assertParams(assert, []interface{}{6}, params)

	query, params = Build(NewQuery("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		Where("Users.ID = ?", 6).
		Select("Users.Name", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products LEFT JOIN Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?", query)
	assertParams(assert, []interface{}{6}, params)

	query, params = Build(NewQuery("Products").
		RightJoin("Users", "Products.TenantID = Users.TenantID").
		Where("Users.ID = ?", 6).
		Select("Users.Name", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products RIGHT JOIN Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?", query)
	assertParams(assert, []interface{}{6}, params)

	query, params = Build(NewQuery("Products").
		InnerJoin("Users", "Products.TenantID = Users.TenantID").
		Where("Users.ID = ?", 6).
		Select("Users.Name", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products INNER JOIN Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?", query)
	assertParams(assert, []interface{}{6}, params)

	query, params = Build(NewQuery("Products").
		NaturalJoin("Users", "Products.TenantID = Users.TenantID").
		Where("Users.ID = ?", 6).
		Select("Users.Name", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products NATURAL JOIN Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?", query)
	assertParams(assert, []interface{}{6}, params)

	query, params = Build(NewQuery("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		RightJoin("Posts", "Products.TenantID = Posts.TenantID").
		Where("Users.ID = ?", 6).
		Select("Users.Name", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products LEFT JOIN Users ON (Products.TenantID = Users.TenantID) RIGHT JOIN Posts ON (Products.TenantID = Posts.TenantID) WHERE Users.ID = ?", query)
	assertParams(assert, []interface{}{6}, params)
}

func TestJoinAlias(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery(NewAlias("Products", "p")).
		LeftJoin(NewAlias("Users", "u"), "p.TenantID = u.TenantID").
		Where("u.ID = ?", 6).
		Select("u.Name", "p.ProductName"))
	assertEqual(assert, "SELECT u.Name, p.ProductName FROM Products AS p LEFT JOIN Users AS u ON (p.TenantID = u.TenantID) WHERE u.ID = ?", query)
	assertParams(assert, []interface{}{6}, params)
}

func TestJoinWhere(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		OrJoinWhere("Users.TenantID = ?", 5).
		Select("Users.Name", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products LEFT JOIN Users ON (Products.TenantID = Users.TenantID OR Users.TenantID = ?)", query)
	assertParams(assert, []interface{}{5}, params)

	query, params = Build(NewQuery("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		JoinWhere("Users.Username = ?", "Wow").
		Select("Users.Name", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products LEFT JOIN Users ON (Products.TenantID = Users.TenantID AND Users.Username = ?)", query)
	assertParams(assert, []interface{}{"Wow"}, params)

	query, params = Build(NewQuery("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		RightJoin("Posts", "Products.TenantID = Posts.TenantID").
		JoinWhere("Posts.Username = ?", "Wow").
		JoinWhere("Users.Username = ?", "Wow").
		Select("Users.Name", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products LEFT JOIN Users ON (Products.TenantID = Users.TenantID AND Users.Username = ?) RIGHT JOIN Posts ON (Products.TenantID = Posts.TenantID AND Posts.Username = ?)", query)
	assertParams(assert, []interface{}{"Wow", "Wow"}, params)
}

//=======================================================
// Sub Query
//=======================================================

func TestSubQuerySelect(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewQuery("Products").Where("Quantity > ?", 2).Select("UserID")
	query, params := Build(NewQuery("Users").Where("ID IN ?", subQuery).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID IN (SELECT UserID FROM Products WHERE Quantity > ?)", query)
	assertParams(assert, []interface{}{2}, params)
}

func TestSubQueryInsert(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewQuery("Users").Where("ID = ?", 6).SelectOne("Name")
	query, params := Build(NewQuery("Products").Insert(H{
		"ProductName": "測試商品",
		"UserID":      subQuery,
		"LastUpdated": NewExpr("NOW()"),
	}))
	assertEqual(assert, "INSERT INTO Products (LastUpdated, ProductName, UserID) VALUES (NOW(), ?, (SELECT Name FROM Users WHERE ID = ? LIMIT 1))", query)
	assertParams(assert, []interface{}{"測試商品", 6}, params)
}

func TestSubQueryJoin(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewQuery("Users").As("Users").Where("Active = ?", 1).Select()
	query, params := Build(NewQuery("Products").
		LeftJoin(subQuery, "Products.UserID = Users.ID").
		Select("Users.Username", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Username, Products.ProductName FROM Products LEFT JOIN (SELECT * FROM Users WHERE Active = ?) AS Users ON (Products.UserID = Users.ID)", query)
	assertParams(assert, []interface{}{1}, params)
}

func TestSubQueryJoinWhere(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewQuery("Users").As("Users").Where("Active = ?", 1).Select()
	query, params := Build(NewQuery("Products").
		LeftJoin(subQuery, "Products.UserID = Users.ID").
		JoinWhere("Users.Username = ?", "Hello").
		Select("Users.Username", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Username, Products.ProductName FROM Products LEFT JOIN (SELECT * FROM Users WHERE Active = ?) AS Users ON (Products.UserID = Users.ID AND Users.Username = ?)", query)
	assertParams(assert, []interface{}{1, "Hello"}, params)
}

func TestSubQueryExists(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewQuery("Users").Where("Company = ?", "測試公司").Select("UserID")
	query, params := Build(NewQuery("Products").Where("EXISTS ?", subQuery).Select())
	assertEqual(assert, "SELECT * FROM Products WHERE EXISTS (SELECT UserID FROM Users WHERE Company = ?)", query)
	assertParams(assert, []interface{}{"測試公司"}, params)

	subQuery = NewQuery("Products").Where("Quantity > ?", 2).Select("UserID")
	query, params = Build(NewQuery("Users").Where("ID IN ?", subQuery).Exists())
	assertEqual(assert, "SELECT EXISTS(SELECT * FROM Users WHERE ID IN (SELECT UserID FROM Products WHERE Quantity > ?))", query)
	assertParams(assert, []interface{}{2}, params)
}

func TestSubQueryRawQuery(t *testing.T) {
	assert := assert.New(t)
	rawQuery := NewRawQuery("SELECT UserID FROM Users WHERE Company = ?", "測試公司")
	query, params := Build(NewQuery("Products").Where("EXISTS ?", rawQuery).Select())
	assertEqual(assert, "SELECT * FROM Products WHERE EXISTS (SELECT UserID FROM Users WHERE Company = ?)", query)
	assertParams(assert, []interface{}{"測試公司"}, params)
}

func TestSubQueryRawQueryReplacement(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewQuery("Locations").Where("Username = ?", "YamiOdymel").Select()
	rawQuery := NewRawQuery("SELECT UserID FROM Users WHERE EXISTS (?)", subQuery)
	query, params := Build(NewQuery("Products").Where("EXISTS ?", rawQuery).Select())
	assertEqual(assert, "SELECT * FROM Products WHERE EXISTS (SELECT UserID FROM Users WHERE EXISTS (SELECT * FROM Locations WHERE Username = ?))", query)
	assertParams(assert, []interface{}{"YamiOdymel"}, params)
}

//=======================================================
// Copy
//=======================================================

func TestCopy(t *testing.T) {
	assert := assert.New(t)
	q := NewQuery("Users").Where("?? = ? AND ?? = ?", "user_id", 30, "nickname", "yamiodymel").Where("?? = ?", "username", "hello").Select()

	query, params := Build(q.Copy().Limit(100))
	assertEqual(assert, "SELECT * FROM Users WHERE `user_id` = ? AND `nickname` = ? AND `username` = ? LIMIT 100", query)
	assertParams(assert, []interface{}{30, "yamiodymel", "hello"}, params)

	query, params = Build(q)
	assertEqual(assert, "SELECT * FROM Users WHERE `user_id` = ? AND `nickname` = ? AND `username` = ?", query)
	assertParams(assert, []interface{}{30, "yamiodymel", "hello"}, params)
}

//=======================================================
// Others
//=======================================================

func TestComplexQueries(t *testing.T) {
	assert := assert.New(t)
	jobHistories := NewQuery("JobHistories").
		Where("DepartmentID BETWEEN ? AND ?", 50, 100).
		Select("JobID")
	jobs := NewQuery("Jobs").
		Where("JobID IN ?", jobHistories).
		GroupBy("JobID").
		Select("JobID", "AVG(MinSalary) AS MyAVG")
	maxAverage := NewQuery(jobs).
		As("SS").
		Select("MAX(MyAVG)")
	employees := NewQuery("Employees").
		GroupBy("JobID").
		Having("AVG(Salary) < ?", maxAverage).
		Select("JobID", "AVG(Salary)")
	query, params := Build(employees)

	assertEqual(assert, "SELECT JobID, AVG(Salary) FROM Employees HAVING AVG(Salary) < (SELECT MAX(MyAVG) FROM (SELECT JobID, AVG(MinSalary) AS MyAVG FROM Jobs WHERE JobID IN (SELECT JobID FROM JobHistories WHERE DepartmentID BETWEEN ? AND ?) GROUP BY JobID) AS SS) GROUP BY JobID", query)
	assertParams(assert, []interface{}{50, 100}, params)
	assertParamOrders(assert, []interface{}{50, 100}, params)

	// Example Source: https://www.w3resource.com/sql/subqueries/nested-subqueries.php
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

	agents := NewQuery("Agents").
		Where("Commission < ?", 0.12).
		Select()
	customers := NewQuery("Customers").
		Where("Grade = ?", 3).
		Where("CustomerCountry <> ?", "India").
		Where("OpeningAmount < ?", 7000).
		Where("EXISTS ?", agents).
		Select("OutstandingAmount")
	orders := NewQuery("Orders").
		Where("OrderAmount > ?", 2000).
		Where("OrderDate < ?", "01-SEP-08").
		Where("AdvanceAmount < ANY (?)", customers).
		Select("OrderNum", "OrderDate", "OrderAmount", "AdvanceAmount")
	query, params = Build(orders)

	assertEqual(assert, "SELECT OrderNum, OrderDate, OrderAmount, AdvanceAmount FROM Orders WHERE OrderAmount > ? AND OrderDate < ? AND AdvanceAmount < ANY (SELECT OutstandingAmount FROM Customers WHERE Grade = ? AND CustomerCountry <> ? AND OpeningAmount < ? AND EXISTS (SELECT * FROM Agents WHERE Commission < ?))", query)
	assertParams(assert, []interface{}{2000, "01-SEP-08", 3, "India", 7000, 0.12}, params)
	assertParamOrders(assert, []interface{}{2000, "01-SEP-08", 3, "India", 7000, 0.12}, params)

	// Example Source: https://www.w3resource.com/sql/subqueries/nested-subqueries.php
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
}

func TestSetQueryOption(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").SetQueryOption("FOR UPDATE").Select("Username"))
	assertEqual(assert, "SELECT Username FROM Users FOR UPDATE", query)
}
