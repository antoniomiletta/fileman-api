package ports

import (
	"context"

	"github.com/antoniomiletta/fileman/internal/domain"
)

type AuthRepository interface {
	Register(ctx context.Context, user *domain.User) error
	Login(ctx context.Context, email, password string) (string, error)
}

type FileRepository interface {
	Create(ctx context.Context, file *domain.File) error
	ListFromFolder(ctx context.Context, folderID string) ([]*domain.File, error)
	GetByID(ctx context.Context, id string) (*domain.File, error)
	Move(ctx context.Context, id, newFolderID string) error
	Delete(ctx context.Context, id string) error
}

type FolderRepository interface {
	Create(ctx context.Context, folder *domain.Folder) error
	ListChildren(ctx context.Context, folderID string) ([]*domain.Folder, error)
	GetByID(ctx context.Context, id string) (*domain.Folder, error)
	Move(ctx context.Context, id, newParentID string) error
	Delete(ctx context.Context, id string) error
}
