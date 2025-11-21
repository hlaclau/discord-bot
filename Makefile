# Wooper Bot Makefile

.PHONY: run build test test-unit test-integration fmt clean help docker-build

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

run: ## Run the bot
	go run ./...

build: ## Build the binary
	go build -o wooper-bot .

test: ## Run all tests (unit + integration)
	go test ./...

test-unit: ## Run unit tests only
	go test ./internal/...

test-integration: ## Run integration tests only
	go test ./tests/integration/...

fmt: ## Format code
	go fmt ./...

clean: ## Clean build artifacts
	go clean
	rm -f wooper-bot

# Docker commands
docker-build: ## Build Docker image
	docker build -t wooper-bot .
