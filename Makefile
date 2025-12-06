.PHONY: dev build run test lint clean docker-build docker-up docker-down

# Development
dev:
	go run cmd/server/main.go

# Build binary
build:
	go build -o bin/server cmd/server/main.go

# Run built binary
run: build
	./bin/server

# Run tests
test:
	go test ./... -v

# Run tests with coverage
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Lint code
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Docker commands
docker-build:
	docker build -t zoomment-server-go .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f zoomment-server

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Generate swagger docs (after installing swag)
swagger:
	swag init -g cmd/server/main.go -o docs

