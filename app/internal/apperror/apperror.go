package apperror

import (
	"encoding/json"
	"fmt"
)

var (
	ErrNotFound = NewAppError("SS-000404", "not found", "not found")
)

type ErrorFields map[string]string
type ErrorParams map[string]string

type AppError struct {
	Err              error       `json:"-"`
	Code             string      `json:"code,omitempty"`
	Message          string      `json:"message,omitempty"`
	DeveloperMessage string      `json:"developer_message,omitempty"`
	Fields           ErrorFields `json:"fields,omitempty"`
	Params           ErrorParams `json:"params,omitempty"`
}

func NewAppError(code, message, developerMessage string) *AppError {
	return &AppError{
		Err:              fmt.Errorf(message),
		Code:             code,
		Message:          message,
		DeveloperMessage: developerMessage,
	}
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Marshal() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return bytes
}

func (e *AppError) WithFields(fields ErrorFields) {
	e.Fields = fields
}

func (e *AppError) WithParams(params ErrorParams) {
	e.Params = params
}

func BadRequestError(message string) *AppError {
	return NewAppError("SS-000400", message, "something wrong with user data")
}

func SystemError(developerMessage string) *AppError {
	return NewAppError("SS-000418", "internal system error", developerMessage)
}
