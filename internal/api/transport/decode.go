package transport

import (
	"encoding/json"
	"net/http"
)

func DecodeJSON[T any](r *http.Request, dest *T) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dest); err != nil {
		return err
	}

	return nil
}
