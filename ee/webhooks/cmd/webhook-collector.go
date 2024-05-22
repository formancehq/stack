package cmd

import (
	"context"
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/licence"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"

	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/uptrace/bun"

	"github.com/formancehq/webhooks/cmd/flag"

	httpclient "github.com/formancehq/webhooks/internal/services/httpclient"

	component "github.com/formancehq/webhooks/internal/components/commons"
	webhookCollector "github.com/formancehq/webhooks/internal/components/webhook_collector"
	storage "github.com/formancehq/webhooks/internal/services/storage/postgres"

	fxmodules "github.com/formancehq/webhooks/cmd/fx-modules"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)





func newWebhookWCollectorCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "collector",
		Aliases: []string{"collect"},
		Short:   "Run webhook Collector",
		RunE:    webhookCollectorRun,
		PreRunE: handleAutoMigrate,
	}
}


func webhookCollectorRun(cmd *cobra.Command, _ []string) error {

	connectionOptions, err := bunconnect.ConnectionOptionsFromFlags(cmd.Context())
	if err != nil {
		return err
	}
	options := []fx.Option{
		auth.CLIAuthModule(),
		licence.CLIModule(ServiceName),
		otlptraces.CLITracesModule(),
		bunconnect.Module(*connectionOptions),
		fx.Provide(

			func() *http.Client {
				return fxmodules.FxProvideHttpClient()
			},
			func (client *http.Client) *httpclient.DefaultHttpClient {
				defaultClient := httpclient.NewDefaultHttpClient(client)
				return &defaultClient
			},
			
			func() *component.RunnerParams {
				
				runnerParams := flag.LoadRunnerParams()

				return &runnerParams
			},

			func (db *bun.DB) *storage.PostgresStore {
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
		),

		fx.Invoke(func(*webhookCollector.Collector){}),
	}


	return service.New(cmd.OutOrStdout(), options...).Run(cmd.Context())
}
