package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/antoniomiletta/fileman/internal/domain"
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
	var file domain.File

	if err := json.NewDecoder(r.Body).Decode(&file); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := h.svc.Create(r.Context(), &file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

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
