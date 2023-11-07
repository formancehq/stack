package wise

import "github.com/pkg/errors"

var (
	// ErrMissingTask is returned when the task is missing.
	ErrMissingTask = errors.New("task is not implemented")

	// ErrMissingAPIKey is returned when the api key is missing from config.
	ErrMissingAPIKey = errors.New("missing apiKey from config")

	// ErrMissingName is returned when the name is missing from config.
	ErrMissingName = errors.New("missing name from config")
)
