package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/licence"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/uptrace/bun"

	"github.com/formancehq/webhooks/cmd/flag"

	"github.com/formancehq/webhooks/internal/commons"
	httpclient "github.com/formancehq/webhooks/internal/services/httpclient"

	component "github.com/formancehq/webhooks/internal/components/commons"
	webhookCollector "github.com/formancehq/webhooks/internal/components/webhook_collector"
	whController "github.com/formancehq/webhooks/internal/components/webhook_controller"
	webhookworker "github.com/formancehq/webhooks/internal/components/webhook_worker"

	fxmodules "github.com/formancehq/webhooks/cmd/fx-modules"
	httpserver "github.com/formancehq/webhooks/internal/services/httpserver"
	storage "github.com/formancehq/webhooks/internal/services/storage/postgres"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func newAllInOneCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "start",
		Aliases: []string{"strt"},
		Short:   "Run StandAlone Webhook",
		RunE:    allInOneRun,
		PreRunE: handleAutoMigrate,
	}
}

func allInOneRun(cmd *cobra.Command, _ []string) error {

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
		publish.CLIPublisherModule(ServiceName),
		fx.Provide(

			func() *http.Client {
				return fxmodules.FxProvideHttpClient()
			},
			func(client *http.Client) *httpclient.DefaultHttpClient {
				defaultClient := httpclient.NewDefaultHttpClient(client)
				return &defaultClient
			},

			func() *component.RunnerParams {

				runnerParams := flag.LoadRunnerParams()

				return &runnerParams
			},
			func(auth auth.Auth, logger logging.Logger) httpserver.DefaultServerParams {

				serverParams := httpserver.DefaultServerParams{}
				serverParams.Addr = viper.GetString(flag.Listen)
				serverParams.Auth = auth
				serverParams.Info = serviceInfo
				serverParams.Logger = logger

				return serverParams
			},

			func(db *bun.DB) *storage.PostgresStore {
				database := storage.NewPostgresStoreProvider(db)
				return &database
			},
			func(lc fx.Lifecycle,
				database *storage.PostgresStore,
				runnerParams *component.RunnerParams,
				client *httpclient.DefaultHttpClient,
			) *webhookCollector.Collector {

				Collector := webhookCollector.NewCollector(*runnerParams, database, client)
				Collector.Init()

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						Collector.Run()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						Collector.Stop()
						return nil
					},
				})

				return Collector
			},
			func(lc fx.Lifecycle,
				database *storage.PostgresStore,
				runnerParams *component.RunnerParams,
				client *httpclient.DefaultHttpClient,
				r *message.Router,
				subscriber message.Subscriber,
				topics []string,
			) *webhookworker.Worker {

				Worker := webhookworker.NewWorker(*runnerParams, database, client)
				Worker.Init()

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						for _, topic := range topics {
							r.AddNoPublisherHandler(fmt.Sprintf("messages-%s", topic), topic, subscriber, Worker.HandleMessage)

						}
						return nil
					},
					OnStop: func(ctx context.Context) error {
						Worker.Stop()
						return nil
					},
				})

				return Worker
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

		fx.Invoke(func(*webhookCollector.Collector, *webhookworker.Worker, *httpserver.DefaultHTTPServer) {}),
	}
	return service.New(cmd.OutOrStdout(), options...).Run(cmd.Context())

}
