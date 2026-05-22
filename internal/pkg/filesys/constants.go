package filesys

import "regexp"

type FSValidationError string

const (
	ErrorTypeRequired         FSValidationError = "REQUIRED"
	ErrorTypeTooLong          FSValidationError = "TOO_LONG"
	ErrorTypePathTraversal    FSValidationError = "PATH_TRAVERSAL"
	ErrorTypeIllegalChars     FSValidationError = "ILLEGAL_CHARS"
	ErrorTypeReservedName     FSValidationError = "RESERVED_NAME"
	ErrorTypeIllegalTrailling FSValidationError = "ILLEGAL_TRAILLING"
)

const (
	MaxFileNameLength  = 255
	MaxFileSizeAllowed = 1000 * 1024 * 1024 // 1GB
)

var (
	illegalCharsRegex  = regexp.MustCompile(`[<>:"/\\|?*%\x00-\x1F]`)
	reservedNamesRegex = regexp.MustCompile(`^(?i)(CON|PRN|AUX|NUL|COM[1-9]|LPT[1-9])(\..*)?$`)
)
