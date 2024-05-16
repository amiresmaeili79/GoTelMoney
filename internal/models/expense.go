package models

import (
	"time"

	"gorm.io/gorm"
)

type ExpenseType struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint
	User   User `gorm:"foreignkey:UserID"`
	Name   string
}

func (t *ExpenseType) ToDisplay() string {
	return t.Name
}

type Expense struct {
	gorm.Model
	Date          time.Time
	Description   string
	Amount        float32
	ExpenseType   ExpenseType `gorm:"foreignkey:ExpenseTypeID"`
	ExpenseTypeID uint
	UserID        uint
	User          User `gorm:"foreignkey:UserID"`
}
