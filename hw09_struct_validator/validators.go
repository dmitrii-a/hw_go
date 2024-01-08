package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
)

func IsDigit(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func validateLenTag(name string, length int, fieldValue reflect.Value) error {
	if fieldValue.Kind() != reflect.String {
		return ValidatorError{Field: name, Err: ErrIncorrectType}
	}
	var validationErrors ValidationErrors
	if fieldValue.Len() != length {
		validationErrors = append(
			validationErrors,
			ValidationError{
				Field: name, Err: fmt.Errorf("%w%v", ErrNotEqualLength, length),
			},
		)
	}
	return validationErrors
}

func validateIntInTag(name string, minValue int, maxValue int, fieldValue reflect.Value) error {
	if fieldValue.Kind() != reflect.String && fieldValue.Kind() != reflect.Int {
		return ValidatorError{Field: name, Err: ErrIncorrectType}
	}
	var validationErrors ValidationErrors
	if fieldValue.Len() >= minValue && fieldValue.Len() <= maxValue {
		validationErrors = append(
			validationErrors,
			ValidationError{
				Field: name, Err: fmt.Errorf("%w%v %v", ErrValueNotIn, minValue, maxValue),
			},
		)
	}
	return validationErrors
}

func validateStringInTag(name string, data []string, fieldValue reflect.Value) error {
	if fieldValue.Kind() != reflect.String && fieldValue.Kind() != reflect.Int {
		return ValidatorError{Field: name, Err: ErrIncorrectType}
	}
	var validationErrors ValidationErrors
	isMatch := false
	for _, v := range data {
		if fmt.Sprint(fieldValue) == v {
			isMatch = true
		}
	}
	if !isMatch {
		validationErrors = append(
			validationErrors,
			ValidationError{
				Field: name, Err: fmt.Errorf("%w%v", ErrValueNotIn, data),
			},
		)
	}
	return validationErrors
}

func validateRegexpTag(name string, re *regexp.Regexp, fieldValue reflect.Value) error {
	if fieldValue.Kind() != reflect.String {
		return ValidatorError{Field: name, Err: ErrIncorrectType}
	}
	var validationErrors ValidationErrors
	match := re.MatchString(fieldValue.String())
	if !match {
		validationErrors = append(
			validationErrors,
			ValidationError{
				Field: name, Err: fmt.Errorf("%w%v", ErrRegexpNotMatch, re.String()),
			},
		)
	}
	return validationErrors
}

func validateMinTag(name string, length int, fieldValue reflect.Value) error {
	if fieldValue.Kind() != reflect.Int {
		return ValidatorError{Field: name, Err: ErrIncorrectType}
	}
	var validationErrors ValidationErrors
	if val, ok := fieldValue.Interface().(int); ok && val < length {
		validationErrors = append(
			validationErrors,
			ValidationError{
				Field: name, Err: fmt.Errorf("%w%v", ErrValueLessThanMin, length),
			},
		)
	}
	return validationErrors
}

func validateMaxTag(name string, length int, fieldValue reflect.Value) error {
	if fieldValue.Kind() != reflect.Int {
		return ValidatorError{Field: name, Err: ErrIncorrectType}
	}
	var validationErrors ValidationErrors
	if val, ok := fieldValue.Interface().(int); ok && val > length {
		validationErrors = append(
			validationErrors,
			ValidationError{
				Field: name, Err: fmt.Errorf("%w%v", ErrValueMoreThanMax, length),
			},
		)
	}
	return validationErrors
}
