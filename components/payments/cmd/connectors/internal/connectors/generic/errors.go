package generic

import "errors"

var (
	// ErrMissingTask is returned when the task is missing.
	ErrMissingTask = errors.New("task is not implemented")

	// ErrMissingEndpoint is returned when the endpoint is missing.
	ErrMissingEndpoint = errors.New("missing endpoint from config")

	// ErrMissingName is returned when the name is missing.
	ErrMissingName = errors.New("missing name from config")
)
