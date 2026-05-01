package ports

import (
	"context"

	"github.com/antoniomiletta/fileman/internal/domain"
)

type FileRepository interface {
	Create(ctx context.Context, file *domain.File) error
	GetByID(ctx context.Context, id string) (*domain.File, error)
	ListByFolder(ctx context.Context, folderID string) ([]*domain.File, error)
	Move(ctx context.Context, id string, newFolderID string) error
	Delete(ctx context.Context, id string) error
}

type FolderRepository interface {
	Create(ctx context.Context, folder *domain.Folder) error
	GetByID(ctx context.Context, id string) (*domain.Folder, error)
	ListChildren(ctx context.Context, id string) ([]*domain.Folder, error)
	Move(ctx context.Context, id string, newParentID string) error
	Delete(ctx context.Context, id string) error
}
