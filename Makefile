.PHONY: build run test clean docker-up docker-down docker-logs help

# Variables
APP_NAME=go-auth-api
DOCKER_COMPOSE=docker-compose

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build the Go application"
	@echo "  run          - Run the application locally"
	@echo "  test         - Run tests"
	@echo "  test-api     - Test API endpoints"
	@echo "  docker-up    - Start Docker containers"
	@echo "  docker-down  - Stop Docker containers"
	@echo "  docker-logs  - Show Docker logs"
	@echo "  clean        - Clean build artifacts"

# Build the application
build:
	go mod tidy
	go build -o bin/$(APP_NAME) ./cmd/server

# Run the application locally
run:
	go run ./cmd/server

# Run tests
test:
	go test -v ./...

# Test API endpoints
test-api:
	./test_api.sh

# Start Docker containers
docker-up:
	$(DOCKER_COMPOSE) up --build -d

# Stop Docker containers
docker-down:
	$(DOCKER_COMPOSE) down

# Show Docker logs
docker-logs:
	$(DOCKER_COMPOSE) logs -f

# Clean build artifacts
clean:
	rm -rf bin/
	$(DOCKER_COMPOSE) down -v
	docker system prune -f
