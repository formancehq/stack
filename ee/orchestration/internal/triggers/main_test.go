package triggers

import (
	"testing"

	"github.com/formancehq/stack/libs/go-libs/testing/docker"
	"github.com/formancehq/stack/libs/go-libs/testing/utils"
	"github.com/stretchr/testify/require"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.temporal.io/sdk/testsuite"

	"github.com/formancehq/stack/libs/go-libs/testing/platform/pgtesting"
)

var (
	srv       *pgtesting.PostgresServer
	devServer *testsuite.DevServer
)

func TestMain(m *testing.M) {
	utils.WithTestMain(func(t *utils.TestingTForMain) int {
		srv = pgtesting.CreatePostgresServer(t, docker.NewPool(t, logging.Testing()))

		var err error
		devServer, err = testsuite.StartDevServer(logging.TestingContext(), testsuite.DevServerOptions{})
		require.NoError(t, err)

		t.Cleanup(func() {
			require.NoError(t, devServer.Stop())
		})

		return m.Run()
	})
}
