.DEFAULT_GOAL := build

.PHONY:build run clean

build:
	go build ./cmd/server/main.go

run: build
	go run ./cmd/server/main.go

clean:
	go clean
