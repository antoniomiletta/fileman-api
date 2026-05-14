package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/antoniomiletta/fileman/internal/domain"
	"github.com/antoniomiletta/fileman/internal/services"
	"github.com/google/uuid"
)

type FolderHandler struct {
	svc *services.FolderService
}

func NewFolderHandler(svc *services.FolderService) *FolderHandler {
	return &FolderHandler{
		svc: svc,
	}
}

// refactor for req type and New() for Domain mapping
func (h *FolderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var folder domain.Folder

	if err := json.NewDecoder(r.Body).Decode(&folder); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := h.svc.Create(r.Context(), &folder); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *FolderHandler) ListChildren(w http.ResponseWriter, r *http.Request) {
	folderID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
	}

	content, err := h.svc.ListChildren(r.Context(), folderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(content)
}

func (h *FolderHandler) Move(w http.ResponseWriter, r *http.Request) {
	folderID, err := uuid.Parse(r.PathValue("id"))
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

	if err := h.svc.Move(r.Context(), folderID, newParentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *FolderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	folderID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
	}

	if err := h.svc.Delete(r.Context(), folderID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
