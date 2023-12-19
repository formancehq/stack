package connectors

import "errors"

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrInvalidConfig  = errors.New("invalid config")
)
