package repositories

import "github.com/amir79esmaeili/go-tel-money/internal/models"

type ExpenseRepository interface {
	All(userId uint, page int, pageSize int) (int, []models.Expense, bool)
	Create(*models.Expense) error
}
