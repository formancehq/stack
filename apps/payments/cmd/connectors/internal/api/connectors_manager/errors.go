package connectors_manager

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrAlreadyInstalled  = errors.New("already installed")
	ErrNotInstalled      = errors.New("not installed")
	ErrNotEnabled        = errors.New("not enabled")
	ErrAlreadyRunning    = errors.New("already running")
	ErrConnectorNotFound = errors.New("connector not found")
	ErrValidation        = errors.New("validation error")
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
