package transport

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

func WriteStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func WriteError(w http.ResponseWriter, err error) {
	mapped := MapError(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(mapped.Status)

	json.NewEncoder(w).Encode(map[string]string{
		"error": mapped.Message,
	})
}
