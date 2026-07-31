package local

import (
	"context"
	"fmt"
	"io"
)

func (s *LocalStorage) Upload(ctx context.Context, fileKey string, r io.Reader, size int64) error {
	fmt.Println("Local.Upload()")
	return nil
}
func (s *LocalStorage) Download(ctx context.Context, fileKey string) (io.ReadCloser, error) {
	fmt.Println("Local.Download()")
	return nil, nil
}
func (s *LocalStorage) Delete(ctx context.Context, fileKey string) error {
	fmt.Println("Local.Delete()")
	return nil
}
func (s *LocalStorage) Exists(ctx context.Context, fileKey string) (bool, error) {
	fmt.Println("Local.Exists()")
	return true, nil
}
