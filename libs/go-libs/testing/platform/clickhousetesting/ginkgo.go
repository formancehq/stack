package clickhousetesting

import (
	. "github.com/formancehq/stack/libs/go-libs/testing/docker/ginkgo"
	. "github.com/formancehq/stack/libs/go-libs/testing/utils"
	. "github.com/onsi/ginkgo/v2"
)

func WithClickhouse(fn func(d *Deferred[*Server])) {
	Context("with clickhouse", func() {
		ret := NewDeferred[*Server]()
		BeforeEach(func() {
			ret.Reset()
			ret.SetValue(CreateServer(ActualDockerPool()))
		})
		fn(ret)
	})
}

func WithNewDatabase(srv *Deferred[*Server], fn func(d *Deferred[*Database])) {
	Context("with new database", func() {
		ret := NewDeferred[*Database]()
		BeforeEach(func() {
			ret.Reset()
			ret.SetValue(srv.GetValue().NewDatabase(GinkgoT()))
		})
		fn(ret)
	})
}
