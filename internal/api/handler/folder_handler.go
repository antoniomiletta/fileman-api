package handler

import (
	"encoding/json"
	"net/http"

	"github.com/antoniomiletta/fileman/internal/api/dto"
	"github.com/antoniomiletta/fileman/internal/api/transport"
	"github.com/antoniomiletta/fileman/internal/pkg/reqctx"
	"github.com/antoniomiletta/fileman/internal/service"
	"github.com/google/uuid"
)

type FolderHandler struct {
	svc *service.FolderService
}

func NewFolderHandler(svc *service.FolderService) *FolderHandler {
	return &FolderHandler{
		svc: svc,
	}
}

func (h *FolderHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.CreateFolderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		transport.WriteError(w, transport.ErrMalformedJSON)
		return
	}
	defer r.Body.Close()

	callerID, err := reqctx.CallerIDFrom(ctx)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	if err := h.svc.CreateFolder(ctx, service.CreateFolderInput{
		OwnerID:  callerID,
		ParentID: req.ParentID,
		Name:     req.Name,
	}); err != nil {
		transport.WriteError(w, err)
		return
	}

	transport.WriteStatus(w, http.StatusCreated)
}

func (h *FolderHandler) ListChildren(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	folderID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		transport.WriteError(w, transport.ErrInvalidPathParam)
		return
	}

	content, err := h.svc.ListChildren(ctx, folderID)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	transport.WriteJSON(w, http.StatusOK, content)
}

func (h *FolderHandler) Move(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	folderID, err := uuid.Parse(r.PathValue("id"))
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

	if err := h.svc.MoveFolder(ctx, folderID, newParentID); err != nil {
		transport.WriteError(w, err)
		return
	}

	transport.WriteStatus(w, http.StatusNoContent)
}

func (h *FolderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	folderID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		transport.WriteError(w, transport.ErrInvalidPathParam)
		return
	}

	if err := h.svc.DeleteFolder(ctx, folderID); err != nil {
		transport.WriteError(w, err)
		return
	}

	transport.WriteStatus(w, http.StatusNoContent)
}
