package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/antoniomiletta/fileman/internal/domain/file"
	"github.com/antoniomiletta/fileman/internal/domain/folder"
	"github.com/antoniomiletta/fileman/internal/ports"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type FolderRepository struct {
	db ports.Querier
}

func NewFolderRepository(q ports.Querier) *FolderRepository {
	return &FolderRepository{db: q}
}

func (r *FolderRepository) CreateRoot(ctx context.Context, ownerID uuid.UUID) error {
	const query = `
		INSERT INTO folders (id, owner_id, parent_id, name, created_at, updated_at)
		VALUES ($1, $2, NULL, 'root', NOW(), NOW())
		`

	_, err := r.db.Exec(ctx, query, uuid.New(), ownerID)
	if err != nil {
		return fmt.Errorf("postgres: create root folder: %w", err)
	}

	return nil
}

func (r *FolderRepository) Create(ctx context.Context, f *folder.Folder) error {
	const query = `
		INSERT INTO folders (id, owner_id, parent_id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		`

	_, err := r.db.Exec(ctx, query,
		f.ID,
		f.OwnerID,
		f.ParentID,
		f.Name)
	if err != nil {
		if isUniqueViolation(err) {
			return folder.ErrFolderNameConflict
		}
		return fmt.Errorf("postgres: create folder: %w", err)
	}

	return nil
}

func (r *FolderRepository) ListChildren(ctx context.Context, folderID uuid.UUID) (*folder.FolderContent, error) {
	subfolders, err := listChildFolders(ctx, r.db, folderID)
	if err != nil {
		return nil, err
	}

	files, err := listChildFiles(ctx, r.db, folderID)
	if err != nil {
		return nil, err
	}

	return &folder.FolderContent{
		Subfolders: subfolders,
		Files:      files,
	}, nil
}

func listChildFolders(ctx context.Context, db ports.Querier, parentID uuid.UUID) ([]*folder.Folder, error) {
	const query = `
		SELECT id, owner_id, parent_id, name, created_at, updated_at
		FROM folders
		WHERE parent_id = $1
		ORDER BY name ASC
		`

	rows, err := db.Query(ctx, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("postgres: list child folders: %w", err)
	}
	defer rows.Close()

	var subfolders []*folder.Folder
	for rows.Next() {
		f := &folder.Folder{}
		if err := rows.Scan(
			&f.ID,
			&f.OwnerID,
			&f.ParentID,
			&f.Name,
			&f.CreatedAt,
			&f.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("postgres: scan child folder: %w", err)
		}
		subfolders = append(subfolders, f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: iterate child folders: %w", err)
	}

	return subfolders, nil
}
func listChildFiles(ctx context.Context, db ports.Querier, parentID uuid.UUID) ([]*file.File, error) {
	const query = `
		SELECT id, owner_id, parent_id, name, mime_type, size, storage_key, upload_status, created_at, updated_at
		FROM files
		WHERE parent_id = $1
		ORDER BY name ASC
		`

	rows, err := db.Query(ctx, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("postgres: list child files: query files: %w", err)
	}
	defer rows.Close()

	var files []*file.File
	for rows.Next() {
		f := &file.File{}
		if err := rows.Scan(
			&f.ID,
			&f.OwnerID,
			&f.ParentID,
			&f.Name,
			&f.MIMEType,
			&f.Size,
			&f.StorageKey,
			&f.UploadStatus,
			&f.CreatedAt,
			&f.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("postgres: list child files: scan files: %w", err)
		}
		files = append(files, f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: list child files: iterate files: %w", err)
	}

	return files, nil
}

func (r *FolderRepository) Move(ctx context.Context, id, newParentID uuid.UUID) error {
	const query = `
		UPDATE folders 
		SET parent_id = $1, updated_at = NOW()
		WHERE id = $2
		`

	tag, err := r.db.Exec(ctx, query, newParentID, id)
	if err != nil {
		if isUniqueViolation(err) {
			return folder.ErrFolderNameConflict
		}
		return fmt.Errorf("postgres: move folder: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return folder.ErrFolderNotFound
	}

	return nil
}

func (r *FolderRepository) Rename(ctx context.Context, id uuid.UUID, newName string) error {
	const query = `
		UPDATE folders
		SET name = $1, updated_at = NOW()
		WHERE id = $2
		`

	tag, err := r.db.Exec(ctx, query, newName, id)
	if err != nil {
		if isUniqueViolation(err) {
			return folder.ErrFolderNameConflict
		}
		return fmt.Errorf("postgres: rename folder: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return folder.ErrFolderNotFound
	}

	return nil
}

func (r *FolderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM folders
		WHERE id = $1
		`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("postgres: delete folder: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return folder.ErrFolderNotFound
	}

	return nil
}

func (r *FolderRepository) FindByID(ctx context.Context, id uuid.UUID) (*folder.Folder, error) {
	const query = `
		SELECT id, owner_id, parent_id, name, created_at, updated_at
		FROM folders
		WHERE id = $1
		`

	f := &folder.Folder{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&f.ID,
		&f.OwnerID,
		&f.ParentID,
		&f.Name,
		&f.CreatedAt,
		&f.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, folder.ErrFolderNotFound
		}
		return nil, fmt.Errorf("postgres: find folder by id: %w", err)
	}

	return f, nil
}

func (r *FolderRepository) FindByNameInParent(ctx context.Context, name string, parentID uuid.UUID) (*folder.Folder, error) {
	const query = `
		SELECT id, owner_id, parent_id, name, created_at, updated_at
		FROM folders
		WHERE name = $1
		AND parent_id = $2
		`

	f := &folder.Folder{}
	err := r.db.QueryRow(ctx, query, name, parentID).Scan(
		&f.ID,
		&f.OwnerID,
		&f.ParentID,
		&f.Name,
		&f.CreatedAt,
		&f.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, folder.ErrFolderNotFound
		}
		return nil, fmt.Errorf("postgres: find folder by name in parent: %w", err)
	}

	return f, nil
}

func (r *FolderRepository) ListDescendantFileKeys(ctx context.Context, id uuid.UUID) ([]string, error) {
	const query = `
		WITH RECURSIVE folder_tree AS (
			SELECT id FROM folders WHERE id = $1
			UNION ALL
			SELECT f.id
			FROM folders f
			JOIN folder_tree ft ON f.parent_id = ft.id
		)
		SELECT files.storage_key
		FROM files
		WHERE files.parent_id IN (SELECT id FROM folder_tree)
		`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("postgres: list file keys: %w", err)
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string

		if err := rows.Scan(&key); err != nil {
			return nil, fmt.Errorf("postgres: scan file key: %w", err)
		}

		keys = append(keys, key)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: iterate file keys: %w", err)
	}

	return keys, nil
}
