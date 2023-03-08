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
		ginkgo.Fail("should have received an event")
	}
	panic("cannot happen")
}

func ChanClosed[T any](ch chan T) func() bool {
	return func() bool {
		select {
		case _, alive := <-ch:
			return !alive
		default:
			return false
		}
	}
}
