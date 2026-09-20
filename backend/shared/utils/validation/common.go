package utils_validation

import (
	"fmt"
	"strings"
)

func ValidateName(v string) error {
	if !Between(v, 2, 32) {
		return fmt.Errorf("'name' must be between 2 and 32 symbols")
	} else if !OnlyCyrillic(v) {
		return fmt.Errorf("'name' must contain only cyrillic symbols")
	} else if !IsFirstLetterUppercase(v) {
		return fmt.Errorf("'name' first letter must be uppercase")
	} else if !IsRestLettersLowercase(v) {
		return fmt.Errorf("'name' letters after first must be lowercase")
	} else if !CyrillicNameRe.MatchString(v) {
		return fmt.Errorf("'name' contains a letter that is not supported")
	}

	return nil
}

func ValidateSurname(v string) error {
	if !Between(v, 2, 32) {
		return fmt.Errorf("'surname' must be between 2 and 32 symbols")
	} else if !OnlyCyrillic(v) {
		return fmt.Errorf("'surname' must contain only cyrillic symbols")
	} else if !IsFirstLetterUppercase(v) {
		return fmt.Errorf("'surname' first letter must be uppercase")
	} else if !IsRestLettersLowercase(v) {
		return fmt.Errorf("'surname' letters after first must be lowercase")
	} else if !CyrillicNameRe.MatchString(v) {
		return fmt.Errorf("'surname' contains a letter that is not supported")
	}

	return nil
}

func ValidatePhoneNumber(v string) error {
	if !Between(v, 12, 15) {
		return fmt.Errorf("'phone number' must be between 12 and 15 symbols")
	} else if !strings.HasPrefix(v, "+") {
		return fmt.Errorf("'phone number' must start with '+'")
	} else if !OnlyDigits(v[1:]) {
		return fmt.Errorf("'phone number' must contain only digits after the '+'")
	} else if !PhoneNumberRe.MatchString(v) {
		return fmt.Errorf("'phone number' contains a symbol that is not supported")
	}

	return nil
}

func ValidateTelegramID(v int64) error {
	if v < 100000 || v > 1000000000 {
		return fmt.Errorf("'telegram id must contains from 6 to 10 digits")
	}

	return nil
}
