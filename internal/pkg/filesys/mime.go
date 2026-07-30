package filesys

import (
	"io"
	"net/http"
)

// DetectMIME reads from r and returns the detected MIME type.
// Defaults to 'application/octet-stream' if it cannot detect a more specific one.
func DetectMIME(r io.ReaderAt) (string, error) {
	buf := make([]byte, 512)

	n, err := r.ReadAt(buf, 0)
	if err != nil && err != io.EOF {
		return "", err
	}

	return http.DetectContentType(buf[:n]), nil
}
