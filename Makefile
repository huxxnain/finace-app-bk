.PHONY: help build run test clean install-deps lint docker-build docker-run

help:
	@echo "Finance App - Available Commands"
	@echo "================================="
	@echo "make build           - Build the application"
	@echo "make run             - Run the application"
	@echo "make dev             - Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)"
	@echo "make test            - Run tests"
	@echo "make clean           - Clean build artifacts"
	@echo "make install-deps    - Download dependencies"
	@echo "make lint            - Run linter (requires golangci-lint)"
	@echo "make docker-build    - Build Docker image"
	@echo "make docker-run      - Run in Docker"
	@echo "make fmt             - Format code"

install-deps:
	go mod download
	go mod tidy

build:
	go build -o finance-app cmd/server/main.go

run: build
	./finance-app

dev:
	air

test:
	go test ./... -v

clean:
	rm -f finance-app
	go clean

lint:
	golangci-lint run

fmt:
	go fmt ./...

docker-build:
	docker build -t finance-app:latest .

docker-run:
	docker run -p 3000:3000 --env-file .env finance-app:latest
