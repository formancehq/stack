package cmd

import (
	"net/http"

	"github.com/formancehq/webhooks/pkg/storage/postgres"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/licence"

	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"

	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/formancehq/webhooks/pkg/backoff"
	"github.com/formancehq/webhooks/pkg/otlp"
	"github.com/formancehq/webhooks/pkg/worker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func newWorkerCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "worker",
		Short:   "Run webhooks worker",
		RunE:    runWorker,
		PreRunE: handleAutoMigrate,
	}
}

func runWorker(cmd *cobra.Command, _ []string) error {
	connectionOptions, err := bunconnect.ConnectionOptionsFromFlags(cmd.Context())
	if err != nil {
		return err
	}

	return service.New(
		cmd.OutOrStdout(),
		otlp.HttpClientModule(),
		licence.CLIModule(ServiceName),
		postgres.NewModule(*connectionOptions),
		fx.Provide(worker.NewWorkerHandler),
		fx.Invoke(func(lc fx.Lifecycle, h http.Handler) {
			lc.Append(httpserver.NewHook(h, httpserver.WithAddress(viper.GetString(flag.Listen))))
		}),
		otlptraces.CLITracesModule(),
		worker.StartModule(
			ServiceName,
			viper.GetDuration(flag.RetryPeriod),
			backoff.NewExponential(
				viper.GetDuration(flag.MinBackoffDelay),
				viper.GetDuration(flag.MaxBackoffDelay),
				viper.GetDuration(flag.AbortAfter),
			),
		),
	).Run(cmd.Context())
}
