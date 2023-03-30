package repository

import (
	"fmt"

	"gorm.io/gorm"
)

func WithPagination(page int, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return db.Limit(limit).Offset(offset)
	}
}

func WithSearch(value string, columns []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if value == "" {
			return db
		}
		for _, column := range columns {
			db.Where(fmt.Sprintf("%s LIKE ?", column), fmt.Sprintf("%%%s%%", value))
		}
		return db
	}
}
