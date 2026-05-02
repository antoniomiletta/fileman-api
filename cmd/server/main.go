package main

import (
	"log"

	"github.com/antoniomiletta/fileman/config"
	"github.com/antoniomiletta/fileman/internal/adapters/repository/postgres"
	"github.com/antoniomiletta/fileman/internal/adapters/storage/local"
	"github.com/antoniomiletta/fileman/internal/ports"
)

func main() {
	cfg := config.Load()

	db, err := postgres.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	var storage ports.StorageBackend
	var storageErr error

	switch storage {
	default:
		storage, storageErr = local.New(cfg.Storage)
	}

	if storageErr != nil {
		log.Fatalf("failed to initialize storage backend: %v", err)
	}
}
