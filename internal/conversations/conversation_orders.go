package conversations

type ConvType int64
type AddExpenseTypeOrder int64

const (
	AddExpenseType ConvType = iota
	AddExpense
	Report
)

var Commands map[string]ConvType

const (
	Start AddExpenseTypeOrder = iota
	AskName
	Submit
)

func init() {
	Commands = map[string]ConvType{
		"Add Expense Type": AddExpenseType,
		"Add Expense":      AddExpense,
		"Report":           Report,
	}
}
