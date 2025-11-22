package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mutsumi-bot/internal/bot"
	"mutsumi-bot/internal/config"
	"mutsumi-bot/internal/handlers"
	"mutsumi-bot/internal/logger"
	"mutsumi-bot/internal/services"

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

	logger.Logger.Info("Starting mutsumi-bot")

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

	// Start health check HTTP server
	healthPort := os.Getenv("HEALTH_PORT")
	if healthPort == "" {
		healthPort = "8089"
	}

	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/health", healthHandler(databaseService))

	healthServer := &http.Server{
		Addr:         ":" + healthPort,
		Handler:      healthMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Start health server in a goroutine
	go func() {
		logger.Logger.Info("Starting health check server", zap.String("port", healthPort))
		if err := healthServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatal("health server error", zap.Error(err))
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start bot in a goroutine
	go func() {
		if err := b.StartWithCommands(ctx, commands); err != nil {
			logger.Logger.Fatal("run error", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	logger.Logger.Info("Shutting down...")

	// Gracefully shutdown health server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := healthServer.Shutdown(shutdownCtx); err != nil {
		logger.Logger.Error("Error shutting down health server", zap.Error(err))
	}
}

// healthHandler returns a handler function for the /health endpoint
func healthHandler(dbService *services.DatabaseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check database connection
		dbHealthy := true
		if err := dbService.Ping(); err != nil {
			dbHealthy = false
			logger.Logger.Warn("Database health check failed", zap.Error(err))
		}

		status := "healthy"
		statusCode := http.StatusOK
		if !dbHealthy {
			status = "unhealthy"
			statusCode = http.StatusServiceUnavailable
		}

		response := map[string]interface{}{
			"status":    status,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"database": map[string]interface{}{
				"connected": dbHealthy,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
	}
}
