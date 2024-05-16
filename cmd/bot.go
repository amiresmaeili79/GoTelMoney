package cmd

import (
	"log"
	"log/slog"

	"github.com/amir79esmaeili/go-tel-money/internal/cfg"
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
	slog.Info("Authorized on account %s", bot.Self.UserName)

	return bot
}

func startGoTelMoney(cmd *cobra.Command) {
	configPath, _ := cmd.Flags().GetString("cfg")
	config := cfg.ParseConfig(configPath)

	bot := setUpTelegram(&config)

	db, err := postgres.ConnectToDB(&config)
	if err != nil {
		panic(err)
	}

	userRepository := repositories.NewUserRepository(db)
	expenseTypeRepository := repositories.NewExpenseTypeRepository(db)
	expenseRepository := repositories.NewExpenseRepository(db)
	registry := repositories.NewRegistry(userRepository, expenseTypeRepository, expenseRepository)
	inMemState := state.NewInMemoryState()
	telegramService := services.NewTelegramService(bot, registry, inMemState)

	err = telegramService.Listen()
	if err != nil {
		slog.Error("Could not start telegram bot", err)
	}
}
