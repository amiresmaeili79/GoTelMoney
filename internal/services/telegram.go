package services

import (
	"fmt"
	"log"

	"github.com/amir79esmaeili/go-tel-money/internal/conversations"
	"github.com/amir79esmaeili/go-tel-money/internal/models"
	"github.com/amir79esmaeili/go-tel-money/internal/repositories"
	"github.com/amir79esmaeili/go-tel-money/internal/state"
	"github.com/amir79esmaeili/go-tel-money/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService struct {
	userRepository        repositories.UserRepository
	expenseTypeRepository repositories.ExpenseTypeRepository
	inMemState            *state.InMemoryState
}

func NewTelegramService(userRepo repositories.UserRepository, expenseTypeRepository repositories.ExpenseTypeRepository,
	memState *state.InMemoryState) *TelegramService {
	return &TelegramService{
		userRepository:        userRepo,
		expenseTypeRepository: expenseTypeRepository,
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

	textMessage := fmt.Sprintf(conversations.WelcomeMsg, update.Message.From.FirstName)
	reply := tgbotapi.NewMessage(update.Message.Chat.ID, textMessage)
	reply.ReplyMarkup = conversations.MainMenu
	return &reply, nil
}

// Handler to add a new type to database
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
	case int(conversations.StartAddExpenseType):
		replyMsg = "Please add a type:\n"
		conversation.State += 1
		all_types := tService.expenseTypeRepository.All(user.ID)
		replyMsg += utils.PrettyPrintExpenseTypes(all_types)

	case int(conversations.AskNameAddExpenseType):
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
