package sqlerrors

import (
	"fmt"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

var (
	ErrConfigurationNotFound = errors.New("configuration not found")
)

type Code string

const (
	ConstraintFailed Code = "CONSTRAINT_FAILED"
	TooManyClient    Code = "TOO_MANY_CLIENT"
)

type Error struct {
	Code          Code
	OriginalError error
}

func (e Error) Is(err error) bool {
	storageErr, ok := err.(*Error)
	if !ok {
		return false
	}
	if storageErr.Code == "" {
		return true
	}
	return storageErr.Code == e.Code
}

func (e Error) Error() string {
	return fmt.Sprintf("%s [%s]", e.OriginalError, e.Code)
}

func NewError(code Code, originalError error) *Error {
	return &Error{
		Code:          code,
		OriginalError: originalError,
	}
}

func IsError(err error) bool {
	return IsErrorCode(err, "")
}

func IsErrorCode(err error, code Code) bool {
	return errors.Is(err, &Error{
		Code: code,
	})
}

// PostgresError is a helper to wrap postgres errors into storage errors
func PostgresError(err error) error {
	if err != nil {
		switch pge := err.(type) {
		case *pq.Error:
			switch pge.Code {
			case "23505":
				return NewError(ConstraintFailed, err)
			case "53300":
				return NewError(TooManyClient, err)
			}
		}
	}

	return err
}
