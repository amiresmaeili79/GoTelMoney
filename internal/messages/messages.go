package messages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	WelcomeMsg = `Hello dear %s
We are happy to see you here 🎉.`
	NewDescription = "✏️ The description of your expense?"
	NewAmount      = "✏️ The amount of your expense?"
	NewType        = "Name of the new type?"
	NewDate        = "The date of your expense? Select the date or just simply type it (e.g. 2024-12-19)."

	Types                 = "🗂️ Your current types 👇\n"
	TypeRow               = "◽ %s\n"
	TypeAddedSuccessfully = "New type '%s' added!"
	CancelMessage         = "Okay! Back to the main menu."
	NewExpenseSaved       = "Done! You're good to go."
	AskType               = "Please select the type of expense!"
)

// Error Messages
const (
	UserCreationFailed = "Failed to create a user. Please try again later."
	Invalid            = "I don't know man!"
	UserNotFound       = "You are not registered in the bot. please type '/start'. Then try again!"
	InvalidAmount      = "This is not a valid number! Try again. (e.g. 23.5)"
	TypeAlreadyAdded   = "Failed to add %s! It is already added."
	TypeAddingFailed   = "Failed to add %s! Try again later please."
	InvalidDate        = "Invalid date! Try again (e.g. 2024-05-25 14:35)"
	InvalidType        = "Invalid type! Please select one from menu!"
	ExpenseSaveFailed  = "Failed to save the expense!"
)

var MainMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("💵 Add Expense"),
		tgbotapi.NewKeyboardButton("🔎 View Report"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("🗂️ Add Expense Type"),
	),
)

var CancelKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Cancel ❌"),
	),
)

func init() {
}
