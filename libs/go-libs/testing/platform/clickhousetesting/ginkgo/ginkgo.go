package ginkgo

import (
	. "github.com/formancehq/stack/libs/go-libs/testing/docker/ginkgo"
	"github.com/formancehq/stack/libs/go-libs/testing/platform/clickhousetesting"
	. "github.com/onsi/ginkgo/v2"
)

var (
	actualClickhouseServer = new(clickhousetesting.Server)
)

func ActualClickhouseServer() *clickhousetesting.Server {
	return actualClickhouseServer
}

func WithClickhouse(fn func()) {
	Context("with clickhouse", func() {
		BeforeEach(func() {
			*actualClickhouseServer = *clickhousetesting.CreateServer(ActualDockerPool())
		})
		fn()
	})
}
