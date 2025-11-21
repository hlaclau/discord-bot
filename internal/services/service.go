package services

// ContentService defines the interface for services that provide content retrieval
type ContentService interface {
	// GetRandomContent returns a random content string for the given command
	GetRandomContent(command string) string

	// GetContentCount returns the number of content entries for a command
	GetContentCount(command string) int

	// GetAvailableCategories returns all available commands
	GetAvailableCategories() []string

	// HasCategory checks if a command exists
	HasCategory(command string) bool
}
