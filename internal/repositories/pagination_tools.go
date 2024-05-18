package repositories

import "gorm.io/gorm"

// getPaginator creates a db connection that handles offset and limit
func getPaginator(db *gorm.DB, page, pageSize int) *gorm.DB {
	if page == -1 { // Don't need pagination. Just all!
		return db
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	return db.Offset(offset).Limit(pageSize)
}
