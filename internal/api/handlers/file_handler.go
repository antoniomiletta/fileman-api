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
func (h *FileHandler) Create(w http.ResponseWriter, r *http.Request)
func (h *FileHandler) Move(w http.ResponseWriter, r *http.Request)
func (h *FileHandler) Delete(w http.ResponseWriter, r *http.Request)
