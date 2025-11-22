# Discord Bot

A Discord bot built in Go that returns text content from a PostgreSQL database based on registered commands.

## Features

- **Database-Driven Commands**: Store and retrieve text content from PostgreSQL database
- **Slash Commands with Autocomplete**: Modern Discord slash commands with command autocomplete
- **PostgreSQL Integration**: Fast, reliable content retrieval from your own PostgreSQL database
- **Dynamic Command Discovery**: Automatically discovers commands from database entries
- **Help System**: Built-in help command to list available commands and entry counts
- **Comprehensive Logging**: Structured logging with Zap for command tracking, user metrics, and performance monitoring
- **Clean Architecture**: Modular design with separate packages for config, services, handlers, and bot logic
- **Environment Configuration**: Support for `.env` files and environment variables
- **Graceful Shutdown**: Proper signal handling for clean shutdowns

## Commands

### Slash Commands (Recommended)
- `/command command:<command>` - Returns random text content for the specified command with autocomplete
  - Example: `/command command:wooper`
  - The command parameter will show available options with autocomplete

### Legacy Text Commands
- `!<command>` - Returns random text content for the specified command (e.g., `!wooper`, `!cats`, `!dogs`)
- `!help` or `!list` - Shows all available commands and entry counts

## Database Setup

The bot uses a PostgreSQL database to store commands and their associated content. The database table is automatically created on first run.

### Database Schema

The bot automatically creates a `commands` table with the following structure:

```sql
CREATE TABLE commands (
    id SERIAL PRIMARY KEY,
    command VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_command ON commands(command);
```

### Adding Commands and Content

To add a new command with content, insert a row into the `commands` table:

```sql
INSERT INTO commands (command, content) VALUES ('wooper', 'This is some wooper content!');
INSERT INTO commands (command, content) VALUES ('wooper', 'Another wooper message');
INSERT INTO commands (command, content) VALUES ('cats', 'Meow meow!');
```

**Notes:**
- Multiple entries with the same `command` value are allowed - the bot will randomly select one
- The `content` field can contain any text (Discord markdown is supported)
- Commands are case-sensitive and should match exactly what users type (e.g., `!wooper` matches command `wooper`)

### Example Database Setup

```sql
-- Create database (if not exists)
CREATE DATABASE mutsumi_bot;

-- Connect to the database
\c mutsumi_bot

-- The bot will create the table automatically, or you can create it manually:
CREATE TABLE IF NOT EXISTS commands (
    id SERIAL PRIMARY KEY,
    command VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_command ON commands(command);

-- Add some example commands
INSERT INTO commands (command, content) VALUES
    ('wooper', '🐸 Wooper is the best!'),
    ('wooper', 'Wooper wooper wooper!'),
    ('help', 'Use !help to see available commands'),
    ('cats', '🐱 Cats are cute!');
```

## Logging

The bot includes comprehensive structured logging using Zap. Logs include:

- **Command Tracking**: Every command execution with user info, timing, and results
- **Performance Metrics**: Response times for database queries and content retrieval
- **User Analytics**: Who is using which commands and when
- **Error Tracking**: Detailed error information for debugging
- **System Events**: Bot startup, database service initialization, etc.

### Log Levels

Set the log level using the `LOG_LEVEL` environment variable:
- `debug` - All messages including debug info
- `info` - General information (default)
- `warn` - Warnings and errors
- `error` - Errors only

### Example Log Output

```json
{
  "timestamp": "2024-01-20T10:30:45Z",
  "level": "info",
  "message": "Command received",
  "command": "!wooper",
  "category": "wooper",
  "user": "username",
  "user_id": "123456789",
  "channel_id": "987654321",
  "guild_id": "111222333"
}
```

### Adding New Commands

To add a new command (e.g., `dogs`):

1. Connect to your PostgreSQL database
2. Insert content for the command:
   ```sql
   INSERT INTO commands (command, content) VALUES ('dogs', 'Woof woof!');
   ```
3. The bot will automatically detect the new command (no restart needed for new entries)
4. Use the new command: `!dogs`

The bot dynamically discovers commands from the database, so new entries are available immediately.

## Setup

### Prerequisites

- Go 1.19 or later (for local development)
- Docker and Docker Compose (for containerized deployment)
- A Discord bot token
- PostgreSQL database (hosted on your VPS or elsewhere)

### Installation

#### Option 1: Docker Deployment with Dokploy (Recommended)

1. **In Dokploy dashboard:**
   - Create a new project
   - Connect your Git repository
   - Set the build context to the root directory
   - Use the provided `Dockerfile`

2. **Environment Variables in Dokploy:**
   - `DISCORD_BOT_TOKEN`: Your Discord bot token
   - `DATABASE_CONNECTION`: PostgreSQL connection string (e.g., `postgres://user:password@host:port/database`)
   - `LOG_LEVEL`: Optional, defaults to `info`

3. **Deploy:**
   - Dokploy will build and deploy your container from the Dockerfile
   - The bot will be available and running on your VPS

4. **Updating:**
   - Push changes to your repository
   - Dokploy can be configured to auto-rebuild or you can manually trigger builds

#### Option 2: Local Development

1. **Clone the repository**
   ```bash
   git clone git@github.com:hlaclau/mutsumi-discord-bot.git
   cd mutsumi-bot
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Configure environment**

   **Option A: Using .env file (recommended)**
   ```bash
   cp .env.example .env
   # Edit .env and add your Discord bot token and database connection string
   ```

   **Option B: Using environment variables**
   ```bash
   export DISCORD_BOT_TOKEN=your_bot_token_here
   export DATABASE_CONNECTION=postgres://user:password@host:port/database
   ```

4. **Run the bot**
   ```bash
   make run
   ```

   Or use the traditional Go command:
   ```bash
   go run ./...
   ```

## Creating a Discord Bot

1. Go to the [Discord Developer Portal](https://discord.com/developers/applications)
2. Click "New Application" and give it a name
3. Go to the "Bot" section
4. Click "Add Bot"
5. Copy the bot token (this is your `DISCORD_BOT_TOKEN`)
6. Under "Privileged Gateway Intents", enable "Message Content Intent"
7. Go to the "OAuth2" > "URL Generator" section
8. Select "bot" scope and "Send Messages" permission
9. Use the generated URL to invite your bot to a server

## Project Structure

```
mutsumi-bot/
├── .env.example          # Environment template
├── .env                  # Your environment variables (gitignored)
├── go.mod               # Go module file
├── go.sum               # Go module checksums
├── main.go              # Application entry point
├── README.md            # This file
├── internal/            # Internal packages
│   ├── bot/             # Discord bot wrapper
│   │   ├── bot.go
│   │   └── bot_test.go
│   ├── config/          # Configuration management
│   │   ├── config.go
│   │   └── config_test.go
│   ├── handlers/        # Message event handlers
│   │   ├── messages.go
│   │   ├── messages_test.go
│   │   ├── interactions.go
│   │   ├── mock_service.go
│   │   └── interactions_test.go
│   ├── logger/          # Structured logging with Zap
│   │   ├── logger.go
│   │   └── logger_test.go
│   └── services/        # Business logic services
│       ├── database.go  # PostgreSQL database service
│       └── service.go   # Content service interface
└── tests/               # Test files
    └── integration/     # Integration tests
        └── integration_test.go
```

## Architecture

The bot follows a clean, layered architecture:

- **`internal/config`**: Environment variable loading with `.env` support
- **`internal/logger`**: Structured logging configuration and initialization
- **`internal/services`**: Business logic for database content management and command discovery
  - **`database.go`**: PostgreSQL database service for storing and retrieving command content
  - **`service.go`**: Content service interface for abstraction
- **`internal/handlers`**: Discord message event processing and slash command interactions with dynamic command support and comprehensive logging
- **`internal/bot`**: Discord session management and lifecycle
- **`main.go`**: Dependency injection and application startup

## Dependencies

- **discordgo**: Discord API client for Go
- **godotenv**: Environment variable loading from `.env` files
- **zap**: High-performance structured logging
- **pgx/v5**: PostgreSQL driver for Go

## Development

### Makefile Commands

The project includes a simple Makefile for common development tasks:

```bash
make help              # Show available commands
make run               # Run the bot
make build             # Build the binary
make test              # Run all tests (unit + integration)
make test-unit         # Run unit tests only
make test-integration  # Run integration tests only
make fmt               # Format code
make clean             # Clean build artifacts
```

### Building

```bash
make build
# or
go build -o mutsumi-bot .
```

### Running Tests

```bash
make test              # Run all tests (unit + integration)
make test-unit         # Run unit tests only
make test-integration  # Run integration tests only
# or
go test ./...                    # All tests
go test ./internal/...          # Unit tests only
go test ./tests/integration/... # Integration tests only
```

### Code Quality

The project follows Go best practices:
- Clean architecture with separated concerns
- Interface-based design for testability
- Proper error handling and context usage
- Graceful shutdown handling

## Troubleshooting

### Bot doesn't respond to commands

1. Ensure the bot has "Message Content Intent" enabled in Discord Developer Portal
2. Check that the bot has "Send Messages" permission in the server
3. Verify the bot token is correct in your `.env` file
4. Verify the database connection string is correct in your `.env` file
5. Use `!help` to see available commands

### No content found for a command

- Check that the command exists in the database:
  ```sql
  SELECT * FROM commands WHERE command = 'yourcommand';
  ```
- Verify the database connection is working
- Check database logs for connection errors
- Ensure the bot has read permissions on the database

### Bot shows "no commands available"

- Verify the database connection string is correct
- Check that the `commands` table exists in your database
- Ensure there are entries in the `commands` table:
  ```sql
  SELECT COUNT(*) FROM commands;
  ```
- Check bot logs for database connection errors

### Database connection errors

- Verify your PostgreSQL database is running and accessible
- Check that the connection string format is correct: `postgres://user:password@host:port/database`
- Ensure the database user has permissions to create tables and read/write data
- Check firewall rules if connecting to a remote database

## License

This project is open source and available under the MIT License.
