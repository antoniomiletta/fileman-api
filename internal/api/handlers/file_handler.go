package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/antoniomiletta/fileman/internal/api/dto"
	"github.com/antoniomiletta/fileman/internal/api/transport"
	"github.com/antoniomiletta/fileman/internal/domain/file"
	"github.com/antoniomiletta/fileman/internal/pkg/authenticator"
	"github.com/antoniomiletta/fileman/internal/pkg/filesys"
	"github.com/antoniomiletta/fileman/internal/services"
	"github.com/google/uuid"
)

type FileHandler struct {
	svc *services.FileService
}

func NewFileHandler(svc *services.FileService) *FileHandler {
	return &FileHandler{
		svc: svc,
	}
}

func (h *FileHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		transport.WriteError(w, transport.ErrMalformedJSON)
		return
	}
	defer r.Body.Close()

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	fileContent, fileHeader, err := r.FormFile(transport.MultipartFieldFile)
	if err != nil {
		transport.WriteError(w, err)
		return
	}
	defer fileContent.Close()

	file := file.NewFile(file.NewFileParams{
		OwnerID:  authenticator.GetIDFromToken(r.Header.Get("Authorization")),
		ParentID: req.ParentID,
		Name:     fileHeader.Filename,
		MIMEType: filesys.DetectMIME(fileContent),
		Size:     fileHeader.Size,
	})

	if err := h.svc.Create(r.Context(), &file, fileContent); err != nil {
		transport.WriteError(w, err)
	}

	transport.WriteStatus(w, http.StatusCreated)
}

func (h *FileHandler) Move(w http.ResponseWriter, r *http.Request) {
	fileID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		transport.WriteError(w, transport.ErrMalformedToken)
	}

	var newParentIdStr string
	if err := json.NewDecoder(r.Body).Decode(&newParentIdStr); err != nil {
		transport.WriteError(w, transport.ErrMalformedJSON)
		return
	}

	newParentID, err := uuid.Parse(newParentIdStr)
	if err != nil {
		transport.WriteError(w, transport.ErrMalformedToken)
	}

	if err := h.svc.Move(r.Context(), fileID, newParentID); err != nil {
		transport.WriteError(w, err)
	}
}

func (h *FileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	fileID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		transport.WriteError(w, err)
	}

	if err := h.svc.Delete(r.Context(), fileID); err != nil {
		transport.WriteError(w, err)
	}
}
