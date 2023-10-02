package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func IsDigit(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func prevIsNotEscape(prevValue string, prevIsString bool) bool {
	return prevValue != `\` || prevIsString
}

func addSymbolInResult(result *strings.Builder, value string, isString bool) {
	if isString {
		result.WriteString(value)
	}
}

func Unpack(s string) (string, error) {
	var (
		prevValue    string
		prevIsString bool
		result       strings.Builder
	)
	data := []rune(s)
	for i := 0; i < len(data); i++ {
		currValue := string(data[i])
		if IsDigit(currValue) && !prevIsString && prevValue != `\` {
			return "", ErrInvalidString
		}
		if prevValue == `\` && !(IsDigit(currValue) || currValue == `\`) {
			return "", ErrInvalidString
		}
		newValue := prevValue
		if IsDigit(currValue) {
			digit, _ := strconv.Atoi(currValue)
			newValue = strings.Repeat(prevValue, digit)
		}
		addSymbolInResult(&result, newValue, prevIsString)
		isString := true
		if (IsDigit(currValue) || currValue == `\`) && prevIsNotEscape(prevValue, prevIsString) {
			isString = false
		}
		prevIsString = isString
		prevValue = currValue
	}
	addSymbolInResult(&result, prevValue, prevIsString)
	return result.String(), nil
}

func UnpackFirst(s string) (string, error) {
	data := []rune(s)
	var prevValue string
	var result strings.Builder
	for i := 0; i < len(data); i++ {
		currValue := string(data[i])
		if IsDigit(currValue) && (IsDigit(prevValue) || prevValue == "") {
			return "", ErrInvalidString
		}
		newValue := prevValue
		if IsDigit(currValue) {
			digit, _ := strconv.Atoi(currValue)
			newValue = strings.Repeat(prevValue, digit)
			currValue = ""
		}
		result.WriteString(newValue)
		prevValue = currValue
	}
	result.WriteString(prevValue)
	return result.String(), nil
}
