package local

import (
	"github.com/antoniomiletta/fileman/config"
)

type LocalStorage struct {
	basePath string
}

func New(cfg config.StorageConfig) (*LocalStorage, error) {
	return &LocalStorage{
		basePath: cfg.LocalStoragePath,
	}, nil
}
