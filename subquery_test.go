package rushia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var subQuery SubQuery

func TestSubQueryGetx(t *testing.T) {
	assert := assert.New(t)
	subQuery = NewSubQuery().Table("Users").Get()
	assertEqual(assert, "SELECT * FROM Users", subQuery.query.query)

	subQuery = NewSubQuery().Table("Users").Get("Username", "Password")
	assertEqual(assert, "SELECT Username, Password FROM Users", subQuery.query.query)
}

func TestSubQueryWhere(t *testing.T) {
	assert := assert.New(t)
	subQuery = NewSubQuery().Table("Users").Where("ID", 1).Where("Username", "admin").Get()
	assertEqual(assert, "SELECT * FROM Users WHERE ID = ? AND Username = ?", subQuery.query.query)
}

func TestSubQueryOrWhere(t *testing.T) {
	assert := assert.New(t)
	subQuery = NewSubQuery().Table("Users").Where("FirstName", "John").OrWhere("FirstName", "Peter").Get()
	assertEqual(assert, "SELECT * FROM Users WHERE FirstName = ? OR FirstName = ?", subQuery.query.query)
	subQuery = NewSubQuery().Table("Users").Where("A = B").OrWhere("(A = C OR A = D)").Get()
	assertEqual(assert, "SELECT * FROM Users WHERE A = B OR (A = C OR A = D)", subQuery.query.query)
}

func TestSubQueryWhereHaving(t *testing.T) {
	assert := assert.New(t)
	subQuery = NewSubQuery().Table("Users").Where("ID", 1).Having("Username", "admin").Get()
	assertEqual(assert, "SELECT * FROM Users WHERE ID = ? HAVING Username = ?", subQuery.query.query)
	subQuery = NewSubQuery().Table("Users").Where("ID", 1).Having("Username", "admin").OrHaving("Password", "test").Get()
	assertEqual(assert, "SELECT * FROM Users WHERE ID = ? HAVING Username = ? OR Password = ?", subQuery.query.query)
}

func TestSubQueryLimit(t *testing.T) {
	assert := assert.New(t)
	subQuery = NewSubQuery().Table("Users").Limit(10).Get()
	assertEqual(assert, "SELECT * FROM Users LIMIT 10", subQuery.query.query)
}

func TestSubQueryOrderBy(t *testing.T) {
	assert := assert.New(t)
	subQuery = NewSubQuery().Table("Users").OrderBy("ID", "ASC").OrderBy("Login", "DESC").OrderBy("RAND()").Get()
	assertEqual(assert, "SELECT * FROM Users ORDER BY ID ASC, Login DESC, RAND()", subQuery.query.query)
}

func TestSubQueryOrderByField(t *testing.T) {
	assert := assert.New(t)
	subQuery = NewSubQuery().Table("Users").OrderBy("UserGroup", "ASC", "SuperUser", "Admin", "Users").Get()
	assertEqual(assert, "SELECT * FROM Users ORDER BY FIELD (UserGroup, ?, ?, ?) ASC", subQuery.query.query)
}

func TestSubQueryGroupBy(t *testing.T) {
	assert := assert.New(t)
	subQuery = NewSubQuery().Table("Users").GroupBy("Name").Get()
	assertEqual(assert, "SELECT * FROM Users GROUP BY Name", subQuery.query.query)
	subQuery = NewSubQuery().Table("Users").GroupBy("Name", "ID").Get()
	assertEqual(assert, "SELECT * FROM Users GROUP BY Name, ID", subQuery.query.query)
}

func TestSubQueryJoinx(t *testing.T) {
	assert := assert.New(t)
	subQuery = NewSubQuery().
		Table("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		Where("Users.ID", 6).
		Get("Users.Name", "Products.ProductName")
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products LEFT JOIN Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?", subQuery.query.query)

	query, _ := builder.
		Table("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		Where("Users.ID", 6).
		Get("Users.Name", "Products.ProductName")
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products LEFT JOIN Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?", query)

	subQuery = NewSubQuery().
		Table("Products").
		RightJoin("Users", "Products.TenantID = Users.TenantID").
		Where("Users.ID", 6).
		Get("Users.Name", "Products.ProductName")
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products RIGHT JOIN Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?", subQuery.query.query)

	subQuery = NewSubQuery().
		Table("Products").
		InnerJoin("Users", "Products.TenantID = Users.TenantID").
		Where("Users.ID", 6).
		Get("Users.Name", "Products.ProductName")
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products INNER JOIN Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?", subQuery.query.query)

	subQuery = NewSubQuery().
		Table("Products").
		NaturalJoin("Users", "Products.TenantID = Users.TenantID").
		Where("Users.ID", 6).
		Get("Users.Name", "Products.ProductName")
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products NATURAL JOIN Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?", subQuery.query.query)

	subQuery = NewSubQuery().
		Table("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		RightJoin("Posts", "Products.TenantID = Posts.TenantID").
		Where("Users.ID", 6).
		Get("Users.Name", "Products.ProductName")
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products RIGHT JOIN Posts ON (Products.TenantID = Posts.TenantID) LEFT JOIN Users ON (Products.TenantID = Users.TenantID) WHERE Users.ID = ?", subQuery.query.query)
}

func TestSubQueryJoinWhere(t *testing.T) {
	assert := assert.New(t)
	subQuery = NewSubQuery().
		Table("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		JoinOrWhere("Users", "Users.TenantID", 5).
		Get("Users.Name", "Products.ProductName")
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products LEFT JOIN Users ON (Products.TenantID = Users.TenantID OR Users.TenantID = ?)", subQuery.query.query)
	subQuery = NewSubQuery().
		Table("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		JoinWhere("Users", "Users.Username", "Wow").
		Get("Users.Name", "Products.ProductName")
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products LEFT JOIN Users ON (Products.TenantID = Users.TenantID AND Users.Username = ?)", subQuery.query.query)
	subQuery = NewSubQuery().
		Table("Products").
		LeftJoin("Users", "Products.TenantID = Users.TenantID").
		RightJoin("Posts", "Products.TenantID = Posts.TenantID").
		JoinWhere("Posts", "Posts.Username", "Wow").
		JoinWhere("Users", "Users.Username", "Wow").
		Get("Users.Name", "Products.ProductName")
	assertEqual(assert, "SELECT Users.Name, Products.ProductName FROM Products LEFT JOIN Users ON (Products.TenantID = Users.TenantID AND Users.Username = ?) RIGHT JOIN Posts ON (Products.TenantID = Posts.TenantID AND Posts.Username = ?)", subQuery.query.query)
}
