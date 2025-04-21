package utils

import (
	"strconv"

	"gorm.io/gorm"
)

// Paginate handles pagination
func Paginate(page, pageSize string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageNum, _ := strconv.Atoi(page)
		if pageNum == 0 {
			pageNum = 1
		}

		limit, _ := strconv.Atoi(pageSize)
		if limit == 0 {
			limit = 10
		}

		offset := (pageNum - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
} 