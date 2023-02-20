package cmd

import (
	"fmt"
	"syscall"

	"github.com/formancehq/stack/libs/go-libs/app"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/formancehq/webhooks/pkg/otlp"
	"github.com/formancehq/webhooks/pkg/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run webhooks server",
	RunE:  RunServer,
}

func RunServer(cmd *cobra.Command, _ []string) error {
	logging.FromContext(cmd.Context()).Debugf(
		"starting webhooks server module: env variables: %+v viper keys: %+v",
		syscall.Environ(), viper.AllKeys())

	ctx := app.DefaultLoggingContext(cmd, viper.GetBool(flag.Debug))

	app := fx.New(
		otlp.HttpClientModule(),
		server.StartModule(
			viper.GetString(flag.HttpBindAddressServer)))

	if err := app.Start(ctx); err != nil {
		return fmt.Errorf("fx.App.Start: %w", err)
	}

	<-app.Done()

	if err := app.Stop(ctx); err != nil {
		return fmt.Errorf("fx.App.Stop: %w", err)
	}

	return nil
}
