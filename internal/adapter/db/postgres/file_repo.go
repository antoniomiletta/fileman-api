package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/antoniomiletta/fileman/internal/domain/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type FileRepository struct {
	db Querier
}

func NewFileRepository(q Querier) *FileRepository {
	return &FileRepository{db: q}
}

// sql queries
func (r *FileRepository) Create(ctx context.Context, f *file.File) error {
	const query = `
		INSERT INTO files (id, owner_id, parent_id, name, mime_type, size, storage_key, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		`

	_, err := r.db.Exec(ctx, query,
		f.ID,
		f.OwnerID,
		f.ParentID,
		f.Name,
		f.MIMEType,
		f.Size,
		f.StorageKey,
		f.Status)
	if err != nil {
		if isUniqueViolation(err) {
			return file.ErrFileNameConflict
		}

		return fmt.Errorf("postgres: create file: %w", err)
	}

	return nil
}

func (r *FileRepository) Move(ctx context.Context, id, newParentID uuid.UUID) error {
	const query = `
		UPDATE files
		SET parent_id = $1, updated_at = NOW()
		WHERE id = $2
		`

	tag, err := r.db.Exec(ctx, query, newParentID, id)
	if err != nil {
		if isUniqueViolation(err) {
			return file.ErrFileNameConflict
		}

		return fmt.Errorf("postgres: move file: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return file.ErrFileNotFound
	}

	return nil
}

func (r *FileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM files
		WHERE id = $1
		`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("postgres: delete file: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return file.ErrFileNotFound
	}

	return nil
}

func (r *FileRepository) FindByID(ctx context.Context, id uuid.UUID) (*file.File, error) {
	const query = `
		SELECT id, owner_id, parent_id, name, mime_type, size, storage_key, status, created_at, updated_at
		FROM files
		WHERE id = $1
		`

	f := &file.File{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&f.ID,
		&f.OwnerID,
		&f.ParentID,
		&f.Name,
		&f.MIMEType,
		&f.Size,
		&f.StorageKey,
		&f.Status,
		&f.CreatedAt,
		&f.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, file.ErrFileNotFound
		}
		return nil, fmt.Errorf("postgres: find file by id: %w", err)
	}

	return f, nil
}

func (r *FileRepository) FindByNameInParent(ctx context.Context, name string, parentID uuid.UUID) (*file.File, error) {
	const query = `
		SELECT id, owner_id, parent_id, name, mime_type, size, storage_key, status, created_at, updated_at
		FROM files
		WHERE name = $1
		AND parent_id = $2
		`

	f := &file.File{}
	err := r.db.QueryRow(ctx, query, name, parentID).Scan(
		&f.ID,
		&f.OwnerID,
		&f.ParentID,
		&f.Name,
		&f.MIMEType,
		&f.Size,
		&f.StorageKey,
		&f.Status,
		&f.CreatedAt,
		&f.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, file.ErrFileNotFound
		}
		return nil, fmt.Errorf("postgres: find file by name in parent: %w", err)
	}

	return f, nil
}
