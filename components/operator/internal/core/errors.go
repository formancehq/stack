package core

import (
	"fmt"

	"github.com/pkg/errors"
)

type ApplicationError struct {
	message string
}

func (e *ApplicationError) Error() string {
	return e.message
}

func (e *ApplicationError) Is(err error) bool {
	_, ok := err.(*ApplicationError)
	return ok
}

func NewApplicationError(msg string, args ...any) *ApplicationError {
	return &ApplicationError{
		message: fmt.Sprintf(msg, args...),
	}
}

func NewStackNotFoundError() *ApplicationError {
	return NewApplicationError("stack not found")
}

func NewPendingError() *ApplicationError {
	return NewApplicationError("pending")
}

func NewMissingSettingsError(msg string) *ApplicationError {
	return NewApplicationError(msg)
}

func IsApplicationError(err error) bool {
	return errors.Is(err, &ApplicationError{})
}
