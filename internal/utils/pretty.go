package utils

import (
	"fmt"

	"github.com/amir79esmaeili/go-tel-money/internal/messages"
	"github.com/amir79esmaeili/go-tel-money/internal/models"
)

func PrettyPrintExpenseTypes(types []models.ExpenseType) string {
	fullMsg := messages.Types

	for _, t := range types {
		fullMsg += fmt.Sprintf(messages.TypeRow, t.Name)
	}

	return fullMsg
}

func PrettyPrintExpenses(expenses []models.Expense) string {
	fullMsg := messages.ReportHead

	for _, e := range expenses {
		fullMsg += fmt.Sprintf(messages.ReportRow,
			e.Amount,
			e.Description,
			e.Date.Format("Jan 2, 2006"),
			e.ExpenseType.Name,
		)
	}

	return fullMsg
}
