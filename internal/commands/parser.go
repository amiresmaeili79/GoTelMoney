package commands

import (
	"fmt"
	"strings"

	"github.com/amir79esmaeili/go-tel-money/internal/messages"
)

func WhatIsCommand(msg string) (messages.ConvType, error) {
	for key, value := range Commands {
		if strings.Compare(key, msg) == 0 {
			return value, nil
		}
	}
	return -1, fmt.Errorf("unknown command %s", msg)
}

var Commands map[string]messages.ConvType

func init() {
	Commands = map[string]messages.ConvType{
		"Add Expense Type": messages.AddExpenseType,
		"💵 Add Expense":    messages.AddExpense,
		"Report":           messages.Report,
	}
}
