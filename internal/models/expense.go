package models

import (
	"fmt"
	"github.com/amir79esmaeili/go-tel-money/internal/utils"
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

func (t *ExpenseType) StringID() string {
	return fmt.Sprintf("%d", t.ID)
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

func (e *Expense) ToDisplay() string {
	return fmt.Sprintf("🗓️ %s / %s / 💵 %.2f / %s",
		e.Date.Format("Jan 2, 2006"),
		utils.TruncateString(e.Description, 50),
		e.Amount,
		e.ExpenseType.Name,
	)
}

func (e *Expense) StringID() string {
	return fmt.Sprintf("%d", e.ID)
}
