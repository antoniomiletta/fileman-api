package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/antoniomiletta/fileman/internal/api/dto"
	"github.com/antoniomiletta/fileman/internal/api/transport"
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
		transport.WriteError(w, transport.ErrMalformedJSON)
		return
	}
	defer r.Body.Close()

	if err := h.svc.CreateFolder(r.Context(), services.CreateFolderInput{
		OwnerID:  authenticator.GetIDFromToken(r.Header.Get("Authorization")),
		ParentID: req.ParentID,
		Name:     req.Name,
	}); err != nil {
		transport.WriteError(w, err)
	}

	transport.WriteStatus(w, http.StatusCreated)
}

func (h *FolderHandler) ListChildren(w http.ResponseWriter, r *http.Request) {
	folderID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		transport.WriteError(w, transport.ErrInvalidPathParam)
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
		transport.WriteError(w, transport.ErrInvalidPathParam)
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

	if err := h.svc.MoveFolder(r.Context(), folderID, newParentID); err != nil {
		transport.WriteError(w, err)
	}
}

func (h *FolderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	folderID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		transport.WriteError(w, transport.ErrInvalidPathParam)
	}

	if err := h.svc.DeleteFolder(r.Context(), folderID); err != nil {
		transport.WriteError(w, err)
	}
}
