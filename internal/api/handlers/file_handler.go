package handlers

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/antoniomiletta/fileman/internal/api/dto"
	"github.com/antoniomiletta/fileman/internal/api/transport"
	"github.com/antoniomiletta/fileman/internal/domain/file"
	"github.com/antoniomiletta/fileman/internal/pkg/authenticator"
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

	// ownerId :=
	// size :=
	fileExt := filepath.Ext(req.Name)
	fileName := req.Name[:len(req.Name)-len(fileExt)]

	file := file.NewFile(file.NewFileParams{
		OwnerID:   authenticator.GetIDFromToken(r.Header.Get("Authorization")),
		ParentID:  req.ParentID,
		Name:      fileName,
		Extension: fileExt,
		Size:      0,
	})

	if err := h.svc.Create(r.Context(), &file); err != nil {
		transport.WriteError(w, err)
	}

	transport.WriteStatus(w, http.StatusCreated)
}

// refactor all from here ->
func (h *FileHandler) Move(w http.ResponseWriter, r *http.Request) {
	fileID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
	}

	var newParentIdStr string
	if err := json.NewDecoder(r.Body).Decode(&newParentIdStr); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	newParentID, err := uuid.Parse(newParentIdStr)
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
	}

	if err := h.svc.Move(r.Context(), fileID, newParentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *FileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	fileID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
	}

	if err := h.svc.Delete(r.Context(), fileID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
