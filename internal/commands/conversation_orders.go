package commands

type ConvType int
type AddExpenseTypeOrder int
type AddExpenseOrder int

const (
	AddExpenseType ConvType = iota
	AddExpense
	Report
	Cancel
)

const (
	StartAddExpenseType AddExpenseTypeOrder = iota
	AskNameAddExpenseType
)

const (
	StartAddExpense AddExpenseOrder = iota
	AskAmountAddExpense
	AskDescriptionAddExpense
	AskDateAddExpense
	AskTypeAddExpense
)
