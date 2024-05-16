package repositories

import "github.com/amir79esmaeili/go-tel-money/internal/models"

type ExpenseRepository interface {
	All(userId uint) []models.Expense
	Create(*models.Expense) error
}
