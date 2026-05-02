package handlers

import (
	"net/http"

	"github.com/antoniomiletta/fileman/internal/services"
)

type FileHandler struct {
	svc *services.FileService
}

func NewFileHandler(svc *services.FileService) *FileHandler {
	return &FileHandler{
		svc: svc,
	}
}

// redirect to service
func (h *FileHandler) CreateFile(w http.ResponseWriter, r *http.Request)
func (h *FileHandler) ListFromFolder(w http.ResponseWriter, r *http.Request)
