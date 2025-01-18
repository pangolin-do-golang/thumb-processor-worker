package errutil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBusinessError(t *testing.T) {
	originalErr := errors.New("original error")
	businessErr := NewBusinessError(originalErr, "business error occurred")

	assert.Equal(t, "business error occurred", businessErr.Message)
	assert.Equal(t, "BUSINESS", businessErr.Type)
	assert.Equal(t, originalErr, businessErr.originalError)
}

func TestNewInputError(t *testing.T) {
	originalErr := errors.New("original error")
	inputErr := NewInputError(originalErr)

	assert.Equal(t, "INPUT", inputErr.Type)
	assert.Equal(t, originalErr, inputErr.originalError)
}

func TestError_Error(t *testing.T) {
	err := &Error{Message: "some error message"}
	assert.Equal(t, "some error message", err.Error())
}

func TestErrRecordNotFound(t *testing.T) {
	assert.EqualError(t, ErrRecordNotFound, "record not found")
}
