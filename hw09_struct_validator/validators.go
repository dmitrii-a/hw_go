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

func validateTypeTag(name string, validator string, fieldValue reflect.Value) *ValidatorError {
	fieldType := fieldValue.Kind()
	if fieldType == reflect.Slice {
		fieldType = fieldValue.Type().Elem().Kind()
	}
	err := &ValidatorError{Field: name, Err: ErrIncorrectType}
	switch {
	case validator == lenTagValidator:
		if fieldType != reflect.String {
			return err
		}
	case validator == regexpTagValidator:
		if fieldType != reflect.String {
			return err
		}
	case validator == inTagValidator:
		if fieldType != reflect.String && fieldType != reflect.Int {
			return err
		}
	case validator == minTagValidator:
		if fieldType != reflect.Int {
			return err
		}
	case validator == maxTagValidator:
		if fieldType != reflect.Int {
			return err
		}
	}
	return nil
}

func validateLenTag(name string, length int, value string) *ValidationError {
	if len([]rune(value)) != length {
		return &ValidationError{Field: name, Err: fmt.Errorf("%w%v", ErrNotEqualLength, length)}
	}
	return nil
}

func validateIntInTag(name string, minValue int64, maxValue int64, value int64) *ValidationError {
	if value >= minValue && value <= maxValue {
		return &ValidationError{Field: name, Err: fmt.Errorf("%w%v %v", ErrValueNotIn, minValue, maxValue)}
	}
	return nil
}

func validateStringInTag(name string, data []string, value string) *ValidationError {
	isMatch := false
	for _, v := range data {
		if value == v {
			isMatch = true
		}
	}
	if !isMatch {
		return &ValidationError{Field: name, Err: fmt.Errorf("%w%v", ErrValueNotIn, data)}
	}
	return nil
}

func validateRegexpTag(name string, re *regexp.Regexp, value string) *ValidationError {
	match := re.MatchString(value)
	if !match {
		return &ValidationError{Field: name, Err: fmt.Errorf("%w%v", ErrRegexpNotMatch, re.String())}
	}
	return nil
}

func validateMinTag(name string, length int64, value int64) *ValidationError {
	if value < length {
		return &ValidationError{Field: name, Err: fmt.Errorf("%w%v", ErrValueLessThanMin, length)}
	}
	return nil
}

func validateMaxTag(name string, length int64, value int64) *ValidationError {
	if value > length {
		return &ValidationError{Field: name, Err: fmt.Errorf("%w%v", ErrValueMoreThanMax, length)}
	}
	return nil
}
