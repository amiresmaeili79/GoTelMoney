package repositories

import (
	"github.com/amir79esmaeili/go-tel-money/internal/models"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db}
}

func (repo *UserRepositoryImpl) Create(user *models.User) error {
	return repo.db.Create(user).Error
}

func (repo *UserRepositoryImpl) GetById(id uint) (*models.User, error) {
	var user models.User
	if err := repo.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepositoryImpl) GetByTID(tid int64) (*models.User, error) {
	var user models.User
	if err := repo.db.Where("user_telegram_id = ?", tid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
