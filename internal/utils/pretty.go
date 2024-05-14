package utils

import (
	"fmt"

	"github.com/amir79esmaeili/go-tel-money/internal/conversations"
	"github.com/amir79esmaeili/go-tel-money/internal/models"
)

func PrettyPrintExpenseTypes(types []models.ExpenseType) string {
	fullMsg := conversations.Types

	for _, t := range types {
		fullMsg += fmt.Sprintf(conversations.TypeRow, t.Name)
	}

	return fullMsg
}
