package cmd

import (
	"fmt"
	"syscall"

	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/formancehq/webhooks/pkg/server"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run webhooks server",
	RunE:  RunServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func RunServer(cmd *cobra.Command, _ []string) error {
	sharedlogging.GetLogger(cmd.Context()).Debugf(
		"starting webhooks server module: env variables: %+v viper keys: %+v",
		syscall.Environ(), viper.AllKeys())

	app := fx.New(
		server.StartModule(
			viper.GetString(flag.HttpBindAddressServer)))

	if err := app.Start(cmd.Context()); err != nil {
		return fmt.Errorf("fx.App.Start: %w", err)
	}

	<-app.Done()

	if err := app.Stop(cmd.Context()); err != nil {
		return fmt.Errorf("fx.App.Stop: %w", err)
	}

	return nil
}
