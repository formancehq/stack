package cmd

import (
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/aws/iam"
	"github.com/formancehq/stack/libs/go-libs/publish"

	"github.com/formancehq/webhooks/pkg/storage/postgres"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/licence"

	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"

	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/formancehq/webhooks/pkg/backoff"
	"github.com/formancehq/webhooks/pkg/otlp"
	"github.com/formancehq/webhooks/pkg/worker"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func newWorkerCommand() *cobra.Command {
	ret := &cobra.Command{
		Use:     "worker",
		Short:   "Run webhooks worker",
		RunE:    runWorker,
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

func runWorker(cmd *cobra.Command, _ []string) error {
	connectionOptions, err := bunconnect.ConnectionOptionsFromFlags(cmd)
	if err != nil {
		return err
	}

	retryPeriod, _ := cmd.Flags().GetDuration(flag.RetryPeriod)
	minBackOffDelay, _ := cmd.Flags().GetDuration(flag.MinBackoffDelay)
	maxBackOffDelay, _ := cmd.Flags().GetDuration(flag.MaxBackoffDelay)
	abortAfter, _ := cmd.Flags().GetDuration(flag.AbortAfter)
	topics, _ := cmd.Flags().GetStringSlice(flag.KafkaTopics)
	listen, _ := cmd.Flags().GetString(flag.Listen)

	return service.New(
		cmd.OutOrStdout(),
		otlp.HttpClientModule(),
		licence.FXModuleFromFlags(cmd, ServiceName),
		postgres.NewModule(*connectionOptions, service.IsDebug(cmd)),
		fx.Provide(worker.NewWorkerHandler),
		fx.Invoke(func(lc fx.Lifecycle, h http.Handler) {
			lc.Append(httpserver.NewHook(h, httpserver.WithAddress(listen)))
		}),
		otlptraces.FXModuleFromFlags(cmd),
		worker.StartModule(
			cmd,
			retryPeriod,
			backoff.NewExponential(
				minBackOffDelay,
				maxBackOffDelay,
				abortAfter,
			),
			service.IsDebug(cmd),
			topics,
		),
	).Run(cmd)
}
