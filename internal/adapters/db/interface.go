package db

import (
	"database/sql"
	"gorm.io/gorm"
)

type IDB interface {
	Model(value interface{}) (tx *gorm.DB)
	Select(query interface{}, args ...interface{}) (tx *gorm.DB)
	Where(query interface{}, args ...interface{}) (tx *gorm.DB)
	Order(value interface{}) (tx *gorm.DB)
	Limit(limit int) (tx *gorm.DB)
	Create(value interface{}) (tx *gorm.DB)
	CreateInBatches(value interface{}, batchSize int) (tx *gorm.DB)
	Save(value interface{}) (tx *gorm.DB)
	First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Take(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Last(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Find(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Update(column string, value interface{}) (tx *gorm.DB)
	Delete(value interface{}, conds ...interface{}) (tx *gorm.DB)
	Count(count *int64) (tx *gorm.DB)
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	Scan(dest interface{}) (tx *gorm.DB)
	Exec(sql string, values ...interface{}) (tx *gorm.DB)
	Raw(sql string, values ...interface{}) (tx *gorm.DB)
	ScanRows(rows *sql.Rows, dest interface{}) error
}
