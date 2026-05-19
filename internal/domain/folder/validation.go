package folder

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// Windows/Linux/macOS forbidden characters and control characters
	illegalCharsRegex = regexp.MustCompile(`[<>:"/\\|?*%\x00-\x1F]`)
	// Windows reserved system names
	reservedNamesRegex = regexp.MustCompile(`(?i)^(CON|PRN|AUX|NUL|COM[1-9]|LPT[1-9])(\..*)?$`)
)

func ValidateFolderName(name string) error {
	if name == "" {
		return ErrFolderNameRequired
	}

	if len(name) > MaxFolderNameLength {
		return fmt.Errorf("%w: name cannot be longer than 255 characters", ErrFolderNameTooLong)
	}

	if illegalCharsRegex.MatchString(name) {
		return fmt.Errorf("%w: name cannot contain control characters", ErrFolderNameInvalid)
	}

	if reservedNamesRegex.MatchString(name) {
		return fmt.Errorf("%w: \"%s\" is a reserved system name", ErrFolderNameInvalid, name)
	}

	if strings.HasSuffix(name, ".") {
		return fmt.Errorf("%w: name cannot end with a period", ErrFolderNameInvalid)
	}

	if name == "." || name == ".." {
		return fmt.Errorf("%w: folder name cannot be \"%s\"", ErrFolderNameInvalid, name)
	}

	return nil
}

