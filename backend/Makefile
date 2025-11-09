.PHONY: help build run test test-unit test-int test-e2e test-coverage clean docker-build docker-up docker-down migrate seed fmt vet lint deps swagger

help:
	@echo "Mini Attendance Backend - Available Commands"
	@echo ""
	@echo "Development:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make clean          - Clean build artifacts"
	@echo ""
	@echo "Testing:"
	@echo "  make test           - Run all tests"
	@echo "  make test-unit      - Run unit tests only"
	@echo "  make test-int       - Run integration tests"
	@echo "  make test-e2e       - Run E2E tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make test-race      - Run tests with race detector"
	@echo "  make test-bench     - Run benchmarks"
	@echo ""
	@echo "Code Quality:"
	@echo "  make fmt            - Format code"
	@echo "  make vet            - Run go vet"
	@echo "  make lint           - Run linter"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-up      - Start Docker containers"
	@echo "  make docker-down    - Stop Docker containers"
	@echo "  make docker-logs    - View Docker logs"
	@echo ""
	@echo "Database:"
	@echo "  make migrate        - Run database migrations"
	@echo "  make seed           - Seed database with sample data"
	@echo ""
	@echo "Dependencies:"
	@echo "  make deps           - Download and tidy dependencies"
	@echo "  make swagger        - Generate Swagger documentation"
	@echo ""

build:
	@echo "Building application..."
	go build -o attendance-api ./cmd/attendance-api/main.go
	@echo "Build complete: attendance-api"

run:
	@echo "Running application..."
	go run ./cmd/attendance-api/main.go

test:
	@echo "Running all tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo ""
	@echo "Test Summary:"
	@go tool cover -func=coverage.out | tail -1

test-unit:
	@echo "Running unit tests..."
	go test -v -short -race ./...

test-int:
	@echo "Running integration tests..."
	@echo "Make sure test database is running..."
	go test -v -race -tags=integration ./tests/integration/...

test-e2e:
	@echo "Running E2E tests..."
	@echo "Make sure application server is running on port 8080..."
	go test -v -tags=e2e ./tests/e2e/...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo ""
	@echo "Generating HTML coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report saved to coverage.html"
	@echo ""
	@echo "Coverage Summary:"
	@go tool cover -func=coverage.out | tail -1

test-coverage-web:
	@echo "Running tests and opening coverage report in browser..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

test-race:
	@echo "Running tests with race detector..."
	go test -v -race ./...

test-bench:
	@echo "Running benchmarks..."
	go test -v -bench=. -benchmem ./...

test-specific:
	@echo "Usage: make test-specific TEST=TestName"
	@echo "Example: make test-specific TEST=TestLogin"
	go test -v -run $(TEST) ./...

test-watch:
	@echo "Watching for changes and running tests..."
	@which entr > /dev/null || (echo "entr not installed. Install with: brew install entr" && exit 1)
	find . -name '*.go' | entr -r go test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -f attendance-api attendance-worker
	go clean

docker-build:
	@echo "Building Docker image..."
	docker build -t mini-attendance:latest .

docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d
	@echo "Services started. Access:"
	@echo "  API: http://localhost:8080"
	@echo "  Grafana: http://localhost:3000 (admin/admin)"
	@echo "  Health: http://localhost:8080/health"

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

docker-logs:
	docker-compose logs -f backend

docker-logs-all:
	docker-compose logs -f

migrate:
	@echo "Running migrations..."
	go run ./cmd/attendance-api/main.go

seed:
	@echo "Seeding database..."
	go run ./cmd/attendance-api/main.go

fmt:
	@echo "Formatting code..."
	go fmt ./...

vet:
	@echo "Running go vet..."
	go vet ./...

lint:
	@echo "Running linter..."
	golangci-lint run

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/attendance-api/main.go

.DEFAULT_GOAL := help
