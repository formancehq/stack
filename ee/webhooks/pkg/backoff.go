package webhooks

import "time"

type BackoffPolicy interface {
	GetRetryDelay(attemptNumber int) (time.Duration, error)
}
