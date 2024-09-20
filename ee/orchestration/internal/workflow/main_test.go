package workflow

import (
	"context"
	"testing"

	"github.com/formancehq/go-libs/testing/docker"
	"github.com/formancehq/go-libs/testing/utils"
	"github.com/stretchr/testify/require"

	"github.com/formancehq/go-libs/logging"
	"go.temporal.io/sdk/testsuite"

	"github.com/formancehq/go-libs/testing/platform/pgtesting"
)

var (
	srv       *pgtesting.PostgresServer
	devServer *testsuite.DevServer
)

func TestMain(m *testing.M) {
	utils.WithTestMain(func(t *utils.TestingTForMain) int {
		srv = pgtesting.CreatePostgresServer(t, docker.NewPool(t, logging.Testing()))

		var err error
		devServer, err = testsuite.StartDevServer(context.Background(), testsuite.DevServerOptions{})
		require.NoError(t, err)

		t.Cleanup(func() {
			require.NoError(t, devServer.Stop())
		})

		return m.Run()
	})
}
