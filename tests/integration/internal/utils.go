package internal

import (
	"time"

	"github.com/onsi/ginkgo/v2"
)

func WaitOnChanWithTimeout[T any](ch chan T, timeout time.Duration) T {
	select {
	case t := <-ch:
		return t
	case <-time.After(timeout):
		ginkgo.Fail("should have received a created transaction event")
	}
	panic("cannot happen")
}
