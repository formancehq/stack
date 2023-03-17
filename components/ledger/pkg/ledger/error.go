package ledger

import "github.com/pkg/errors"

type ValidationError struct {
	Msg string
}

func (v ValidationError) Error() string {
	return v.Msg
}

func (v ValidationError) Is(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

func NewValidationError(msg string) *ValidationError {
	return &ValidationError{
		Msg: msg,
	}
}

func IsValidationError(err error) bool {
	return errors.Is(err, &ValidationError{})
}

type ConflictError struct{}

func (e ConflictError) Error() string {
	return "conflict error on reference"
}

func (e ConflictError) Is(err error) bool {
	_, ok := err.(*ConflictError)
	return ok
}

func NewConflictError() *ConflictError {
	return &ConflictError{}
}

func IsConflictError(err error) bool {
	return errors.Is(err, &ConflictError{})
}

type NotFoundError struct {
	Msg string
}

func (v NotFoundError) Error() string {
	return v.Msg
}

func (v NotFoundError) Is(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{
		Msg: msg,
	}
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, &NotFoundError{})
}
