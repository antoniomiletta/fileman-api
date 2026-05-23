package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/antoniomiletta/fileman/internal/api/dto"
	"github.com/antoniomiletta/fileman/internal/api/transport"
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
		transport.WriteError(w, transport.ErrMalformedJSON)
		return
	}
	defer r.Body.Close()

	if err := h.svc.CreateUser(r.Context(), services.CreateUserInput{
		Email:    req.Email,
		Password: req.Password,
	}); err != nil {
		transport.WriteError(w, err)
		return
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		transport.WriteError(w, transport.ErrMalformedJSON)
		return
	}
	defer r.Body.Close()

	token, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	transport.WriteJSON(w, http.StatusOK, token)
}
