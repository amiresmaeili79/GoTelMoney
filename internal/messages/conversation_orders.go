package messages

type ConvType int
type AddExpenseTypeOrder int
type AddExpenseOrder int

const (
	AddExpenseType ConvType = iota
	AddExpense
	Report
)

const (
	StartAddExpenseType AddExpenseTypeOrder = iota
	AskNameAddExpenseType
	SubmitAddExpenseType
)

const (
	StartAddExpense AddExpenseOrder = iota
	AskAmountAddExpense
	AskDescriptionAddExpense
	AskDateAddExpense
	AskTypeAddExpense
	SubmitAddExpense
)
