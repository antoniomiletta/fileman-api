package main

import (
	"log"

	"github.com/antoniomiletta/fileman/config"
	"github.com/antoniomiletta/fileman/internal/adapters/repository/postgres"
	"github.com/antoniomiletta/fileman/internal/adapters/storage/s3"
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

	storage, err := s3.New(cfg.Storage)
	if err != nil {
		log.Fatalf("failed to initialize storage backend: %v", err)
	}

	authRepo := postgres.NewAuthRepository(db)
	folderRepo := postgres.NewFolderRepository(db)
	fileRepo := postgres.NewFileRepository(db)

	authSvc := services.NewAuthService(authRepo)
	folderSvc := services.NewFolderService(folderRepo)
	fileSvc := services.NewFileService(fileRepo, folderRepo, storage)

	server := api.NewServer(api.NewServerInput{
		AuthSvc:   authSvc,
		FolderSvc: folderSvc,
		FileSvc:   fileSvc,
		Cfg:       cfg.Server,
	})

	server.Start()
}
