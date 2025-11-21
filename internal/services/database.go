package services

import (
	"database/sql"
	"fmt"

	"wooper-bot/internal/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type DatabaseService struct {
	db *sql.DB
}

// NewDatabaseService creates a new database service with a PostgreSQL connection
func NewDatabaseService(connectionString string) (*DatabaseService, error) {
	logger.Logger.Info("Initializing database service")

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		logger.Logger.Error("Failed to open database connection", zap.Error(err))
		return nil, fmt.Errorf("open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		logger.Logger.Error("Failed to ping database", zap.Error(err))
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	// Ensure the table exists
	if err := ensureTableExists(db); err != nil {
		logger.Logger.Error("Failed to ensure table exists", zap.Error(err))
		db.Close()
		return nil, fmt.Errorf("ensure table exists: %w", err)
	}

	service := &DatabaseService{db: db}

	// Log available commands
	commands := service.GetAvailableCategories()
	logger.Logger.Info("Database service initialized successfully",
		zap.Int("total_commands", len(commands)))

	return service, nil
}

// ensureTableExists creates the commands table if it doesn't exist
func ensureTableExists(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS commands (
		id SERIAL PRIMARY KEY,
		command VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_command ON commands(command);
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	logger.Logger.Info("Commands table ensured")
	return nil
}

// getRandomContentInternal retrieves a random content entry for a given command (internal method)
func (s *DatabaseService) getRandomContentInternal(command string) (string, error) {
	query := `
		SELECT content FROM commands
		WHERE command = $1
		ORDER BY RANDOM()
		LIMIT 1
	`

	var content string
	err := s.db.QueryRow(query, command).Scan(&content)
	if err == sql.ErrNoRows {
		logger.Logger.Debug("No content found for command", zap.String("command", command))
		return "", nil
	}
	if err != nil {
		logger.Logger.Error("Failed to query database", zap.String("command", command), zap.Error(err))
		return "", fmt.Errorf("query database: %w", err)
	}

	logger.Logger.Debug("Retrieved content for command",
		zap.String("command", command),
		zap.String("content", content))

	return content, nil
}

// getContentCountInternal returns the number of content entries for a command (internal method)
func (s *DatabaseService) getContentCountInternal(command string) (int, error) {
	query := `SELECT COUNT(*) FROM commands WHERE command = $1`

	var count int
	err := s.db.QueryRow(query, command).Scan(&count)
	if err != nil {
		logger.Logger.Error("Failed to count content", zap.String("command", command), zap.Error(err))
		return 0, fmt.Errorf("count content: %w", err)
	}

	return count, nil
}

// GetAvailableCategories returns all unique commands from the database
// Implements ContentService interface
func (s *DatabaseService) GetAvailableCategories() []string {
	query := `SELECT DISTINCT command FROM commands ORDER BY command`

	rows, err := s.db.Query(query)
	if err != nil {
		logger.Logger.Error("Failed to query commands", zap.Error(err))
		return []string{}
	}
	defer rows.Close()

	var commands []string
	for rows.Next() {
		var command string
		if err := rows.Scan(&command); err != nil {
			logger.Logger.Error("Failed to scan command", zap.Error(err))
			continue
		}
		commands = append(commands, command)
	}

	if err := rows.Err(); err != nil {
		logger.Logger.Error("Error iterating rows", zap.Error(err))
		return []string{}
	}

	return commands
}

// HasCategory checks if a command exists in the database
func (s *DatabaseService) HasCategory(command string) bool {
	count, err := s.getContentCountInternal(command)
	if err != nil {
		logger.Logger.Error("Failed to check category", zap.String("command", command), zap.Error(err))
		return false
	}
	return count > 0
}

// Close closes the database connection
func (s *DatabaseService) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// ContentService interface methods

// GetRandomContent returns a random content string for the given command
func (s *DatabaseService) GetRandomContent(command string) string {
	content, err := s.getRandomContentInternal(command)
	if err != nil {
		logger.Logger.Error("Failed to get random content", zap.String("command", command), zap.Error(err))
		return ""
	}
	return content
}

// GetContentCount returns the number of content entries for a command
func (s *DatabaseService) GetContentCount(command string) int {
	count, err := s.getContentCountInternal(command)
	if err != nil {
		return 0
	}
	return count
}
