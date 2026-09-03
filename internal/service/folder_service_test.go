package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/antoniomiletta/fileman/internal/domain/auth"
	"github.com/antoniomiletta/fileman/internal/domain/file"
	"github.com/antoniomiletta/fileman/internal/domain/folder"
	"github.com/antoniomiletta/fileman/internal/pkg/reqctx"
	"github.com/antoniomiletta/fileman/internal/ports"
	"github.com/antoniomiletta/fileman/internal/service"
)

type fakeFolderRepo struct {
	folders map[uuid.UUID]*folder.Folder
}

func newFakeFolderRepo() *fakeFolderRepo {
	return &fakeFolderRepo{folders: make(map[uuid.UUID]*folder.Folder)}
}

type fakeTxRunner struct {
	t *testing.T
}

func newFakeTxRunner(t *testing.T) *fakeTxRunner {
	return &fakeTxRunner{t: t}
}

func (f *fakeTxRunner) RunTx(ctx context.Context, fn func(ports.Querier) error) error {
	f.t.Fatal("txRunner.RunTx should not be called in unit tests.")
	return nil
}

func (f *fakeFolderRepo) seed(fol *folder.Folder) {
	f.folders[fol.ID] = fol
}

func (f *fakeFolderRepo) CreateRoot(ctx context.Context, ownerID uuid.UUID) error {
	root := &folder.Folder{
		ID:       uuid.New(),
		OwnerID:  ownerID,
		ParentID: nil,
		Name:     "root",
	}
	f.folders[root.ID] = root
	return nil
}

func (f *fakeFolderRepo) Create(ctx context.Context, fol *folder.Folder) error {
	f.folders[fol.ID] = fol
	return nil
}

func (f *fakeFolderRepo) ListChildren(ctx context.Context, folderID uuid.UUID) (*folder.FolderContent, error) {
	var subfolders []*folder.Folder
	for _, fol := range f.folders {
		if fol.ParentID != nil && *fol.ParentID == folderID {
			subfolders = append(subfolders, fol)
		}
	}
	return &folder.FolderContent{Subfolders: subfolders, Files: []*file.File{}}, nil
}

func (f *fakeFolderRepo) Move(ctx context.Context, id, newParentID uuid.UUID) error {
	fol, ok := f.folders[id]
	if !ok {
		return folder.ErrFolderNotFound
	}
	fol.ParentID = &newParentID
	return nil
}

func (f *fakeFolderRepo) Rename(ctx context.Context, id uuid.UUID, newName string) error {
	fol, ok := f.folders[id]
	if !ok {
		return folder.ErrFolderNotFound
	}
	fol.Name = newName
	return nil
}

func (f *fakeFolderRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if _, ok := f.folders[id]; !ok {
		return folder.ErrFolderNotFound
	}
	delete(f.folders, id)
	return nil
}

func (f *fakeFolderRepo) FindByID(ctx context.Context, id uuid.UUID) (*folder.Folder, error) {
	fol, ok := f.folders[id]
	if !ok {
		return nil, folder.ErrFolderNotFound
	}
	return fol, nil
}

func (f *fakeFolderRepo) FindByNameInParent(ctx context.Context, name string, parentID uuid.UUID) (*folder.Folder, error) {
	for _, fol := range f.folders {
		if fol.Name == name && fol.ParentID != nil && *fol.ParentID == parentID {
			return fol, nil
		}
	}
	return nil, folder.ErrFolderNotFound
}

func (f *fakeFolderRepo) ListDescendantFileKeys(ctx context.Context, id uuid.UUID) ([]string, error) {
	return nil, nil
}

// Tests
func TestFolderService_CreateFolder(t *testing.T) {
	ownerA := uuid.New()
	ownerB := uuid.New()

	setup := func() (*fakeFolderRepo, *folder.Folder) {
		repo := newFakeFolderRepo()
		root := &folder.Folder{ID: uuid.New(), OwnerID: ownerA, ParentID: nil, Name: "root"}
		repo.seed(root)
		return repo, root
	}

	t.Run("rejects empty name", func(t *testing.T) {
		repo, root := setup()
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.CreateFolder(ctx, &root.ID, "")
		if err == nil {
			t.Fatal("expected an error for empty name, got nil")
		}
	})

	t.Run("rejects creating a new root folder", func(t *testing.T) {
		repo, _ := setup()
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.CreateFolder(ctx, nil, "new-root")
		if !errors.Is(err, folder.ErrCannotCreateNewRootFolder) {
			t.Fatalf("expected ErrCannotCreateNewRootFolder, got: %v", err)
		}
	})

	t.Run("propagates error when parent does not exist", func(t *testing.T) {
		repo, _ := setup()
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		missingParent := uuid.New()
		err := svc.CreateFolder(ctx, &missingParent, "docs")
		if !errors.Is(err, folder.ErrFolderNotFound) {
			t.Fatalf("expected ErrFolderNotFound, got: %v", err)
		}
	})

	t.Run("rejects creating in a folder owned by someone else", func(t *testing.T) {
		repo, root := setup()
		root.OwnerID = ownerB
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.CreateFolder(ctx, &root.ID, "docs")
		if !errors.Is(err, auth.ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got: %v", err)
		}
	})

	t.Run("rejects name conflict in parent", func(t *testing.T) {
		repo, root := setup()
		existing := &folder.Folder{ID: uuid.New(), OwnerID: ownerA, ParentID: &root.ID, Name: "docs"}
		repo.seed(existing)
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.CreateFolder(ctx, &root.ID, "docs")
		if !errors.Is(err, folder.ErrFolderNameConflict) {
			t.Fatalf("expected ErrFolderNameConflict, got: %v", err)
		}
	})

	t.Run("allows same name with different casing in the same parent", func(t *testing.T) {
		repo, root := setup()
		existing := &folder.Folder{ID: uuid.New(), OwnerID: ownerA, ParentID: &root.ID, Name: "docs"}
		repo.seed(existing)
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.CreateFolder(ctx, &root.ID, "Docs")
		if err != nil {
			t.Fatalf("expected success, got: %v", err)
		}

		created, err := repo.FindByNameInParent(ctx, "Docs", root.ID)
		if err != nil {
			t.Fatalf("expected created folder to be findable, got: %v", err)
		}
		if created.Name != "Docs" {
			t.Fatalf("expected name %q, got %q", "Docs", created.Name)
		}
	})

	t.Run("successful create", func(t *testing.T) {
		repo, root := setup()
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.CreateFolder(ctx, &root.ID, "docs")
		if err != nil {
			t.Fatalf("expected success, got: %v", err)
		}

		created, err := repo.FindByNameInParent(ctx, "docs", root.ID)
		if err != nil {
			t.Fatalf("expected created folder to be findable, got: %v", err)
		}
		if created.OwnerID != ownerA {
			t.Fatalf("expected owner %s, got %s", ownerA, created.OwnerID)
		}
		if created.ParentID == nil || *created.ParentID != root.ID {
			t.Fatalf("expected parent %s, got %v", root.ID, created.ParentID)
		}
	})
}

func TestFolderService_MoveFolder(t *testing.T) {
	ownerA := uuid.New()
	ownerB := uuid.New()

	setup := func() (*fakeFolderRepo, *folder.Folder, *folder.Folder, *folder.Folder) {
		repo := newFakeFolderRepo()
		root := &folder.Folder{ID: uuid.New(), OwnerID: ownerA, ParentID: nil, Name: "root"}
		docs := &folder.Folder{ID: uuid.New(), OwnerID: ownerA, ParentID: &root.ID, Name: "docs"}
		photos := &folder.Folder{ID: uuid.New(), OwnerID: ownerA, ParentID: &root.ID, Name: "photos"}
		repo.seed(root)
		repo.seed(docs)
		repo.seed(photos)
		return repo, root, docs, photos
	}

	t.Run("cannot move root folder", func(t *testing.T) {
		repo, root, docs, _ := setup()
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.MoveFolder(ctx, root.ID, docs.ID)
		if !errors.Is(err, folder.ErrCannotMoveRootFolder) {
			t.Fatalf("expected ErrCannotMoveRootFolder, got: %v", err)
		}
	})

	t.Run("cannot move folder into itself", func(t *testing.T) {
		repo, _, docs, _ := setup()
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.MoveFolder(ctx, docs.ID, docs.ID)
		if !errors.Is(err, folder.ErrCannotMoveIntoItself) {
			t.Fatalf("expected ErrCannotMoveIntoItself, got: %v", err)
		}
	})

	t.Run("cannot move folder into its own descendant", func(t *testing.T) {
		repo, _, docs, _ := setup()
		nested := &folder.Folder{ID: uuid.New(), OwnerID: ownerA, ParentID: &docs.ID, Name: "nested"}
		repo.seed(nested)
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.MoveFolder(ctx, docs.ID, nested.ID)
		if !errors.Is(err, folder.ErrCannotMoveIntoDescendant) {
			t.Fatalf("expected ErrCannotMoveIntoDescendant, got: %v", err)
		}
	})

	t.Run("cannot move into current parent folder", func(t *testing.T) {
		repo, root, docs, _ := setup()
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.MoveFolder(ctx, docs.ID, root.ID)
		if !errors.Is(err, folder.ErrAlreadyInDestination) {
			t.Fatalf("expected ErrAlreadyInDestination, got: %v", err)
		}
	})

	t.Run("cannot move someone else's folder", func(t *testing.T) {
		repo, _, docs, photos := setup()
		docs.OwnerID = ownerB
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.MoveFolder(ctx, docs.ID, photos.ID)
		if !errors.Is(err, auth.ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got: %v", err)
		}
	})

	t.Run("cannot move into a folder owned by someone else", func(t *testing.T) {
		repo, _, docs, photos := setup()
		photos.OwnerID = ownerB
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.MoveFolder(ctx, docs.ID, photos.ID)
		if !errors.Is(err, auth.ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got: %v", err)
		}
	})

	t.Run("cannot move into a folder with a name conflict", func(t *testing.T) {
		repo, _, docs, photos := setup()
		clash := &folder.Folder{ID: uuid.New(), OwnerID: ownerA, ParentID: &photos.ID, Name: docs.Name}
		repo.seed(clash)
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.MoveFolder(ctx, docs.ID, photos.ID)
		if !errors.Is(err, folder.ErrFolderNameConflict) {
			t.Fatalf("expected ErrFolderNameConflict, got: %v", err)
		}
	})

	t.Run("successful move", func(t *testing.T) {
		repo, _, docs, photos := setup()
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		if err := svc.MoveFolder(ctx, docs.ID, photos.ID); err != nil {
			t.Fatalf("expected success, got: %v", err)
		}
		moved, _ := repo.FindByID(ctx, docs.ID)
		if moved.ParentID == nil || *moved.ParentID != photos.ID {
			t.Fatalf("expected docs to be moved under photos, got parent: %v", moved.ParentID)
		}
	})
}

func TestFolderService_ListChildren(t *testing.T) {
	ownerA := uuid.New()
	ownerB := uuid.New()

	setup := func() (*fakeFolderRepo, *folder.Folder) {
		repo := newFakeFolderRepo()
		root := &folder.Folder{ID: uuid.New(), OwnerID: ownerA, ParentID: nil, Name: "root"}
		repo.seed(root)
		return repo, root
	}

	t.Run("cannot list children of someone else's folder", func(t *testing.T) {
		repo, root := setup()
		root.OwnerID = ownerB
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		_, err := svc.ListChildren(ctx, root.ID)
		if !errors.Is(err, auth.ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got: %v", err)
		}
	})

	t.Run("propagates error when folder does not exist", func(t *testing.T) {
		repo, _ := setup()
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		_, err := svc.ListChildren(ctx, uuid.New())
		if !errors.Is(err, folder.ErrFolderNotFound) {
			t.Fatalf("expected ErrFolderNotFound, got: %v", err)
		}
	})

	t.Run("successful list returns subfolders", func(t *testing.T) {
		repo, root := setup()
		child := &folder.Folder{ID: uuid.New(), OwnerID: ownerA, ParentID: &root.ID, Name: "docs"}
		repo.seed(child)
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		content, err := svc.ListChildren(ctx, root.ID)
		if err != nil {
			t.Fatalf("expected success, got: %v", err)
		}
		if len(content.Subfolders) != 1 || content.Subfolders[0].ID != child.ID {
			t.Fatalf("expected 1 subfolder (%s), got: %+v", child.ID, content.Subfolders)
		}
	})
}

func TestFolderService_RenameFolder(t *testing.T) {
	ownerA := uuid.New()
	ownerB := uuid.New()

	setup := func() (*fakeFolderRepo, *folder.Folder) {
		repo := newFakeFolderRepo()
		root := &folder.Folder{ID: uuid.New(), OwnerID: ownerA, ParentID: nil, Name: "root"}
		docs := &folder.Folder{ID: uuid.New(), OwnerID: ownerA, ParentID: &root.ID, Name: "docs"}
		repo.seed(root)
		repo.seed(docs)
		return repo, docs
	}

	t.Run("cannot rename someone else's folder", func(t *testing.T) {
		repo, docs := setup()
		docs.OwnerID = ownerB
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.RenameFolder(ctx, docs.ID, "documents")
		if !errors.Is(err, auth.ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got: %v", err)
		}
	})

	t.Run("rejects invalid name", func(t *testing.T) {
		repo, docs := setup()
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.RenameFolder(ctx, docs.ID, "")
		if err == nil {
			t.Fatal("expected an error for empty name, got nil")
		}
	})

	t.Run("successful rename", func(t *testing.T) {
		repo, docs := setup()
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.RenameFolder(ctx, docs.ID, "documents")
		if err != nil {
			t.Fatalf("expected success, got: %v", err)
		}

		renamed, err := repo.FindByID(ctx, docs.ID)
		if err != nil {
			t.Fatalf("expected folder to still exist: %v", err)
		}
		if renamed.Name != "documents" {
			t.Fatalf("expected name %q, got %q", "documents", renamed.Name)
		}
	})
}

func TestFolderService_DeleteFolder(t *testing.T) {
	ownerA := uuid.New()
	ownerB := uuid.New()

	setup := func() (*fakeFolderRepo, *folder.Folder, *folder.Folder) {
		repo := newFakeFolderRepo()
		root := &folder.Folder{ID: uuid.New(), OwnerID: ownerA, ParentID: nil, Name: "root"}
		docs := &folder.Folder{ID: uuid.New(), OwnerID: ownerA, ParentID: &root.ID, Name: "docs"}
		repo.seed(root)
		repo.seed(docs)
		return repo, root, docs
	}

	t.Run("cannot delete someone else's folder", func(t *testing.T) {
		repo, _, docs := setup()
		docs.OwnerID = ownerB
		svc := service.NewFolderService(repo, newFakeTxRunner(t)) // fails the test if RunTx is ever reached
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.DeleteFolder(ctx, docs.ID)
		if !errors.Is(err, auth.ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got: %v", err)
		}
	})

	t.Run("cannot delete root folder", func(t *testing.T) {
		repo, root, _ := setup()
		svc := service.NewFolderService(repo, newFakeTxRunner(t))
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.DeleteFolder(ctx, root.ID)
		if !errors.Is(err, folder.ErrCannotDeleteRootFolder) {
			t.Fatalf("expected ErrCannotDeleteRootFolder, got: %v", err)
		}
	})

	// Cascade-cleanup (descendant files key enqueueing and folders deletetion)
	// is not covered inside unit tests. It should be verified against a real Postgres
	// inside an integration test via test containers.
}
