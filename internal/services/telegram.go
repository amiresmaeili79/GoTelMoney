package services

import (
	"fmt"

	"github.com/amir79esmaeili/go-tel-money/internal/models"
	"github.com/amir79esmaeili/go-tel-money/internal/repositories"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService struct {
	userRepository repositories.UserRepository
}

func NewTelegramService(userRepo repositories.UserRepository) *TelegramService {
	return &TelegramService{
		userRepository: userRepo,
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

	textMessage := fmt.Sprintf(WelcomeMsg, update.Message.From.FirstName)
	reply := tgbotapi.NewMessage(update.Message.Chat.ID, textMessage)
	reply.ReplyMarkup = MainMenu
	return &reply, nil
}
