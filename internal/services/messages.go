package services

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	WelcomeMsg = `Hello dear %s
We are happy to see you here 🎉.`
)

var MainMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("💵 Add Expense"),
		tgbotapi.NewKeyboardButton("🔎 View Report"),
	),
)
