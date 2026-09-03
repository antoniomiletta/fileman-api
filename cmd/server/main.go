package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/antoniomiletta/fileman/config"
	"github.com/antoniomiletta/fileman/internal/adapters/db/pg"
	"github.com/antoniomiletta/fileman/internal/adapters/storage/awss3"
	"github.com/antoniomiletta/fileman/internal/adapters/storage/local"
	"github.com/antoniomiletta/fileman/internal/api"
	"github.com/antoniomiletta/fileman/internal/api/transport"
	"github.com/antoniomiletta/fileman/internal/pkg/authenticator"
	"github.com/antoniomiletta/fileman/internal/ports"
	"github.com/antoniomiletta/fileman/internal/service"
	"github.com/antoniomiletta/fileman/internal/workers"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db := initDb(cfg.DB)
	defer db.Close()

	store := initStorage(cfg.Storage)

	authRepo := pg.NewAuthRepository(db)
	folderRepo := pg.NewFolderRepository(db)
	fileRepo := pg.NewFileRepository(db)
	txRunner := pg.NewTxRunner(db.Pool())
	cleanupRepo := pg.NewCleanupJobRepository(db)

	authn := authenticator.NewAuthenticator(cfg.Auth)
	resp := transport.NewResponder(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	authSvc := service.NewAuthService(authRepo, txRunner, authn)
	folderSvc := service.NewFolderService(folderRepo, txRunner)
	fileSvc := service.NewFileService(fileRepo, folderRepo, store, txRunner, cleanupRepo)

	svr := api.NewServer(api.ServerDeps{
		Cfg:       cfg.Server,
		AuthSvc:   authSvc,
		FolderSvc: folderSvc,
		FileSvc:   fileSvc,
		Authn:     authn,
		Resp:      resp,
	})

	// Sweep inconsitent uploads
	uploadReconciler := workers.NewUploadReconciler(fileRepo, store)
	if err := uploadReconciler.Sweep(context.Background(), time.Hour); err != nil {
		log.Println(err)
	}

	// Start cleanup worker
	cleanupWorker := workers.NewCleanupWorker(cleanupRepo, store, cfg.Cleanup)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup
	defer wg.Wait() // runs before db.Close()
	wg.Go(func() {
		cleanupWorker.Run(ctx)
	})

	// Start server
	chServeErr := make(chan error, 1)
	go func() {
		chServeErr <- svr.Serve(cfg.Server)
	}()

	select {
	case <-ctx.Done():
		log.Printf("shutdown signal received, stopping server...")
		// Server shutdown uses a fresh context, otherwise the cancellation signal
		// would make the shutdown itself fail.
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*10)
		defer shutdownCancel()
		if err := svr.Shutdown(shutdownCtx); err != nil {
			log.Printf("failed to shutdown server: %v", err)
		}
	case err := <-chServeErr:
		if err != nil {
			log.Fatalf("server error: %v", err)
		}
	}
}

func initDb(cfg config.DBConfig) *pg.DB {
	db, err := pg.Connect(cfg)
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
		if err := os.MkdirAll(cfg.LocalConfig.LocalRoot, local.DataDirPerm); err != nil {
			log.Fatalf("failed to create storage root: %v", err)
		}
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
