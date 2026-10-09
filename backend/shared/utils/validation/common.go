package utils_validation

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func ValidateUsername(v string) error {
	if !LenBetween(v, 3, 32) {
		return fmt.Errorf("'username' must be between 3 and 32 symbols")
	} else if !UsernameRe.MatchString(v) {
		return fmt.Errorf("'username' can contain only letter, number, '.', '-' or '_'")
	}

	return nil
}

func ValidateName(v string) error {
	if !LenBetween(v, 2, 32) {
		return fmt.Errorf("'name' must be between 2 and 32 symbols")
	} else if !IsOnlyCyrillic(v) {
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
	if !LenBetween(v, 2, 32) {
		return fmt.Errorf("'surname' must be between 2 and 32 symbols")
	} else if !IsOnlyCyrillic(v) {
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
	if !LenBetween(v, 12, 15) {
		return fmt.Errorf("'phone number' must be between 12 and 15 symbols")
	} else if !strings.HasPrefix(v, "+") {
		return fmt.Errorf("'phone number' must start with '+'")
	} else if !IsOnlyDigits(v[1:]) {
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

func ValidateID(v uuid.UUID) error {
	if v == uuid.Nil {
		return fmt.Errorf("id must not be nil")
	}

	return nil
}

func ValidatePassword(v string) error {
	if !LenBetween(v, 8, 64) {
		return fmt.Errorf("'password' must be between 8 and 64 symbols")
	} else if !IsOnlyASCII(v) {
		return fmt.Errorf("'password' must contains only ASCII symbols")
	} else if !ContainsUppercaseLetter(v) {
		return fmt.Errorf("'password' must contains at least one uppercase letter")
	} else if !ContainsLowercaseLetter(v) {
		return fmt.Errorf("'password' must contains at least one lowercase letter")
	} else if !ContainsNumber(v) {
		return fmt.Errorf("'password' must contains at least one number")
	} else if !ContainsSpecialSymbol(v) {
		return fmt.Errorf("'password' must contains at least one symbol")
	}

	return nil
}
