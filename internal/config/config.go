package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordBotToken    string
	DatabaseConnection string
}

// Load reads configuration from environment variables and validates required fields.
// If a .env file is present in the working directory, it will be loaded first.
func Load() (Config, error) {
	// Load .env if present; ignore error when file does not exist
	_ = godotenv.Load()

	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		return Config{}, errors.New("missing DISCORD_BOT_TOKEN env var")
	}

	dbConn := os.Getenv("DATABASE_CONNECTION")
	if dbConn == "" {
		return Config{}, errors.New("missing DATABASE_CONNECTION env var")
	}

	return Config{
		DiscordBotToken:    token,
		DatabaseConnection: dbConn,
	}, nil
}
