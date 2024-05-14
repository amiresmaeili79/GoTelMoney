package repositories

import "github.com/amir79esmaeili/go-tel-money/internal/models"

type ExpenseTypeRepository interface {
	All(userId uint) []models.ExpenseType
	Create(*models.ExpenseType) error
	FindByName(name string, userID uint) (models.ExpenseType, error)
}
