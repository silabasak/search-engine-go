.PHONY: build run test clean mock-server help

# Variables
BINARY_NAME=search-engine-service
MOCK_SERVER=mock-server
BUILD_DIR=build

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(shell git describe --tags --always --dirty)"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the main application
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

build-mock: ## Build the mock server
	@echo "Building $(MOCK_SERVER)..."
	$(GOBUILD) -o $(BUILD_DIR)/$(MOCK_SERVER) ./cmd/mock-server

build-all: build build-mock ## Build all applications

run: build ## Run the main application
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

run-mock: build-mock ## Run the mock server
	@echo "Running $(MOCK_SERVER)..."
	./$(BUILD_DIR)/$(MOCK_SERVER)

dev: ## Run in development mode
	@echo "Running in development mode..."
	$(GOCMD) run ./cmd/server

dev-mock: ## Run mock server in development mode
	@echo "Running mock server in development mode..."
	$(GOCMD) run ./cmd/mock-server

test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -v -cover ./...

test-benchmark: ## Run benchmark tests
	@echo "Running benchmark tests..."
	$(GOTEST) -bench=. ./...

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	$(GOMOD) get -u ./...

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

install: deps build ## Install dependencies and build

setup: ## Setup the project (install dependencies, build, create directories)
	@echo "Setting up the project..."
	mkdir -p $(BUILD_DIR)
	mkdir -p logs
	$(GOMOD) download
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server
	$(GOBUILD) -o $(BUILD_DIR)/$(MOCK_SERVER) ./cmd/mock-server
	@echo "Setup complete!"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(BINARY_NAME)

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

format: ## Format code
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOCMD) vet ./...

# Database commands
db-migrate: ## Run database migrations
	@echo "Running database migrations..."
	$(GOCMD) run ./cmd/server --migrate

db-seed: ## Seed database with sample data
	@echo "Seeding database..."
	$(GOCMD) run ./cmd/server --seed

# Development helpers
logs: ## Show application logs
	@echo "Showing application logs..."
	tail -f logs/app.log

monitor: ## Monitor application resources
	@echo "Monitoring application resources..."
	watch -n 1 'ps aux | grep $(BINARY_NAME)'

# Quick start commands
start: setup run ## Quick start: setup and run

start-all: setup ## Start all services
	@echo "Starting all services..."
	./$(BUILD_DIR)/$(MOCK_SERVER) &
	./$(BUILD_DIR)/$(BINARY_NAME) &
	@echo "All services started!"
	@echo "Main service: http://localhost:8080"
	@echo "Mock server: http://localhost:3001"
	@echo "Dashboard: http://localhost:8080/dashboard"

stop: ## Stop all services
	@echo "Stopping all services..."
	pkill -f $(BINARY_NAME) || true
	pkill -f $(MOCK_SERVER) || true
	@echo "All services stopped!" 