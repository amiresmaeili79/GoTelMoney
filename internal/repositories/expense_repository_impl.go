package repositories

import (
	"github.com/amir79esmaeili/go-tel-money/internal/models"
	"gorm.io/gorm"
)

type ExpenseRepositoryImpl struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) *ExpenseRepositoryImpl {
	return &ExpenseRepositoryImpl{
		db: db,
	}
}

func (r *ExpenseRepositoryImpl) All(userId uint) []models.Expense {
	var expenses []models.Expense
	r.db.Find(&expenses, "user_id = ?", userId)
	return expenses
}

func (r *ExpenseRepositoryImpl) Create(e *models.Expense) error {
	return r.db.Create(e).Error
}
