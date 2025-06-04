APP_NAME=distributed-rate-limiter
PORT?=8080
REDIS_URL?=redis://localhost:6379

.PHONY: all build run docker clean test

all: build

build:
	go build -o bin/$(APP_NAME) ./cmd/server

run:
	REDIS_URL=$(REDIS_URL) go run ./cmd/server

test:
	go test ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/

# Start everything using Docker Compose
up:
	docker-compose up --build

# Stop and remove containers
down:
	docker-compose down

# Rebuild container without cache
rebuild:
	docker-compose build --no-cache

# View logs
logs:
	docker-compose logs -f rlaas


