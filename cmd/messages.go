package cmd

import (
	"fmt"
	"net/http"
	"syscall"

	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/formancehq/webhooks/pkg/worker/messages"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var messagesCmd = &cobra.Command{
	Use:   "messages",
	Short: "Run webhooks worker messages",
	RunE:  RunWorkerMessages,
}

func RunWorkerMessages(cmd *cobra.Command, _ []string) error {
	sharedlogging.GetLogger(cmd.Context()).Debugf(
		"starting webhooks worker messages module: env variables: %+v viper keys: %+v",
		syscall.Environ(), viper.AllKeys())

	app := fx.New(
		messages.StartModule(
			viper.GetString(flag.HttpBindAddressWorkerMessages),
			http.DefaultClient,
			retriesSchedule))

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
	workerCmd.AddCommand(messagesCmd)
}
