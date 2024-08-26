package cmd

import (
	"github.com/formancehq/stack/libs/go-libs/aws/iam"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/publish"

	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/licence"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/formancehq/webhooks/pkg/backoff"
	"github.com/formancehq/webhooks/pkg/otlp"
	"github.com/formancehq/webhooks/pkg/server"
	"github.com/formancehq/webhooks/pkg/storage/postgres"
	"github.com/formancehq/webhooks/pkg/worker"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func newServeCommand() *cobra.Command {
	ret := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"server"},
		Short:   "Run webhooks server",
		RunE:    serve,
		PreRunE: handleAutoMigrate,
	}
	otlptraces.AddFlags(ret.Flags())
	publish.AddFlags(ServiceName, ret.Flags())
	auth.AddFlags(ret.Flags())
	flag.Init(ret.Flags())
	bunconnect.AddFlags(ret.Flags())
	iam.AddFlags(ret.Flags())
	service.AddFlags(ret.Flags())
	licence.AddFlags(ret.Flags())

	return ret
}

func serve(cmd *cobra.Command, _ []string) error {
	connectionOptions, err := bunconnect.ConnectionOptionsFromFlags(cmd)
	if err != nil {
		return err
	}

	listen, _ := cmd.Flags().GetString(flag.Listen)
	options := []fx.Option{
		fx.Provide(func() server.ServiceInfo {
			return server.ServiceInfo{
				Version: Version,
			}
		}),
		auth.FXModuleFromFlags(cmd),
		postgres.NewModule(*connectionOptions, service.IsDebug(cmd)),
		otlp.HttpClientModule(),
		server.FXModuleFromFlags(cmd, listen, service.IsDebug(cmd)),
		licence.FXModuleFromFlags(cmd, ServiceName),
	}
	isWorker, _ := cmd.Flags().GetBool(flag.Worker)
	if isWorker {
		retryPeriod, _ := cmd.Flags().GetDuration(flag.RetryPeriod)
		minBackOffDelay, _ := cmd.Flags().GetDuration(flag.MinBackoffDelay)
		maxBackOffDelay, _ := cmd.Flags().GetDuration(flag.MaxBackoffDelay)
		abortAfter, _ := cmd.Flags().GetDuration(flag.AbortAfter)
		topics, _ := cmd.Flags().GetStringSlice(flag.KafkaTopics)

		options = append(options, worker.StartModule(
			cmd,
			retryPeriod,
			backoff.NewExponential(
				minBackOffDelay,
				maxBackOffDelay,
				abortAfter,
			),
			service.IsDebug(cmd),
			topics,
		))
	}

	return service.New(cmd.OutOrStdout(), options...).Run(cmd)
}
