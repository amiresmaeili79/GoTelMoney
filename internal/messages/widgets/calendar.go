package widgets

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"time"
)

const BtnPrev = "<"
const BtnNext = ">"

func HandleButton(data string) tgbotapi.InlineKeyboardMarkup {
	config := strings.Split(data, " ")
	year, _ := strconv.Atoi(config[1])
	month, _ := strconv.Atoi(config[2])
	if strings.Compare(config[0], BtnNext) == 0 {
		return handlerNextButton(year, time.Month(month))
	} else {
		return handlerPrevButton(year, time.Month(month))
	}
}

func GenerateCalendar(year int, month time.Month) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	keyboard = addShortcuts(keyboard)
	keyboard = addMonthYearRow(year, month, keyboard)
	keyboard = addDaysNamesRow(keyboard)
	keyboard = generateMonth(year, int(month), keyboard)
	keyboard = addSpecialButtons(keyboard, year, int(month))
	return keyboard
}

func addShortcuts(keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.InlineKeyboardMarkup {
	var row []tgbotapi.InlineKeyboardButton
	now := time.Now()
	today := tgbotapi.NewInlineKeyboardButtonData("Today", now.Format("2006-01-02"))
	yesterdayDate := now.AddDate(0, 0, -1)
	yesterday := tgbotapi.NewInlineKeyboardButtonData("Yesterday", yesterdayDate.Format("2006-01-02"))
	row = append(row, today, yesterday)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	return keyboard
}

func handlerPrevButton(year int, month time.Month) tgbotapi.InlineKeyboardMarkup {
	if month != 1 {
		month--
	} else {
		month = 12
		year--
	}
	return GenerateCalendar(year, month)
}

func handlerNextButton(year int, month time.Month) tgbotapi.InlineKeyboardMarkup {
	if month != 12 {
		month++
	} else {
		year++
	}
	return GenerateCalendar(year, month)
}

func addMonthYearRow(year int, month time.Month, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.InlineKeyboardMarkup {
	var row []tgbotapi.InlineKeyboardButton
	btn := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s %v", month, year), "1")
	row = append(row, btn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	return keyboard
}

func addDaysNamesRow(keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.InlineKeyboardMarkup {
	days := [7]string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"}
	var rowDays []tgbotapi.InlineKeyboardButton
	for _, day := range days {
		btn := tgbotapi.NewInlineKeyboardButtonData(day, day)
		rowDays = append(rowDays, btn)
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowDays)
	return keyboard
}

func generateMonth(year int, month int, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.InlineKeyboardMarkup {
	firstDay := date(year, month, 0)
	amountDaysInMonth := date(year, month+1, 0).Day()

	weekday := int(firstDay.Weekday())
	var rowDays []tgbotapi.InlineKeyboardButton
	for i := 1; i <= weekday; i++ {
		btn := tgbotapi.NewInlineKeyboardButtonData(" ", string(rune(i)))
		rowDays = append(rowDays, btn)
	}

	amountWeek := weekday
	for i := 1; i <= amountDaysInMonth; i++ {
		if amountWeek == 7 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowDays)
			amountWeek = 0
			rowDays = []tgbotapi.InlineKeyboardButton{}
		}

		day := strconv.Itoa(i)
		if len(day) == 1 {
			day = fmt.Sprintf("0%v", day)
		}
		monthStr := strconv.Itoa(int(month))
		if len(monthStr) == 1 {
			monthStr = fmt.Sprintf("0%v", monthStr)
		}

		currDate := date(year, month, i)

		btnText := fmt.Sprintf("%v", i)
		if time.Now().Day() == i {
			btnText = fmt.Sprintf("%v!", i)
		}
		btn := tgbotapi.NewInlineKeyboardButtonData(btnText, currDate.Format("2006-01-02"))
		rowDays = append(rowDays, btn)
		amountWeek++
	}
	for i := 1; i <= 7-amountWeek; i++ {
		btn := tgbotapi.NewInlineKeyboardButtonData(" ", string(rune(i)))
		rowDays = append(rowDays, btn)
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowDays)

	return keyboard
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func addSpecialButtons(keyboard tgbotapi.InlineKeyboardMarkup, year int, month int) tgbotapi.InlineKeyboardMarkup {
	var rowDays []tgbotapi.InlineKeyboardButton
	btnPrev := tgbotapi.NewInlineKeyboardButtonData(BtnPrev, fmt.Sprintf("%s %d %d", BtnPrev, year, month))
	btnNext := tgbotapi.NewInlineKeyboardButtonData(BtnNext, fmt.Sprintf("%s %d %d", BtnNext, year, month))
	rowDays = append(rowDays, btnPrev, btnNext)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowDays)
	return keyboard
}
