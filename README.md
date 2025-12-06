# Zoomment Server (Go)

> ⚠️ **BETA VERSION** - This project is currently in beta and not ready for production use. Use at your own risk.

An open-source, self-hosted comment system API built with Go. This is a Go port of the original [Node.js Zoomment server](https://github.com/zoomment/zoomment-server).

## Features

- 💬 Threaded comments with replies
- 👍 Emoji reactions
- 🔐 Passwordless authentication (magic links)
- 📧 Email notifications
- 🌐 Multi-site support
- 🚀 High performance (Go)

## Requirements

- [Go](https://golang.org/) 1.22+
- [MongoDB](https://www.mongodb.com/) 4.4+

## Project Structure

```
zoomment-server-go/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Environment configuration
│   ├── database/
│   │   └── mongodb.go           # MongoDB connection
│   ├── models/
│   │   ├── comment.go           # Comment model
│   │   ├── user.go              # User model
│   │   ├── site.go              # Site model
│   │   └── reaction.go          # Reaction model
│   ├── handlers/
│   │   ├── comments.go          # Comment endpoints
│   │   ├── users.go             # User/auth endpoints
│   │   ├── sites.go             # Site management endpoints
│   │   └── reactions.go         # Reaction endpoints
│   ├── middleware/
│   │   └── auth.go              # JWT authentication & authorization
│   ├── routes/
│   │   └── routes.go            # API route definitions
│   ├── repository/
│   │   └── comments.go          # Database queries (with aggregation)
│   ├── services/
│   │   ├── mailer/              # Email service (SMTP)
│   │   └── metadata/            # HTML scraper for site verification
│   ├── validators/
│   │   └── validators.go        # Request validation
│   ├── errors/
│   │   └── errors.go            # Structured error handling
│   ├── logger/
│   │   └── logger.go            # Structured logging (zerolog)
│   └── utils/
│       ├── email.go             # Email utilities
│       ├── name.go              # Name sanitization
│       ├── gravatar.go          # Gravatar hash generation
│       ├── secret.go            # Secret token generation
│       └── sanitizer.go         # XSS protection
├── docs/
│   └── docs.go                  # Swagger API documentation
├── Dockerfile                   # Multi-stage Docker build
├── docker-compose.yml           # Docker Compose configuration
├── Makefile                     # Build commands
├── go.mod                       # Go module dependencies
└── env.example                  # Example environment variables
```

## Quick Start

### 1. Clone and setup

```bash
git clone https://github.com/zoomment/zoomment-server-go.git
cd zoomment-server-go

# Copy environment file
cp env.example .env

# Edit .env with your settings
vim .env
```

### 2. Configure environment variables

```env
# Server
PORT=8080
MONGODB_URI=mongodb://localhost:27017/zoomment

# Authentication
JWT_SECRET=your-super-secret-key-change-this

# Dashboard
DASHBOARD_URL=http://localhost:3000
BRAND_NAME=Zoomment

# Email (SMTP) - Optional
BOT_EMAIL_ADDR=your-email@gmail.com
BOT_EMAIL_PASS=your-app-password
BOT_EMAIL_HOST=smtp.gmail.com
BOT_EMAIL_PORT=465
```

### 3. Run the server

#### Development (with hot reload)

```bash
go run cmd/server/main.go
```

#### Using Make

```bash
make dev      # Run in development mode
make build    # Build binary
make run      # Build and run
make test     # Run tests
```

#### Using Docker

```bash
make docker-build   # Build Docker image
make docker-up      # Start with docker-compose
make docker-down    # Stop containers
make docker-logs    # View logs
```

## API Endpoints

| Method        | Endpoint                      | Auth  | Description              |
| ------------- | ----------------------------- | ----- | ------------------------ |
| GET           | `/health`                     | -     | Health check             |
| GET           | `/swagger/*`                  | -     | API documentation        |
| **Comments**  |
| GET           | `/api/comments?pageId=xxx`    | -     | List comments for a page |
| POST          | `/api/comments`               | -     | Add a comment            |
| DELETE        | `/api/comments/:id`           | ✓     | Delete a comment         |
| GET           | `/api/comments/sites/:siteId` | Admin | List comments by site    |
| **Users**     |
| POST          | `/api/users/auth`             | -     | Request magic link       |
| GET           | `/api/users/profile`          | ✓     | Get user profile         |
| DELETE        | `/api/users`                  | ✓     | Delete account           |
| **Sites**     |
| GET           | `/api/sites`                  | Admin | List user's sites        |
| POST          | `/api/sites`                  | Admin | Register a site          |
| DELETE        | `/api/sites/:id`              | Admin | Remove a site            |
| **Reactions** |
| GET           | `/api/reactions?pageId=xxx`   | -     | Get reactions            |
| POST          | `/api/reactions`              | -     | Add/toggle reaction      |

## API Documentation

Swagger UI is available at: `http://localhost:8080/swagger/index.html`

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

## Production Deployment

> ⚠️ This project is in **BETA**. Review the code before deploying to production.

### Docker (Recommended)

```bash
# Build production image
docker build -t zoomment-server-go .

# Run with docker-compose
docker-compose up -d
```

## Contributing

Contributions are welcome! Please read the contributing guidelines first.

## License

MIT License - see [LICENSE](LICENSE) for details.

---

**Note:** This is a beta version. Please report any issues on GitHub.
