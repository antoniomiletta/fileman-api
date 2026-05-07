.DEFAULT_GOAL := build

.PHONY:build run clean docker-up docker-down

build:
	go build ./cmd/server

run:
	go run ./cmd/server

clean:
	go clean

docker-up:
	docker compose up -d

docker-down:
	docker compose down -d
