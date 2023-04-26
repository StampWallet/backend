package database

import (
	"context"
	"database/sql"
	. "gorm.io/gorm"
	clause "gorm.io/gorm/clause"
)

type GormDB interface {
	AddError(err error) error
	Assign(attrs ...interface{}) (tx GormDB)
	Association(column string) *Association
	Attrs(attrs ...interface{}) (tx GormDB)
	AutoMigrate(dst ...interface{}) error
	Begin(opts ...*sql.TxOptions) GormDB
	//tldr we cant use callbacks now
	//Callback() *callbacks
	Clauses(conds ...clause.Expression) (tx GormDB)
	Commit() GormDB
	Connection(fc func(tx GormDB) error) (err error)
	Count(count *int64) (tx GormDB)
	Create(value interface{}) (tx GormDB)
	CreateInBatches(value interface{}, batchSize int) (tx GormDB)
	DB() (*sql.DB, error)
	Debug() (tx GormDB)
	Delete(value interface{}, conds ...interface{}) (tx GormDB)
	Distinct(args ...interface{}) (tx GormDB)
	Exec(sql string, values ...interface{}) (tx GormDB)
	Find(dest interface{}, conds ...interface{}) (tx GormDB)
	FindInBatches(dest interface{}, batchSize int, fc func(tx GormDB, batch int) error) GormDB
	First(dest interface{}, conds ...interface{}) (tx GormDB)
	FirstOrCreate(dest interface{}, conds ...interface{}) (tx GormDB)
	FirstOrInit(dest interface{}, conds ...interface{}) (tx GormDB)
	Get(key string) (interface{}, bool)
	Group(name string) (tx GormDB)
	Having(query interface{}, args ...interface{}) (tx GormDB)
	InnerJoins(query string, args ...interface{}) (tx GormDB)
	InstanceGet(key string) (interface{}, bool)
	InstanceSet(key string, value interface{}) GormDB
	Joins(query string, args ...interface{}) (tx GormDB)
	Last(dest interface{}, conds ...interface{}) (tx GormDB)
	Limit(limit int) (tx GormDB)
	Migrator() Migrator
	Model(value interface{}) (tx GormDB)
	Not(query interface{}, args ...interface{}) (tx GormDB)
	Offset(offset int) (tx GormDB)
	Omit(columns ...string) (tx GormDB)
	Or(query interface{}, args ...interface{}) (tx GormDB)
	Order(value interface{}) (tx GormDB)
	Pluck(column string, dest interface{}) (tx GormDB)
	Preload(query string, args ...interface{}) (tx GormDB)
	Raw(sql string, values ...interface{}) (tx GormDB)
	Rollback() GormDB
	RollbackTo(name string) GormDB
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	Save(value interface{}) (tx GormDB)
	SavePoint(name string) GormDB
	Scan(dest interface{}) (tx GormDB)
	ScanRows(rows *sql.Rows, dest interface{}) error
	Scopes(funcs ...func(GormDB) GormDB) (tx GormDB)
	Select(query interface{}, args ...interface{}) (tx GormDB)
	Session(config *Session) GormDB
	Set(key string, value interface{}) GormDB
	SetupJoinTable(model interface{}, field string, joinTable interface{}) error
	Table(name string, args ...interface{}) (tx GormDB)
	Take(dest interface{}, conds ...interface{}) (tx GormDB)
	ToSQL(queryFn func(tx GormDB) GormDB) string
	Transaction(fc func(tx GormDB) error, opts ...*sql.TxOptions) (err error)
	Unscoped() (tx GormDB)
	Update(column string, value interface{}) (tx GormDB)
	UpdateColumn(column string, value interface{}) (tx GormDB)
	UpdateColumns(values interface{}) (tx GormDB)
	Updates(values interface{}) (tx GormDB)
	Use(plugin Plugin) error
	Where(query interface{}, args ...interface{}) (tx GormDB)
	WithContext(ctx context.Context) GormDB

	GetConfig() *Config
	GetError() error
	GetRowsAffected() int64
	GetStatement() *Statement
}

type GormDBImpl struct {
	Db *DB
}

//a little copying is better than a little dependency or whatever the fuck
//maybe writing a generator here would be a better idea. didn't have time.
//i couldn't find anything that would generate a proxy just like this, perhaps i'm searching wrong
//i probably should've just given up at this point and set up the db for unit tests

func (self *GormDBImpl) AddError(err error) error {
	return self.Db.AddError(err)
}

func (self *GormDBImpl) Association(column string) *Association {
	return self.Db.Association(column)
}

func (self *GormDBImpl) AutoMigrate(dst ...interface{}) error {
	return self.Db.AutoMigrate(dst...)
}

func (self *GormDBImpl) Connection(fc func(tx GormDB) error) (err error) {
	return self.Db.Connection(func(tx *DB) error {
		return fc(&GormDBImpl{tx})
	})
}

func (self *GormDBImpl) DB() (*sql.DB, error) {
	return self.Db.DB()
}

func (self *GormDBImpl) Get(key string) (interface{}, bool) {
	return self.Db.Get(key)
}

func (self *GormDBImpl) InstanceGet(key string) (interface{}, bool) {
	return self.Db.InstanceGet(key)
}

func (self *GormDBImpl) Migrator() Migrator {
	return self.Db.Migrator()
}

func (self *GormDBImpl) Row() *sql.Row {
	return self.Db.Row()
}

func (self *GormDBImpl) Rows() (*sql.Rows, error) {
	return self.Db.Rows()
}

func (self *GormDBImpl) ScanRows(rows *sql.Rows, dest interface{}) error {
	return self.Db.ScanRows(rows, dest)
}

func (self *GormDBImpl) SetupJoinTable(model interface{}, field string, joinTable interface{}) error {
	return self.Db.SetupJoinTable(model, field, joinTable)
}

func (self *GormDBImpl) ToSQL(queryFn func(tx GormDB) GormDB) string {
	return self.Db.ToSQL(func(tx *DB) *DB {
		return queryFn(&GormDBImpl{tx}).(*GormDBImpl).Db
	})
}

func (self *GormDBImpl) Transaction(fc func(tx GormDB) error, opts ...*sql.TxOptions) (err error) {
	return self.Db.Transaction(func(tx *DB) error {
		return fc(&GormDBImpl{tx})
	}, opts...)
}

func (self *GormDBImpl) Use(plugin Plugin) error {
	return self.Db.Use(plugin)
}

func (self *GormDBImpl) Assign(attrs ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Assign(attrs...)}
}

func (self *GormDBImpl) Attrs(attrs ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Attrs(attrs...)}
}

func (self *GormDBImpl) Begin(opts ...*sql.TxOptions) GormDB {
	return &GormDBImpl{self.Db.Begin(opts...)}
}

func (self *GormDBImpl) Clauses(conds ...clause.Expression) (tx GormDB) {
	return &GormDBImpl{self.Db.Clauses(conds...)}
}

func (self *GormDBImpl) Commit() GormDB {
	return &GormDBImpl{self.Db.Commit()}
}

func (self *GormDBImpl) Count(count *int64) (tx GormDB) {
	return &GormDBImpl{self.Db.Count(count)}
}

func (self *GormDBImpl) Create(value interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Create(value)}
}

func (self *GormDBImpl) CreateInBatches(value interface{}, batchSize int) (tx GormDB) {
	return &GormDBImpl{self.Db.CreateInBatches(value, batchSize)}
}

func (self *GormDBImpl) Debug() (tx GormDB) {
	return &GormDBImpl{self.Db.Debug()}
}

func (self *GormDBImpl) Delete(value interface{}, conds ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Delete(value, conds...)}
}

func (self *GormDBImpl) Distinct(args ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Distinct(args...)}
}

func (self *GormDBImpl) Exec(sql string, values ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Exec(sql, values...)}
}

func (self *GormDBImpl) Find(dest interface{}, conds ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Find(dest, conds...)}
}

func (self *GormDBImpl) FindInBatches(dest interface{}, batchSize int, fc func(tx GormDB, batch int) error) GormDB {
	return &GormDBImpl{self.Db.FindInBatches(dest, batchSize, func(tx *DB, batch int) error {
		return fc(&GormDBImpl{tx}, batch)
	})}
}

func (self *GormDBImpl) First(dest interface{}, conds ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.First(dest, conds...)}
}

func (self *GormDBImpl) FirstOrCreate(dest interface{}, conds ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.FirstOrCreate(dest, conds...)}
}

func (self *GormDBImpl) FirstOrInit(dest interface{}, conds ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.FirstOrInit(dest, conds...)}
}

func (self *GormDBImpl) Group(name string) (tx GormDB) {
	return &GormDBImpl{self.Db.Group(name)}
}

func (self *GormDBImpl) Having(query interface{}, args ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Having(query, args...)}
}

func (self *GormDBImpl) InnerJoins(query string, args ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.InnerJoins(query, args...)}
}

func (self *GormDBImpl) InstanceSet(key string, value interface{}) GormDB {
	return &GormDBImpl{self.Db.InstanceSet(key, value)}
}

func (self *GormDBImpl) Joins(query string, args ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Joins(query, args)}
}

func (self *GormDBImpl) Last(dest interface{}, conds ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Last(dest, conds...)}
}

func (self *GormDBImpl) Limit(limit int) (tx GormDB) {
	return &GormDBImpl{self.Db.Limit(limit)}
}

func (self *GormDBImpl) Model(value interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Model(value)}
}

func (self *GormDBImpl) Not(query interface{}, args ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Not(query, args...)}
}

func (self *GormDBImpl) Offset(offset int) (tx GormDB) {
	return &GormDBImpl{self.Db.Offset(offset)}
}

func (self *GormDBImpl) Omit(columns ...string) (tx GormDB) {
	return &GormDBImpl{self.Db.Omit(columns...)}
}

func (self *GormDBImpl) Or(query interface{}, args ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Or(query, args...)}
}

func (self *GormDBImpl) Order(value interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Order(value)}
}

func (self *GormDBImpl) Pluck(column string, dest interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Pluck(column, dest)}
}

func (self *GormDBImpl) Preload(query string, args ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Preload(query, args...)}
}

func (self *GormDBImpl) Raw(sql string, values ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Raw(sql, values...)}
}

func (self *GormDBImpl) Rollback() GormDB {
	return &GormDBImpl{self.Db.Rollback()}
}

func (self *GormDBImpl) RollbackTo(name string) GormDB {
	return &GormDBImpl{self.Db.RollbackTo(name)}
}

func (self *GormDBImpl) Save(value interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Save(value)}
}

func (self *GormDBImpl) SavePoint(name string) GormDB {
	return &GormDBImpl{self.Db.SavePoint(name)}
}

func (self *GormDBImpl) Scan(dest interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Scan(dest)}
}

func (self *GormDBImpl) Scopes(funcs ...func(GormDB) GormDB) (tx GormDB) {
	var args []func(db *DB) *DB
	for _, f := range funcs {
		args = append(args, func(db *DB) *DB {
			return f(&GormDBImpl{db}).(*GormDBImpl).Db
		})
	}
	return &GormDBImpl{self.Db.Scopes(args...)}
}

func (self *GormDBImpl) Select(query interface{}, args ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Select(query, args...)}
}

func (self *GormDBImpl) Session(config *Session) GormDB {
	return &GormDBImpl{self.Db.Session(config)}
}

func (self *GormDBImpl) Set(key string, value interface{}) GormDB {
	return &GormDBImpl{self.Db.Set(key, value)}
}

func (self *GormDBImpl) Table(name string, args ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Table(name, args...)}
}

func (self *GormDBImpl) Take(dest interface{}, conds ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Take(dest, conds...)}
}

func (self *GormDBImpl) Unscoped() (tx GormDB) {
	return &GormDBImpl{self.Db.Unscoped()}
}

func (self *GormDBImpl) Update(column string, value interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Update(column, value)}
}

func (self *GormDBImpl) UpdateColumn(column string, value interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.UpdateColumn(column, value)}
}

func (self *GormDBImpl) UpdateColumns(values interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.UpdateColumns(values)}
}

func (self *GormDBImpl) Updates(values interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Updates(values)}
}

func (self *GormDBImpl) Where(query interface{}, args ...interface{}) (tx GormDB) {
	return &GormDBImpl{self.Db.Where(query, args...)}
}

func (self *GormDBImpl) WithContext(ctx context.Context) GormDB {
	return &GormDBImpl{self.Db.WithContext(ctx)}
}

func (self *GormDBImpl) GetConfig() *Config {
	return self.Db.Config
}

func (self *GormDBImpl) GetError() error {
	return self.Db.Error
}

func (self *GormDBImpl) GetRowsAffected() int64 {
	return self.Db.RowsAffected
}

func (self *GormDBImpl) GetStatement() *Statement {
	return self.Db.Statement
}
