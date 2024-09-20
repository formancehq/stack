package natstesting

import (
	"github.com/formancehq/stack/libs/go-libs/logging"
	. "github.com/formancehq/stack/libs/go-libs/testing/utils"
	. "github.com/onsi/ginkgo/v2"
)

func WithNewNatsServer(logger logging.Logger, fn func(p *Deferred[*NatsServer])) bool {
	return Context("With new postgres server", func() {
		ret := NewDeferred[*NatsServer]()
		BeforeEach(func() {
			ret.Reset()
			ret.SetValue(CreateServer(
				GinkgoT(),
				true,
				logger,
			))
		})
		fn(ret)
	})
}
