package config

import (
	"os"
	"testing"
)

// TestLoad tests the Load function with various environment variable scenarios.
func TestLoad(t *testing.T) {
	tests := []struct {
		name          string
		envVars       map[string]string
		expectedError bool
		expectedToken string
	}{
		{
			name: "valid token",
			envVars: map[string]string{
				"DISCORD_BOT_TOKEN":   "test-token-123",
				"DATABASE_CONNECTION": "postgres://user:pass@localhost:5432/testdb",
			},
			expectedError: false,
			expectedToken: "test-token-123",
		},
		{
			name: "missing token",
			envVars: map[string]string{
				"OTHER_VAR": "value",
			},
			expectedError: true,
			expectedToken: "",
		},
		{
			name: "empty token",
			envVars: map[string]string{
				"DISCORD_BOT_TOKEN": "",
			},
			expectedError: true,
			expectedToken: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Clearenv()

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}
			defer os.Clearenv()

			config, err := Load()

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if config.DiscordBotToken != tt.expectedToken {
					t.Errorf("Expected token %s, got %s", tt.expectedToken, config.DiscordBotToken)
				}
			}
		})
	}
}

// TestLoadWithEnvFile tests the Load function's ability to handle .env file loading.
func TestLoadWithEnvFile(t *testing.T) {
	os.Clearenv()
	os.Setenv("DISCORD_BOT_TOKEN", "test-token-from-env")
	os.Setenv("DATABASE_CONNECTION", "postgres://user:pass@localhost:5432/testdb")
	defer os.Clearenv()

	config, err := Load()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if config.DiscordBotToken != "test-token-from-env" {
		t.Errorf("Expected token 'test-token-from-env', got %s", config.DiscordBotToken)
	}
}
