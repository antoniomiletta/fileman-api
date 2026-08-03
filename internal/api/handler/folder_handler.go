package handler

import (
	"net/http"

	"github.com/antoniomiletta/fileman/internal/api/dto"
	"github.com/antoniomiletta/fileman/internal/api/transport"
	"github.com/antoniomiletta/fileman/internal/pkg/authenticator"
	"github.com/antoniomiletta/fileman/internal/service"
	"github.com/google/uuid"
)

type FolderHandler struct {
	svc  *service.FolderService
	resp *transport.Responder
}

func NewFolderHandler(svc *service.FolderService, resp *transport.Responder) *FolderHandler {
	return &FolderHandler{
		svc:  svc,
		resp: resp,
	}
}

func (h *FolderHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateFolderRequest
	if err := transport.DecodeJSON(r, &req); err != nil {
		h.resp.WriteError(w, transport.ErrInvalidJSON)
		return
	}
	defer r.Body.Close()

	if err := h.svc.CreateFolder(ctx, req.ParentID, req.Name); err != nil {
		h.resp.WriteError(w, err)
		return
	}

	h.resp.WriteStatus(w, http.StatusCreated)
}

func (h *FolderHandler) ListChildren(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	folderID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.resp.WriteError(w, transport.ErrInvalidPathParam)
		return
	}

	content, err := h.svc.ListChildren(ctx, folderID)
	if err != nil {
		h.resp.WriteError(w, err)
		return
	}

	h.resp.WriteJSON(w, http.StatusOK, content)
}

func (h *FolderHandler) Move(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	folderID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.resp.WriteError(w, transport.ErrInvalidPathParam)
		return
	}

	var req dto.MoveFolderRequest
	if err := transport.DecodeJSON(r, &req); err != nil {
		h.resp.WriteError(w, transport.ErrInvalidJSON)
		return
	}

	newParentID, err := uuid.Parse(req.NewParentID.String())
	if err != nil {
		h.resp.WriteError(w, authenticator.ErrInvalidToken)
		return
	}

	if err := h.svc.MoveFolder(ctx, folderID, newParentID); err != nil {
		h.resp.WriteError(w, err)
		return
	}

	h.resp.WriteStatus(w, http.StatusNoContent)
}

func (h *FolderHandler) Rename(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	folderID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.resp.WriteError(w, transport.ErrInvalidPathParam)
		return
	}

	var req dto.RenameFolderRequest
	if err := transport.DecodeJSON(r, &req); err != nil {
		h.resp.WriteError(w, transport.ErrInvalidJSON)
		return
	}

	if err := h.svc.RenameFolder(ctx, folderID, req.NewName); err != nil {
		h.resp.WriteError(w, err)
		return
	}

	h.resp.WriteStatus(w, http.StatusNoContent)
}

func (h *FolderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	folderID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.resp.WriteError(w, transport.ErrInvalidPathParam)
		return
	}

	if err := h.svc.DeleteFolder(ctx, folderID); err != nil {
		h.resp.WriteError(w, err)
		return
	}

	h.resp.WriteStatus(w, http.StatusNoContent)
}
