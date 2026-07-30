package handler

import (
	"encoding/json"
	"net/http"

	"github.com/antoniomiletta/fileman/internal/api/dto"
	"github.com/antoniomiletta/fileman/internal/api/transport"
	"github.com/antoniomiletta/fileman/internal/pkg/filesys"
	"github.com/antoniomiletta/fileman/internal/pkg/reqctx"
	"github.com/antoniomiletta/fileman/internal/service"
	"github.com/google/uuid"
)

type FileHandler struct {
	svc *service.FileService
}

func NewFileHandler(svc *service.FileService) *FileHandler {
	return &FileHandler{
		svc: svc,
	}
}

func (h *FileHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		transport.WriteError(w, err)
		return
	}

	parentID, err := uuid.Parse(r.FormValue(dto.MultipartFieldParentID))
	if err != nil {
		transport.WriteError(w, transport.ErrInvalidPathParam)
		return
	}

	fileContent, fileHeader, err := r.FormFile(dto.MultipartFieldFile)
	if err != nil {
		transport.WriteError(w, err)
		return
	}
	defer fileContent.Close()

	callerID, err := reqctx.CallerIDFrom(ctx)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	mimeType, err := filesys.DetectMIME(fileContent)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	if err := h.svc.CreateFile(ctx, service.CreateFileInput{
		OwnerID:  callerID,
		ParentID: parentID,
		Name:     fileHeader.Filename,
		MIMEType: mimeType,
		Size:     fileHeader.Size,
	}, fileContent); err != nil {
		transport.WriteError(w, err)
		return
	}

	transport.WriteStatus(w, http.StatusCreated)
}

func (h *FileHandler) Move(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fileID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		transport.WriteError(w, transport.ErrInvalidPathParam)
		return
	}

	var newParentIdStr string
	if err := json.NewDecoder(r.Body).Decode(&newParentIdStr); err != nil {
		transport.WriteError(w, transport.ErrMalformedJSON)
		return
	}

	newParentID, err := uuid.Parse(newParentIdStr)
	if err != nil {
		transport.WriteError(w, transport.ErrMalformedToken)
		return
	}

	if err := h.svc.MoveFile(ctx, fileID, newParentID); err != nil {
		transport.WriteError(w, err)
		return
	}

	transport.WriteStatus(w, http.StatusNoContent)
}

func (h *FileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fileID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		transport.WriteError(w, transport.ErrInvalidPathParam)
		return
	}

	if err := h.svc.DeleteFile(ctx, fileID); err != nil {
		transport.WriteError(w, err)
		return
	}

	transport.WriteStatus(w, http.StatusNoContent)
}
