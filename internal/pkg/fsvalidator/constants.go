package fsvalidator

import "regexp"

var (
	illegalCharsRegex  = regexp.MustCompile(`[<>:"/\\|?*%\x00-\x1F]`)
	reservedNamesRegex = regexp.MustCompile(`^(?i)(CON|PRN|AUX|NUL|COM[1-9]|LPT[1-9])(\..*)?$`)
)

const (
	MaxFileNameLength  = 255
	MaxFileSizeAllowed = 1000 * 1024 * 1024 // 1GB
)
