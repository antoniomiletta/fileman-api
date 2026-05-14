package transport

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, err error) error {
	httpErr := MapError(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.Status)

	return json.NewEncoder(w).Encode(map[string]string{
		"error": httpErr.Message,
	})
}
