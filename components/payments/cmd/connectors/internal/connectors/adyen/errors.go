package adyen

import "errors"

var (
	// ErrMissingTask is returned when the task is missing.
	ErrMissingTask = errors.New("task is not implemented")

	// ErrMissingAPIKey is returned when the apiKey is missing.
	ErrMissingAPIKey = errors.New("missing apiKey from config")

	// ErrMissingEndpoint is returned when the endpoint is missing.
	ErrMissingLiveEndpointPrefix = errors.New("missing live endpoint prefix from config")

	// ErrMissingName is returned when the name is missing.
	ErrMissingName = errors.New("missing name from config")

	// ErrMissingHMACKey is returned when the hmacKey is missing.
	ErrMissingHMACKey = errors.New("missing hmacKey from config")
)
