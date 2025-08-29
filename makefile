# Variables
BINARY_NAME=go-currency-api
DOCKER_COMPOSE_FILE=docker-compose.yml
DB_URL=postgresql://currency_user:currency_pass@localhost:5432/currency_db?sslmode=disable

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
NC=\033[0m # No Color

.PHONY: help build run test clean docker-up docker-down migrate-up migrate-down

# Default target
help:
	@echo "Available targets:"
	@echo "  $(GREEN)build$(NC)          - Build the application"
	@echo "  $(GREEN)run$(NC)            - Run the application"
	@echo "  $(GREEN)test$(NC)           - Run tests"
	@echo "  $(GREEN)test-coverage$(NC)  - Run tests with coverage"
	@echo "  $(GREEN)clean$(NC)          - Clean build artifacts"
	@echo "  $(GREEN)docker-up$(NC)      - Start Docker services"
	@echo "  $(GREEN)docker-down$(NC)    - Stop Docker services"
	@echo "  $(GREEN)docker-logs$(NC)    - View Docker logs"
	@echo "  $(GREEN)migrate-up$(NC)     - Run database migrations"
	@echo "  $(GREEN)migrate-down$(NC)   - Rollback database migrations"
	@echo "  $(GREEN)dev$(NC)            - Start development environment"
	  @echo "  $(GREEN)docker-compose-check$(NC) - Check Docker Compose version"

# Build the application
build:
	@echo "$(YELLOW)Building $(BINARY_NAME)...$(NC)"
	go build -o bin/$(BINARY_NAME) cmd/api/main.go
	@echo "$(GREEN)Build completed!$(NC)"

# Run the application
run:
	@echo "$(YELLOW)Running $(BINARY_NAME)...$(NC)"
	go run cmd/api/main.go

# Run tests
test:
	@echo "$(YELLOW)Running tests...$(NC)"
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "$(YELLOW)Running tests with coverage...$(NC)"
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

# Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning...$(NC)"
	go clean
	rm -f bin/$(BINARY_NAME)
	rm -f coverage.out coverage.html
	@echo "$(GREEN)Clean completed!$(NC)"

# Start Docker services
docker-up:
	@echo "$(YELLOW)Starting Docker services...$(NC)"
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose -f $(DOCKER_COMPOSE_FILE) up -d; \
	else \
		docker compose -f $(DOCKER_COMPOSE_FILE) up -d; \
	fi
	@echo "$(GREEN)Docker services started!$(NC)"
	@echo "$(YELLOW)PostgreSQL:$(NC) localhost:5432"
	@echo "$(YELLOW)Redis:$(NC) localhost:6379"

# Start Docker services with tools (pgAdmin, Redis Commander)
docker-up-tools:
	@echo "$(YELLOW)Starting Docker services with tools...$(NC)"
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose -f $(DOCKER_COMPOSE_FILE) --profile tools up -d; \
	else \
		docker compose -f $(DOCKER_COMPOSE_FILE) --profile tools up -d; \
	fi
	@echo "$(GREEN)Docker services with tools started!$(NC)"
	@echo "$(YELLOW)PostgreSQL:$(NC) localhost:5432"
	@echo "$(YELLOW)Redis:$(NC) localhost:6379"
	@echo "$(YELLOW)pgAdmin:$(NC) http://localhost:5050 (admin@admin.com / admin)"
	@echo "$(YELLOW)Redis Commander:$(NC) http://localhost:8081"

# Stop Docker services
docker-down:
	@echo "$(YELLOW)Stopping Docker services...$(NC)"
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose -f $(DOCKER_COMPOSE_FILE) down; \
	else \
		docker compose -f $(DOCKER_COMPOSE_FILE) down; \
	fi
	@echo "$(GREEN)Docker services stopped!$(NC)"

# View Docker logs
docker-logs:
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f; \
	else \
		docker compose -f $(DOCKER_COMPOSE_FILE) logs -f; \
	fi

# Install migrate tool if not present
install-migrate:
	@if ! command -v migrate >/dev/null 2>&1; then \
		echo "$(YELLOW)Installing golang-migrate...$(NC)"; \
		OS=$(uname -s | tr '[:upper:]' '[:lower:]'); \
		ARCH=$(uname -m); \
		case $ARCH in \
			x86_64) ARCH="amd64" ;; \
			arm64|aarch64) ARCH="arm64" ;; \
			*) echo "$(RED)Unsupported architecture: $ARCH$(NC)"; exit 1 ;; \
		esac; \
		case $OS in \
			darwin) OS="darwin" ;; \
			linux) OS="linux" ;; \
			*) echo "$(RED)Unsupported OS: $OS$(NC)"; exit 1 ;; \
		esac; \
		MIGRATE_VERSION="v4.16.2"; \
		DOWNLOAD_URL="https://github.com/golang-migrate/migrate/releases/download/$MIGRATE_VERSION/migrate.$OS-$ARCH.tar.gz"; \
		echo "$(YELLOW)Downloading migrate for $OS-$ARCH...$(NC)"; \
		curl -L $DOWNLOAD_URL | tar xvz; \
		chmod +x migrate; \
		if [ -w /usr/local/bin ]; then \
			mv migrate /usr/local/bin/; \
		else \
			echo "$(YELLOW)Moving migrate to /usr/local/bin (requires sudo)...$(NC)"; \
			sudo mv migrate /usr/local/bin/; \
		fi; \
		echo "$(GREEN)golang-migrate installed successfully$(NC)"; \
	else \
		echo "$(GREEN)golang-migrate is already installed$(NC)"; \
	fi

# Run database migrations up
migrate-up: install-migrate
	@echo "$(YELLOW)Running database migrations...$(NC)"
	migrate -path migrations -database "$(DB_URL)" -verbose up
	@echo "$(GREEN)Migrations completed!$(NC)"

# Rollback database migrations
migrate-down: install-migrate
	@echo "$(YELLOW)Rolling back database migrations...$(NC)"
	migrate -path migrations -database "$(DB_URL)" -verbose down 1
	@echo "$(GREEN)Rollback completed!$(NC)"

# Alternative: Run migrations using Go (without installing migrate binary)
migrate-go-up:
	@echo "$(YELLOW)Running database migrations using Go...$(NC)"
	@if [ ! -d "migrations" ]; then \
		echo "$(YELLOW)Creating migrations directory...$(NC)"; \
		mkdir -p migrations; \
	fi
	@if [ -z "$(ls -A migrations 2>/dev/null)" ]; then \
		echo "$(YELLOW)No migration files found. Skipping migrations...$(NC)"; \
	else \
		go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
			-path migrations \
			-database "$(DB_URL)" \
			-verbose up; \
	fi
	@echo "$(GREEN)Migrations completed!$(NC)"

# Alternative: Rollback migrations using Go
migrate-go-down:
	@echo "$(YELLOW)Rolling back database migrations using Go...$(NC)"
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
		-path migrations \
		-database "$(DB_URL)" \
		-verbose down 1
	@echo "$(GREEN)Rollback completed!$(NC)"

# Create a new migration using Go
migrate-go-create:
	@read -p "Enter migration name: " name; \
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
		create -ext sql -dir migrations $name
	@echo "$(GREEN)Migration files created!$(NC)"

# Start development environment
dev: docker-up
	@echo "$(YELLOW)Waiting for services to be ready...$(NC)"
	sleep 5
	@echo "$(YELLOW)Running migrations...$(NC)"
	make migrate-go-up
	@echo "$(GREEN)Development environment ready!$(NC)"
	@echo "$(YELLOW)Starting application...$(NC)"
	make run

# Alternative dev command using binary migrate (if you prefer)
dev-binary: docker-up
	@echo "$(YELLOW)Waiting for services to be ready...$(NC)"
	sleep 5
	@echo "$(YELLOW)Running migrations...$(NC)"
	make migrate-up
	@echo "$(GREEN)Development environment ready!$(NC)"
	@echo "$(YELLOW)Starting application...$(NC)"
	make run

# Install dependencies
deps:
	@echo "$(YELLOW)Installing dependencies...$(NC)"
	go mod download
	go mod tidy
	@echo "$(GREEN)Dependencies installed!$(NC)"

# Check Docker Compose version and availability
docker-compose-check:
	@echo "$(YELLOW)Checking Docker Compose availability...$(NC)"
	@if command -v docker-compose >/dev/null 2>&1; then \
		echo "$(GREEN)Legacy docker-compose found:$(NC)"; \
		docker-compose --version; \
	elif docker compose version >/dev/null 2>&1; then \
		echo "$(GREEN)Modern docker compose found:$(NC)"; \
		docker compose version; \
	else \
		echo "$(RED)Neither docker-compose nor docker compose found!$(NC)"; \
		echo "$(YELLOW)Please install Docker Desktop or Docker Compose$(NC)"; \
		exit 1; \
	fi

# Run linter (requires golangci-lint)
lint:
	@echo "$(YELLOW)Running linter...$(NC)"
	golangci-lint run
	@echo "$(GREEN)Linting completed!$(NC)"

# Format code
fmt:
	@echo "$(YELLOW)Formatting code...$(NC)"
	go fmt ./...
	@echo "$(GREEN)Code formatted!$(NC)"

# Run security check
security:
	@echo "$(YELLOW)Running security check...$(NC)"
	gosec ./...
	@echo "$(GREEN)Security check completed!$(NC)"