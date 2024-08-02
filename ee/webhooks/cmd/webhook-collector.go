package cmd

import (
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/licence"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"

	"github.com/formancehq/stack/libs/go-libs/service"

	fxmodules "github.com/formancehq/webhooks/cmd/fx-modules"

	"github.com/spf13/cobra"

	"go.uber.org/fx"
)

func newWebhookWCollectorCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "collector",
		Aliases: []string{"collect"},
		Short:   "Run webhook Collector",
		RunE:    webhookCollectorRun,
		PreRunE: handleAutoMigrate,
	}
}

func webhookCollectorRun(cmd *cobra.Command, _ []string) error {

	connectionOptions, err := bunconnect.ConnectionOptionsFromFlags(cmd.Context())
	if err != nil {
		return err
	}
	options := []fx.Option{
		auth.CLIAuthModule(),
		licence.CLIModule(ServiceName),
		otlptraces.CLITracesModule(),
		bunconnect.Module(*connectionOptions),
		fxmodules.ProvideHttpClient(),
		fxmodules.ProvideCacheParams(),
		fxmodules.ProvideDatabase(),
		fxmodules.InvokeCollector(),
	}

	return service.New(cmd.OutOrStdout(), options...).Run(cmd.Context())
}
