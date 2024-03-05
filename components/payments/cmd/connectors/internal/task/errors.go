package task

import "github.com/pkg/errors"

var (
	// ErrRetryable will be sent by the task if we can retry the task,
	// e.g. if the task failed because of a temporary network issue.
	ErrRetryable = errors.New("retryable error")

	// ErrNonRetryable will be sent by the task if we can't retry the task,
	// e.g. if the task failed because of a validation error.
	ErrNonRetryable = errors.New("non-retryable error")
)
