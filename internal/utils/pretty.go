package utils

import (
	"fmt"
	"unicode/utf8"

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

// TruncateString truncates a string to the specified number of characters
func TruncateString(str string, limit int) string {
	if limit <= 0 {
		return "" // return empty string if limit is not positive
	}

	if utf8.RuneCountInString(str) <= limit {
		return str // return the original string if it's shorter than the limit
	}

	// Truncate the string to the specified number of characters
	runes := []rune(str)
	return string(runes[:limit])
}
