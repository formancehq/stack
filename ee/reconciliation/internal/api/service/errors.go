package service

import (
	"errors"
	"fmt"
)

var (
	ErrValidation = errors.New("validation error")
	ErrInvalidID  = errors.New("invalid id")
)

type storageError struct {
	err error
	msg string
}

func (e *storageError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err)
}

func (e *storageError) Is(err error) bool {
	_, ok := err.(*storageError)
	return ok
}

func (e *storageError) Unwrap() error {
	return e.err
}

func newStorageError(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &storageError{
		err: err,
		msg: msg,
	}
}
