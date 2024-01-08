package hw09structvalidator

import (
	"errors"
	"fmt"
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
