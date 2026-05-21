package folder

import (
	"fmt"

	"github.com/antoniomiletta/fileman/internal/pkg/fsvalidator"
)

func ValidateFolderName(name string) error {
	reason, ok := fsvalidator.ValidateName(name)
	if !ok {
		switch reason {
		case fsvalidator.ErrorTypeRequired:
			return ErrFolderNameRequired

		case fsvalidator.ErrorTypeTooLong:
			return fmt.Errorf("%w: name cannot be longer than 255 characters", ErrFolderNameTooLong)

		case fsvalidator.ErrorTypePathTraversal:
			return fmt.Errorf("%w: name cannot contain path structures", ErrFolderNameInvalid)

		case fsvalidator.ErrorTypeIllegalChars:
			return fmt.Errorf("%w: name cannot contain control characters", ErrFolderNameInvalid)

		case fsvalidator.ErrorTypeReservedName:
			return fmt.Errorf("%w: \"%s\" is a reserved system name", ErrFolderNameInvalid, name)

		case fsvalidator.ErrorTypeIllegalTrailling:
			return fmt.Errorf("%w: name cannot end with a period or space", ErrFolderNameInvalid)
		}
	}

	return nil
}
