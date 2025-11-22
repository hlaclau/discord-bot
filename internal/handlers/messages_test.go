package handlers

import (
	"testing"

	"mutsumi-bot/internal/logger"
)

// setupTestHandler creates a message handler with mock content service
func setupTestHandler(t *testing.T) *MessageHandler {
	// Initialize logger for tests
	err := logger.Init()
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	t.Cleanup(func() {
		logger.Close()
	})

	// Create mock content service
	mockService := newMockContentService()
	mockService.addCommand("mutsumi", "Mutsumi content 1", "Mutsumi content 2")
	mockService.addCommand("cats", "Cats content 1", "Cats content 2")

	// Create message handler
	handler := NewMessageHandler(mockService)

	return handler
}

// TestNewMessageHandler tests the NewMessageHandler constructor function.
func TestNewMessageHandler(t *testing.T) {
	// Create a mock content service
	mockService := newMockContentService()

	handler := NewMessageHandler(mockService)

	if handler == nil {
		t.Fatalf("Expected handler but got nil")
		return
	}

	if handler.ContentService != mockService {
		t.Errorf("Expected content service %v but got %v", mockService, handler.ContentService)
	}
}

// TestMessageHandler_ContentServiceIntegration tests the integration between MessageHandler and ContentService.
func TestMessageHandler_ContentServiceIntegration(t *testing.T) {
	handler := setupTestHandler(t)

	// Test that the handler has access to the content service
	if handler.ContentService == nil {
		t.Errorf("Expected content service but got nil")
	}

	// Test that the content service has commands
	commands := handler.ContentService.GetAvailableCategories()
	if len(commands) == 0 {
		t.Errorf("Expected commands but got none")
	}

	// Test that we can get random content
	for _, command := range commands {
		content := handler.ContentService.GetRandomContent(command)
		if content == "" {
			t.Errorf("Expected content for command %s but got empty", command)
		}
	}
}
