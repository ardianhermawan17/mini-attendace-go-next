.PHONY: help build run test clean docker-build docker-up docker-down migrate seed

help:
	@echo "Mini Attendance Backend - Available Commands"
	@echo ""
	@echo "Development:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"
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

build:
	@echo "Building application..."
	go build -o attendance-api ./cmd/attendance-api/main.go
	@echo "Build complete: attendance-api"

run:
	@echo "Running application..."
	go run ./cmd/attendance-api/main.go

test:
	@echo "Running tests..."
	go test -v -cover ./...

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
