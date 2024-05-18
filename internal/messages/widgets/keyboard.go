package widgets

import (
	"github.com/amir79esmaeili/go-tel-money/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math"
)

// GetDynamicInlineKeyboard Creates dynamic menu based on given data
func GetDynamicInlineKeyboard(items []models.InlineKeyboardItem, perRow int) *tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	totalChunks := int(math.Ceil(float64(len(items)) / float64(perRow)))

	for i := 0; i < totalChunks; i++ {
		start := i * perRow
		end := start + perRow

		if end > len(items) {
			end = len(items)
		}

		oneRow := tgbotapi.NewInlineKeyboardRow()
		for _, i := range items[start:end] {
			oneRow = append(oneRow, tgbotapi.NewInlineKeyboardButtonData(
				i.ToDisplay(), i.StringID()))
		}

		rows = append(rows, oneRow)
	}

	inlineKB := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &inlineKB
}
