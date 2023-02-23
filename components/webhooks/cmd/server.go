package cmd

import (
	"syscall"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/formancehq/webhooks/pkg/otlp"
	"github.com/formancehq/webhooks/pkg/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	return service.New(cmd.OutOrStdout(),
		otlp.HttpClientModule(),
		server.StartModule(viper.GetString(flag.HttpBindAddressServer)),
	).Run(cmd.Context())
}
