package awss3

import (
	"context"
	"io"
)

func (s *S3Storage) Upload(ctx context.Context, fileKey string, r io.Reader, size int64) error
func (s *S3Storage) Download(ctx context.Context, fileKey string) (io.ReadCloser, error)
func (s *S3Storage) Delete(ctx context.Context, fileKey string) error
func (s *S3Storage) Exists(ctx context.Context, fileKey string) (bool, error)
