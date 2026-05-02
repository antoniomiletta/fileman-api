package handlers

import (
	"net/http"

	"github.com/antoniomiletta/fileman/internal/services"
)

type FolderHandler struct {
	svc *services.FolderService
}

func NewFolderHandler(svc *services.FolderService) *FolderHandler {
	return &FolderHandler{
		svc: svc,
	}
}

// redirect to service
func (h *FolderHandler) CreateFolder(w http.ResponseWriter, r *http.Request)
func (h *FolderHandler) ListChildren(w http.ResponseWriter, r *http.Request)
