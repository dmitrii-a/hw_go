package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func prevIsNotEscape(prevValue rune, prevIsString bool) bool {
	return prevValue != '\\' || prevIsString
}

func IsDigit(s rune) bool {
	if s < '0' || s > '9' {
		return false
	}
	return true
}

func Unpack(s string) (string, error) {
	var (
		prevValue    rune
		prevIsString bool
		result       strings.Builder
	)
	backslash := '\\'
	for _, value := range s {
		if IsDigit(value) && !prevIsString && prevValue != backslash {
			return "", ErrInvalidString
		}
		if prevValue == backslash && !(IsDigit(value) || value == backslash) {
			return "", ErrInvalidString
		}
		newValue := []rune{prevValue}
		if IsDigit(value) {
			size := int(value - '0')
			newValue = make([]rune, size)
			for i := 0; i < size; i++ {
				newValue[i] = prevValue
			}
		}
		if prevIsString {
			for _, r := range newValue {
				result.WriteRune(r)
			}
		}
		isString := true
		if (IsDigit(value) || value == backslash) && prevIsNotEscape(prevValue, prevIsString) {
			isString = false
		}
		prevIsString = isString
		prevValue = value
	}
	if prevValue == backslash {
		return "", ErrInvalidString
	}
	if prevIsString {
		result.WriteRune(prevValue)
	}
	return result.String(), nil
}
