package folder

import (
	"fmt"

	"github.com/antoniomiletta/fileman/internal/pkg/filesys"
)

func ValidateFolderName(name string) error {
	reason, ok := filesys.ValidateName(name)
	if !ok {
		switch reason {
		case filesys.ErrorTypeRequired:
			return ErrFolderNameRequired

		case filesys.ErrorTypeTooLong:
			return fmt.Errorf("%w: name cannot be longer than 255 characters", ErrFolderNameTooLong)

		case filesys.ErrorTypePathTraversal:
			return fmt.Errorf("%w: name cannot contain path structures", ErrFolderNameInvalid)

		case filesys.ErrorTypeIllegalChars:
			return fmt.Errorf("%w: name cannot contain control characters", ErrFolderNameInvalid)

		case filesys.ErrorTypeReservedName:
			return fmt.Errorf("%w: \"%s\" is a reserved system name", ErrFolderNameInvalid, name)

		case filesys.ErrorTypeIllegalTrailling:
			return fmt.Errorf("%w: name cannot end with a period or space", ErrFolderNameInvalid)
		}
	}

	return nil
}
