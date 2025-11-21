package handlers

import "wooper-bot/internal/services"

// mockContentService is a mock implementation of ContentService for testing
type mockContentService struct {
	commands map[string][]string // command -> list of content entries
}

func newMockContentService() *mockContentService {
	return &mockContentService{
		commands: make(map[string][]string),
	}
}

func (m *mockContentService) addCommand(command string, content ...string) {
	m.commands[command] = content
}

func (m *mockContentService) GetRandomContent(command string) string {
	contents, exists := m.commands[command]
	if !exists || len(contents) == 0 {
		return ""
	}
	// Return first content for deterministic testing
	return contents[0]
}

func (m *mockContentService) GetContentCount(command string) int {
	if contents, exists := m.commands[command]; exists {
		return len(contents)
	}
	return 0
}

func (m *mockContentService) GetAvailableCategories() []string {
	var commands []string
	for cmd := range m.commands {
		commands = append(commands, cmd)
	}
	return commands
}

func (m *mockContentService) HasCategory(command string) bool {
	_, exists := m.commands[command]
	return exists && len(m.commands[command]) > 0
}

// Ensure mockContentService implements ContentService
var _ services.ContentService = (*mockContentService)(nil)
