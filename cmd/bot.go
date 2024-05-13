package cmd

import (
	"fmt"
	"log"

	"github.com/amir79esmaeili/go-tel-money/internal/cfg"
	"github.com/amir79esmaeili/go-tel-money/internal/postgres"
	"github.com/amir79esmaeili/go-tel-money/internal/repositories"
	"github.com/amir79esmaeili/go-tel-money/internal/services"
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
	telegramService := services.NewTelegramService(userRepository)

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
		}
	}
}
