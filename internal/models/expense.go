package models

import (
	"time"

	"gorm.io/gorm"
)

type ExpenseType struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type Expense struct {
	gorm.Model
	Date          time.Time
	Description   string
	Amount        float32
	ExpenseType   ExpenseType `gorm:"foreignkey:ExpenseTypeID"`
	ExpenseTypeID uint
}
