package cmd

import (
	"syscall"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/formancehq/webhooks/pkg/otlp"
	"github.com/formancehq/webhooks/pkg/server"
	"github.com/formancehq/webhooks/pkg/storage/postgres"
	"github.com/formancehq/webhooks/pkg/worker"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func newServeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run webhooks server",
		RunE:  serve,
	}
}

func serve(cmd *cobra.Command, _ []string) error {
	logging.FromContext(cmd.Context()).Debugf(
		"starting webhooks server module: env variables: %+v viper keys: %+v",
		syscall.Environ(), viper.AllKeys())

	options := []fx.Option{
		postgres.NewModule(viper.GetString(flag.StoragePostgresConnString)),
		otlp.HttpClientModule(),
		server.StartModule(viper.GetString(flag.Listen)),
	}

	if viper.GetBool(flag.Worker) {
		options = append(options, worker.StartModule(
			ServiceName,
			viper.GetDuration(flag.RetriesCron),
			retriesSchedule))
	}

	return errors.Wrap(service.New(cmd.OutOrStdout(), options...).Run(cmd.Context()), "staging service")
}
