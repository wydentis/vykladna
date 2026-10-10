package utils_validation

import (
	"unicode"
	"unicode/utf8"
)

func LenBetween(v string, from, to int) bool {
	length := len([]rune(v))
	return length >= from && length <= to
}

func IsOnlyASCII(v string) bool {
	for _, r := range v {
		if r > unicode.MaxASCII {
			return false
		}
	}

	return true
}

func IsOnlyCyrillic(v string) bool {
	for _, r := range v {
		if !unicode.IsLetter(r) || !unicode.Is(unicode.Cyrillic, r) {
			return false
		}
	}

	return true
}

func IsOnlyDigits(v string) bool {
	for _, r := range v {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}

func ContainsUppercaseLetter(v string) bool {
	for _, r := range v {
		if unicode.IsUpper(r) {
			return true
		}
	}

	return false
}

func ContainsLowercaseLetter(v string) bool {
	for _, r := range v {
		if unicode.IsLower(r) {
			return true
		}
	}

	return false
}

func ContainsNumber(v string) bool {
	for _, r := range v {
		if unicode.IsNumber(r) {
			return true
		}
	}

	return false
}

func ContainsSpecialSymbol(v string) bool {
	for _, r := range v {
		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			return true
		}
	}

	return false
}

func IsFirstLetterUppercase(v string) bool {
	r, _ := utf8.DecodeRuneInString(v)
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
