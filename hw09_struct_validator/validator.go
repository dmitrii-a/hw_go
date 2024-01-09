package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type validateFieldFunc func(value reflect.Value) *ValidationError

func getFieldValidateFunc(
	name string,
	validator string,
	validatorValue string,
	fieldValue reflect.Value,
) (validateFieldFunc, *ValidatorError) {
	if vrError := validateTypeTag(name, validator, fieldValue); vrError != nil {
		return nil, vrError
	}
	switch {
	case validator == lenTagValidator:
		length, err := strconv.Atoi(validatorValue)
		if err != nil {
			return nil, &ValidatorError{Field: name, Err: ErrInvalidLength}
		}
		return func(value reflect.Value) *ValidationError {
			return validateLenTag(name, length, value.String())
		}, nil
	case validator == regexpTagValidator:
		re, err := regexp.Compile(validatorValue)
		if err != nil {
			return nil, &ValidatorError{Field: name, Err: ErrInvalidRegexp}
		}
		return func(value reflect.Value) *ValidationError {
			return validateRegexpTag(name, re, value.String())
		}, nil
	case validator == inTagValidator:
		data := strings.Split(validatorValue, valueSeparator)
		if len(data) == 2 && IsDigit(data[0]) && IsDigit(data[1]) {
			minValue, err := strconv.Atoi(data[0])
			if err != nil {
				return nil, &ValidatorError{Field: name, Err: fmt.Errorf("%w%v", ErrInvalidValue, data[0])}
			}
			maxValue, err := strconv.Atoi(data[1])
			if err != nil {
				return nil, &ValidatorError{Field: name, Err: fmt.Errorf("%w%v", ErrInvalidValue, data[1])}
			}
			return func(value reflect.Value) *ValidationError {
				return validateIntInTag(name, int64(minValue), int64(maxValue), value.Int())
			}, nil
		}
		return func(value reflect.Value) *ValidationError {
			return validateStringInTag(name, data, fmt.Sprint(fieldValue))
		}, nil
	case validator == minTagValidator:
		length, err := strconv.Atoi(validatorValue)
		if err != nil {
			return nil, &ValidatorError{Field: name, Err: ErrInvalidMin}
		}
		return func(value reflect.Value) *ValidationError {
			return validateMinTag(name, int64(length), value.Int())
		}, nil
	case validator == maxTagValidator:
		length, err := strconv.Atoi(validatorValue)
		if err != nil {
			return nil, &ValidatorError{Field: name, Err: ErrInvalidMax}
		}
		return func(value reflect.Value) *ValidationError {
			return validateMaxTag(name, int64(length), value.Int())
		}, nil
	}
	return nil, nil
}

func validateField(fieldValue reflect.Value, name string, tags string) (ValidationErrors, ValidatorErrors) {
	var validationErrors ValidationErrors
	var validatorErrors ValidatorErrors
	for _, tag := range strings.Split(tags, separator) {
		splitTag := strings.Split(tag, tagSeparator)
		if len(splitTag) < 2 {
			validatorErrors = append(validatorErrors, ValidatorError{Field: name, Err: ErrIncorrectValidator})
			continue
		}
		fieldValidateFunc, validatorError := getFieldValidateFunc(name, splitTag[0], splitTag[1], fieldValue)
		if validatorError != nil {
			validatorErrors = append(validatorErrors, *validatorError)
			continue
		}
		if fieldValue.Kind() == reflect.Slice {
			for i := 0; i < fieldValue.Len(); i++ {
				if err := fieldValidateFunc(fieldValue.Index(i)); err != nil {
					validationErrors = append(validationErrors, *err)
				}
			}
		} else {
			if err := fieldValidateFunc(fieldValue); err != nil {
				validationErrors = append(validationErrors, *err)
			}
		}
	}
	return validationErrors, validatorErrors
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
		vnErrors, vrErrors := validateField(fieldValue, fieldType.Name, tag)
		if vrErrors != nil {
			return vrErrors
		}
		validationErrors = append(validationErrors, vnErrors...)
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}
