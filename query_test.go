package rushia

import (
	"reflect"
	"strings"
	"testing"

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
	return
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
	query, _ := Build(NewQuery("Users").Insert(u))
	assertEqual(assert, "INSERT INTO Users (Username, Password) VALUES (?, ?)", query)
}

func TestInsert(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Insert(H{
		"Username": "YamiOdymel",
		"Password": "test",
	}))
	assertEqual(assert, "INSERT INTO Users (Username, Password) VALUES (?, ?)", query)
}

func TestInsertMap(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Insert(map[string]interface{}{
		"Username": "YamiOdymel",
		"Password": "test",
	}))
	assertEqual(assert, "INSERT INTO Users (Username, Password) VALUES (?, ?)", query)
}

func TestInsertParams(t *testing.T) {
	assert := assert.New(t)
	query, params := Build(NewQuery("Users").Insert(H{
		"Username": "YamiOdymel",
		"Password": "test",
	}))
	assertEqual(assert, "INSERT INTO Users (Username, Password) VALUES (?, ?)", query)
	assert.Len(params, 2)
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
	query, _ := Build(NewQuery("Users").Insert(data))
	assertEqual(assert, "INSERT INTO Users (Password, Username) VALUES (?, ?), (?, ?)", query)
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
	query, _ := Build(NewQuery("Users").Insert(data))
	assertEqual(assert, "INSERT INTO Users (Password, Username) VALUES (?, ?), (?, ?)", query)
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
	query, _ := Build(NewQuery("Users").Omit("Username").Insert(u))
	assertEqual(assert, "INSERT INTO Users (Password) VALUES (?)", query)
}

func TestInsertStructTagOmit(t *testing.T) {
	u := struct {
		Username string `rushia:"-"`
		Password string
	}{
		Username: "YamiOdymel",
		Password: "test",
	}
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Insert(u))
	assertEqual(assert, "INSERT INTO Users (Password) VALUES (?)", query)
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
	query, _ := Build(NewQuery("Users").Insert(u))
	assertEqual(assert, "INSERT INTO Users (user_name, Password) VALUES (?, ?)", query)
}

func TestInsertOmit(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Omit("Username").Insert(H{
		"Username": "YamiOdymel",
		"Password": "test",
	}))
	assertEqual(assert, "INSERT INTO Users (Password) VALUES (?)", query)
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
	query, _ := Build(NewQuery("Users").Omit("Username").Insert(data))
	assertEqual(assert, "INSERT INTO Users (Password) VALUES (?), (?)", query)
}

func TestInsertExpr(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Insert(H{
		"Username":  "YamiOdymel",
		"Password":  NewExpr("SHA1(?)", "secretpassword+salt"),
		"Expires":   NewExpr("NOW() + INTERVAL 1 YEAR"),
		"CreatedAt": NewExpr("NOW()"),
	}))
	assertEqual(assert, "INSERT INTO Users (CreatedAt, Expires, Password, Username) VALUES (NOW(), NOW() + INTERVAL 1 YEAR, SHA1(?), ?)", query)
}

func TestInsertSubQueryExpr(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewQuery("Salaries").Where("Username", "YamiOdymel").Select("Salary")
	query, _ := Build(NewQuery("Users").Insert(H{
		"Username":  "YamiOdymel",
		"AvgSalary": NewExpr("SUM((?))", subQuery),
	}))
	assertEqual(assert, "INSERT INTO Users (Username, AvgSalary) VALUES (?, SUM((SELECT Salary FROM Salaries WHERE Username = ?)))", query)

	subQuery = NewQuery("Salaries").Where("Username", "YamiOdymel").Select("Salary")
	query, _ = Build(NewQuery("Users").Insert(H{
		"Username": "YamiOdymel",
		"Salary":   subQuery,
	}))
	assertEqual(assert, "INSERT INTO Users (Username, Salary) VALUES (?, (SELECT Salary FROM Salaries WHERE Username = ?))", query)
}

func TestInsertSelect(t *testing.T) {
	assert := assert.New(t)
	from := NewQuery("AdditionalUsers").WhereLike("Name", "ABC%").Select("ID", "Username", "Nickname")
	query, _ := Build(NewQuery("Users").InsertSelect(from, "ID", "Username", "Nickname"))
	assertEqual(assert, "INSERT INTO Users (ID, Username, Nickname) SELECT ID, Username, Nickname FROM AdditionalUsers WHERE Username LIKE ?", query)
}

func TestOnDuplicateInsert(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").OnDuplicate(H{
		"UpdatedAt": NewExpr("VALUES(UpdatedAt)"), // Deprecated in MySQL 8.0.20
		"ID":        NewExpr("LAST_INSERT_ID(ID)"),
	}).Insert(H{
		"Username":  "YamiOdymel",
		"Password":  "test",
		"UpdatedAt": NewExpr("NOW()"),
	}))
	assertEqual(assert, "INSERT INTO Users (Password, UpdatedAt, Username) VALUES (?, NOW(), ?) ON DUPLICATE KEY UPDATE ID = LAST_INSERT_ID(ID), UpdatedAt = VALUES(UpdatedAt)", query)

	query, _ = Build(NewQuery("Users").As("New").OnDuplicate(H{
		"UpdatedAt": NewExpr("New.UpdatedAt"),
		"ID":        NewExpr("LAST_INSERT_ID(ID)"),
	}).Insert(H{
		"Username":  "YamiOdymel",
		"Password":  "test",
		"UpdatedAt": NewExpr("NOW()"),
	}))
	assertEqual(assert, "INSERT INTO Users (Password, UpdatedAt, Username) VALUES (?, NOW(), ?) AS New ON DUPLICATE KEY UPDATE ID = LAST_INSERT_ID(ID), UpdatedAt = New.UpdatedAt", query)
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
	query, _ := Build(NewQuery("Users").Replace(H{
		"Username": "YamiOdymel",
		"Password": "test",
	}))
	assertEqual(assert, "REPLACE INTO Users (Password, Username) VALUES (?, ?)", query)
}

//=======================================================
// Update
//=======================================================

func TestUpdateOmit(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("Username", "YamiOdymel").Omit("Username").Update(H{
		"Username": "Karisu",
		"Password": "123456",
	}))
	assertEqual(assert, "UPDATE Users SET Password = ? WHERE Username = ?", query)
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
	query, _ := Build(NewQuery("Users").Where("Username", "YamiOdymel").Omit("Username").Update(u))
	assertEqual(assert, "UPDATE Users SET Password = ? WHERE Username = ?", query)
}

func TestUpdate(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("Username", "YamiOdymel").Update(H{
		"Username": "",
		"Password": "123456",
	}))
	assertEqual(assert, "UPDATE Users SET Password = ?, Username = ? WHERE Username = ?", query)
}

func TestUpdateCase(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("Username", "YamiOdymel").Update(H{
		"Username": "",
		"Password": NewExpr(`CASE Username WHEN "YamiOdymel" THEN 123 WHEN "Foobar" THEN 456 ELSE 789 END`),
	}))
	assertEqual(assert, `UPDATE Users SET Password = CASE Username WHEN "YamiOdymel" THEN 123 WHEN "Foobar" THEN 456 ELSE 789 END, Username = ? WHERE Username = ?`, query)
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
	query, _ := Build(NewQuery("Users").Where("Username", "YamiOdymel").Update(u))
	assertEqual(assert, "UPDATE Users SET Password = ?, Username = ? WHERE Username = ?", query)
}

func TestLimitUpdate(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Limit(10).Update(H{
		"Username": "Karisu",
		"Password": "123456",
	}))
	assertEqual(assert, "UPDATE Users SET Password = ?, Username = ? LIMIT 10", query)
}

//=======================================================
// Patch
//=======================================================

func TestPatch(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("Username", "YamiOdymel").Patch(H{
		"Username": "",
		"Age":      0,
		"Height":   183,
		"Password": "123456",
	}))
	assertEqual(assert, "UPDATE Users SET Password = ?, Height = ? WHERE Username = ?", query)
}

func TestPatchExcludeTypes(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("Username", "YamiOdymel").Exclude(reflect.String).Patch(H{
		"Username": "",
		"Age":      0,
		"Height":   183,
		"Password": "123456",
	}))
	assertEqual(assert, "UPDATE Users SET Password = ?, Username = ?, Height = ? WHERE Username = ?", query)
}

func TestPatchExcludeColumns(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("Username", "YamiOdymel").Exclude("Username", "Age").Patch(H{
		"Username": "",
		"Age":      0,
		"Height":   0,
		"Password": "123456",
	}))
	assertEqual(assert, "UPDATE Users SET Password = ?, Username = ?, Age = ? WHERE Username = ?", query)
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
	query, _ := Build(NewQuery("Users").Where("Username", "YamiOdymel").Patch(u))
	assertEqual(assert, "UPDATE Users SET Password = ? WHERE Username = ?", query)
}

//=======================================================
// Exists
//=======================================================

func TestExists(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("Username", "YamiOdymel").Exists())
	assertEqual(assert, "SELECT EXISTS(SELECT * FROM Users WHERE Username = ?)", query)
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

func TestLimitMultiSelect(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Limit(10, 20).Select())
	assertEqual(assert, "SELECT * FROM Users LIMIT 10, 20", query)
}

func TestGetColumns(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Select("Username", "Nickname"))
	assertEqual(assert, "SELECT Username, Nickname FROM Users", query)

	query, _ = Build(NewQuery("Users").Select("COUNT(*) AS Count"))
	assertEqual(assert, "SELECT COUNT(*) AS Count FROM Users", query)

	query, _ = Build(NewQuery("Users").Select(NewExpr("COUNT(*) AS Count")))
	assertEqual(assert, "SELECT COUNT(*) AS Count FROM Users", query)

	query, _ = Build(NewQuery("Users").Select(NewExpr("SUM(ID)"), NewExpr("COUNT(*) AS Count")))
	assertEqual(assert, "SELECT SUM(ID), COUNT(*) AS Count FROM Users", query)
}

func TestSelectOne(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("ID", 1).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID = ?", query)

	query, _ = Build(NewQuery("Users").Limit(1).Select())
	assertEqual(assert, "SELECT * FROM Users LIMIT 1", query)

	query, _ = Build(NewQuery("Users").SelectOne())
	assertEqual(assert, "SELECT * FROM Users LIMIT 1", query)

	query, _ = Build(NewQuery("Users").SelectOne("Username", "Nickname"))
	assertEqual(assert, "SELECT Username, Nickname FROM Users", query)
}

//=======================================================
// Raw Query
//=======================================================

func TestRawQuery(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewRawQuery("SELECT * FROM Users WHERE ID >= ?", 10))
	assertEqual(assert, "SELECT * FROM Users WHERE ID >= ?", query)
}

//=======================================================
// Where
//=======================================================

func TestWhere(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("ID", 1).Where("Username", "admin").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID = ? AND Username = ?", query)

	query, _ = Build(NewQuery("Users").Where("ID", 1).OrWhere("Username", "admin").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID = ? OR Username = ?", query)
}

func TestWhereQuery(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("(ID = ?)", 1).Where("Username", "admin").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE (ID = ?) AND Username = ?", query)
	query, _ = Build(NewQuery("Users").Where("(ID = ? OR Password = SHA(?))", 1, "password").Where("Username", "admin").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE (ID = ? OR Password = SHA(?)) AND Username = ?", query)
	query, _ = Build(NewQuery("Users").Where("(ID = ? OR Password = SHA(?))", 1, "password").OrWhere("Username", "admin").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE (ID = ? OR Password = SHA(?)) OR Username = ?", query)
	query, _ = Build(NewQuery("Users").Where("(ID = ? OR Password = SHA(?) OR Username = ?)", 1, "password", "Hello").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE (ID = ? OR Password = SHA(?) OR Username = ?)", query)
}

func TestWhereExpr(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where(NewExpr("(ID = ?)", 1)).Where("Username", "admin").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE (ID = ?) AND Username = ?", query)
	query, _ = Build(NewQuery("Users").Where(NewExpr("(ID = ? OR Password = SHA(?))", 1, "password")).Where("Username", "admin").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE (ID = ? OR Password = SHA(?)) AND Username = ?", query)
}

func TestWhereHaving(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("ID", 1).Having("Username", "admin").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID = ? HAVING Username = ?", query)
	query, _ = Build(NewQuery("Users").Where("ID", 1).Having("Username", "admin").OrHaving("Password", "test").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID = ? HAVING Username = ? OR Password = ?", query)
}

func TestWhereColumns(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("LastLogin = CreatedAt").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE LastLogin = CreatedAt", query)

	query, _ = Build(NewQuery("Users").WhereColumn("LastLogin", "=", "CreatedAt").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE LastLogin = CreatedAt", query)
}

func TestWhereOperator(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("ID", ">=", 50).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID >= ?", query)
}

func TestWhereLike(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").WhereLike("ID", 50).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID LIKE ?", query)

	query, _ = Build(NewQuery("Users").Where("ID", "LIKE", 50).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID LIKE ?", query)

	query, _ = Build(NewQuery("Users").WhereNotLike("ID", 50).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID NOT LIKE ?", query)

	query, _ = Build(NewQuery("Users").Where("ID", "NOT LIKE", 50).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID NOT LIKE ?", query)
}

func TestWhereBetween(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").WhereBetween("ID", 0, 20).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID BETWEEN ? AND ?", query)

	query, _ = Build(NewQuery("Users").Where("ID", "BETWEEN", 0, 20).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID BETWEEN ? AND ?", query)

	query, _ = Build(NewQuery("Users").WhereNotBetween("ID", 0, 20).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID NOT BETWEEN ? AND ?", query)

	query, _ = Build(NewQuery("Users").Where("ID", "NOT BETWEEN", 0, 20).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID NOT BETWEEN ? AND ?", query)
}

func TestWhereIn(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("ID", "IN", 1, 5, 27, -1, "d").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID IN (?, ?, ?, ?, ?)", query)

	query, _ = Build(NewQuery("Users").WhereIn("ID", 1, 5, 27, -1, "d").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID IN (?, ?, ?, ?, ?)", query)

	query, _ = Build(NewQuery("Users").Where("ID", "NOT IN", 1, 5, 27, -1, "d").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID NOT IN (?, ?, ?, ?, ?)", query)

	query, _ = Build(NewQuery("Users").WhereNotIn("ID", 1, 5, 27, -1, "d").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID NOT IN (?, ?, ?, ?, ?)", query)
}

func TestOrWhere(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("FirstName", "John").OrWhere("FirstName", "Peter").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE FirstName = ? OR FirstName = ?", query)

	query, _ = Build(NewQuery("Users").Where("A = B").OrWhere("(A = C OR A = D)").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE A = B OR (A = C OR A = D)", query)
}

func TestWhereNull(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("LastName", "IS", nil).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE LastName IS NULL", query)

	query, _ = Build(NewQuery("Users").WhereIsNull("LastName").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE LastName IS NULL", query)

	query, _ = Build(NewQuery("Users").Where("LastName", "IS NOT", nil).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE LastName IS NOT NULL", query)

	query, _ = Build(NewQuery("Users").WhereIsNotNull("LastName").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE LastName IS NOT NULL", query)
}

func TestWhereExists(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewQuery("Products").Select()

	query, _ := Build(NewQuery("Users").Where("EXISTS", subQuery).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE EXISTS (SELECT * FROM Products)", query)

	query, _ = Build(NewQuery("Users").WhereExists(subQuery).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE EXISTS (SELECT * FROM Products)", query)

	query, _ = Build(NewQuery("Users").Where("NOT EXISTS", subQuery).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE NOT EXISTS (SELECT * FROM Products)", query)

	query, _ = Build(NewQuery("Users").WhereNotExists(subQuery).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE NOT EXISTS (SELECT * FROM Products)", query)
}

func TestRawWhere(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("ID != CompanyID").Where("DATE(CreatedAt) = DATE(LastLogin)").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID != CompanyID AND DATE(CreatedAt) = DATE(LastLogin)", query)
	query, _ = Build(NewQuery("Users").WhereRaw("ID != CompanyID").WhereRaw("DATE(CreatedAt) = DATE(LastLogin)").WhereRaw("Username = ?", "YamiOdymel").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID != CompanyID AND DATE(CreatedAt) = DATE(LastLogin) AND Username = ?", query)
}

func TestRawWhereParams(t *testing.T) {
	assert := assert.New(t)
	query, p := Build(NewQuery("Users").Where("(ID = ? OR ID = ?)", 6, 2).Where("Login", "Mike").Select())
	assertEqual(assert, "SELECT * FROM Users WHERE (ID = ? OR ID = ?) AND Login = ?", query)
	assert.Len(p, 3)
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

	query, _ = Build(NewQuery(NewQuery("Users").Union(tableQuery).Select()).Where("Username", "YamiOdymel").Select())
	assertEqual(assert, "SELECT * FROM (SELECT * FROM Users UNION SELECT * FROM Locations) WHERE Username = ?)", query)
}

func TestUnionAll(t *testing.T) {
	assert := assert.New(t)
	tableQuery := NewQuery("Locations").Select()

	query, _ := Build(NewQuery("Users").UnionAll(tableQuery).Select())
	assertEqual(assert, "SELECT * FROM Users UNION ALL SELECT * FROM Locations", query)

	query, _ = Build(NewQuery(
		NewQuery("Users").UnionAll(tableQuery).Select(),
	).As("Result").Where("Username", "YamiOdymel").Select())
	assertEqual(assert, "SELECT * FROM (SELECT * FROM Users UNION ALL SELECT * FROM Locations) AS Result WHERE Username = ?", query)
}

//=======================================================
// Delete
//=======================================================

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").Where("ID", 1).Delete())
	assertEqual(assert, "DELETE FROM Users WHERE ID = ?", query)
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
	query, _ := Build(NewQuery("Users").OrderByField("UserGroup", "SuperUser", "Admin", "Users").Select())
	assertEqual(assert, "SELECT * FROM Users ORDER BY FIELD (UserGroup, ?, ?, ?)", query)
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
	query, _ := Build(NewQuery("Products").
		CrossJoin("Users", "Products.TenantID = Users.TenantID").
		Where("Users.ID", 6).
		Select("Users.Name", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products CROSS JOIN Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?", query)

	query, _ = Build(NewQuery("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		Where("Users.ID", 6).
		Select("Users.Name", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products LEFT JOIN Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?", query)

	query, _ = Build(NewQuery("Products").
		RightJoin("Users", "Products.TenantID = Users.TenantID").
		Where("Users.ID", 6).
		Select("Users.Name", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products RIGHT JOIN Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?", query)

	query, _ = Build(NewQuery("Products").
		InnerJoin("Users", "Products.TenantID = Users.TenantID").
		Where("Users.ID", 6).
		Select("Users.Name", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products INNER JOIN Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?", query)

	query, _ = Build(NewQuery("Products").
		NaturalJoin("Users", "Products.TenantID = Users.TenantID").
		Where("Users.ID", 6).
		Select("Users.Name", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products NATURAL JOIN Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?", query)

	query, _ = Build(NewQuery("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		RightJoin("Posts", "Products.TenantID = Posts.TenantID").
		Where("Users.ID", 6).
		Select("Users.Name", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products LEFT JOIN Users ON (Products.TenantID = Users.TenantID) RIGHT JOIN Posts ON (Products.TenantID = Posts.TenantID) WHERE Users.ID = ?", query)
}

func TestJoinWhere(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		OrJoinWhere("Users.TenantID", 5).
		Select("Users.Name", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products LEFT JOIN Users ON (Products.TenantID = Users.TenantID OR Users.TenantID = ?)", query)

	query, _ = Build(NewQuery("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		JoinWhere("Users.Username", "Wow").
		Select("Users.Name", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products LEFT JOIN Users ON (Products.TenantID = Users.TenantID AND Users.Username = ?)", query)

	query, _ = Build(NewQuery("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		RightJoin("Posts", "Products.TenantID = Posts.TenantID").
		JoinWhere("Posts.Username", "Wow").
		JoinWhere("Users.Username", "Wow").
		Select("Users.Name", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products LEFT JOIN Users ON (Products.TenantID = Users.TenantID AND Users.Username = ?) RIGHT JOIN Posts ON (Products.TenantID = Posts.TenantID AND Posts.Username = ?)", query)
}

TEST PARAMS TRUE LA

//=======================================================
// Sub Query
//=======================================================

func TestSubQuerySelect(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewQuery("Products").Where("Quantity", ">", 2).Select("UserID")
	query, _ := Build(NewQuery("Users").WhereIn("ID", subQuery).Select())
	assertEqual(assert, "SELECT * FROM Users WHERE ID IN (SELECT UserID FROM Products WHERE Quantity > ?)", query)
}

func TestSubQueryInsert(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewQuery("Users").Where("ID", 6).Select("Name")
	query, _ := Build(NewQuery("Products").Insert(H{
		"ProductName": "測試商品",
		"UserID":      subQuery,
		"LastUpdated": NewExpr("NOW()"),
	}))
	assertEqual(assert, "INSERT INTO Products (LastUpdated, ProductName, UserID) VALUES (NOW(), ?, (SELECT Name FROM Users WHERE ID = ?))", query)
}

func TestSubQueryJoin(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewQuery("Users").As("Users").Where("Active", 1).Select()
	query, _ := Build(NewQuery("Products").
		LeftJoin(subQuery, "Products.UserID = Users.ID").
		Select("Users.Username", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Username, Products.ProductName FROM Products LEFT JOIN (SELECT * FROM Users WHERE Active = ?) AS Users ON (Products.UserID = Users.ID)", query)
}

func TestSubQueryJoinWhere(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewQuery("Users").As("Users").Where("Active", 1).Select()
	query, _ := Build(NewQuery("Products").
		LeftJoin(subQuery, "Products.UserID = Users.ID").
		JoinWhere("Users", "Users.Username", "Hello").
		Select("Users.Username", "Products.ProductName"))
	assertEqual(assert, "SELECT Users.Username, Products.ProductName FROM Products LEFT JOIN (SELECT * FROM Users WHERE Active = ?) AS Users ON (Products.UserID = Users.ID AND Users.Username = ?)", query)
}

func TestSubQueryExists(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewQuery("Users").Where("Company", "測試公司").Select("UserID")
	query, _ := Build(NewQuery("Products").WhereExists(subQuery).Select())
	assertEqual(assert, "SELECT * FROM Products WHERE EXISTS (SELECT UserID FROM Users WHERE Company = ?)", query)

	subQuery = NewQuery("Products").Where("Quantity", ">", 2).Select("UserID")
	query, _ = Build(NewQuery("Users").Where("ID", "IN", subQuery).Exists())
	assertEqual(assert, "SELECT EXISTS(SELECT * FROM Users WHERE ID IN (SELECT UserID FROM Products WHERE Quantity > ?))", query)
}

func TestSubQueryRawQuery(t *testing.T) {
	assert := assert.New(t)
	rawQuery := NewRawQuery("SELECT UserID FROM Users WHERE Company = ?", "測試公司")
	query, _ := Build(NewQuery("Products").WhereExists(rawQuery).Select())
	assertEqual(assert, "SELECT * FROM Products WHERE EXISTS (SELECT UserID FROM Users WHERE Company = ?)", query)
}

func TestSubQueryRawQueryReplacement(t *testing.T) {
	assert := assert.New(t)
	subQuery := NewQuery("Locations").Where("Username", "YamiOdymel").Select()
	rawQuery := NewRawQuery("SELECT UserID FROM Users WHERE EXISTS (?)", subQuery)
	query, _ := Build(NewQuery("Products").WhereExists(rawQuery).Select())
	assertEqual(assert, "SELECT * FROM Products WHERE EXISTS (SELECT UserID FROM Users WHERE EXISTS (SELECT * FROM Locations WHERE Username = ?))", query)
}

//=======================================================
// Others
//=======================================================

func TestComplexQueries(t *testing.T) {
	assert := assert.New(t)
	ids := NewQuery("JobHistories").WhereBetween("DepartmentID", 50, 100).Select("JobID")
	avgs := NewQuery("Jobs").WhereIn("JobID", ids).GroupBy("JobID").Select("JobID", "AVG(MinSalary) AS MyAVG")
	maxAvg := NewQuery(avgs).As("SS").Select("MAX(MyAVG)")
	query, _ := Build(NewQuery("Employees").GroupBy("JobID").Having("Avg(Salary)", "<", maxAvg).Select("JobID", "AVG(Salary)"))

	assertEqual(assert, "SELECT JobID, AVG(Salary) FROM Employees HAVING AVG(Salary) < (SELECT MAX(MyAVG) FROM (SELECT JobID, AVG(MinSalary) AS MyAVG FROM Jobs WHERE JobID IN (SELECT JobID FROM JobHistories WHERE DepartmentID BETWEEN ? AND ?) GROUP BY JobID) AS SS) GROUP BY JobID", query)

	// Example Source: https://www.w3resource.com/sql/subqueries/nested-subqueries.php
	// SELECT job_id,
	//        Avg(salary)
	// FROM   employees
	// HAVING Avg(salary) < (SELECT Max(myavg)
	//                       FROM   (SELECT job_id,
	//                                      Avg(min_salary) AS myavg
	//                               FROM   jobs
	//                               WHERE  job_id IN (SELECT job_id
	//                                                 FROM   job_history
	//                                                 WHERE  department_id BETWEEN 50
	//                                                        AND 100
	//                                                )
	//                               GROUP  BY job_id) ss)
	// GROUP  BY job_id;

	// Example Source: https://www.w3resource.com/sql/subqueries/nested-subqueries.php
	// SELECT ord_num,
	//        ord_date,
	//        ord_amount,
	//        advance_amount
	// FROM   orders
	// WHERE  ord_amount > 2000
	//        AND ord_date < '01-SEP-08'
	//        AND advance_amount < ANY (SELECT outstanding_amt
	//                                  FROM   customer
	//                                  WHERE  grade = 3
	//                                         AND cust_country <> 'India'
	//                                         AND opening_amt < 7000
	//                                         AND EXISTS (SELECT *
	//                                                     FROM   agents
	//                                                     WHERE  commission < .12));
}

func TestSetQueryOption(t *testing.T) {
	assert := assert.New(t)
	query, _ := Build(NewQuery("Users").SetQueryOption("FOR UPDATE").Select("Username"))
	assertEqual(assert, "SELECT Username FROM Users FOR UPDATE", query)
}
