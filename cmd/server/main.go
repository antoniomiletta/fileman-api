package main

import (
	"log"

	"github.com/antoniomiletta/fileman/config"
	"github.com/antoniomiletta/fileman/internal/adapters/repository/postgres"
	"github.com/antoniomiletta/fileman/internal/adapters/storage/local"
	"github.com/antoniomiletta/fileman/internal/api"
	"github.com/antoniomiletta/fileman/internal/services"
)

func main() {
	cfg := config.Load()

	db, err := postgres.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	storage, err := local.New(cfg.Storage)
	if err != nil {
		log.Fatalf("failed to initialize storage backend: %v", err)
	}

	authRepo := postgres.NewAuthRepository(db)
	fileRepo := postgres.NewFileRepository(db)
	folderRepo := postgres.NewFolderRepository(db)

	authSvc := services.NewAuthService(authRepo)
	fileSvc := services.NewFileService(fileRepo, storage)
	folderSvc := services.NewFolderService(folderRepo)

	server := api.NewServer(authSvc, fileSvc, folderSvc, cfg.Server)

	server.Start()
}
