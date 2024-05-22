package fxmodules

import (
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/formancehq/webhooks/internal/app/webhook_server/api/router"
	apiutils "github.com/formancehq/webhooks/internal/app/webhook_server/api/utils"
	httpclient "github.com/formancehq/webhooks/internal/services/httpclient"
	storage "github.com/formancehq/webhooks/internal/services/storage/postgres"

	"go.uber.org/fx"
)

func InvokeServer() fx.Option {
	return fx.Invoke(
		func(
			lc fx.Lifecycle,
			healthcontroller *health.HealthController,
			database *storage.PostgresStore,
			serverParams *apiutils.DefaultServerParams,
			client *httpclient.DefaultHttpClient,
		) {
			router := router.NewRouter(database, client, healthcontroller, serverParams.Auth, serverParams.Info)
			lc.Append(httpserver.NewHook(router, httpserver.WithAddress(serverParams.Addr)))
		})

}
