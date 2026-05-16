package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/antoniomiletta/fileman/internal/api/dto"
	"github.com/antoniomiletta/fileman/internal/api/transport"
	"github.com/antoniomiletta/fileman/internal/domain"
	"github.com/antoniomiletta/fileman/internal/pkg/authenticator"
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

func (h *FolderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFolderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		transport.WriteError(w, err)
		return
	}
	defer r.Body.Close()

	token := r.Header.Get("Authorization")

	folder := domain.NewFolder(domain.NewFolderParams{
		OwnerID:  authenticator.GetIDFromToken(token),
		ParentID: req.ParentID,
		Name:     req.Name,
	})

	if err := h.svc.Create(r.Context(), &folder); err != nil {
		transport.WriteError(w, err)
	}

	transport.WriteStatus(w, http.StatusCreated)
}

func (h *FolderHandler) ListChildren(w http.ResponseWriter, r *http.Request) {
	folderID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		transport.WriteError(w, err)
	}

	content, err := h.svc.ListChildren(r.Context(), folderID)
	if err != nil {
		transport.WriteError(w, err)
	}

	transport.WriteJSON(w, http.StatusOK, content)
}

func (h *FolderHandler) Move(w http.ResponseWriter, r *http.Request) {
	folderID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		transport.WriteError(w, domain.ErrInvalidToken)
	}

	var newParentIdStr string
	if err := json.NewDecoder(r.Body).Decode(&newParentIdStr); err != nil {
		transport.WriteError(w, err)
		return
	}

	newParentID, err := uuid.Parse(newParentIdStr)
	if err != nil {
		transport.WriteError(w, err)
	}

	if err := h.svc.Move(r.Context(), folderID, newParentID); err != nil {
		transport.WriteError(w, err)
	}
}

func (h *FolderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	folderID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		transport.WriteError(w, err)
	}

	if err := h.svc.Delete(r.Context(), folderID); err != nil {
		transport.WriteError(w, err)
	}
}
