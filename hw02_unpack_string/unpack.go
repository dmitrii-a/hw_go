package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func prevIsNotEscape(prevValue rune, prevIsString bool) bool {
	return prevValue != '\\' || prevIsString
}

func addSymbolInResult(result *strings.Builder, value []rune, isString bool) {
	if isString {
		for _, r := range value {
			result.WriteRune(r)
		}
	}
}

func Unpack(s string) (string, error) {
	var (
		prevValue    rune
		prevIsString bool
		result       strings.Builder
	)
	backslash := '\\'
	for _, value := range s {
		if unicode.IsDigit(value) && !prevIsString && prevValue != backslash {
			return "", ErrInvalidString
		}
		if prevValue == backslash && !(unicode.IsDigit(value) || value == backslash) {
			return "", ErrInvalidString
		}
		var newValue []rune
		if unicode.IsDigit(value) {
			for i := 0; i < int(value-'0'); i++ {
				newValue = append(newValue, prevValue)
			}
		} else {
			newValue = append(newValue, prevValue)
		}
		addSymbolInResult(&result, newValue, prevIsString)
		isString := true
		if (unicode.IsDigit(value) || value == backslash) && prevIsNotEscape(prevValue, prevIsString) {
			isString = false
		}
		prevIsString = isString
		prevValue = value
	}
	if prevValue == backslash {
		return "", ErrInvalidString
	}
	addSymbolInResult(&result, []rune{prevValue}, prevIsString)
	return result.String(), nil
}

//func IsDigit(s string) bool {
//	for _, c := range s {
//		if c < '0' || c > '9' {
//			return false
//		}
//	}
//	return true
//}
//func UnpackFirst(s string) (string, error) {
//	var prevValue string
//	var result strings.Builder
//	for _, v := range s {
//		currValue := string(v)
//		if IsDigit(currValue) && (IsDigit(prevValue) || prevValue == "") {
//			return "", ErrInvalidString
//		}
//		newValue := prevValue
//		if IsDigit(currValue) {
//			digit, _ := strconv.Atoi(currValue)
//			newValue = strings.Repeat(prevValue, digit)
//			currValue = ""
//		}
//		result.WriteString(newValue)
//		prevValue = currValue
//	}
//	result.WriteString(prevValue)
//	return result.String(), nil
//}
