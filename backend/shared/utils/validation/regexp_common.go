package utils_validation

import "regexp"

var (
	CyrillicNameRe = regexp.MustCompile(`^[^\P{Cyrillic}\P{Lu}][^\P{Cyrillic}\P{Ll}]{1,31}$`)
	PhoneNumberRe  = regexp.MustCompile(`^\+?[0-9]{9,15}$`)
)
