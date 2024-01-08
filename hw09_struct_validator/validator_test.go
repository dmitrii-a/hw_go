package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name           string
		in             interface{}
		expectedErr    error
		inCorrectError bool
	}{
		{
			name:        "not struct",
			in:          "string, not struct",
			expectedErr: ErrTypeNotStruct,
		},
		{
			name: "correct User",
			in: User{
				ID:   "123456789012345678901234567890123456",
				Name: "John", Age: 20,
				Email:  "test@test.ru",
				Role:   "admin",
				Phones: []string{"89999999999"},
				meta:   json.RawMessage(`{"test": "test"}`),
			},
			expectedErr: nil,
		},
		{
			name: "incorrect User",
			in: User{
				ID:     "123",
				Name:   "John",
				Age:    20,
				Email:  "test@test.ru",
				Role:   "admin",
				Phones: []string{"89999999999"},
				meta:   json.RawMessage(`{"test": "test"}`),
			},
			expectedErr: ValidationErrors{{Field: "ID", Err: ErrNotEqualLength}},
		},
		{
			name: "incorrect User role",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John",
				Age:    20,
				Email:  "test@test.ru",
				Role:   "test",
				Phones: []string{"89999999999"},
				meta:   json.RawMessage(`{"test": "test"}`),
			},
			expectedErr: ValidationErrors{{Field: "Role", Err: ErrValueNotIn}},
		},
		{
			name: "incorrect all fields in User",
			in: User{
				ID:     "123",
				Name:   "John",
				Age:    2,
				Email:  "test",
				Role:   "test",
				Phones: []string{"89999999999999999999"},
				meta:   json.RawMessage(`{"test": "test"}`),
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrNotEqualLength},
				{Field: "Age", Err: ErrValueLessThanMin},
				{Field: "Email", Err: ErrRegexpNotMatch},
				{Field: "Role", Err: ErrValueNotIn},
				{Field: "Phones", Err: ErrNotEqualLength},
			},
		},
		{
			name:        "correct length of App",
			in:          App{Version: "1.0.0"},
			expectedErr: nil,
		},
		{
			name:        "incorrect length of App",
			in:          App{Version: "1.0.0.0"},
			expectedErr: ValidationErrors{{Field: "Version", Err: ErrNotEqualLength}},
		},
		{
			name: "Token without validation",
			in: Token{
				Header:    nil,
				Payload:   nil,
				Signature: nil,
			},
			expectedErr: nil,
		},
		{
			name: "correct Response",
			in: Response{
				Code: 200,
				Body: "",
			},
			expectedErr: nil,
		},
		{
			name: "incorrect Response",
			in: Response{
				Code: 700,
				Body: "",
			},
			expectedErr: ValidationErrors{{Field: "Code", Err: ErrValueNotIn}},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			errFromValidate := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, errFromValidate, errFromValidate)
			}
			if tt.expectedErr != nil && errors.As(errFromValidate, &ValidationErrors{}) {
				var validationErrors ValidationErrors
				if errors.As(errFromValidate, &validationErrors) {
					for i, err := range validationErrors {
						var validationErrorsFromValidate ValidationErrors
						if errors.As(errFromValidate, &validationErrorsFromValidate) {
							require.Contains(t, validationErrorsFromValidate[i].Field, err.Field)
							require.ErrorIs(t, validationErrorsFromValidate[i].Err, err.Err)
						}
					}
				}
			}
		})
	}
}
