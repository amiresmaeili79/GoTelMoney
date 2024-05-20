package models

import (
	"fmt"
	"github.com/amir79esmaeili/go-tel-money/internal/messages"
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
	var date string
	if e.Date.Year() == time.Now().Year() {
		date = e.Date.Format("Jan 2")
	} else {
		date = e.Date.Format("Jan 2, 2006")
	}

	return fmt.Sprintf(messages.ExpenseRow,
		e.ID,
		date,
		e.Amount,
		e.ExpenseType.Name,
		utils.TruncateString(e.Description, 50),
	)
}

func (e *Expense) StringID() string {
	return fmt.Sprintf("%d", e.ID)
}

func PrettyPrintExpenseTypes(types []ExpenseType) string {
	fullMsg := messages.Types

	for _, t := range types {
		fullMsg += fmt.Sprintf(messages.TypeRow, t.Name)
	}

	return fullMsg
}

func PrettyPrintExpenses(expenses []Expense, page, pages int) string {
	fullMsg := fmt.Sprintf(messages.ReportHead, page, pages)

	var total float32
	for _, e := range expenses {
		fullMsg += e.ToDisplay()
		total += e.Amount
	}

	fullMsg += fmt.Sprintf(messages.ExpenseTotal, total)

	return fullMsg
}
