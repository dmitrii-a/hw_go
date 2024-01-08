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
	ErrIncorrectType    = errors.New("incorrect type for validator")
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
	var fieldValidateFunc func(value reflect.Value) ValidationErrors
	for _, tag := range strings.Split(tags, separator) {
		switch {
		case strings.HasPrefix(tag, lenTagValidator):
			length, err := strconv.Atoi(strings.TrimPrefix(tag, lenTagValidator))
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: name, Err: ErrInvalidLength})
				continue
			}
			fieldValidateFunc = func(value reflect.Value) ValidationErrors {
				return validateLenTag(name, length, value)
			}
		case strings.HasPrefix(tag, regexpTagValidator):
			pattern := strings.TrimPrefix(tag, regexpTagValidator)
			re, err := regexp.Compile(pattern)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: name, Err: ErrInvalidRegexp})
				continue
			}
			fieldValidateFunc = func(value reflect.Value) ValidationErrors {
				return validateRegexpTag(name, re, value)
			}
		case strings.HasPrefix(tag, inTagValidator):
			data := strings.Split(strings.TrimPrefix(tag, inTagValidator), valueSeparator)
			isDigitData := len(data) == 2 && IsDigit(data[0]) && IsDigit(data[1])
			fieldValidateFunc = func(value reflect.Value) ValidationErrors {
				return validateInTag(name, data, isDigitData, value)
			}
		case strings.HasPrefix(tag, minTagValidator):
			length, err := strconv.Atoi(strings.TrimPrefix(tag, minTagValidator))
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: name, Err: ErrInvalidMin})
				continue
			}
			fieldValidateFunc = func(value reflect.Value) ValidationErrors {
				return validateMinTag(name, length, value)
			}
		case strings.HasPrefix(tag, maxTagValidator):
			length, err := strconv.Atoi(strings.TrimPrefix(tag, maxTagValidator))
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: name, Err: ErrInvalidMax})
				continue
			}
			fieldValidateFunc = func(value reflect.Value) ValidationErrors {
				return validateMaxTag(name, length, value)
			}
		}
		if fieldValue.Kind() == reflect.Slice {
			for i := 0; i < fieldValue.Len(); i++ {
				validationErrors = append(validationErrors, fieldValidateFunc(fieldValue.Index(i))...)
			}
		} else {
			validationErrors = append(validationErrors, fieldValidateFunc(fieldValue)...)
		}
	}
	return validationErrors
}

func validateLenTag(name string, length int, fieldValue reflect.Value) ValidationErrors {
	if fieldValue.Kind() != reflect.String {
		return ValidationErrors{{Field: name, Err: ErrIncorrectType}}
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

func validateInTag(name string, data []string, isDigitData bool, fieldValue reflect.Value) ValidationErrors {
	if fieldValue.Kind() != reflect.String && fieldValue.Kind() != reflect.Int {
		return ValidationErrors{{Field: name, Err: ErrIncorrectType}}
	}
	var validationErrors ValidationErrors
	isMatch := false
	if isDigitData {
		minValue, err := strconv.Atoi(data[0])
		if err != nil {
			validationErrors = append(
				validationErrors, ValidationError{
					Field: name,
					Err:   fmt.Errorf("%w%v", ErrInvalidValue, data[0]),
				},
			)
		}
		maxValue, err := strconv.Atoi(data[1])
		if err != nil {
			validationErrors = append(
				validationErrors, ValidationError{
					Field: name, Err: fmt.Errorf("%w%v", ErrInvalidValue, data[1]),
				},
			)
		}
		if fieldValue.Len() >= minValue && fieldValue.Len() <= maxValue {
			isMatch = true
		}
	}
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

func validateRegexpTag(name string, re *regexp.Regexp, fieldValue reflect.Value) ValidationErrors {
	if fieldValue.Kind() != reflect.String {
		return ValidationErrors{{Field: name, Err: ErrIncorrectType}}
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

func validateMinTag(name string, length int, fieldValue reflect.Value) ValidationErrors {
	if fieldValue.Kind() != reflect.Int {
		return ValidationErrors{{Field: name, Err: ErrIncorrectType}}
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

func validateMaxTag(name string, length int, fieldValue reflect.Value) ValidationErrors {
	if fieldValue.Kind() != reflect.Int {
		return ValidationErrors{{Field: name, Err: ErrIncorrectType}}
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
		err := validateField(fieldValue, fieldType.Name, tag)
		validationErrors = append(validationErrors, err...)
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}
