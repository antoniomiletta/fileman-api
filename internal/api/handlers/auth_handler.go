package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/antoniomiletta/fileman/internal/api/dto"
	"github.com/antoniomiletta/fileman/internal/domain"
	"github.com/antoniomiletta/fileman/internal/services"
)

type AuthHandler struct {
	svc *services.AuthService
}

func NewAuthHandler(svc *services.AuthService) *AuthHandler {
	return &AuthHandler{
		svc: svc,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user := domain.NewUser(domain.NewUserParams{
		Email:    req.Email,
		Password: req.Password,
	})

	if err := h.svc.Register(r.Context(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := h.svc.Login(r.Context(), req.Email, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
