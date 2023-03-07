package cmd

import (
	"net/http"
	"syscall"

	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/formancehq/webhooks/pkg/otlp"
	"github.com/formancehq/webhooks/pkg/storage/postgres"
	"github.com/formancehq/webhooks/pkg/worker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func newWorkerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "worker",
		Short: "Run webhooks worker",
		RunE:  runWorker,
	}
}

func runWorker(cmd *cobra.Command, _ []string) error {
	logging.FromContext(cmd.Context()).Debugf(
		"starting webhooks worker module: env variables: %+v viper keys: %+v",
		syscall.Environ(), viper.AllKeys())

	return service.New(
		cmd.OutOrStdout(),
		otlp.HttpClientModule(),
		postgres.NewModule(viper.GetString(flag.StoragePostgresConnString)),
		fx.Provide(worker.NewWorkerHandler),
		fx.Invoke(func(lc fx.Lifecycle, h http.Handler) {
			lc.Append(httpserver.NewHook(viper.GetString(flag.Listen), h))
		}),
		worker.StartModule(
			ServiceName,
			viper.GetDuration(flag.RetriesCron),
			retriesSchedule),
	).Run(cmd.Context())
}
