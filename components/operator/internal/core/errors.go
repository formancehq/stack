package core

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	CodePending         = "pending"
	CodeStackNotFound   = "stack not found"
	CodeMissingSettings = "missing settings"
)

type ApplicationError struct {
	code    string
	message interface{}
}

func (e *ApplicationError) Error() string {
	if e.message == "" {
		return e.code
	}
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

func (e *ApplicationError) Is(err error) bool {
	_, ok := err.(*ApplicationError)
	return ok
}

func NewApplicationError(code string, msg string) *ApplicationError {
	return &ApplicationError{
		code:    code,
		message: msg,
	}
}

func NewStackNotFoundError() *ApplicationError {
	return NewApplicationError(CodeStackNotFound, "")
}

func NewPendingError() *ApplicationError {
	return NewApplicationError(CodePending, "")
}

func NewMissingSettingsError(msg string) *ApplicationError {
	return NewApplicationError(CodeMissingSettings, msg)
}

func IsApplicationError(err error) bool {
	return errors.Is(err, &ApplicationError{})
}
