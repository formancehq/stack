package cmd

import (
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/licence"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/publish"

	"github.com/formancehq/stack/libs/go-libs/service"

	fxmodules "github.com/formancehq/webhooks/cmd/fx-modules"
	apiutils "github.com/formancehq/webhooks/internal/app/webhook_server/api/utils"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func newWebhookWorkerCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "worker",
		Aliases: []string{"work", "wrk"},
		Short:   "Run webhook Worker",
		RunE:    webhookWorkerRun,
		PreRunE: handleAutoMigrate,
	}
}

func webhookWorkerRun(cmd *cobra.Command, _ []string) error {

	connectionOptions, err := bunconnect.ConnectionOptionsFromFlags(cmd.Context())
	if err != nil {
		return err
	}
	options := []fx.Option{
		auth.CLIAuthModule(),
		licence.CLIModule(ServiceName),
		otlptraces.CLITracesModule(),
		bunconnect.Module(*connectionOptions),
		publish.CLIPublisherModule(ServiceName),
		fx.Provide(
			func() apiutils.ServiceInfo {
				return apiutils.ServiceInfo{Name: ServiceName, Version: Version}
			},
		),
		fxmodules.ProvideHttpClient(),
		fxmodules.ProvideDatabase(),
		fxmodules.ProvideCacheParams(),
		fxmodules.ProvideTopics(),
		fxmodules.InvokeWorker(),
	}

	return service.New(cmd.OutOrStdout(), options...).Run(cmd.Context())
}
