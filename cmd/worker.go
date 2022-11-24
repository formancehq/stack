package cmd

import (
	"fmt"
	"syscall"

	"github.com/formancehq/go-libs/sharedlogging"
	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/formancehq/webhooks/pkg/worker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Run webhooks worker",
	RunE:  RunWorker,
}

func RunWorker(cmd *cobra.Command, _ []string) error {
	sharedlogging.GetLogger(cmd.Context()).Debugf(
		"starting webhooks worker module: env variables: %+v viper keys: %+v",
		syscall.Environ(), viper.AllKeys())

	app := fx.New(
		worker.StartModule(
			viper.GetString(flag.HttpBindAddressWorker),
			viper.GetDuration(flag.RetryCron),
			retrySchedule))

	if err := app.Start(cmd.Context()); err != nil {
		return fmt.Errorf("fx.App.Start: %w", err)
	}

	<-app.Done()

	if err := app.Stop(cmd.Context()); err != nil {
		return fmt.Errorf("fx.App.Stop: %w", err)
	}

	return nil
}
