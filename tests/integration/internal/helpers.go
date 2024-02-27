package internal

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

func Given(text string, body func()) bool {
	return Context(text, body)
}

func With(text string, body func()) bool {
	return Context(text, body)
}

func Then(text string, body func()) bool {
	return Context(text, body)
}

func NeedModule(m *Module) {
	BeforeEach(func() {
		Expect(currentTest.loadModule(ctx, m)).To(Succeed())
	})
	AfterEach(func() {
		Expect(currentTest.unloadModule(ctx, m)).To(Succeed())
	})
}

func WithModules(modules []*Module, callback func()) bool {
	Context("with modules", func() {
		for _, m := range modules {
			NeedModule(m)
		}
		callback()
	})
	return true
}

func WaitOnChanWithTimeout[T any](ch chan T, timeout time.Duration) T {
	select {
	case t := <-ch:
		return t
	case <-time.After(timeout):
		Fail("should have received an event", 1)
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
