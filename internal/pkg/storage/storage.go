package storage

import "github.com/google/uuid"

type FileKeyParams struct {
	OwnerID      uuid.UUID
	ResourceType string
	Filename     string
}

func GenerateFileKey(p FileKeyParams) string
func ResourceTypeFrom(mimeType string) string
