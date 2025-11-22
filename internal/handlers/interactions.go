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

type InteractionHandler struct {
	ContentService services.ContentService
}

func NewInteractionHandler(contentService services.ContentService) *InteractionHandler {
	return &InteractionHandler{ContentService: contentService}
}

func (h *InteractionHandler) OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == "command" {
		h.handleCommand(s, i)
	}
}

func (h *InteractionHandler) handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	startTime := time.Now()

	// Get the category from the command options
	var category string
	if len(i.ApplicationCommandData().Options) > 0 {
		category = i.ApplicationCommandData().Options[0].StringValue()
	}

	// Log the interaction
	logger.Logger.Info("Slash command received",
		zap.String("command", "command"),
		zap.String("category", category),
		zap.String("user", i.Member.User.Username),
		zap.String("user_id", i.Member.User.ID),
		zap.String("channel_id", i.ChannelID),
		zap.String("guild_id", i.GuildID))

	// Check if category exists
	if !h.ContentService.HasCategory(category) {
		availableCategories := h.ContentService.GetAvailableCategories()
		message := fmt.Sprintf("Category '%s' not found. Available categories: %s",
			category, strings.Join(availableCategories, ", "))

		logger.Logger.Warn("Invalid category requested",
			zap.String("category", category),
			zap.Strings("available_categories", availableCategories),
			zap.String("user", i.Member.User.Username))

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
			},
		})
		return
	}

	// Get random content
	content := h.ContentService.GetRandomContent(category)
	if content == "" {
		logger.Logger.Warn("No content available for command",
			zap.String("command", category),
			zap.String("user", i.Member.User.Username))

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("No content available for `%s`", category),
			},
		})
		return
	}

	// Send the content
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})

	duration := time.Since(startTime)

	if err != nil {
		logger.Logger.Error("Failed to send content",
			zap.String("command", category),
			zap.String("user", i.Member.User.Username),
			zap.Duration("duration", duration),
			zap.Error(err))
	} else {
		logger.Logger.Info("Content sent successfully via slash command",
			zap.String("command", category),
			zap.String("user", i.Member.User.Username),
			zap.String("user_id", i.Member.User.ID),
			zap.String("channel_id", i.ChannelID),
			zap.Duration("duration", duration))
	}
}
