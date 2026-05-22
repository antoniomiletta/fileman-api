package file

import (
	"fmt"

	"github.com/antoniomiletta/fileman/internal/pkg/filesys"
)

func ValidateFileName(name string) error {
	reason, ok := filesys.ValidateName(name)
	if !ok {
		switch reason {
		case filesys.ErrorTypeRequired:
			return ErrFileNameRequired

		case filesys.ErrorTypeTooLong:
			return fmt.Errorf("%w: name cannot be longer than 255 characters", ErrFileNameTooLong)

		case filesys.ErrorTypePathTraversal:
			return fmt.Errorf("%w: name cannot contain path structures", ErrFileNameInvalid)

		case filesys.ErrorTypeIllegalChars:
			return fmt.Errorf("%w: name cannot contain control characters", ErrFileNameInvalid)

		case filesys.ErrorTypeReservedName:
			return fmt.Errorf("%w: \"%s\" is a reserved system name", ErrFileNameInvalid, name)

		case filesys.ErrorTypeIllegalTrailling:
			return fmt.Errorf("%w: name cannot end with a period or space", ErrFileNameInvalid)
		}
	}

	return nil
}

func IsAllowedMIME(mimeType string) bool {
	return nil
}
