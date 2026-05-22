package filesys

import (
	"path/filepath"
	"strings"
)

func ValidateName(name string) (FSValidationError, bool) {
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
