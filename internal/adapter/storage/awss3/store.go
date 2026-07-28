package awss3

import (
	"context"
	"io"
)

func (s *S3Backend) Upload(ctx context.Context, fileKey string, r io.Reader, size int64) error
func (s *S3Backend) Download(ctx context.Context, fileKey string) (io.ReadCloser, error)
func (s *S3Backend) Delete(ctx context.Context, fileKey string) error
func (s *S3Backend) Exists(ctx context.Context, fileKey string) (bool, error)
