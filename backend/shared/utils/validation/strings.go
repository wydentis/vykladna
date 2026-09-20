package utils_validation

import (
	"unicode"
	"unicode/utf8"
)

func Between(v string, from, to int) bool {
	length := len([]rune(v))
	return length >= from && length <= to
}

func OnlyCyrillic(v string) bool {
	for _, r := range v {
		if !unicode.IsLetter(r) || !unicode.Is(unicode.Cyrillic, r) {
			return false
		}
	}

	return true
}

func OnlyDigits(v string) bool {
	for _, r := range v {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}

func IsFirstLetterUppercase(v string) bool {
	r, _ := utf8.DecodeRuneInString(v)
	if r == utf8.RuneError {
		return false
	}
	return unicode.IsUpper(r)
}

func IsRestLettersLowercase(v string) bool {
	_, size := utf8.DecodeRuneInString(v)
	for _, r := range v[size:] {
		if !unicode.IsLower(r) {
			return false
		}
	}

	return true
}
