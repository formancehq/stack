package task

import "github.com/pkg/errors"

var (
	// ErrRetryableError will be sent by the task if we can retry the task,
	// e.g. if the task failed because of a temporary network issue.
	ErrRetryableError = errors.New("retryable error")

	// ErrNonRetryableError will be sent by the task if we can't retry the task,
	// e.g. if the task failed because of a validation error.
	ErrNonRetryableError = errors.New("non-retryable error")
)
