package s3

import (
	"github.com/antoniomiletta/fileman/config"
)

type LocalStorage struct {
	bucket string
}

func New(cfg config.StorageConfig) (*LocalStorage, error) {
	return &LocalStorage{
		bucket: cfg.S3Bucket,
	}, nil
}
