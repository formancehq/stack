package cmd

import (
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/formancehq/stack/libs/go-libs/licence"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"

	"github.com/formancehq/stack/libs/go-libs/service"

	fxmodules "github.com/formancehq/webhooks/cmd/fx-modules"
	apiutils "github.com/formancehq/webhooks/internal/app/webhook_server/api/utils"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func newWebhookControllerCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "serve",
		Aliases: []string{"server"},
		Short:   "Run webhook controller server",
		RunE:    webhookControllerRun,
		PreRunE: handleAutoMigrate,
	}
}

func webhookControllerRun(cmd *cobra.Command, _ []string) error {

	connectionOptions, err := bunconnect.ConnectionOptionsFromFlags(cmd.Context())
	if err != nil {
		return err
	}

	options := []fx.Option{
		auth.CLIAuthModule(),
		licence.CLIModule(ServiceName),
		otlptraces.CLITracesModule(),
		bunconnect.Module(*connectionOptions),
		health.Module(),
		fx.Provide(
			func() apiutils.ServiceInfo {
				return apiutils.ServiceInfo{Name: ServiceName, Version: Version}
			},
		),
		fxmodules.ProvideDatabase(),
		fxmodules.ProvideHttpClient(),
		fxmodules.ProvideCacheParams(),
		fxmodules.ProvideServerParams(),
		fxmodules.InvokeServer(),
	}

	return service.New(cmd.OutOrStdout(), options...).Run(cmd.Context())
}
