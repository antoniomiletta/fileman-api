package ports

import (
	"context"

	"github.com/antoniomiletta/fileman/internal/domain"
	"github.com/google/uuid"
)

type AuthRepository interface {
	Register(ctx context.Context, user *domain.User) error
	Login(ctx context.Context, email, password string) (string, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

type FileRepository interface {
	Create(ctx context.Context, file *domain.File) error
	Move(ctx context.Context, id, newFolderID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type FolderRepository interface {
	Create(ctx context.Context, folder *domain.Folder) error
	ListChildren(ctx context.Context, folderID uuid.UUID) ([]*domain.Folder, error)
	Move(ctx context.Context, id, newParentID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}
