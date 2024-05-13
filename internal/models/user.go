package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserTelegramID int64
	FirstName      string
	LastName       string
}
