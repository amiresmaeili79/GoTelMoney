package repositories

import "github.com/amir79esmaeili/go-tel-money/internal/models"

type ExpenseRepository interface {
	All(page int, pageSize int) (int, []models.Expense, bool)
	Filter(page, pageSize int, args ...any) (int, []models.Expense, bool)
	Create(*models.Expense) error
}
