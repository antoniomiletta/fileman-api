package ports

import (
	"context"
	"io"
)

type StorageBackend interface {
	Upload(ctx context.Context, fileKey string, r io.Reader, size int64) error
	Download(ctx context.Context, fileKey string) (io.ReadCloser, error)
	Delete(ctx context.Context, fileKey string) error
	Exists(ctx context.Context, fileKey string) (bool, error)
}
