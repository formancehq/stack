package cmd

import (
	"syscall"

	"github.com/formancehq/stack/libs/go-libs/app"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/formancehq/webhooks/pkg/otlp"
	"github.com/formancehq/webhooks/pkg/worker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Run webhooks worker",
	RunE:  RunWorker,
}

func RunWorker(cmd *cobra.Command, _ []string) error {
	logging.FromContext(cmd.Context()).Debugf(
		"starting webhooks worker module: env variables: %+v viper keys: %+v",
		syscall.Environ(), viper.AllKeys())

	return app.New(
		cmd.OutOrStdout(),
		otlp.HttpClientModule(),
		worker.StartModule(
			viper.GetString(flag.HttpBindAddressWorker),
			ServiceName,
			viper.GetDuration(flag.RetriesCron),
			retriesSchedule),
	).Run(cmd.Context())
}
