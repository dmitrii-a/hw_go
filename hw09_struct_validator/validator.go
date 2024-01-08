package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type validateFieldFunc func(value reflect.Value) error

func getFieldValidateFunc(
	name string,
	validator,
	validatorValue string,
) (validateFieldFunc, *ValidatorError) {
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
	fieldValidateFunc validateFieldFunc,
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
		splitTag := strings.Split(tag, tagSeparator)
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
