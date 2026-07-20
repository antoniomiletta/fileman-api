package main

import (
	"log"

	"github.com/antoniomiletta/fileman/config"
	"github.com/antoniomiletta/fileman/internal/adapter/db/postgres"
	"github.com/antoniomiletta/fileman/internal/adapter/storage/s3"
	"github.com/antoniomiletta/fileman/internal/api"
	"github.com/antoniomiletta/fileman/internal/pkg/authenticator"
	"github.com/antoniomiletta/fileman/internal/service"
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

	authenticator := authenticator.NewAuthenticator(cfg.Auth)

	authSvc := service.NewAuthService(authRepo, authenticator)
	folderSvc := service.NewFolderService(folderRepo)
	fileSvc := service.NewFileService(fileRepo, folderRepo, storage)

	server := api.NewServer(api.ServerDeps{
		Cfg:       cfg.Server,
		AuthSvc:   authSvc,
		FolderSvc: folderSvc,
		FileSvc:   fileSvc,
	})

	server.Start()
}
