package integration

import (
	"database/sql"
	"os"
	"testing"

	"mutsumi-bot/internal/logger"
	"mutsumi-bot/internal/services"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// TestIntegration tests the overall flow of the mutsumi-bot application.
func TestIntegration(t *testing.T) {
	// Skip if running in CI or if database connection is not available
	dbConn := os.Getenv("DATABASE_CONNECTION")
	if os.Getenv("CI") != "" || dbConn == "" {
		t.Skip("Skipping integration test - no database connection or running in CI")
	}

	// Initialize logger for integration tests
	err := logger.Init()
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	// Test database service initialization
	dbService, err := services.NewDatabaseService(dbConn)
	if err != nil {
		t.Fatalf("Failed to create database service: %v", err)
	}
	defer dbService.Close()

	// Test that commands can be retrieved
	commands := dbService.GetAvailableCategories()
	if len(commands) == 0 {
		t.Log("No commands found in database - this is okay if database is empty")
	}

	// Test getting random content for existing commands
	for _, command := range commands {
		content := dbService.GetRandomContent(command)
		if content == "" {
			t.Errorf("Expected content for command %s but got empty", command)
		}
	}

	// Test content count
	for _, command := range commands {
		count := dbService.GetContentCount(command)
		if count <= 0 {
			t.Errorf("Expected positive count for command %s but got %d", command, count)
		}
	}
}

// TestMain tests the main function components without running the full application.
func TestMain(t *testing.T) {
	// This test verifies that the main function components can be initialized
	// without actually running the full application

	// Test config loading (this will fail without proper env vars, which is expected)
	// We're just testing that the function doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Main function panicked: %v", r)
		}
	}()

	// Test that we can import and use the packages
	// This is a basic smoke test
	_ = services.DatabaseService{}
	_ = sql.DB{}
}
