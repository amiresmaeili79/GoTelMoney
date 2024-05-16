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
