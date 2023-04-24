package database

import (
	"context"
	"database/sql"
	. "gorm.io/gorm"
	clause "gorm.io/gorm/clause"
)

type GormDB interface { 
    AddError(err error) error
    Assign(attrs ...interface{}) (tx *DB)
    Association(column string) *Association
    Attrs(attrs ...interface{}) (tx *DB)
    AutoMigrate(dst ...interface{}) error
    Begin(opts ...*sql.TxOptions) *DB
	//tldr we cant use callbacks now
    //Callback() *callbacks
    Clauses(conds ...clause.Expression) (tx *DB)
    Commit() *DB
    Connection(fc func(tx *DB) error) (err error)
    Count(count *int64) (tx *DB)
    Create(value interface{}) (tx *DB)
    CreateInBatches(value interface{}, batchSize int) (tx *DB)
    DB() (*sql.DB, error)
    Debug() (tx *DB)
    Delete(value interface{}, conds ...interface{}) (tx *DB)
    Distinct(args ...interface{}) (tx *DB)
    Exec(sql string, values ...interface{}) (tx *DB)
    Find(dest interface{}, conds ...interface{}) (tx *DB)
    FindInBatches(dest interface{}, batchSize int, fc func(tx *DB, batch int) error) *DB
    First(dest interface{}, conds ...interface{}) (tx *DB)
    FirstOrCreate(dest interface{}, conds ...interface{}) (tx *DB)
    FirstOrInit(dest interface{}, conds ...interface{}) (tx *DB)
    Get(key string) (interface{}, bool)
    Group(name string) (tx *DB)
    Having(query interface{}, args ...interface{}) (tx *DB)
    InnerJoins(query string, args ...interface{}) (tx *DB)
    InstanceGet(key string) (interface{}, bool)
    InstanceSet(key string, value interface{}) *DB
    Joins(query string, args ...interface{}) (tx *DB)
    Last(dest interface{}, conds ...interface{}) (tx *DB)
    Limit(limit int) (tx *DB)
    Migrator() Migrator
    Model(value interface{}) (tx *DB)
    Not(query interface{}, args ...interface{}) (tx *DB)
    Offset(offset int) (tx *DB)
    Omit(columns ...string) (tx *DB)
    Or(query interface{}, args ...interface{}) (tx *DB)
    Order(value interface{}) (tx *DB)
    Pluck(column string, dest interface{}) (tx *DB)
    Preload(query string, args ...interface{}) (tx *DB)
    Raw(sql string, values ...interface{}) (tx *DB)
    Rollback() *DB
    RollbackTo(name string) *DB
    Row() *sql.Row
    Rows() (*sql.Rows, error)
    Save(value interface{}) (tx *DB)
    SavePoint(name string) *DB
    Scan(dest interface{}) (tx *DB)
    ScanRows(rows *sql.Rows, dest interface{}) error
    Scopes(funcs ...func(*DB) *DB) (tx *DB)
    Select(query interface{}, args ...interface{}) (tx *DB)
    Session(config *Session) *DB
    Set(key string, value interface{}) *DB
    SetupJoinTable(model interface{}, field string, joinTable interface{}) error
    Table(name string, args ...interface{}) (tx *DB)
    Take(dest interface{}, conds ...interface{}) (tx *DB)
    ToSQL(queryFn func(tx *DB) *DB) string
    Transaction(fc func(tx *DB) error, opts ...*sql.TxOptions) (err error)
    Unscoped() (tx *DB)
    Update(column string, value interface{}) (tx *DB)
    UpdateColumn(column string, value interface{}) (tx *DB)
    UpdateColumns(values interface{}) (tx *DB)
    Updates(values interface{}) (tx *DB)
    Use(plugin Plugin) error
    Where(query interface{}, args ...interface{}) (tx *DB)
    WithContext(ctx context.Context) *DB
}
