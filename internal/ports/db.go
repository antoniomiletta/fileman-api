package ports

import (
	"context"
	"time"

	"github.com/antoniomiletta/fileman/internal/domain/auth"
	"github.com/antoniomiletta/fileman/internal/domain/file"
	"github.com/antoniomiletta/fileman/internal/domain/folder"
	"github.com/antoniomiletta/fileman/internal/jobs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type AuthRepository interface {
	SignUp(ctx context.Context, user *auth.User) error
	FindByEmail(ctx context.Context, email string) (*auth.User, error)
}

type FileRepository interface {
	Create(ctx context.Context, file *file.File) error
	Move(ctx context.Context, id, newParentID uuid.UUID) error
	Rename(ctx context.Context, id uuid.UUID, newName string) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*file.File, error)
	FindByNameInParent(ctx context.Context, name string, parentID uuid.UUID) (*file.File, error)
	MarkUploaded(ctx context.Context, id uuid.UUID) error
}

type FolderRepository interface {
	CreateRoot(ctx context.Context, ownerID uuid.UUID) error
	Create(ctx context.Context, folder *folder.Folder) error
	ListChildren(ctx context.Context, folderID uuid.UUID) (*folder.FolderContent, error)
	Move(ctx context.Context, id, newParentID uuid.UUID) error
	Rename(ctx context.Context, id uuid.UUID, newName string) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*folder.Folder, error)
	FindByNameInParent(ctx context.Context, name string, parentID uuid.UUID) (*folder.Folder, error)
	ListDescendantFileKeys(ctx context.Context, id uuid.UUID) ([]string, error)
}

type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type TxRunner interface {
	RunTx(ctx context.Context, fn func(Querier) error) error
}

type CleanupJobRepository interface {
	Enqueue(ctx context.Context, storageKey string) error
	ClaimBatch(ctx context.Context, limit int) ([]jobs.CleanupJob, error)
	MarkDone(ctx context.Context, id uuid.UUID) error
	MarkFailed(ctx context.Context, id uuid.UUID, cause error, maxRetries int) error
	Reclaim(ctx context.Context, staleTime time.Duration) (int64, error)
}
