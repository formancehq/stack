package plugins

import "errors"

var (
	ErrNotImplemented  = errors.New("not implemented")
	ErrNotYetInstalled = errors.New("not yet installed")
)
