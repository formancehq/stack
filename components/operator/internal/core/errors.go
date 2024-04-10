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

func (e *ApplicationError) WithMessage(msg string, args ...any) *ApplicationError {
	e.message = fmt.Sprintf(msg, args...)
	return e
}

func NewApplicationError() *ApplicationError {
	return &ApplicationError{}
}

func NewStackNotFoundError() *ApplicationError {
	return NewApplicationError().WithMessage("stack not found")
}

func NewPendingError() *ApplicationError {
	return NewApplicationError().WithMessage("pending")
}

func NewMissingSettingsError(msg string) *ApplicationError {
	return NewApplicationError().WithMessage(msg)
}

func IsApplicationError(err error) bool {
	return errors.Is(err, &ApplicationError{})
}
