package errutil

import "errors"

var ErrRecordNotFound = errors.New("record not found")

type Error struct {
	Message       string
	Type          string
	originalError error
}

func (e *Error) Error() string {
	return e.Message
}

func NewBusinessError(err error, mes string) *Error {
	return &Error{
		Message:       mes,
		originalError: err,
		Type:          "BUSINESS",
	}
}

func NewInputError(err error) *Error {
	return &Error{
		originalError: err,
		Type:          "INPUT",
	}
}
