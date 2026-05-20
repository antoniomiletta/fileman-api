package ports

import (
	"context"

	"github.com/antoniomiletta/fileman/internal/domain/auth"
	"github.com/antoniomiletta/fileman/internal/domain/file"
	"github.com/antoniomiletta/fileman/internal/domain/folder"
	"github.com/google/uuid"
)

type AuthRepository interface {
	Register(ctx context.Context, user *auth.User) error
	Login(ctx context.Context, email, password string) (string, error)
	FindByEmail(ctx context.Context, email string) (*auth.User, error)
}

type FileRepository interface {
	Create(ctx context.Context, file *file.File) error
	Move(ctx context.Context, id, newFolderID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type FolderRepository interface {
	Create(ctx context.Context, folder *folder.Folder) error
	FindByID(ctx context.Context, id uuid.UUID) (*folder.Folder, error)
	ListChildren(ctx context.Context, folderID uuid.UUID) (*folder.FolderContent, error)
	Move(ctx context.Context, id, newParentID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}
