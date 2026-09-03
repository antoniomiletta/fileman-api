package local

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/antoniomiletta/fileman/internal/pkg/storage"
)

const DataDirPerm os.FileMode = 0o755 // rwxr-xr-x

// Upload writes r into fileKey, creating any necessary parent directories.
// The write is atomic: content is staged in a temp file and rename only happens
// after a full, successful write. Crashes or cancels never leave a partially written file.
func (s *LocalStorage) Upload(ctx context.Context, fileKey string, r io.Reader, size int64) error {
	dest := s.fullPath(fileKey)
	dir := filepath.Dir(dest)

	if err := os.MkdirAll(dir, DataDirPerm); err != nil {
		return fmt.Errorf("local storage: upload: create directory: %w", err)
	}

	tmp, err := os.CreateTemp(dir, ".upload-*.tmp")
	if err != nil {
		return fmt.Errorf("local storage: upload: create temp file: %w", err)
	}
	// If returning early, cleanup temp file. After rename is successful,
	// tmp.Name() no longer exists, so this just becomes a no-op.
	defer os.Remove(tmp.Name())

	written, err := io.Copy(tmp, r)
	if err != nil {
		return fmt.Errorf("local storage: upload: write file: %w", err)
	}

	if written != size {
		return fmt.Errorf("%w: expected %d bytes, wrote %d", storage.ErrFileSizeMismatch, size, written)
	}

	if err := tmp.Close(); err != nil {
		return fmt.Errorf("local storage: upload: close temp file: %w", err)
	}

	if err := os.Rename(tmp.Name(), dest); err != nil {
		return fmt.Errorf("local storage: upload: rename tmp file: %w", err)
	}

	return nil
}

// Download returns a reader for the content at fileKey.
// The caller is responsible for closing it.
func (s *LocalStorage) Download(ctx context.Context, fileKey string) (io.ReadCloser, error) {
	f, err := os.Open(s.fullPath(fileKey))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("local storage: download: open file: %w", storage.ErrObjectNotFound)
		}
		return nil, fmt.Errorf("local storage: download: open file: %w", err)
	}

	return f, nil
}

// Delete remove the object at fileKey. Deleting a key that does not exist is not an error.
func (s *LocalStorage) Delete(ctx context.Context, fileKey string) error {
	if err := os.Remove(s.fullPath(fileKey)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("local storage: delete: remove file: %w", err)
	}

	return nil
}

func (s *LocalStorage) Exists(ctx context.Context, fileKey string) (bool, error) {
	_, err := os.Stat(s.fullPath(fileKey))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("local storage: exists: stat file: %w", err)
	}

	return true, nil
}
