# Notification Service

A Go-based notification service that allows you to send messages to different messaging platforms through a unified API.

## Features

- Send notifications to multiple messaging platforms (currently DingTalk and Telegram)
- Manage multiple bot configurations for each platform
- Monitor message delivery status and history
- Simple API for sending notifications
- Web interface for configuration and management
- Support for multiple database types (SQLite, MySQL, PostgreSQL)

## System Requirements

- Go 1.20 or later
- Web browser (Chrome, Firefox, Safari, etc.)
- Database (SQLite, MySQL, or PostgreSQL)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/username/yuyan.git
cd yuyan
```

2. Download dependencies:

```bash
go mod download
```

3. Build the application:

```bash
go build -o yuyan cmd/server/main.go
```

## Configuration

The application uses a YAML configuration file located at `config/config.yaml`. You can modify this file directly or use the web interface to update settings.

Default configuration:

```yaml
server:
  port: 8080
  mode: debug

database:
  type: sqlite # sqlite, mysql, postgres
  sqlite:
    path: ./data/yuyan.db
  mysql:
    host: localhost
    port: 3306
    username: root
    password: password
    dbname: yuyan
    params: charset=utf8mb4&parseTime=True&loc=Local
  postgres:
    host: localhost
    port: 5432
    username: postgres
    password: password
    dbname: yuyan
    sslmode: disable
```

## Usage

### Running the service

To start the service, run:

```bash
./yuyan
```

The web interface will be available at `http://localhost:8080`.

### API endpoints

The service provides the following API endpoints:

#### Bot Management

- `GET /api/bots` - Get all bot configurations
- `GET /api/bots/:id` - Get bot configuration by ID
- `POST /api/bots` - Create a new bot configuration
- `PUT /api/bots/:id` - Update an existing bot configuration
- `DELETE /api/bots/:id` - Delete a bot configuration

#### Message Management

- `GET /api/messages` - Get all messages (with optional filtering)
- `GET /api/messages/:id` - Get message by ID
- `POST /api/messages` - Send a new message

#### System Configuration

- `GET /api/config` - Get current system configuration
- `PUT /api/config` - Update system configuration

### Sending a notification

To send a notification, use the `/api/messages` endpoint:

Basic example:
```bash
curl -X POST \
  http://localhost:8090/api/messages \
  -H 'Content-Type: application/json' \
  -d '{
    "bot_id": 4,
    "content": "Hello, this is a test message!",
    "format": "text"
  }'
```

Telegram specific examples with @mentions:

1. Message with user mentions:
```bash
curl -X POST \
  http://localhost:8090/api/messages \
  -H 'Content-Type: application/json' \
  -d '{
    "bot_id": 4,
    "content": "Hey @wolf1688, please check this! cc: @learnjie",
    "format": "html"
  }'
```

2. Message with @all mention:
```bash
curl -X POST \
  http://localhost:8090/api/messages \
  -H 'Content-Type: application/json' \
  -d '{
    "bot_id": 4,
    "content": "Important announcement @all: Team meeting at 3 PM!",
    "format": "html"
  }'
```

3. Message with phone number mention:
```bash
curl -X POST \
  http://localhost:8090/api/messages \
  -H 'Content-Type: application/json' \
  -d '{
    "bot_id": 4,
    "content": "Please contact @1365754902 for support.",
    "format": "html"
  }'
```

Note: When using Telegram bot with mentions:
- User mentions (@username) will be converted to clickable links
- @all will be replaced with a highlighted announcement message
- Phone number mentions will be highlighted but not clickable
- Set "format": "html" when using mentions to ensure proper formatting

## Bot Configuration

### DingTalk

To configure a DingTalk bot, you'll need:

1. A webhook URL from the DingTalk custom robot configuration
2. (Optional) A secret for signing requests

### Telegram

To configure a Telegram bot, you'll need:

1. A bot token from the BotFather
2. A chat ID where messages will be sent

## Development

### Project structure

```
yuyan/
├── cmd/
│   └── server/
│       └── main.go         # Entry point
├── config/
│   └── config.yaml         # Configuration file
├── internal/
│   ├── api/                # API handlers
│   ├── bot/                # Bot implementations
│   ├── database/           # Database access
│   ├── models/             # Data models
│   ├── service/            # Business logic
│   └── utils/              # Utility functions
├── web/
│   ├── static/             # Static files (CSS, JS)
│   └── templates/          # HTML templates
└── go.mod                  # Go module file
```

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 