package transport

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Responder struct {
	logger *slog.Logger
}

func NewResponder(logger *slog.Logger) *Responder {
	return &Responder{
		logger: logger,
	}
}

func (r *Responder) WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

func (r *Responder) WriteStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func (r *Responder) WriteError(w http.ResponseWriter, err error) {
	httpErr := MapError(err)

	if httpErr.Status >= 500 {
		r.logger.Error("server error", "error", err.Error())
	} else {
		r.logger.Info("request failed", "error", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.Status)

	json.NewEncoder(w).Encode(map[string]string{
		"error": httpErr.Message,
	})
}
