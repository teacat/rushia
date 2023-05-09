package rushiadb

import (
	"github.com/teacat/rushia/v3"
	"gorm.io/gorm"
)

// NewGorm returns a new Rushia DB instance by a existing Gorm Connection.
func NewGorm(db *gorm.DB) rushia.DB {
	return &gormDB{gorm: db}
}

// gormDB
type gormDB struct {
	gorm *gorm.DB
}

// New
func (e *gormDB) New() rushia.DB {
	return NewGorm(e.gorm)
}

// Query
func (e *gormDB) Query(query string, params []interface{}, dest interface{}) error {
	return e.gorm.Raw(query, params...).Scan(dest).Error
}

// Exec
func (e *gormDB) Exec(query string, params []interface{}) error {
	result := e.gorm.Exec(query, params...)
	return result.Error
}

// Transaction
func (e *gormDB) Transaction(handler func(tx *rushia.Query) error) error {
	return e.gorm.Transaction(func(tx *gorm.DB) error {
		q := rushia.SetDB(rushia.NewQuery(nil), NewGorm(tx))
		return handler(q)
	})
}

// Begin
func (e *gormDB) Begin() *rushia.Query {
	q := rushia.SetDB(rushia.NewQuery(nil), NewGorm(e.gorm.Begin()))
	return q
}

// Rollback
func (e *gormDB) Rollback() error {
	return e.gorm.Rollback().Error
}

// Commit
func (e *gormDB) Commit() error {
	return e.gorm.Commit().Error
}

// RollbackTo
func (e *gormDB) RollbackTo(name string) error {
	return e.gorm.RollbackTo(name).Error
}

// RollbackTo
func (e *gormDB) SavePoint(name string) error {
	return e.gorm.SavePoint(name).Error
}
