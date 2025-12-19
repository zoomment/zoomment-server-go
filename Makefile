.PHONY: dev dev-air build run test lint clean docker-build docker-up docker-down

# Development (manual - requires rebuild on changes)
dev:
	go run cmd/server/main.go

# Development with auto-reload (requires: go install github.com/air-verse/air@latest)
dev-air:
	@if command -v air &> /dev/null; then \
		air; \
	elif [ -f "$(HOME)/go/bin/air" ]; then \
		$(HOME)/go/bin/air; \
	else \
		echo "❌ Air not found!"; \
		echo ""; \
		echo "Install Air with:"; \
		echo "  go install github.com/air-verse/air@latest"; \
		echo ""; \
		echo "Then add to your PATH (add to ~/.zshrc):"; \
		echo "  export PATH=\$$PATH:\$$HOME/go/bin"; \
		echo ""; \
		echo "Or run directly:"; \
		echo "  ~/go/bin/air"; \
		exit 1; \
	fi

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

# Install Air for auto-reload
install-air:
	@echo "Installing Air..."
	go install github.com/air-verse/air@latest
	@echo "✅ Air installed!"
	@echo ""
	@echo "Add to your PATH (add to ~/.zshrc):"
	@echo "  export PATH=\$$PATH:\$$HOME/go/bin"
	@echo ""
	@echo "Or run directly: ~/go/bin/air"

# Generate swagger docs (after installing swag)
swagger:
	swag init -g cmd/server/main.go -o docs

