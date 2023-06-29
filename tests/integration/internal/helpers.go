package internal

import (
	"github.com/onsi/ginkgo/v2"
)

func Given(text string, body func()) bool {
	return ginkgo.Context(text, body)
}

func With(text string, body func()) bool {
	return ginkgo.Context(text, body)
}

func Then(text string, body func()) bool {
	return ginkgo.Context(text, body)
}
