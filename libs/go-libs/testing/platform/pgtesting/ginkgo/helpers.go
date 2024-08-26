package ginkgo

import (
	"reflect"

	. "github.com/formancehq/stack/libs/go-libs/testing/docker/ginkgo"
	"github.com/formancehq/stack/libs/go-libs/testing/platform/pgtesting"
	. "github.com/onsi/ginkgo/v2"
)

var actualServer = new(pgtesting.PostgresServer)

func WithNewPostgresServer(fn func()) bool {
	return Context("With new postgres server", func() {
		BeforeEach(func() {
			*actualServer = *pgtesting.CreatePostgresServer(
				GinkgoT(),
				ActualDockerPool(),
			)
		})
		fn()
	})
}

func ActualServer() *pgtesting.PostgresServer {
	if reflect.ValueOf(*actualServer).IsZero() {
		Fail("server is not configured")
	}
	return actualServer
}

var actualDatabase = new(pgtesting.Database)

func ActualDatabase() *pgtesting.Database {
	return actualDatabase
}

func WithNewPostgresDatabase(fn func()) {
	Context("With new postgres database", func() {
		BeforeEach(func() {
			*actualDatabase = *ActualServer().NewDatabase()
		})
		fn()
	})
}
