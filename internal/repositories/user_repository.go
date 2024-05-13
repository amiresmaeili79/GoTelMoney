package repositories

import "github.com/amir79esmaeili/go-tel-money/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	GetById(id uint) (*models.User, error)
	GetByTID(tid int64) (*models.User, error)
}
