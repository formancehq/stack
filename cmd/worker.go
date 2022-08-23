package cmd

import (
	"net/http"
	"syscall"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/constants"
	"github.com/numary/webhooks/pkg/worker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Run webhooks worker",
	RunE:  RunWorker,
}

func RunWorker(cmd *cobra.Command, args []string) error {
	sharedlogging.GetLogger(cmd.Context()).Debugf(
		"starting webhooks worker module: env variables: %+v viper keys: %+v",
		syscall.Environ(), viper.AllKeys())

	app := fx.New(
		worker.StartModule(
			http.DefaultClient, viper.GetString(constants.HttpBindAddressWorkerFlag)))

	if err := app.Start(cmd.Context()); err != nil {
		return err
	}

	<-app.Done()

	if err := app.Stop(cmd.Context()); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
