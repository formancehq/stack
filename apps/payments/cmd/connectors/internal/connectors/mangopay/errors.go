package mangopay

import "errors"

var (
	// ErrMissingTask is returned when the task is missing.
	ErrMissingTask = errors.New("task is not implemented")

	// ErrMissingClientID is returned when the clientID is missing.
	ErrMissingClientID = errors.New("missing clientID from config")

	// ErrMissingAPIKey is returned when the apiKey is missing.
	ErrMissingAPIKey = errors.New("missing apiKey from config")

	// ErrMissingEndpoint is returned when the endpoint is missing.
	ErrMissingEndpoint = errors.New("missing endpoint from config")

	// ErrMissingName is returned when the name is missing.
	ErrMissingName = errors.New("missing name from config")
)
