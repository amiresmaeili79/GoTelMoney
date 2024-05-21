package services

import (
	"fmt"
	"github.com/amir79esmaeili/go-tel-money/internal/commands"
	"github.com/amir79esmaeili/go-tel-money/internal/messages"
	"github.com/amir79esmaeili/go-tel-money/internal/messages/widgets"
	"github.com/amir79esmaeili/go-tel-money/internal/models"
	"github.com/amir79esmaeili/go-tel-money/internal/repositories"
	"github.com/amir79esmaeili/go-tel-money/internal/state"
	"github.com/amir79esmaeili/go-tel-money/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type TelegramService struct {
	bot        *tgbotapi.BotAPI
	registry   *repositories.Registry
	inMemState *state.InMemoryState
}

func NewTelegramService(bot *tgbotapi.BotAPI,
	registry *repositories.Registry,
	memState *state.InMemoryState) *TelegramService {
	return &TelegramService{
		bot:        bot,
		registry:   registry,
		inMemState: memState,
	}
}

func (t *TelegramService) Cancel(u *tgbotapi.Update) {
	t.inMemState.DeleteState(u.Message.Chat.ID)
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, messages.CancelMessage)
	msg.ReplyMarkup = messages.MainMenu
	t.bot.Send(msg)
}

func (t *TelegramService) handleCommand(u *tgbotapi.Update) {
	switch u.Message.Command() {
	// Extract the command from the Message.
	case "start":
		t.start(u)
	default:
		t.bot.Send(tgbotapi.NewMessage(u.Message.Chat.ID, messages.Invalid))
	}
}

func (t *TelegramService) loadConversation(update *tgbotapi.Update) *models.Conversation {
	var chatID int64
	if update.Message == nil {
		chatID = update.CallbackQuery.From.ID
	} else {
		chatID = update.Message.Chat.ID
	}

	var command commands.ConvType
	var conversation *models.Conversation

	if session, err := t.inMemState.GetSession(chatID); err != nil {
		// This is a new session
		command, err = commands.WhatIsCommand(update.Message.Text)
		if err == nil {
			conversation = &models.Conversation{
				Type:  command,
				State: 0,
			}
			t.inMemState.CreateSession(chatID, conversation)
		}
	} else {
		// Continue a session
		conversation = session.(*models.Conversation)
		command = conversation.Type
	}
	return conversation
}
func (t *TelegramService) handleMsg(u *tgbotapi.Update) {

	conversation := t.loadConversation(u)
	var text string
	if u.Message != nil {
		text = u.Message.Text
	} else if u.CallbackQuery != nil {
		text = u.CallbackQuery.Data
	} else {
		slog.Error("What is this %v\n", u)
		return
	}

	if cmd, _ := commands.WhatIsCommand(text); cmd == commands.Cancel {
		t.Cancel(u)
		return
	}

	switch conversation.Type {
	case commands.AddExpenseType:
		t.addExpenseType(u, conversation)
	case commands.AddExpense:
		t.addExpense(u, conversation)
	case commands.Report:
		t.getReport(u, conversation)
	default:
		chatId := getChatID(u)
		t.bot.Send(tgbotapi.NewMessage(chatId, messages.Invalid))
	}
}

func (t *TelegramService) Listen() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				t.handleCommand(&update)
				continue
			}
		}
		t.handleMsg(&update)
	}
	return nil
}

// Start Starts the bot and registers the user
func (t *TelegramService) start(update *tgbotapi.Update) {
	telegramId := update.Message.From.ID
	// test if user is already available or not
	_, err := t.registry.UserRepository().GetByTID(telegramId)
	if err != nil {
		newUser := &models.User{
			UserTelegramID: update.Message.From.ID,
			FirstName:      update.Message.From.FirstName,
			LastName:       update.Message.From.LastName,
		}
		err := t.registry.UserRepository().Create(newUser)
		if err != nil {
			slog.Error("Could not create user", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messages.UserCreationFailed)
			t.bot.Send(msg)
		}
	}

	textMessage := fmt.Sprintf(messages.WelcomeMsg, update.Message.From.FirstName)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, textMessage)
	msg.ReplyMarkup = messages.MainMenu
	t.bot.Send(msg)
}

// AddExpenseType Handler to add a new type to database
func (t *TelegramService) addExpenseType(u *tgbotapi.Update, c *models.Conversation) {
	chatId := getChatID(u)
	// Check if user exists
	user, err := t.fetchUser(chatId)
	if err != nil {
		t.inMemState.DeleteState(chatId)
		return
	}

	switch c.State {
	case int(commands.StartAddExpenseType):
		allTypes := t.registry.ExpenseTypeRepository().All(user.ID, -1, -1)
		msg := tgbotapi.NewMessage(chatId, messages.NewType)
		msg.ReplyMarkup = messages.CancelKeyboard
		t.bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, models.PrettyPrintExpenseTypes(allTypes))
		t.bot.Send(msg)
		c.State += 1
	case int(commands.AskNameAddExpenseType):
		name := u.Message.Text
		// Check if type is already added or not!?
		_, err = t.registry.ExpenseTypeRepository().GetByName(name, user.ID)
		if err == nil {
			// this type exists
			log.Printf("Failed to add a new type, %v. It does exist\n", err)
			t.bot.Send(tgbotapi.NewMessage(chatId, fmt.Sprintf(messages.TypeAlreadyAdded, name)))
			return
		}
		expenseType := &models.ExpenseType{Name: name, UserID: user.ID}
		if err := t.registry.ExpenseTypeRepository().Create(expenseType); err != nil {
			msg := tgbotapi.NewMessage(chatId,
				fmt.Sprintf(messages.TypeAddingFailed, expenseType.Name))
			t.bot.Send(msg)
			slog.Error("Failed to add a new type, %v\n", err)
			return
		}
		msg := tgbotapi.NewMessage(chatId, fmt.Sprintf(messages.TypeAddedSuccessfully, expenseType.Name))
		t.bot.Send(msg)
		t.inMemState.DeleteState(chatId)
	}
}

func (t *TelegramService) fetchUser(chatId int64) (*models.User, error) {
	user, err := t.registry.UserRepository().GetByTID(chatId)
	if err != nil {
		msg := tgbotapi.NewMessage(chatId, messages.UserNotFound)
		slog.Error("Failed to find user, %v\n", err)
		t.bot.Send(msg)
		return nil, err
	}
	return user, err
}

// addExpense It is a service that adds new expense
func (t *TelegramService) addExpense(u *tgbotapi.Update, c *models.Conversation) {

	chatId := getChatID(u)
	// Check if user exists
	user, err := t.fetchUser(chatId)
	if err != nil {
		t.inMemState.DeleteState(chatId)
		return
	}

	var expense *models.Expense
	if c.Data != nil {
		expense = c.Data.(*models.Expense)
	} else {
		expense = &models.Expense{
			UserID: user.ID,
		}
		c.Data = expense
	}

	switch c.State {
	case int(commands.StartAddExpense):
		msg := tgbotapi.NewMessage(chatId, messages.NewAmount)
		msg.ReplyMarkup = messages.CancelKeyboard
		t.bot.Send(msg)
		c.State += 1
	case int(commands.AskAmountAddExpense):
		s, err := strconv.ParseFloat(u.Message.Text, 32)
		if err != nil {
			slog.Error("Failed to parse amount, %v\n", err)
			t.bot.Send(tgbotapi.NewMessage(chatId, messages.InvalidAmount))
			return
		}
		expense.Amount = float32(s)
		c.State += 1
		t.bot.Send(tgbotapi.NewMessage(chatId, messages.NewDescription))
	case int(commands.AskDescriptionAddExpense):
		expense.Description = u.Message.Text
		c.State += 1
		msg := tgbotapi.NewMessage(chatId, messages.NewDate)
		msg.ReplyMarkup = widgets.GenerateCalendar(time.Now().Year(), time.Now().Month())
		t.bot.Send(msg)
	case int(commands.AskDateAddExpense):
		var date time.Time
		if u.Message == nil {
			if strings.ContainsAny(u.CallbackQuery.Data, widgets.BtnNext+widgets.BtnPrev) {
				cal := widgets.HandleButton(u.CallbackQuery.Data)
				msg := tgbotapi.NewEditMessageReplyMarkup(chatId, u.CallbackQuery.Message.MessageID, cal)
				t.bot.Send(msg)
				return
			} else {
				date, err = time.Parse("2006-01-02", u.CallbackQuery.Data)
				callback := tgbotapi.NewCallback(u.CallbackQuery.ID, u.CallbackQuery.Data)
				t.bot.Request(callback)
			}
		} else {
			date, err = time.Parse("2006-01-02", u.Message.Text)
		}
		expense.Date = date

		exTypes := t.registry.ExpenseTypeRepository().All(user.ID, -1, -1)
		var data []models.InlineKeyboardItem
		for _, exType := range exTypes {
			data = append(data, &exType) // converting to struct
		}
		replyMarkup := widgets.GetDynamicInlineKeyboard(data, 2)
		msg := tgbotapi.NewMessage(chatId, messages.AskType)
		msg.ReplyMarkup = replyMarkup
		c.State += 1
		t.bot.Send(msg)
	case int(commands.AskTypeAddExpense):
		id, err := strconv.Atoi(u.CallbackQuery.Data)
		if err != nil {
			slog.Error("Failed to parse id, %v\n", err)
			t.bot.Send(tgbotapi.NewMessage(chatId, messages.InvalidType))
		}
		exType, err := t.registry.ExpenseTypeRepository().GetByID(uint(id), user.ID)
		if err != nil {
			slog.Error("Invalid type %s\n", u.CallbackQuery.Data)
			t.bot.Send(tgbotapi.NewMessage(chatId, messages.InvalidType))
			return
		}
		expense.ExpenseTypeID = exType.ID
		err = t.registry.ExpenseRepository().Create(expense)
		if err != nil {
			slog.Error("Failed to save the expense, %v\n", err)
			t.bot.Send(tgbotapi.NewMessage(chatId, messages.ExpenseSaveFailed))
			return
		}
		t.inMemState.DeleteState(chatId)
		callback := tgbotapi.NewCallback(u.CallbackQuery.ID, u.CallbackQuery.Data)
		t.bot.Request(callback)
		msg := tgbotapi.NewMessage(chatId, messages.NewExpenseSaved)
		msg.ReplyMarkup = messages.MainMenu
		t.bot.Send(msg)
	}
}

func (t *TelegramService) getReport(u *tgbotapi.Update, c *models.Conversation) {
	chatId := getChatID(u)
	user, _ := t.fetchUser(chatId)

	pageSize := 10
	page := 1

	switch c.State {
	case int(commands.StartReport):
		msg := tgbotapi.NewMessage(chatId, messages.StartReport)
		msg.ReplyMarkup = messages.CancelKeyboard
		t.bot.Send(msg)

		msg = tgbotapi.NewMessage(chatId, messages.SelectReportRange)
		msg.ReplyMarkup = messages.ReportRangeKeyboard
		t.bot.Send(msg)
		c.State += 1
	case int(commands.AskReportRange):
		dataRange := u.CallbackQuery.Data
		callback := tgbotapi.NewCallback(u.CallbackQuery.ID, u.CallbackQuery.Data)
		t.bot.Request(callback)
		start, end, err := utils.GetDateRange(dataRange)
		c.Data = [2]time.Time{start, end}
		if err != nil {
			slog.Error("Invalid date range selected!, %v\n", err)
			t.bot.Send(tgbotapi.NewMessage(chatId, messages.InvalidDate))
		}
		pages, expenses, hasNextPage := t.registry.ExpenseRepository().Filter(
			page, pageSize, "user_id = ? AND date >= ? AND date <= ?", user.ID, start, end)
		messageText := models.PrettyPrintExpenses(expenses, page, pages)
		kb := widgets.AddPaginationButtons(nil, page, false, hasNextPage)
		msg := tgbotapi.NewEditMessageTextAndMarkup(chatId, u.CallbackQuery.Message.MessageID, messageText, *kb)
		t.bot.Send(msg)
		c.State += 1
	case int(commands.ViewData):
		dateRange := c.Data.([2]time.Time)
		page, _ = strconv.Atoi(u.CallbackQuery.Data)
		pages, expenses, hasNextPage := t.registry.ExpenseRepository().Filter(
			page, pageSize, "user_id = ? AND date >= ? AND date <= ?", user.ID, dateRange[0], dateRange[1])
		messageText := models.PrettyPrintExpenses(expenses, page, pages)
		kb := widgets.AddPaginationButtons(nil, page, page > 1, hasNextPage)
		msg := tgbotapi.NewEditMessageTextAndMarkup(
			chatId, u.CallbackQuery.Message.MessageID, messageText, *kb,
		)
		t.bot.Send(msg)
	}
}

func getChatID(u *tgbotapi.Update) int64 {
	if u.Message == nil {
		return u.CallbackQuery.From.ID
	}
	return u.Message.Chat.ID
}
