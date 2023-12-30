package core

import "github.com/pkg/errors"

var (
	ErrPending = errors.New("pending")
	ErrDeleted = errors.New("deleted")
)
