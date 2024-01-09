package cmd

import (
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/formancehq/webhooks/pkg/backoff"
	"github.com/formancehq/webhooks/pkg/otlp"
	"github.com/formancehq/webhooks/pkg/server"
	"github.com/formancehq/webhooks/pkg/storage/postgres"
	"github.com/formancehq/webhooks/pkg/worker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func newServeCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "serve",
		Aliases: []string{"server"},
		Short:   "Run webhooks server",
		RunE:    serve,
	}
}

func serve(cmd *cobra.Command, _ []string) error {
	options := []fx.Option{
		fx.Provide(func() server.ServiceInfo {
			return server.ServiceInfo{
				Version: Version,
			}
		}),
		auth.CLIAuthModule(viper.GetViper()),
		postgres.NewModule(viper.GetString(flag.StoragePostgresConnString)),
		otlp.HttpClientModule(),
		server.StartModule(viper.GetString(flag.Listen)),
	}
	if viper.GetBool(flag.Worker) {
		options = append(options, worker.StartModule(
			ServiceName,
			viper.GetDuration(flag.RetryPeriod),
			backoff.NewExponential(
				viper.GetDuration(flag.MinBackoffDelay),
				viper.GetDuration(flag.MaxBackoffDelay),
				viper.GetDuration(flag.AbortAfter),
			),
		))
	}

	return service.New(cmd.OutOrStdout(), options...).Run(cmd.Context())
}
