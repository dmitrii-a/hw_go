package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	validateTag        = "validate"
	separator          = "|"
	valueSeparator     = ","
	lenTagValidator    = "len:"
	inTagValidator     = "in:"
	regexpTagValidator = "regexp:"
	minTagValidator    = "min:"
	maxTagValidator    = "max:"
)

var (
	ErrTypeNotStruct    = errors.New("type is not a struct")
	ErrInvalidLength    = errors.New("invalid length")
	ErrNotEqualLength   = errors.New("length not equal ")
	ErrInvalidRegexp    = errors.New("invalid regexp")
	ErrRegexpNotMatch   = errors.New("regexp not match ")
	ErrInvalidValue     = errors.New("invalid value ")
	ErrValueNotIn       = errors.New("value not in ")
	ErrInvalidMin       = errors.New("invalid value for min validator")
	ErrInvalidMax       = errors.New("invalid value for max validator")
	ErrValueLessThanMin = errors.New("value less than ")
	ErrValueMoreThanMax = errors.New("value more than ")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var validationError string
	for _, err := range v {
		validationError += fmt.Sprintf("%s: %s\n", err.Field, err.Err.Error())
	}
	return validationError
}

func IsDigit(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func validateField(fieldValue reflect.Value, name string, tags string) ValidationErrors {
	var validationErrors ValidationErrors
	for _, tag := range strings.Split(tags, separator) {
		switch {
		case strings.HasPrefix(tag, lenTagValidator):
			validationErrors = append(validationErrors, validateLenTag(name, tag, fieldValue)...)
		case strings.HasPrefix(tag, regexpTagValidator):
			validationErrors = append(validationErrors, validateRegexpTag(name, tag, fieldValue)...)
		case strings.HasPrefix(tag, inTagValidator):
			validationErrors = append(validationErrors, validateInTag(name, tag, fieldValue)...)
		case strings.HasPrefix(tag, minTagValidator):
			validationErrors = append(validationErrors, validateMinTag(name, tag, fieldValue)...)
		case strings.HasPrefix(tag, maxTagValidator):
			validationErrors = append(validationErrors, validateMaxTag(name, tag, fieldValue)...)
		}
	}
	return validationErrors
}

func validateLenTag(name, tag string, fieldValue reflect.Value) ValidationErrors {
	var validationErrors ValidationErrors
	length, err := strconv.Atoi(strings.TrimPrefix(tag, lenTagValidator))
	if err != nil {
		validationErrors = append(validationErrors, ValidationError{Field: name, Err: ErrInvalidLength})
	}
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

func validateInTag(name, tag string, fieldValue reflect.Value) ValidationErrors {
	var validationErrors ValidationErrors
	s := strings.Split(strings.TrimPrefix(tag, inTagValidator), valueSeparator)
	isMatch := false
	if len(s) == 2 && IsDigit(s[0]) && IsDigit(s[1]) {
		minValue, err := strconv.Atoi(s[0])
		if err != nil {
			validationErrors = append(
				validationErrors, ValidationError{
					Field: name,
					Err:   fmt.Errorf("%w%v", ErrInvalidValue, s[0]),
				},
			)
		}
		maxValue, err := strconv.Atoi(s[1])
		if err != nil {
			validationErrors = append(
				validationErrors, ValidationError{
					Field: name, Err: fmt.Errorf("%w%v", ErrInvalidValue, s[1]),
				},
			)
		}
		if fieldValue.Len() >= minValue && fieldValue.Len() <= maxValue {
			isMatch = true
		}
	}
	for _, v := range s {
		if fmt.Sprint(fieldValue) == v {
			isMatch = true
		}
	}
	if !isMatch {
		validationErrors = append(
			validationErrors,
			ValidationError{
				Field: name, Err: fmt.Errorf("%w%v", ErrValueNotIn, s),
			},
		)
	}
	return validationErrors
}

func validateRegexpTag(name, tag string, fieldValue reflect.Value) ValidationErrors {
	var validationErrors ValidationErrors
	pattern := strings.TrimPrefix(tag, regexpTagValidator)
	re, err := regexp.MatchString(pattern, fieldValue.String())
	if err != nil {
		validationErrors = append(validationErrors, ValidationError{Field: name, Err: ErrInvalidRegexp})
	}
	if !re {
		validationErrors = append(
			validationErrors,
			ValidationError{
				Field: name, Err: fmt.Errorf("%w%v", ErrRegexpNotMatch, pattern),
			},
		)
	}
	return validationErrors
}

func validateMinTag(name, tag string, fieldValue reflect.Value) ValidationErrors {
	var validationErrors ValidationErrors
	length, err := strconv.Atoi(strings.TrimPrefix(tag, minTagValidator))
	if err != nil {
		validationErrors = append(validationErrors, ValidationError{Field: name, Err: ErrInvalidMin})
	}
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

func validateMaxTag(name, tag string, fieldValue reflect.Value) ValidationErrors {
	var validationErrors ValidationErrors
	length, err := strconv.Atoi(strings.TrimPrefix(tag, maxTagValidator))
	if err != nil {
		validationErrors = append(validationErrors, ValidationError{Field: name, Err: ErrInvalidMax})
	}
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

func Validate(v interface{}) error {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Struct {
		return ErrTypeNotStruct
	}
	validationErrors := ValidationErrors{}
	for i := 0; i < reflectValue.NumField(); i++ {
		fieldValue := reflectValue.Field(i)
		fieldType := reflectValue.Type().Field(i)
		tag := fieldType.Tag.Get(validateTag)
		if tag == "" {
			continue
		}
		if fieldType.Type.Kind() == reflect.Slice {
			for i := 0; i < fieldValue.Len(); i++ {
				err := validateField(fieldValue.Index(i), fieldType.Name, tag)
				validationErrors = append(validationErrors, err...)
			}
		} else {
			err := validateField(fieldValue, fieldType.Name, tag)
			validationErrors = append(validationErrors, err...)
		}
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}
