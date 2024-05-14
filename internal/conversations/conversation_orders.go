package conversations

type ConvType int
type AddExpenseTypeOrder int

const (
	AddExpenseType ConvType = iota
	AddExpense
	Report
)

var Commands map[string]ConvType

const (
	StartAddExpenseType AddExpenseTypeOrder = iota
	AskNameAddExpenseType
	SubmitAddExpenseType
)

func init() {
	Commands = map[string]ConvType{
		"Add Expense Type": AddExpenseType,
		"Add Expense":      AddExpense,
		"Report":           Report,
	}
}
