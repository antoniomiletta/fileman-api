package storage

import (
	"fmt"
	"mime"

	"github.com/google/uuid"
)

type ResourceType string

const (
	ResourceTypeImage    ResourceType = "images"
	ResourceTypeDocument ResourceType = "documents"
	ResourceTypeArchive  ResourceType = "archives"
	ResourceTypeAudio    ResourceType = "audio"
	ResourceTypeVideo    ResourceType = "video"
)

// mimeResourceMap maps all the system's allowed MIME types
// to their correspondent resource type in storage.
var mimeResourceMap = map[string]ResourceType{
	// images
	"image/jpeg":    ResourceTypeImage,
	"image/png":     ResourceTypeImage,
	"image/gif":     ResourceTypeImage,
	"image/webp":    ResourceTypeImage,
	"image/svg+xml": ResourceTypeImage,

	// documents
	"application/pdf":          ResourceTypeDocument,
	"text/plain":               ResourceTypeDocument,
	"text/csv":                 ResourceTypeDocument,
	"application/msword":       ResourceTypeDocument,
	"application/vnd.ms-excel": ResourceTypeDocument,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ResourceTypeDocument,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":       ResourceTypeDocument,

	// archives
	"application/zip":  ResourceTypeArchive,
	"application/gzip": ResourceTypeArchive,

	// audio/video
	"audio/mpeg": ResourceTypeAudio,
	"video/mp4":  ResourceTypeVideo,
}

// Classify reports whether the provided MIME type maps to a valid
// resource type in storage, returning the resource type if it does.
// It can be used to determine if the MIME type is accepted by the system,
// since every allowed MIME type is mapped and all mapped types are allowed.
func Classify(mediaType string) (ResourceType, bool) {
	mimeType, _, err := mime.ParseMediaType(mediaType)
	if err != nil {
		return "", false
	}
	rt, ok := mimeResourceMap[mimeType]
	return rt, ok
}

type FileKeyParams struct {
	OwnerID      uuid.UUID
	ResourceType ResourceType
	FileID       uuid.UUID
}

// GenerateFileKey builds a {ownerId}/{resourceType}/{fileId} filekey.
// Owner first so all of a user's objects share a common prefix in storage.
func GenerateFileKey(p FileKeyParams) string {
	return fmt.Sprintf("%s/%s/%s", p.OwnerID, p.ResourceType, p.FileID)
}
