package filesys

import (
	"path/filepath"
	"regexp"
	"strings"
)

var (
	illegalCharsRegex  = regexp.MustCompile(`[<>:"/\\|?*%\x00-\x1F]`)
	reservedNamesRegex = regexp.MustCompile(`^(?i)(CON|PRN|AUX|NUL|COM[1-9]|LPT[1-9])(\..*)?$`)
)

const (
	MaxFileNameLength  = 255
	MaxFileSizeAllowed = 1000 * 1024 * 1024 // 1GB
)

type FSValidationError string

const (
	ErrorTypeRequired         FSValidationError = "REQUIRED"
	ErrorTypeTooLong          FSValidationError = "TOO_LONG"
	ErrorTypePathTraversal    FSValidationError = "PATH_TRAVERSAL"
	ErrorTypeIllegalChars     FSValidationError = "ILLEGAL_CHARS"
	ErrorTypeReservedName     FSValidationError = "RESERVED_NAME"
	ErrorTypeIllegalTrailling FSValidationError = "ILLEGAL_TRAILLING"
)

// ValidateName reports whether the name provided is a valid file system name,
// returning the specific FSValidationError if it's not
func ValidateName(name string) (FSValidationError, bool) {
	name = strings.TrimSpace(name)

	if name == "" {
		return ErrorTypeRequired, false
	}

	if len(name) > MaxFileNameLength {
		return ErrorTypeTooLong, false
	}

	if filepath.Base(name) != name {
		return ErrorTypePathTraversal, false
	}

	if illegalCharsRegex.MatchString(name) {
		return ErrorTypeIllegalChars, false
	}

	if reservedNamesRegex.MatchString(name) {
		return ErrorTypeReservedName, false
	}

	if strings.HasSuffix(name, ".") || strings.HasSuffix(name, " ") {
		return ErrorTypeIllegalTrailling, false
	}

	return "", true
}
