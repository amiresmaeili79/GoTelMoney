package repositories

import (
	"github.com/amir79esmaeili/go-tel-money/internal/models"
	"gorm.io/gorm"
	"math"
)

type ExpenseRepositoryImpl struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) *ExpenseRepositoryImpl {
	return &ExpenseRepositoryImpl{
		db: db,
	}
}

func (r *ExpenseRepositoryImpl) All(userId uint, page, pageSize int) (int, []models.Expense, bool) {
	var expenses []models.Expense
	var totalRows int64
	hasNextPage := false
	if page != -1 {
		hasNextPage, totalRows = getPageStat(r.db, models.Expense{}, page, pageSize)
	}
	getPaginator(r.db, page, pageSize).Preload("ExpenseType").Find(&expenses, "user_id = ?", userId)

	totalPages := int(math.Ceil(float64(totalRows) / float64(pageSize)))
	return totalPages, expenses, hasNextPage
}

func (r *ExpenseRepositoryImpl) Create(e *models.Expense) error {
	return r.db.Create(e).Error
}
