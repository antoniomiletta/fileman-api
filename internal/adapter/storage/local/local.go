package local

import (
	"path/filepath"

	"github.com/antoniomiletta/fileman/config"
)

type LocalStorage struct {
	root string
}

func (s *LocalStorage) fullPath(fileKey string) string {
	return filepath.Join(s.root, fileKey)
}

func NewLocalStorage(cfg config.LocalConfig) *LocalStorage {
	return &LocalStorage{
		root: cfg.LocalRoot,
	}
}
