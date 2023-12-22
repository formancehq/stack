package backoff

import (
	"errors"
	"time"

	webhooks "github.com/formancehq/webhooks/pkg"
)

var ErrMaxAttemptsReached = errors.New("max attempts reached")

func NewExponential(minRetryDelay, maxRetryDelay, abortAfterDelay time.Duration) webhooks.BackoffPolicy {
	return &exponential{
		minRetryDelay,
		maxRetryDelay,
		abortAfterDelay,
	}
}

type exponential struct {
	minRetryDelay   time.Duration
	maxRetryDelay   time.Duration
	abortAfterDelay time.Duration
}

func (e *exponential) GetRetryDelay(attemptNumber int) (time.Duration, error) {
	delay := e.minRetryDelay
	sinceFirstAttempt := delay
	for i := 0; i < attemptNumber; i++ {
		delay <<= 1
		if delay > e.maxRetryDelay {
			delay = e.maxRetryDelay
		}
		sinceFirstAttempt += delay
	}
	if sinceFirstAttempt > e.abortAfterDelay {
		return 0, ErrMaxAttemptsReached
	}
	return delay, nil
}
