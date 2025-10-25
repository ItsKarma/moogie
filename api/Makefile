# Makefile for Moogie API

.PHONY: help build run dev test clean migrate-up migrate-down migrate-create deps docs

# Default target
help:
	@echo "Available commands:"
	@echo "  deps         - Download dependencies"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  dev          - Run in development mode with auto-reload"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  migrate-up   - Run database migrations"
	@echo "  migrate-down - Rollback database migrations"
	@echo "  migrate-create NAME=<name> - Create new migration"
	@echo "  docs         - Generate API documentation"

# Download dependencies
deps:
	go mod download
	go mod tidy

# Build the application
build: deps
	go build -o bin/server cmd/server/main.go

# Run the application
run: build
	./bin/server

# Development mode (requires air for auto-reload)
dev:
	@if ! command -v air > /dev/null; then \
		echo "Installing air for hot reload..."; \
		go install github.com/air-verse/air@latest; \
	fi
	air

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Database migrations
migrate-up:
	@if ! command -v goose > /dev/null; then \
		echo "Installing goose..."; \
		go install github.com/pressly/goose/v3/cmd/goose@latest; \
	fi
	goose -dir migrations postgres "$(DB_URL)" up

migrate-down:
	@if ! command -v goose > /dev/null; then \
		echo "Installing goose..."; \
		go install github.com/pressly/goose/v3/cmd/goose@latest; \
	fi
	goose -dir migrations postgres "$(DB_URL)" down

migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@if ! command -v goose > /dev/null; then \
		echo "Installing goose..."; \
		go install github.com/pressly/goose/v3/cmd/goose@latest; \
	fi
	goose -dir migrations create $(NAME) sql

# Generate API documentation
docs:
	@if ! command -v swag > /dev/null; then \
		echo "Installing swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	swag init -g cmd/server/main.go -o docs/

# Docker commands
docker-build:
	docker build -t moogie-api .

docker-run:
	docker run -p 8080:8080 --env-file .env moogie-api

# Environment variables for development
DB_URL ?= postgres://moogie:moogie@localhost:5432/moogie?sslmode=disable
