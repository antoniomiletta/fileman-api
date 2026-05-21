package file

import (
	"fmt"

	"github.com/antoniomiletta/fileman/internal/pkg/fsvalidator"
)

func ValidateFileName(name string) error {
	reason, ok := fsvalidator.ValidateName(name)
	if !ok {
		switch reason {
		case fsvalidator.ErrorTypeRequired:
			return ErrFileNameRequired

		case fsvalidator.ErrorTypeTooLong:
			return fmt.Errorf("%w: name cannot be longer than 255 characters", ErrFileNameTooLong)

		case fsvalidator.ErrorTypePathTraversal:
			return fmt.Errorf("%w: name cannot contain path structures", ErrFileNameInvalid)

		case fsvalidator.ErrorTypeIllegalChars:
			return fmt.Errorf("%w: name cannot contain control characters", ErrFileNameInvalid)

		case fsvalidator.ErrorTypeReservedName:
			return fmt.Errorf("%w: \"%s\" is a reserved system name", ErrFileNameInvalid, name)

		case fsvalidator.ErrorTypeIllegalTrailling:
			return fmt.Errorf("%w: name cannot end with a period or space", ErrFileNameInvalid)
		}
	}

	return nil
}
