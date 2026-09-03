package service_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/antoniomiletta/fileman/internal/domain/auth"
	"github.com/antoniomiletta/fileman/internal/domain/file"
	"github.com/antoniomiletta/fileman/internal/domain/folder"
	"github.com/antoniomiletta/fileman/internal/pkg/reqctx"
	"github.com/antoniomiletta/fileman/internal/pkg/storage"
	"github.com/antoniomiletta/fileman/internal/service"
	"github.com/google/uuid"
)

// Any behavior that depends directly on a real service like
// Postgres (transactions, cascading and skipping locked rows) is not verified here.

type fakeFileRepo struct {
	files map[uuid.UUID]*file.File
}

func newFakeFileRepo() *fakeFileRepo {
	return &fakeFileRepo{files: make(map[uuid.UUID]*file.File)}
}

func (f *fakeFileRepo) seed(fl *file.File) {
	f.files[fl.ID] = fl
}

func (f *fakeFileRepo) Create(ctx context.Context, fl *file.File) error {
	f.files[fl.ID] = fl
	return nil
}

func (f *fakeFileRepo) Move(ctx context.Context, id, newParentID uuid.UUID) error {
	fl, ok := f.files[id]
	if !ok {
		return file.ErrFileNotFound
	}
	fl.ParentID = newParentID
	return nil
}

func (f *fakeFileRepo) Rename(ctx context.Context, id uuid.UUID, newName string) error {
	fl, ok := f.files[id]
	if !ok {
		return file.ErrFileNotFound
	}
	fl.Name = newName
	return nil
}

func (f *fakeFileRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if _, ok := f.files[id]; !ok {
		return file.ErrFileNotFound
	}
	delete(f.files, id)
	return nil
}

func (f *fakeFileRepo) FindByID(ctx context.Context, id uuid.UUID) (*file.File, error) {
	fl, ok := f.files[id]
	if !ok {
		return nil, file.ErrFileNotFound
	}
	return fl, nil
}

func (f *fakeFileRepo) FindByNameInParent(ctx context.Context, name string, parentID uuid.UUID) (*file.File, error) {
	for _, fl := range f.files {
		if fl.Name == name && fl.ParentID == parentID {
			return fl, nil
		}
	}
	return nil, file.ErrFileNotFound
}

func (f *fakeFileRepo) MarkUploaded(ctx context.Context, id uuid.UUID) error {
	fl, ok := f.files[id]
	if !ok {
		return file.ErrFileNotFound
	}
	fl.UploadStatus = file.UploadStatusComplete
	return nil
}

// Not exercised in this layer of tests
func (f *fakeFileRepo) ListStale(ctx context.Context, staleTime time.Duration) ([]*file.File, error) {
	return nil, nil
}

type fakeStorageBackend struct {
	uploaded  map[string][]byte
	uploadErr error
}

func newFakeStorageBackend() *fakeStorageBackend {
	return &fakeStorageBackend{uploaded: make(map[string][]byte)}
}

func (f *fakeStorageBackend) Upload(ctx context.Context, fileKey string, r io.Reader, size int64) error {
	if f.uploadErr != nil {
		return f.uploadErr
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	f.uploaded[fileKey] = data
	return nil
}

func (f *fakeStorageBackend) Download(ctx context.Context, fileKey string) (io.ReadCloser, error) {
	data, ok := f.uploaded[fileKey]
	if !ok {
		return nil, storage.ErrObjectNotFound
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func (f *fakeStorageBackend) Delete(ctx context.Context, fileKey string) error {
	delete(f.uploaded, fileKey)
	return nil
}

func (f *fakeStorageBackend) Exists(ctx context.Context, fileKey string) (bool, error) {
	_, ok := f.uploaded[fileKey]
	return ok, nil
}

// Tests
func TestFileService_CreateFile(t *testing.T) {
	ownerA := uuid.New()
	ownerB := uuid.New()

	setup := func() (*fakeFileRepo, *fakeFolderRepo, *fakeStorageBackend, *folder.Folder) {
		fileRepo := newFakeFileRepo()
		folderRepo := newFakeFolderRepo()
		store := newFakeStorageBackend()
		root := &folder.Folder{ID: uuid.New(), OwnerID: ownerA, ParentID: nil, Name: "root"}
		folderRepo.seed(root)
		return fileRepo, folderRepo, store, root
	}

	newSvc := func(fileRepo *fakeFileRepo, folderRepo *fakeFolderRepo, store *fakeStorageBackend) *service.FileService {
		return service.NewFileService(fileRepo, folderRepo, store, newFakeTxRunner(t), nil)
	}

	t.Run("rejects invalid name", func(t *testing.T) {
		fileRepo, folderRepo, store, root := setup()
		svc := newSvc(fileRepo, folderRepo, store)
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.CreateFile(ctx, service.CreateFileInput{
			ParentID: root.ID,
			Name:     "",
			MIMEType: "application/pdf",
			Size:     100,
		}, strings.NewReader("content"))
		if err == nil {
			t.Fatal("expected an error for empty name, got nil")
		}
	})

	t.Run("propagates error when parent does not exist", func(t *testing.T) {
		fileRepo, folderRepo, store, _ := setup()
		svc := newSvc(fileRepo, folderRepo, store)
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.CreateFile(ctx, service.CreateFileInput{
			ParentID: uuid.New(),
			Name:     "resume.pdf",
			MIMEType: "application/pdf",
			Size:     100,
		}, strings.NewReader("content"))
		if !errors.Is(err, folder.ErrFolderNotFound) {
			t.Fatalf("expected ErrFolderNotFound, got: %v", err)
		}
	})

	t.Run("rejects creating in a folder owned by someone else", func(t *testing.T) {
		fileRepo, folderRepo, store, root := setup()
		root.OwnerID = ownerB
		svc := newSvc(fileRepo, folderRepo, store)
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.CreateFile(ctx, service.CreateFileInput{
			ParentID: root.ID,
			Name:     "resume.pdf",
			MIMEType: "application/pdf",
			Size:     100,
		}, strings.NewReader("content"))
		if !errors.Is(err, auth.ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got: %v", err)
		}
	})

	t.Run("rejects name conflict in parent", func(t *testing.T) {
		fileRepo, folderRepo, store, root := setup()
		existing := &file.File{ID: uuid.New(), OwnerID: ownerA, ParentID: root.ID, Name: "resume.pdf"}
		fileRepo.seed(existing)
		svc := newSvc(fileRepo, folderRepo, store)
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.CreateFile(ctx, service.CreateFileInput{
			ParentID: root.ID,
			Name:     "resume.pdf",
			MIMEType: "application/pdf",
			Size:     100,
		}, strings.NewReader("content"))
		if !errors.Is(err, file.ErrFileNameConflict) {
			t.Fatalf("expected ErrFileNameConflict, got: %v", err)
		}
	})

	t.Run("rejects disallowed MIME type", func(t *testing.T) {
		fileRepo, folderRepo, store, root := setup()
		svc := newSvc(fileRepo, folderRepo, store)
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.CreateFile(ctx, service.CreateFileInput{
			ParentID: root.ID,
			Name:     "malware.exe",
			MIMEType: "application/x-msdownload",
			Size:     100,
		}, strings.NewReader("content"))
		if !errors.Is(err, file.ErrFileTypeNotAllowed) {
			t.Fatalf("expected ErrFileTypeNotAllowed, got: %v", err)
		}
	})

	t.Run("successful create assigns a non-zero ID and marks uploaded", func(t *testing.T) {
		fileRepo, folderRepo, store, root := setup()
		svc := newSvc(fileRepo, folderRepo, store)
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.CreateFile(ctx, service.CreateFileInput{
			ParentID: root.ID,
			Name:     "resume.pdf",
			MIMEType: "application/pdf",
			Size:     7,
		}, strings.NewReader("content"))
		if err != nil {
			t.Fatalf("expected success, got: %v", err)
		}

		created, err := fileRepo.FindByNameInParent(ctx, "resume.pdf", root.ID)
		if err != nil {
			t.Fatalf("expected created file to be findable, got: %v", err)
		}

		if created.ID == uuid.Nil {
			t.Fatal("expected a non-zero file ID, got the zero UUID")
		}

		if created.OwnerID != ownerA {
			t.Fatalf("expected owner %s, got %s", ownerA, created.OwnerID)
		}
		if created.UploadStatus != file.UploadStatusComplete {
			t.Fatalf("expected status %q, got %q", file.UploadStatusComplete, created.UploadStatus)
		}
		if _, ok := store.uploaded[created.StorageKey]; !ok {
			t.Fatalf("expected content to be uploaded under key %q", created.StorageKey)
		}
	})

	t.Run("second create does not collide with the first", func(t *testing.T) {
		fileRepo, folderRepo, store, root := setup()
		svc := newSvc(fileRepo, folderRepo, store)
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err1 := svc.CreateFile(ctx, service.CreateFileInput{
			ParentID: root.ID, Name: "a.pdf", MIMEType: "application/pdf", Size: 1,
		}, strings.NewReader("a"))
		err2 := svc.CreateFile(ctx, service.CreateFileInput{
			ParentID: root.ID, Name: "b.pdf", MIMEType: "application/pdf", Size: 1,
		}, strings.NewReader("b"))

		if err1 != nil || err2 != nil {
			t.Fatalf("expected both creates to succeed, got: %v / %v", err1, err2)
		}

		a, _ := fileRepo.FindByNameInParent(ctx, "a.pdf", root.ID)
		b, _ := fileRepo.FindByNameInParent(ctx, "b.pdf", root.ID)
		if a.ID == b.ID {
			t.Fatal("expected distinct IDs for two separately created files, got the same ID")
		}
	})

	t.Run("upload failure leaves the file pending, not marked uploaded", func(t *testing.T) {
		fileRepo, folderRepo, store, root := setup()
		store.uploadErr = errors.New("simulated storage failure")
		svc := newSvc(fileRepo, folderRepo, store)
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.CreateFile(ctx, service.CreateFileInput{
			ParentID: root.ID, Name: "resume.pdf", MIMEType: "application/pdf", Size: 7,
		}, strings.NewReader("content"))
		if err == nil {
			t.Fatal("expected an error from the failed upload, got nil")
		}

		created, findErr := fileRepo.FindByNameInParent(ctx, "resume.pdf", root.ID)
		if findErr != nil {
			t.Fatalf("expected the row to still exist despite upload failure: %v", findErr)
		}
		if created.UploadStatus != file.UploadStatusPending {
			t.Fatalf("expected status to remain %q, got %q", file.UploadStatusPending, created.UploadStatus)
		}
	})
}

func TestFileService_MoveFile(t *testing.T) {
	ownerA := uuid.New()
	ownerB := uuid.New()

	setup := func() (*fakeFileRepo, *fakeFolderRepo, *folder.Folder, *folder.Folder, *file.File) {
		fileRepo := newFakeFileRepo()
		folderRepo := newFakeFolderRepo()
		root := &folder.Folder{ID: uuid.New(), OwnerID: ownerA, ParentID: nil, Name: "root"}
		photos := &folder.Folder{ID: uuid.New(), OwnerID: ownerA, ParentID: &root.ID, Name: "photos"}
		folderRepo.seed(root)
		folderRepo.seed(photos)
		f := &file.File{ID: uuid.New(), OwnerID: ownerA, ParentID: root.ID, Name: "resume.pdf"}
		fileRepo.seed(f)
		return fileRepo, folderRepo, root, photos, f
	}

	newSvc := func(fileRepo *fakeFileRepo, folderRepo *fakeFolderRepo) *service.FileService {
		return service.NewFileService(fileRepo, folderRepo, newFakeStorageBackend(), newFakeTxRunner(t), nil)
	}

	t.Run("cannot move someone else's file", func(t *testing.T) {
		fileRepo, folderRepo, _, photos, f := setup()
		f.OwnerID = ownerB
		svc := newSvc(fileRepo, folderRepo)
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.MoveFile(ctx, f.ID, photos.ID)
		if !errors.Is(err, auth.ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got: %v", err)
		}
	})

	t.Run("cannot move into already parent folder", func(t *testing.T) {
		fileRepo, folderRepo, root, _, f := setup()
		svc := newSvc(fileRepo, folderRepo)
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.MoveFile(ctx, f.ID, root.ID)
		if !errors.Is(err, folder.ErrAlreadyInDestination) {
			t.Fatalf("expected ErrAlreadyInDestination, got: %v", err)
		}
	})

	t.Run("propagates error when destination folder does not exist", func(t *testing.T) {
		fileRepo, folderRepo, _, _, f := setup()
		svc := newSvc(fileRepo, folderRepo)
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.MoveFile(ctx, f.ID, uuid.New())
		if !errors.Is(err, folder.ErrFolderNotFound) {
			t.Fatalf("expected ErrFolderNotFound, got: %v", err)
		}
	})

	t.Run("cannot move into a folder owned by someone else", func(t *testing.T) {
		fileRepo, folderRepo, _, photos, f := setup()
		photos.OwnerID = ownerB
		svc := newSvc(fileRepo, folderRepo)
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.MoveFile(ctx, f.ID, photos.ID)
		if !errors.Is(err, auth.ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got: %v", err)
		}
	})

	t.Run("cannot move into a folder with a name conflict", func(t *testing.T) {
		fileRepo, folderRepo, _, photos, f := setup()
		clash := &file.File{ID: uuid.New(), OwnerID: ownerA, ParentID: photos.ID, Name: f.Name}
		fileRepo.seed(clash)
		svc := newSvc(fileRepo, folderRepo)
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.MoveFile(ctx, f.ID, photos.ID)
		if !errors.Is(err, file.ErrFileNameConflict) {
			t.Fatalf("expected ErrFileNameConflict, got: %v", err)
		}
	})

	t.Run("successful move", func(t *testing.T) {
		fileRepo, folderRepo, _, photos, f := setup()
		svc := newSvc(fileRepo, folderRepo)
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		if err := svc.MoveFile(ctx, f.ID, photos.ID); err != nil {
			t.Fatalf("expected success, got: %v", err)
		}
		moved, _ := fileRepo.FindByID(ctx, f.ID)
		if moved.ParentID != photos.ID {
			t.Fatalf("expected file to be moved under %s, got parent: %v", photos.ID, moved.ParentID)
		}
	})
}

func TestFileService_RenameFile(t *testing.T) {
	ownerA := uuid.New()
	ownerB := uuid.New()

	setup := func() (*fakeFileRepo, *file.File) {
		fileRepo := newFakeFileRepo()
		root := uuid.New()
		f := &file.File{ID: uuid.New(), OwnerID: ownerA, ParentID: root, Name: "resume.pdf"}
		fileRepo.seed(f)
		return fileRepo, f
	}

	newSvc := func(fileRepo *fakeFileRepo) *service.FileService {
		return service.NewFileService(fileRepo, newFakeFolderRepo(), newFakeStorageBackend(), newFakeTxRunner(t), nil)
	}

	t.Run("cannot rename someone else's file", func(t *testing.T) {
		fileRepo, f := setup()
		f.OwnerID = ownerB
		svc := newSvc(fileRepo)
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.RenameFile(ctx, f.ID, "cv.pdf")
		if !errors.Is(err, auth.ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got: %v", err)
		}
	})

	t.Run("rejects invalid name", func(t *testing.T) {
		fileRepo, f := setup()
		svc := newSvc(fileRepo)
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.RenameFile(ctx, f.ID, "")
		if err == nil {
			t.Fatal("expected an error for empty name, got nil")
		}
	})

	t.Run("successful rename", func(t *testing.T) {
		fileRepo, f := setup()
		svc := newSvc(fileRepo)
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.RenameFile(ctx, f.ID, "cv.pdf")
		if err != nil {
			t.Fatalf("expected success, got: %v", err)
		}

		renamed, err := fileRepo.FindByID(ctx, f.ID)
		if err != nil {
			t.Fatalf("expected file to still exist: %v", err)
		}
		if renamed.Name != "cv.pdf" {
			t.Fatalf("expected name %q, got %q", "cv.pdf", renamed.Name)
		}
	})
}

func TestFileService_DeleteFile_Unit(t *testing.T) {
	ownerA := uuid.New()
	ownerB := uuid.New()

	setup := func() (*fakeFileRepo, *file.File) {
		fileRepo := newFakeFileRepo()
		f := &file.File{ID: uuid.New(), OwnerID: ownerA, ParentID: uuid.New(), Name: "resume.pdf"}
		fileRepo.seed(f)
		return fileRepo, f
	}

	t.Run("cannot delete someone else's file", func(t *testing.T) {
		fileRepo, f := setup()
		f.OwnerID = ownerB
		svc := service.NewFileService(fileRepo, newFakeFolderRepo(), newFakeStorageBackend(), newFakeTxRunner(t), nil)
		ctx := reqctx.WithCallerID(context.Background(), ownerA)

		err := svc.DeleteFile(ctx, f.ID)
		if !errors.Is(err, auth.ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got: %v", err)
		}
	})
}
