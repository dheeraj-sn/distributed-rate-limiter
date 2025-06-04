APP_NAME=distributed-rate-limiter
PORT?=8080
REDIS_URL?=redis://localhost:6379

.PHONY: all build run docker clean test

all: build

build:
	go build -o bin/$(APP_NAME) ./cmd/server

run:
	REDIS_URL=$(REDIS_URL) go run ./cmd/server

docker:
	docker build -t $(APP_NAME):latest .

docker-run:
	docker run -e REDIS_URL=$(REDIS_URL) -p $(PORT):8080 $(APP_NAME):latest

test:
	go test ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/