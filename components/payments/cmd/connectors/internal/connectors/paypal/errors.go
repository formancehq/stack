package paypal

import "errors"

var (
	// ErrMissingTask is returned when the task is missing.
	ErrMissingTask = errors.New("task is not implemented")

	// ErrMissingClientID is returned when the clientID is missing.
	ErrMissingClientID = errors.New("missing clientID from config")

	// ErrMissingSecret is returned when the secret is missing.
	ErrMissingSecret = errors.New("missing secret from config")

	// ErrMissingEndpoint is returned when the endpoint is missing.
	ErrMissingEndpoint = errors.New("missing endpoint from config")

	// ErrMissingName is returned when the name is missing.
	ErrMissingName = errors.New("missing name from config")
)
