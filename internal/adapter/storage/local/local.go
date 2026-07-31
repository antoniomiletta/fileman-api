package local

import "github.com/antoniomiletta/fileman/config"

type LocalStorage struct {
	root string
}

func NewLocalStorage(cfg config.LocalConfig) *LocalStorage {
	return &LocalStorage{
		root: cfg.LocalRoot,
	}
}
