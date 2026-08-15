ifeq (,$(wildcard .env))
$(error .env file not found)
endif

include .env
export

BUILD_DIR = ./bin
BINARY_NAME = fileman
ENTRYPOINT = ./cmd/server/main.go

DB_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATIONS_DIR = db/migrations

.PHONY: dev build run test test-integration clean docker-up docker-down docker-down-v docker-ps migrate-up migrate-down migrate-version

dev:
	go run $(ENTRYPOINT)

build:
	mkdir -p $(BUILD_DIR) && go build -o $(BUILD_DIR)/$(BINARY_NAME) $(ENTRYPOINT)

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

test:
	go test -short ./... | grep -v "no test files"

test-integration:
	go test ./... | grep -v "no test files"

clean:
	rm -rf $(BUILD_DIR)

docker-up:
	docker compose up -d

docker-down:
	docker compose down

# wipe volumes
docker-down-v:
	docker compose down -v

docker-ps:
	docker compose ps

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) down 1

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) version
