ifeq (,$(wildcard .env))
$(error .env file not found)
endif

include .env
export

MIGRATIONS_DIR = db/migrations
BUILD_DIR = ./bin
BINARY_NAME = fileman
ENTRYPOINT = ./cmd/server/main.go

.PHONY:dev build run clean docker-up docker-down docker-ps migrate-up migrate-down

dev:
	go run $(ENTRYPOINT)

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(ENTRYPOINT)

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

clean:
	rm -rf $(BUILD_DIR)

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-ps:
	docker compose ps

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) down 1

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) version
