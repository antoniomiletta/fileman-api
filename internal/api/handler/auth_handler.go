package handler

import (
	"net/http"

	"github.com/antoniomiletta/fileman/internal/api/dto"
	"github.com/antoniomiletta/fileman/internal/api/transport"
	"github.com/antoniomiletta/fileman/internal/service"
)

type AuthHandler struct {
	svc  *service.AuthService
	resp *transport.Responder
}

func NewAuthHandler(svc *service.AuthService, resp *transport.Responder) *AuthHandler {
	return &AuthHandler{
		svc:  svc,
		resp: resp,
	}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.SignUpRequest

	if err := transport.DecodeJSON(r, &req); err != nil {
		h.resp.WriteError(w, transport.ErrInvalidJSON)
		return
	}
	defer r.Body.Close()

	if err := h.svc.CreateUser(ctx, service.CreateUserInput{
		Email:    req.Email,
		Password: req.Password,
	}); err != nil {
		h.resp.WriteError(w, err)
		return
	}
	h.resp.WriteStatus(w, http.StatusOK)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.LoginRequest

	if err := transport.DecodeJSON(r, &req); err != nil {
		h.resp.WriteError(w, transport.ErrInvalidJSON)
		return
	}
	defer r.Body.Close()

	token, err := h.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		h.resp.WriteError(w, err)
		return
	}

	h.resp.WriteJSON(w, http.StatusOK, token)
}
