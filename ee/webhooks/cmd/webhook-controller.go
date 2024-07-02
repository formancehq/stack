package cmd

import (
	"context"
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/licence"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"

	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/uptrace/bun"

	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/formancehq/webhooks/internal/commons"

	httpclient "github.com/formancehq/webhooks/internal/services/httpclient"
	httpserver "github.com/formancehq/webhooks/internal/services/httpserver"
	storage "github.com/formancehq/webhooks/internal/services/storage/postgres"

	whController "github.com/formancehq/webhooks/internal/components/webhook_controller"

	fxmodules "github.com/formancehq/webhooks/cmd/fx-modules"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func newWebhookControllerCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "controller",
		Aliases: []string{"control", "ctrl"},
		Short:   "Run webhook controller server",
		RunE:    webhookControllerRun,
		PreRunE: handleAutoMigrate,
	}
}

func webhookControllerRun(cmd *cobra.Command, _ []string) error {

	connectionOptions, err := bunconnect.ConnectionOptionsFromFlags(cmd.Context())
	if err != nil {
		return err
	}
	serviceInfo := commons.ServiceInfo{Name: ServiceName, Version: Version}

	options := []fx.Option{
		auth.CLIAuthModule(),
		licence.CLIModule(ServiceName),
		otlptraces.CLITracesModule(),
		bunconnect.Module(*connectionOptions),

		fx.Provide(

			func() *http.Client {
				return fxmodules.FxProvideHttpClient()
			},

			func(auth auth.Auth, logger logging.Logger) httpserver.DefaultServerParams {

				serverParams := httpserver.DefaultServerParams{}
				serverParams.Addr = viper.GetString(flag.Listen)
				serverParams.Auth = auth
				serverParams.Info = serviceInfo
				serverParams.Logger = logger

				return serverParams
			},

			func(db *bun.DB) storage.PostgresStore {
				return storage.NewPostgresStoreProvider(db)
			},

			func(client *http.Client) httpclient.DefaultHttpClient {
				return httpclient.NewDefaultHttpClient(client)
			},

			func(lc fx.Lifecycle,
				database storage.PostgresStore,
				serverParams httpserver.DefaultServerParams,
				client httpclient.DefaultHttpClient,
			) *httpserver.DefaultHTTPServer {

				defaultHTTPServer := httpserver.NewDefaultHTTPServer(serverParams.Addr,
					serverParams.Info, serverParams.Auth, serverParams.Logger)

				whController.Init(&defaultHTTPServer, database, &client)

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						defaultHTTPServer.Run(ctx)
						return nil
					},
					OnStop: func(ctx context.Context) error {
						defaultHTTPServer.Stop(ctx)
						return nil
					},
				})

				return &defaultHTTPServer
			},
		),

		fx.Invoke(func(*httpserver.DefaultHTTPServer) {}),
	}

	return service.New(cmd.OutOrStdout(), options...).Run(cmd.Context())
}
