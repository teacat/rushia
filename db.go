package rushia

import (
	"github.com/iancoleman/strcase"
)

type DB interface {
	New() DB
	Query(query string, params []interface{}, dest interface{}) error
	Exec(query string, params []interface{}) (err error)
	Transaction(handler func(tx *Query) error) error
	Begin() *Query
	Rollback() error
	RollbackTo(name string) error
	Commit() error
	SavePoint(name string) error
}

type Config struct {
	ColumnNamer  func(col string) string
	OrderColumns bool
}

func DefaultConfig() *Config {
	return &Config{
		ColumnNamer: func(col string) string {
			return strcase.ToSnake(col)
		},
		OrderColumns: false,
	}
}

// NewDB
func NewDB(db DB, conf *Config) *Query {
	return SetDB(NewQuery(nil).SetConfig(conf), db)
}

// SetDB
func SetDB(q *Query, db DB) *Query {
	q.db = db
	return q
}

// GetDB
func GetDB(q *Query) DB {
	return q.db
}

// SetConfig
func (q *Query) SetConfig(conf *Config) *Query {
	q.conf = conf
	return q
}

// NewQuery
func (q *Query) NewQuery(table interface{}) *Query {
	return SetDB(NewQuery(table), q.db.New())
}

// Query
func (q *Query) Query(dest interface{}) error {
	query, params := Build(q)
	return q.db.Query(query, params, dest)
}

// Exec
func (q *Query) Exec() error {
	query, params := Build(q)
	return q.db.Exec(query, params)
}

// Transaction
func (q *Query) Transaction(handler func(tx *Query) error) error {
	return q.db.Transaction(handler)
}

// Begin
func (q *Query) Begin() *Query {
	return q.db.Begin()
}

// Rollback
func (q *Query) Rollback() error {
	return q.db.Rollback()
}

// RollbackTo
func (q *Query) RollbackTo(name string) error {
	return q.db.RollbackTo(name)
}

// Commit
func (q *Query) Commit() error {
	return q.db.Commit()
}

// SavePoint
func (q *Query) SavePoint(name string) error {
	return q.db.SavePoint(name)
}
