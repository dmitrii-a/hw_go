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
	lenTagValidator    = "len"
	inTagValidator     = "in"
	regexpTagValidator = "regexp"
	minTagValidator    = "min"
	maxTagValidator    = "max"
)

var (
	ErrTypeNotStruct      = errors.New("type is not a struct")
	ErrInvalidLength      = errors.New("invalid length")
	ErrNotEqualLength     = errors.New("length not equal ")
	ErrInvalidRegexp      = errors.New("invalid regexp")
	ErrRegexpNotMatch     = errors.New("regexp not match ")
	ErrInvalidValue       = errors.New("invalid value ")
	ErrValueNotIn         = errors.New("value not in ")
	ErrInvalidMin         = errors.New("invalid value for min validator")
	ErrInvalidMax         = errors.New("invalid value for max validator")
	ErrValueLessThanMin   = errors.New("value less than ")
	ErrValueMoreThanMax   = errors.New("value more than ")
	ErrIncorrectType      = errors.New("incorrect type for validator")
	ErrIncorrectValidator = errors.New("incorrect validator")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidatorError struct {
	Field string
	Err   error
}

func (v ValidatorError) Error() string {
	return fmt.Sprintf("%s: %s\n", v.Field, v.Err.Error())
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var validationError string
	for _, err := range v {
		validationError += fmt.Sprintf("%s: %s\n", err.Field, err.Err.Error())
	}
	return validationError
}

type ValidatorErrors []ValidatorError

func (v ValidatorErrors) Error() string {
	var validatorError string
	for _, err := range v {
		validatorError += fmt.Sprintf("%s: %s\n", err.Field, err.Err.Error())
	}
	return validatorError
}

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

func getFieldValidateFunc(
	name string,
	validator,
	validatorValue string,
) (func(value reflect.Value) error, *ValidatorError) {
	switch {
	case validator == lenTagValidator:
		length, err := strconv.Atoi(validatorValue)
		if err != nil {
			return nil, &ValidatorError{Field: name, Err: ErrInvalidLength}
		}
		return func(value reflect.Value) error {
			return validateLenTag(name, length, value)
		}, nil
	case validator == regexpTagValidator:
		re, err := regexp.Compile(validatorValue)
		if err != nil {
			return nil, &ValidatorError{Field: name, Err: ErrInvalidRegexp}
		}
		return func(value reflect.Value) error {
			return validateRegexpTag(name, re, value)
		}, nil
	case validator == inTagValidator:
		data := strings.Split(validatorValue, valueSeparator)
		isDigitData := len(data) == 2 && IsDigit(data[0]) && IsDigit(data[1])
		if isDigitData {
			minValue, err := strconv.Atoi(data[0])
			if err != nil {
				return nil, &ValidatorError{Field: name, Err: fmt.Errorf("%w%v", ErrInvalidValue, data[0])}
			}
			maxValue, err := strconv.Atoi(data[1])
			if err != nil {
				return nil, &ValidatorError{Field: name, Err: fmt.Errorf("%w%v", ErrInvalidValue, data[1])}
			}
			return func(value reflect.Value) error {
				return validateIntInTag(name, minValue, maxValue, value)
			}, nil
		}
		return func(value reflect.Value) error {
			return validateStringInTag(name, data, value)
		}, nil
	case validator == minTagValidator:
		length, err := strconv.Atoi(validatorValue)
		if err != nil {
			return nil, &ValidatorError{Field: name, Err: ErrInvalidMin}
		}
		return func(value reflect.Value) error {
			return validateMinTag(name, length, value)
		}, nil
	case validator == maxTagValidator:
		length, err := strconv.Atoi(validatorValue)
		if err != nil {
			return nil, &ValidatorError{Field: name, Err: ErrInvalidMax}
		}
		return func(value reflect.Value) error {
			return validateMaxTag(name, length, value)
		}, nil
	}
	return nil, nil
}

func processValidateField(
	fieldValue reflect.Value,
	fieldValidateFunc func(value reflect.Value) error,
	validationErrors ValidationErrors,
	validatorErrors ValidatorErrors,
) (ValidationErrors, ValidatorErrors) {
	err := fieldValidateFunc(fieldValue)
	var vrError ValidatorError
	var vnErrors ValidationErrors
	switch {
	case errors.As(err, &vrError):
		validatorErrors = append(validatorErrors, vrError)
	case errors.As(err, &vnErrors):
		validationErrors = append(validationErrors, vnErrors...)
	}
	return validationErrors, validatorErrors
}

func validateField(fieldValue reflect.Value, name string, tags string) ValidationErrors {
	var validationErrors ValidationErrors
	var validatorErrors ValidatorErrors
	for _, tag := range strings.Split(tags, separator) {
		splitTag := strings.Split(tag, ":")
		if len(splitTag) < 2 {
			validatorErrors = append(validatorErrors, ValidatorError{Field: name, Err: ErrIncorrectValidator})
		}
		validator := splitTag[0]
		validatorValue := splitTag[1]
		fieldValidateFunc, validatorError := getFieldValidateFunc(name, validator, validatorValue)
		if validatorError != nil {
			validatorErrors = append(validatorErrors, *validatorError)
			continue
		}
		if fieldValue.Kind() == reflect.Slice {
			for i := 0; i < fieldValue.Len(); i++ {
				validationErrors, validatorErrors = processValidateField(
					fieldValue.Index(i), fieldValidateFunc, validationErrors, validatorErrors,
				)
			}
		} else {
			validationErrors, validatorErrors = processValidateField(
				fieldValue, fieldValidateFunc, validationErrors, validatorErrors,
			)
		}
	}
	if len(validatorErrors) > 0 {
		panic(validatorErrors)
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
