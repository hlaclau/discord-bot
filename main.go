package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"wooper-bot/internal/bot"
	"wooper-bot/internal/config"
	"wooper-bot/internal/handlers"
	"wooper-bot/internal/logger"
	"wooper-bot/internal/services"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func buildCategoryChoices(contentService services.ContentService) []*discordgo.ApplicationCommandOptionChoice {
	categories := contentService.GetAvailableCategories()

	choices := make([]*discordgo.ApplicationCommandOptionChoice, len(categories))
	for i, category := range categories {
		choices[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  category,
			Value: category,
		}
	}

	return choices
}

func main() {
	// Initialize logging
	if err := logger.Init(); err != nil {
		log.Fatalf("logger init error: %v", err)
	}
	defer logger.Close()

	logger.Logger.Info("Starting wooper-bot")

	cfg, err := config.Load()
	if err != nil {
		logger.Logger.Fatal("config error", zap.Error(err))
	}

	databaseService, err := services.NewDatabaseService(cfg.DatabaseConnection)
	if err != nil {
		logger.Logger.Fatal("database service error", zap.Error(err))
	}
	defer databaseService.Close()

	messageHandler := handlers.NewMessageHandler(databaseService)
	interactionHandler := handlers.NewInteractionHandler(databaseService)

	b, err := bot.New(cfg.DiscordBotToken)
	if err != nil {
		logger.Logger.Fatal("bot error", zap.Error(err))
	}
	b.AddHandler(messageHandler.OnMessageCreate)
	b.AddHandler(interactionHandler.OnInteractionCreate)

	// Register slash commands
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "command",
			Description: "Get content from a registered command",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "command",
					Description: "Command to get content from",
					Required:    true,
					Choices:     buildCategoryChoices(databaseService),
				},
			},
		},
	}

	logger.Logger.Info("Bot initialized successfully")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := b.StartWithCommands(ctx, commands); err != nil {
		logger.Logger.Fatal("run error", zap.Error(err))
	}
}
