package main

import (
	"log"
	"os"

	"github.com/antoniomiletta/fileman/config"
	"github.com/antoniomiletta/fileman/internal/adapter/db/postgres"
	"github.com/antoniomiletta/fileman/internal/adapter/storage/awss3"
	"github.com/antoniomiletta/fileman/internal/adapter/storage/local"
	"github.com/antoniomiletta/fileman/internal/api"
	"github.com/antoniomiletta/fileman/internal/pkg/authenticator"
	"github.com/antoniomiletta/fileman/internal/ports"
	"github.com/antoniomiletta/fileman/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db := initDb(cfg.DB)
	defer db.Close()

	store := initStorage(cfg.Storage)

	authRepo := postgres.NewAuthRepository(db)
	folderRepo := postgres.NewFolderRepository(db)
	fileRepo := postgres.NewFileRepository(db)
	txRunner := postgres.NewTxRunner(db.Pool())

	authn := authenticator.NewAuthenticator(cfg.Auth)

	authSvc := service.NewAuthService(authRepo, txRunner, authn)
	folderSvc := service.NewFolderService(folderRepo)
	fileSvc := service.NewFileService(fileRepo, folderRepo, store)

	svr := api.NewServer(api.ServerDeps{
		Cfg:       cfg.Server,
		AuthSvc:   authSvc,
		FolderSvc: folderSvc,
		FileSvc:   fileSvc,
		Authn:     authn,
	})

	if err := svr.Serve(); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}

func initDb(cfg config.DBConfig) *postgres.DB {
	db, err := postgres.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	return db
}

func initStorage(cfg config.StorageConfig) ports.StorageBackend {
	var store ports.StorageBackend

	switch cfg.Backend {
	case "local":
		if cfg.LocalConfig.LocalRoot == "" {
			log.Fatalf("STORAGE_BACKEND=local but LOCAL_STORAGE_ROOT is not set")
		}
		os.MkdirAll(cfg.LocalConfig.LocalRoot, 0o755) // rwxr-xr-x
		store = local.NewLocalStorage(cfg.LocalConfig)
	case "s3":
		if cfg.S3Config.S3Bucket == "" {
			log.Fatalf("STORAGE_BACKEND=s3 but S3_BUCKET is not set")
		}
		store = awss3.NewS3Storage(cfg.S3Config)
	default:
		log.Fatalf("unexpected storage backend: %q", cfg.Backend)
	}

	return store
}
