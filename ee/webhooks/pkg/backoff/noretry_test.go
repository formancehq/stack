package backoff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoRetry(t *testing.T) {
	policy := NewNoRetry()
	_, err := policy.GetRetryDelay(0)
	assert.ErrorIs(t, err, ErrMaxAttemptsReached)
}
