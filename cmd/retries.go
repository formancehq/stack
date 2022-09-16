package cmd

import (
	"fmt"
	"net/http"
	"syscall"

	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/formancehq/webhooks/pkg/worker/retries"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var retriesCmd = &cobra.Command{
	Use:   "retries",
	Short: "Run webhooks worker retries",
	RunE:  RunWorkerRetries,
}

func RunWorkerRetries(cmd *cobra.Command, _ []string) error {
	sharedlogging.GetLogger(cmd.Context()).Debugf(
		"starting webhooks worker retries module: env variables: %+v viper keys: %+v",
		syscall.Environ(), viper.AllKeys())

	app := fx.New(
		retries.StartModule(
			viper.GetString(flag.HttpBindAddressWorkerRetries),
			http.DefaultClient,
			viper.GetDuration(flag.RetriesCron),
			retriesSchedule,
		))

	if err := app.Start(cmd.Context()); err != nil {
		return fmt.Errorf("fx.App.Start: %w", err)
	}

	<-app.Done()

	if err := app.Stop(cmd.Context()); err != nil {
		return fmt.Errorf("fx.App.Stop: %w", err)
	}

	return nil
}

func init() {
	workerCmd.AddCommand(retriesCmd)
}
