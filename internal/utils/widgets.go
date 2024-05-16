package utils

import (
	"github.com/amir79esmaeili/go-tel-money/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetDynamicInlineKeyboard Creates dynamic menu based on given data
func GetDynamicInlineKeyboard(items []models.Choosable) *tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for i := 0; i < len(items); i += 2 {
		if i+1 < len(items) {
			pair := []tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData(items[i].ToDisplay(), items[i].ToDisplay()),
				tgbotapi.NewInlineKeyboardButtonData(items[i+1].ToDisplay(), items[i+1].ToDisplay()),
			}
			rows = append(rows, pair)
		} else {
			row := []tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData(items[i].ToDisplay(), items[i].ToDisplay()),
			}
			rows = append(rows, row)
		}
	}

	inlineKB := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &inlineKB
}
