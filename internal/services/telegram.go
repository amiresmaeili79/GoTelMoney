package services

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/amir79esmaeili/go-tel-money/internal/messages"
	"github.com/amir79esmaeili/go-tel-money/internal/models"
	"github.com/amir79esmaeili/go-tel-money/internal/repositories"
	"github.com/amir79esmaeili/go-tel-money/internal/state"
	"github.com/amir79esmaeili/go-tel-money/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService struct {
	userRepository        repositories.UserRepository
	expenseTypeRepository repositories.ExpenseTypeRepository
	expenseRepository     repositories.ExpenseRepository
	inMemState            *state.InMemoryState
}

func NewTelegramService(userRepo repositories.UserRepository,
	expenseTypeRepository repositories.ExpenseTypeRepository,
	expenseRepository repositories.ExpenseRepository,
	memState *state.InMemoryState) *TelegramService {
	return &TelegramService{
		userRepository:        userRepo,
		expenseTypeRepository: expenseTypeRepository,
		expenseRepository:     expenseRepository,
		inMemState:            memState,
	}
}

func (tService *TelegramService) Start(update tgbotapi.Update) (*tgbotapi.MessageConfig, error) {
	telegramId := update.Message.From.ID
	// test if user is already available or not
	_, err := tService.userRepository.GetByTID(telegramId)
	if err != nil {
		newUser := &models.User{
			UserTelegramID: update.Message.From.ID,
			FirstName:      update.Message.From.FirstName,
			LastName:       update.Message.From.LastName,
		}
		err := tService.userRepository.Create(newUser)
		if err != nil {
			return nil, err
		}
	}

	textMessage := fmt.Sprintf(messages.WelcomeMsg, update.Message.From.FirstName)
	reply := tgbotapi.NewMessage(update.Message.Chat.ID, textMessage)
	reply.ReplyMarkup = messages.MainMenu
	return &reply, nil
}

// AddExpenseType Handler to add a new type to database
func (tService *TelegramService) AddExpenseType(update tgbotapi.Update) *tgbotapi.MessageConfig {
	var replyMsg string

	chatId := update.Message.Chat.ID
	conv, _ := tService.inMemState.GetSession(chatId)

	tid := update.Message.From.ID
	// Check if user exists
	user, err := tService.userRepository.GetByTID(tid)
	if err != nil {
		replyMsg = "You are not registered in the bot. please type '/start'. Then try again!"
		reply := tgbotapi.NewMessage(chatId, replyMsg)
		log.Printf("Failed to find user, %v\n", err)
		return &reply
	}

	conversation := conv.(*models.Conversation)
	switch conversation.State {
	case int(messages.StartAddExpenseType):
		replyMsg = "Please add a type:\n"
		conversation.State += 1
		all_types := tService.expenseTypeRepository.All(user.ID)
		replyMsg += utils.PrettyPrintExpenseTypes(all_types)

	case int(messages.AskNameAddExpenseType):
		name := update.Message.Text

		// Check if type is already added or not!?
		_, err = tService.expenseTypeRepository.FindByName(name, user.ID)
		if err == nil {
			// this type exists
			replyMsg = fmt.Sprintf("Failed to add %s! It is already added.", name)
			log.Printf("Failed to add a new type, %v. It does exist\n", err)
			break
		}

		expenseType := &models.ExpenseType{
			Name:   name,
			UserID: user.ID,
		}

		replyMsg = fmt.Sprintf("%s added!\n", expenseType.Name)
		if err := tService.expenseTypeRepository.Create(expenseType); err != nil {
			replyMsg = fmt.Sprintf("Failed to add %s! Try again later please.", expenseType.Name)
			log.Printf("Failed to add a new type, %v\n", err)
		}

		tService.inMemState.DeleteState(chatId)
	}

	reply := tgbotapi.NewMessage(chatId, replyMsg)
	return &reply
}

func (tService *TelegramService) AddExpense(update tgbotapi.Update) *tgbotapi.MessageConfig {
	var replyMsg string
	//var replyMarkup interface{}

	chatId := update.Message.Chat.ID
	conv, _ := tService.inMemState.GetSession(chatId)

	tid := update.Message.From.ID
	// Check if user exists
	user, err := tService.userRepository.GetByTID(tid)
	if err != nil {
		replyMsg = "You are not registered in the bot. please type '/start'. Then try again!"
		reply := tgbotapi.NewMessage(chatId, replyMsg)
		log.Printf("Failed to find user, %v\n", err)
		return &reply
	}

	conversation := conv.(*models.Conversation)
	var expense *models.Expense

	if conversation.Data != nil {
		expense = conversation.Data.(*models.Expense)
	} else {
		expense = &models.Expense{}
		conversation.Data = expense
	}

	switch conversation.State {
	case int(messages.StartAddExpense):
		replyMsg = "Amount:\n"
		conversation.State += 1
	case int(messages.AskAmountAddExpense):
		replyMsg = "Description:\n"
		conversation.State += 1
		s, err := strconv.ParseFloat(update.Message.Text, 32)
		if err != nil {
			replyMsg = "This is not a valid number! Try again. (e.g. 23.5)"
			break
		}
		expense.Amount = float32(s)
	case int(messages.AskDescriptionAddExpense):
		replyMsg = "Date:\n"
		conversation.State += 1
		expense.Description = update.Message.Text
	case int(messages.AskDateAddExpense):
		replyMsg = "Type:\n"
		date, err := time.Parse("2006-01-02 15:04:05", update.Message.Text)
		if err != nil {
			replyMsg = "Invalid date! Try again (e.g. 2024-05-25 14:35)"
			break
		}
		conversation.State += 1
		expense.Date = date
		exTypes := tService.expenseTypeRepository.All(user.ID)

		var data []models.Choosable
		for _, exType := range exTypes {
			data = append(data, &exType) // converting to struct
		}
		//replyMarkup = utils.GetDynamicInlineKeyboard(data)
	case int(messages.AskTypeAddExpense):
		replyMsg = "Done!"
		exTypeName := update.Message.Text
		exType, err := tService.expenseTypeRepository.FindByName(exTypeName, user.ID)
		if err != nil {
			replyMsg = "Invalid type! Please select one from menu!"
			log.Printf("Invalid type %s\n", exTypeName)
			break
		}
		expense.ExpenseTypeID = exType.ID
		err = tService.expenseRepository.Create(expense)
		if err != nil {
			replyMsg = "Failed to save the expense!"
			log.Printf("Failed to save the expense, %v\n", err)
			break
		}

	}
	reply := tgbotapi.NewMessage(chatId, replyMsg)
	reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
			tgbotapi.NewInlineKeyboardButtonData("2", "2"),
			tgbotapi.NewInlineKeyboardButtonData("3", "3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("4", "4"),
			tgbotapi.NewInlineKeyboardButtonData("5", "5"),
			tgbotapi.NewInlineKeyboardButtonData("6", "6"),
		),
	)
	return &reply
}
