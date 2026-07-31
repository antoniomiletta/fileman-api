package awss3

import (
	"context"
	"fmt"
	"io"
)

func (s *S3Storage) Upload(ctx context.Context, fileKey string, r io.Reader, size int64) error {
	fmt.Println("Local.Upload()")
	return nil
}
func (s *S3Storage) Download(ctx context.Context, fileKey string) (io.ReadCloser, error) {
	fmt.Println("Local.Upload()")
	return nil, nil
}
func (s *S3Storage) Delete(ctx context.Context, fileKey string) error {
	fmt.Println("Local.Upload()")
	return nil
}
func (s *S3Storage) Exists(ctx context.Context, fileKey string) (bool, error) {
	fmt.Println("Local.Exists()")
	return true, nil

}
