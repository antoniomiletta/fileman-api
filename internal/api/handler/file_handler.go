package handler

import (
	"net/http"

	"github.com/antoniomiletta/fileman/internal/api/dto"
	"github.com/antoniomiletta/fileman/internal/api/transport"
	"github.com/antoniomiletta/fileman/internal/pkg/filesys"
	"github.com/antoniomiletta/fileman/internal/service"
	"github.com/google/uuid"
)

type FileHandler struct {
	svc  *service.FileService
	resp *transport.Responder
}

func NewFileHandler(svc *service.FileService, resp *transport.Responder) *FileHandler {
	return &FileHandler{
		svc:  svc,
		resp: resp,
	}
}

func (h *FileHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.resp.WriteError(w, err)
		return
	}

	parentID, err := uuid.Parse(r.FormValue(dto.MultipartFieldParentID))
	if err != nil {
		h.resp.WriteError(w, transport.ErrInvalidPathParam)
		return
	}

	fileContent, fileHeader, err := r.FormFile(dto.MultipartFieldFile)
	if err != nil {
		h.resp.WriteError(w, err)
		return
	}
	defer fileContent.Close()

	mimeType, err := filesys.DetectMIME(fileContent)
	if err != nil {
		h.resp.WriteError(w, err)
		return
	}

	if err := h.svc.CreateFile(ctx, service.CreateFileInput{
		ParentID: parentID,
		Name:     fileHeader.Filename,
		MIMEType: mimeType,
		Size:     fileHeader.Size,
	}, fileContent); err != nil {
		h.resp.WriteError(w, err)
		return
	}

	h.resp.WriteStatus(w, http.StatusCreated)
}

func (h *FileHandler) Move(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fileID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.resp.WriteError(w, transport.ErrInvalidPathParam)
		return
	}

	var req dto.MoveFileRequest
	if err := transport.DecodeJSON(r, &req); err != nil {
		h.resp.WriteError(w, transport.ErrInvalidJSON)
		return
	}

	newParentID, err := uuid.Parse(req.NewParentID.String())
	if err != nil {
		h.resp.WriteError(w, transport.ErrInvalidJSON)
		return
	}

	if err := h.svc.MoveFile(ctx, fileID, newParentID); err != nil {
		h.resp.WriteError(w, err)
		return
	}

	h.resp.WriteStatus(w, http.StatusNoContent)
}

func (h *FileHandler) Rename(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fileID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.resp.WriteError(w, transport.ErrInvalidPathParam)
		return
	}

	var req dto.RenameFileRequest
	if err := transport.DecodeJSON(r, &req); err != nil {
		h.resp.WriteError(w, transport.ErrInvalidJSON)
		return
	}

	if err := h.svc.RenameFile(ctx, fileID, req.NewName); err != nil {
		h.resp.WriteError(w, err)
		return
	}

	h.resp.WriteStatus(w, http.StatusNoContent)
}

func (h *FileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fileID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.resp.WriteError(w, transport.ErrInvalidPathParam)
		return
	}

	if err := h.svc.DeleteFile(ctx, fileID); err != nil {
		h.resp.WriteError(w, err)
		return
	}

	h.resp.WriteStatus(w, http.StatusNoContent)
}
