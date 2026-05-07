package local

import (
	"context"
	"io"
)

func (s *LocalStorage) Upload(ctx context.Context, fileKey string, r io.Reader, size int64) error
func (s *LocalStorage) Download(ctx context.Context, fileKey string) (io.ReadCloser, error)
func (s *LocalStorage) Delete(ctx context.Context, fileKey string) error
func (s *LocalStorage) Exists(ctx context.Context, fileKey string) (bool, error)
