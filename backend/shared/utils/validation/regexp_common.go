package utils_validation

import "regexp"

var (
	CyrillicNameRe = regexp.MustCompile(`^[\x{0400}-\x{042F}\x{0490}][\x{0430}-\x{045F}\x{0491}]{1,31}$`)
	PhoneNumberRe  = regexp.MustCompile(`^\+?[0-9]{9,15}$`)
)
