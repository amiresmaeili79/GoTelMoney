package repositories

import (
	"github.com/amir79esmaeili/go-tel-money/internal/models"
	"gorm.io/gorm"
)

type ExpenseTypeRepositoryImpl struct {
	db *gorm.DB
}

func NewExpenseTypeRepository(db *gorm.DB) *ExpenseTypeRepositoryImpl {
	return &ExpenseTypeRepositoryImpl{db}
}

func (repo *ExpenseTypeRepositoryImpl) Create(exType *models.ExpenseType) error {
	return repo.db.Create(exType).Error
}

func (repo *ExpenseTypeRepositoryImpl) All(userId uint, page, pageSize int) []models.ExpenseType {
	var types []models.ExpenseType
	getPaginator(repo.db, page, pageSize).Find(&types, "user_id = ?", userId)
	return types
}

func (repo *ExpenseTypeRepositoryImpl) GetByName(name string, userID uint) (models.ExpenseType, error) {
	var exType models.ExpenseType
	err := repo.db.First(&exType, "user_id = ? AND name = ?", userID, name).Error
	return exType, err
}

func (repo *ExpenseTypeRepositoryImpl) GetByID(id uint, userID uint) (models.ExpenseType, error) {
	var exType models.ExpenseType
	err := repo.db.First(&exType, "id = ? AND user_id = ?", id, userID).Error
	return exType, err
}
