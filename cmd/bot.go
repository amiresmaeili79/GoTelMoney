package cmd

import (
	"fmt"
	"log"

	"github.com/amir79esmaeili/go-tel-money/internal/cfg"
	"github.com/amir79esmaeili/go-tel-money/internal/commands"
	"github.com/amir79esmaeili/go-tel-money/internal/conversations"
	"github.com/amir79esmaeili/go-tel-money/internal/models"
	"github.com/amir79esmaeili/go-tel-money/internal/postgres"
	"github.com/amir79esmaeili/go-tel-money/internal/repositories"
	"github.com/amir79esmaeili/go-tel-money/internal/services"
	"github.com/amir79esmaeili/go-tel-money/internal/state"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
)

func addBotCommand(root *cobra.Command) {
	botCmd := &cobra.Command{
		Use:   "bot",
		Short: "Start telegram bot",
		Long:  "Start the telegram bot",
		Run: func(cmd *cobra.Command, args []string) {
			startGoTelMoney(cmd)
		},
	}

	botCmd.Flags().StringP("cfg", "c", ".env", "Config file path")
	root.AddCommand(botCmd)
}

func setUpTelegram(config *cfg.Config) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(config.TelgramToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return bot
}

// Handles chat based commands with respect to session
func handleMessageCommands(bot *tgbotapi.BotAPI, update *tgbotapi.Update,
	inMemState *state.InMemoryState, tlgService *services.TelegramService) {
	chatID := update.Message.Chat.ID

	var command conversations.ConvType
	var conversation *models.Conversation

	if session, err := inMemState.GetSession(chatID); err != nil {
		// This is a new session
		command, err = commands.WhatIsCommand(update.Message.Text)
		if err == nil {
			conversation = &models.Conversation{
				Type:  command,
				State: 0,
			}
			inMemState.CreateSession(chatID, conversation)
		}
	} else {
		// Continue a session
		conversation = session.(*models.Conversation)
		command = conversation.Type
	}

	var reply *tgbotapi.MessageConfig

	switch command {
	case conversations.AddExpenseType:
		reply = tlgService.AddExpenseType(*update)
	default:
		invalidMess := tgbotapi.NewMessage(chatID, conversations.Invalid)
		reply = &invalidMess
	}

	if _, err := bot.Send(reply); err != nil {
		log.Panic(err)
	}
}

func startGoTelMoney(cmd *cobra.Command) {
	configPath, _ := cmd.Flags().GetString("cfg")
	config := cfg.ParseConfig(configPath)

	bot := setUpTelegram(&config)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	db, err := postgres.ConnectToDB(&config)
	if err != nil {
		panic(err)
	}

	userRepository := repositories.NewUserRepository(db)
	expenseTypeRepository := repositories.NewExpenseTypeRepository(db)
	inMemState := state.NewInMemoryState()
	telegramService := services.NewTelegramService(userRepository, expenseTypeRepository, inMemState)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}
		var reply *tgbotapi.MessageConfig

		if update.Message.IsCommand() {
			// Extract the command from the Message.
			switch update.Message.Command() {
			case "start":
				reply, err = telegramService.Start(update)
				if err != nil {
					log.Printf("Failed to start the bot, %v\n", err)
				}
			default:
				fmt.Println("hi")

			}
			if _, err := bot.Send(reply); err != nil {
				log.Panic(err)
			}
		} else {
			handleMessageCommands(bot, &update, inMemState, telegramService)
		}
	}
}
