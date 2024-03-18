package web

import "github.com/pkg/errors"

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}

type Error struct {
	Value  error
	Status int
	Fields []FieldError
}

func NewRequestError(value error, status int) error {
	return &Error{Value: value, Status: status, Fields: nil}
}

func (e *Error) Error() string {
	return e.Value.Error()
}

type shutdown struct {
	Message string
}

func NewShutdown(message string) error {
	return &shutdown{Message: message}
}

func (s *shutdown) Error() string {
	return s.Message
}

func IsShutdown(err error) bool {
	if _, ok := errors.Cause(err).(*shutdown); ok {
		return true
	}

	return false
}
