# Zoomment Server (Go)

> ⚠️ **BETA VERSION** - This project is currently in beta and not ready for production use. Use at your own risk.

An open-source, self-hosted comment system API built with Go. This is a Go port of the original [Node.js Zoomment server](https://github.com/zoomment/zoomment-server).

## Features

- 💬 **Threaded comments** with nested replies
- 👍 **Emoji reactions** for quick feedback
- 🔐 **Passwordless authentication** via magic links
- 📧 **Email notifications** for comments and verification
- 🌐 **Multi-site support** with domain verification
- 🚀 **High performance** built with Go
- 🛡️ **XSS protection** with HTML sanitization
- 📚 **Swagger API documentation** included
- 🔄 **Auto-reload** development workflow

## Requirements

- [Go](https://golang.org/) 1.22+
- [MongoDB](https://www.mongodb.com/) 4.4+

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

Edit `.env` file with your settings:

```env
# Server
PORT=8080
MONGODB_URI=mongodb://localhost:27017/zoomment

# Authentication
JWT_SECRET=your-super-secret-jwt-key-change-this

# Dashboard
DASHBOARD_URL=http://localhost:3000
BRAND_NAME=Zoomment

# Email (SMTP) - Optional but recommended
BOT_EMAIL_ADDR=your-email@gmail.com
BOT_EMAIL_PASS=your-app-password
BOT_EMAIL_HOST=smtp.gmail.com
BOT_EMAIL_PORT=465

# Admin
ADMIN_EMAIL_ADDR=admin@example.com
```

> 💡 **Tip**: For Gmail, use an [App Password](https://support.google.com/accounts/answer/185833) instead of your regular password.

### 3. Run the server

#### Development (Recommended - with auto-reload)

First, install Air for automatic rebuilds:

```bash
# Install Air (one-time)
make install-air
# OR
go install github.com/air-verse/air@latest

# Add to PATH (add to ~/.zshrc)
export PATH=$PATH:$HOME/go/bin
```

Then run with auto-reload:

```bash
make dev-air    # Auto-rebuilds on code changes
```

#### Development (Manual)

```bash
make dev        # Run with go run (rebuilds each time)
go run cmd/server/main.go
```

#### Production Build

```bash
make build      # Build optimized binary
make run        # Run the built binary
```

#### All Make Commands

```bash
make dev        # Run in development mode
make dev-air    # Run with auto-reload (requires Air)
make build      # Build binary
make run        # Build and run
make test       # Run tests
make test-coverage  # Run tests with coverage
make lint       # Run linter
make clean      # Clean build artifacts
make install-air    # Install Air for auto-reload
```

> 📖 **For detailed development workflow, see [DEVELOPMENT.md](DEVELOPMENT.md)**

#### Using Docker

```bash
make docker-build   # Build Docker image
make docker-up      # Start with docker-compose
make docker-down    # Stop containers
make docker-logs    # View logs
```

## API Endpoints

### Health & Documentation

| Method | Endpoint           | Auth | Description         |
|--------|-------------------|------|---------------------|
| GET    | `/health`         | -    | Health check        |
| GET    | `/swagger/*`       | -    | Swagger UI docs     |

### Comments

| Method | Endpoint                          | Auth  | Description                    |
|--------|-----------------------------------|-------|--------------------------------|
| GET    | `/api/comments?pageId=xxx`        | -     | List comments for a page       |
| GET    | `/api/comments?domain=xxx`         | -     | List comments for a domain     |
| POST   | `/api/comments`                   | -     | Add a comment                   |
| DELETE | `/api/comments/:id?secret=xxx`     | ✓     | Delete a comment (auth/secret) |
| GET    | `/api/comments/sites/:siteId`      | Admin | List all comments for a site   |

### Users

| Method | Endpoint              | Auth | Description          |
|--------|----------------------|------|----------------------|
| POST   | `/api/users/auth`     | -    | Request magic link   |
| GET    | `/api/users/profile`  | ✓    | Get user profile     |
| DELETE | `/api/users`          | ✓    | Delete account        |

### Sites

| Method | Endpoint         | Auth  | Description       |
|--------|------------------|-------|-------------------|
| GET    | `/api/sites`      | Admin | List user's sites |
| POST   | `/api/sites`      | Admin | Register a site  |
| DELETE | `/api/sites/:id`  | Admin | Remove a site     |

### Reactions

| Method | Endpoint                    | Auth | Description              |
|--------|----------------------------|------|--------------------------|
| GET    | `/api/reactions?pageId=xxx` | -    | Get reactions for page   |
| POST   | `/api/reactions`            | -    | Add/toggle reaction      |

> 📖 **Full API documentation**: Visit `http://localhost:8080/swagger/index.html` when server is running

## API Documentation

Interactive Swagger UI is available at: `http://localhost:8080/swagger/index.html`

The API documentation is auto-generated from code comments. To regenerate:

```bash
make swagger
```

## Development

### Auto-reload Development

For the best development experience, use Air for automatic rebuilds:

```bash
# Install Air (one-time)
make install-air

# Run with auto-reload
make dev-air
```

Air will automatically:
- Watch for `.go` file changes
- Rebuild the application
- Restart the server

### Code Structure

- **Handlers**: HTTP request handlers (REST endpoints)
- **Models**: Database models (MongoDB documents)
- **Repository**: Database queries and aggregations
- **Middleware**: Authentication, authorization, logging
- **Services**: Business logic (email, metadata scraping)
- **Utils**: Helper functions (sanitization, validation)
- **Constants**: Application-wide constants

### Code Quality

```bash
# Run linter
make lint

# Run tests
make test

# Run tests with coverage
make test-coverage
```

> 📖 **See [DEVELOPMENT.md](DEVELOPMENT.md) for detailed development workflow**

## Testing

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage
# Opens coverage.html in browser
```

## Architecture

### Key Design Decisions

- **Repository Pattern**: Database queries separated from handlers
- **Structured Errors**: Consistent error responses across API
- **Middleware Chain**: Authentication, CORS, logging, recovery
- **Aggregation Queries**: Efficient nested comment fetching (no N+1 problem)
- **XSS Protection**: HTML sanitization for user input
- **Type Safety**: Strong typing with Go structs

### Performance Optimizations

- MongoDB aggregation pipelines for efficient queries
- Single query for comments with nested replies
- Async email sending (non-blocking)
- Connection pooling for MongoDB

## Production Deployment

> ⚠️ This project is in **BETA**. Review the code before deploying to production.

### Build for Production

```bash
# Build optimized binary
make build

# Binary will be in bin/server
./bin/server
```

### Docker Deployment (Recommended)

```bash
# Build production image
make docker-build
# OR
docker build -t zoomment-server-go .

# Run with docker-compose
make docker-up
# OR
docker-compose up -d

# View logs
make docker-logs
# OR
docker-compose logs -f
```

### Environment Variables for Production

Make sure to set:
- Strong `JWT_SECRET` (use a secure random string)
- Production `MONGODB_URI`
- Valid SMTP credentials for emails
- Correct `DASHBOARD_URL` for your frontend

## Troubleshooting

### Server won't start
- Check MongoDB is running: `mongosh` or check connection string
- Verify environment variables in `.env`
- Check port 8080 is not in use

### Air not found
```bash
# Install Air
make install-air

# Add to PATH
export PATH=$PATH:$HOME/go/bin
```

### CORS errors
- CORS is configured in `cmd/server/main.go`
- Make sure frontend origin is allowed
- Check `fingerprint` and `token` headers are in allowed headers

### Email not sending
- Verify SMTP credentials in `.env`
- For Gmail, use App Password (not regular password)
- Check firewall/network allows SMTP connections

## Contributing

Contributions are welcome! 

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

Please ensure:
- Code follows Go conventions
- Tests pass: `make test`
- Linter passes: `make lint`
- Documentation is updated

## License

MIT License - see [LICENSE](LICENSE) for details.

## Related Documentation

- [DEVELOPMENT.md](DEVELOPMENT.md) - Development workflow guide

## License

MIT License - see [LICENSE](LICENSE) for details.

---

**Note:** This is a beta version. Please report any issues on GitHub.

Made with ❤️ using Go
