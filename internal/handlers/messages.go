package handlers

import (
	"fmt"
	"strings"
	"time"

	"mutsumi-bot/internal/logger"
	"mutsumi-bot/internal/services"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type MessageHandler struct {
	ContentService services.ContentService
}

func NewMessageHandler(contentService services.ContentService) *MessageHandler {
	return &MessageHandler{ContentService: contentService}
}

func (h *MessageHandler) OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author == nil || m.Author.Bot {
		return
	}

	content := strings.TrimSpace(m.Content)

	// Log all messages for debugging (can be filtered by log level)
	logger.Logger.Debug("Message received",
		zap.String("user", m.Author.Username),
		zap.String("user_id", m.Author.ID),
		zap.String("channel_id", m.ChannelID),
		zap.String("guild_id", m.GuildID),
		zap.String("content", content))

	// Check if message starts with ! and has a valid category
	if strings.HasPrefix(content, "!") {
		category := strings.TrimPrefix(content, "!")

		// Log command attempt
		logger.Logger.Info("Command received",
			zap.String("command", content),
			zap.String("category", category),
			zap.String("user", m.Author.Username),
			zap.String("user_id", m.Author.ID),
			zap.String("channel_id", m.ChannelID),
			zap.String("guild_id", m.GuildID))

		if h.ContentService.HasCategory(category) {
			startTime := time.Now()

			content := h.ContentService.GetRandomContent(category)
			if content == "" {
				logger.Logger.Warn("No content available for command",
					zap.String("command", category),
					zap.String("user", m.Author.Username))
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("no content available for `!%s`", category))
				return
			}

			_, err := s.ChannelMessageSend(m.ChannelID, content)
			duration := time.Since(startTime)

			if err != nil {
				logger.Logger.Error("Failed to send content",
					zap.String("command", category),
					zap.String("user", m.Author.Username),
					zap.Duration("duration", duration),
					zap.Error(err))
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("failed to send content for `!%s`: %v", category, err))
			} else {
				logger.Logger.Info("Content sent successfully",
					zap.String("command", category),
					zap.String("user", m.Author.Username),
					zap.String("user_id", m.Author.ID),
					zap.String("channel_id", m.ChannelID),
					zap.Duration("duration", duration))
			}
		} else if category == "help" || category == "list" {
			// Show available categories
			logger.Logger.Info("Help command requested",
				zap.String("user", m.Author.Username),
				zap.String("user_id", m.Author.ID))

			categories := h.ContentService.GetAvailableCategories()
			if len(categories) == 0 {
				logger.Logger.Warn("No categories available for help",
					zap.String("user", m.Author.Username))
				_, _ = s.ChannelMessageSend(m.ChannelID, "no commands available")
				return
			}

			message := "Available commands:\n"
			for _, cat := range categories {
				count := h.ContentService.GetContentCount(cat)
				message += fmt.Sprintf("• `!%s` (%d entries)\n", cat, count)
			}

			logger.Logger.Info("Help response sent",
				zap.String("user", m.Author.Username),
				zap.Int("categories_count", len(categories)))

			_, _ = s.ChannelMessageSend(m.ChannelID, message)
		} else {
			// Unknown command
			logger.Logger.Info("Unknown command received",
				zap.String("command", content),
				zap.String("category", category),
				zap.String("user", m.Author.Username),
				zap.String("user_id", m.Author.ID))
		}
	}
}
