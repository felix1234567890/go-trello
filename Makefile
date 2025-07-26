# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Test parameters
TEST_FLAGS=-v -race -coverprofile=coverage.out
E2E_FLAGS=-v -race
UNIT_FLAGS=-v -race -short

# Binary name
BINARY_NAME=go-trello
BINARY_UNIX=$(BINARY_NAME)_unix

.PHONY: all build clean test test-unit test-e2e test-coverage test-integration deps help

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

all: deps test build ## Download dependencies, run tests, and build

build: ## Build the binary
	$(GOBUILD) -o $(BINARY_NAME) -v ./main.go

clean: ## Remove binary and clean cache
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) tidy

test: ## Run all tests
	$(GOTEST) $(TEST_FLAGS) ./...

test-unit: ## Run unit tests only (excludes e2e tests)
	$(GOTEST) $(UNIT_FLAGS) ./handlers ./service ./repository ./models ./utils ./middlewares

test-handlers: ## Run handler unit tests
	$(GOTEST) $(UNIT_FLAGS) ./tests -run "TestUserHandler|TestGroupHandler" 

test-e2e: ## Run end-to-end tests
	$(GOTEST) $(E2E_FLAGS) ./tests -run "TestE2ETestSuite"

test-integration: ## Run all integration tests (including e2e)
	$(GOTEST) $(E2E_FLAGS) ./tests

test-coverage: ## Run tests with coverage report
	$(GOTEST) $(TEST_FLAGS) ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-coverage-func: ## Show test coverage by function
	$(GOTEST) $(TEST_FLAGS) ./...
	$(GOCMD) tool cover -func=coverage.out

run: ## Run the application
	$(GOCMD) run main.go

run-dev: ## Run the application in development mode
	$(GOCMD) run main.go -port 3000

docker-build: ## Build docker image
	docker build -t $(BINARY_NAME) .

docker-run: ## Run docker container
	docker-compose up --build

docker-clean: ## Clean docker containers and images
	docker-compose down
	docker system prune -f

# Cross compilation
build-linux: ## Build for Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./main.go

install: ## Install the application
	$(GOCMD) install

# Linting and formatting
fmt: ## Format Go code
	$(GOCMD) fmt ./...

vet: ## Run go vet
	$(GOCMD) vet ./...

lint: fmt vet ## Run formatting and vetting

# Generate documentation
docs: ## Generate swagger documentation
	swag init -g main.go -o ./docs

# Test specific handlers
test-user-handlers: ## Test user handlers specifically
	$(GOTEST) $(UNIT_FLAGS) ./tests -run "TestUserHandler"

test-group-handlers: ## Test group handlers specifically  
	$(GOTEST) $(UNIT_FLAGS) ./tests -run "TestGroupHandler"

# Database related
migrate: ## Run database migrations (if using a migration tool)
	@echo "Database migration would go here"

# CI/CD helpers
ci-test: deps test-coverage ## Run tests for CI (with coverage)

ci-lint: lint ## Run linting for CI

ci-build: deps ci-lint ci-test build ## Full CI pipeline

# Development helpers
dev-setup: deps ## Setup development environment
	@echo "Setting up development environment..."
	@echo "Installing required tools..."
	$(GOGET) github.com/swaggo/swag/cmd/swag@latest

watch: ## Watch for changes and run tests (requires entr or similar tool)
	find . -name "*.go" | entr -r make test

# Security
security-check: ## Run security checks (requires gosec)
	@which gosec > /dev/null || (echo "Installing gosec..." && $(GOGET) github.com/securecodewarrior/gosec/v2/cmd/gosec@latest)
	gosec ./...

.DEFAULT_GOAL := help