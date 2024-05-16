package commands

import (
	"fmt"
	"strings"
)

func WhatIsCommand(msg string) (ConvType, error) {
	for key, value := range Commands {
		if strings.Compare(key, msg) == 0 {
			return value, nil
		}
	}
	return -1, fmt.Errorf("unknown command %s", msg)
}

var Commands map[string]ConvType

func init() {
	Commands = map[string]ConvType{
		"🗂️ Add Expense Type": AddExpenseType,
		"💵 Add Expense":       AddExpense,
		"Report":              Report,
		"Cancel ❌":            Cancel,
	}
}
