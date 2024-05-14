package conversations

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	WelcomeMsg = `Hello dear %s
We are happy to see you here 🎉.`
	Invalid        = "I don't know man!"
	NewDescription = "✏️ The description of your expense:"
	NewAmount      = "✏️ The amount of your expense:"
)

var MainMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("💵 Add Expense"),
		tgbotapi.NewKeyboardButton("🔎 View Report"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Add Expense Type"),
	),
)

const (
	ExpenseConv = iota
)
