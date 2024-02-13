package cmd

import (
	"github.com/formancehq/stack/components/stargate/internal/server/grpc"
	"github.com/formancehq/stack/components/stargate/internal/server/http"
	"github.com/formancehq/stack/components/stargate/internal/server/http/controllers"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlpmetrics"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/publish"
	app "github.com/formancehq/stack/libs/go-libs/service"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.uber.org/fx"
)

const (
	serverServiceName = "stargate-server"

	serviceGrpcAddrFlag = "grpc-addr"
	serviceHttpAddrFlag = "http-addr"

	natsRequestTimeout         = "nats-request-timeout"
	authIssuerURLFlag          = "auth-issuer-url"
	maxRetriesJWKSFetchingFlag = "max-retries-jwks-fetching"
)

func newServer() *cobra.Command {
	return &cobra.Command{
		Use:          "server",
		Aliases:      []string{"serve"},
		Short:        "Launch server",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.New(cmd.OutOrStdout(), resolveServerOptions(
				viper.GetViper(),
			)...).Run(cmd.Context())
		},
	}
}

func resolveServerOptions(v *viper.Viper, userOptions ...fx.Option) []fx.Option {
	options := make([]fx.Option, 0)
	options = append(options, fx.NopLogger)

	options = append(options,
		otlptraces.CLITracesModule(),
		otlpmetrics.CLIMetricsModule(),

		// HTTP and GRPC APIs need to be started with a NATS conn, so we need to
		// create a NATS conn first.
		publish.NatsModule(
			v.GetString(publish.PublisherNatsURLFlag),
			serverServiceName,
			nats.Name(v.GetString(publish.PublisherNatsClientIDFlag)),
			nats.MaxReconnects(v.GetInt(publish.PublisherNatsMaxReconnectFlag)),
			nats.ReconnectWait(v.GetDuration(publish.PublisherNatsReconnectWaitFlag)),
		),

		fx.Provide(func() controllers.StargateControllerConfig {
			return controllers.NewStargateControllerConfig(
				Version,
				viper.GetDuration(natsRequestTimeout),
			)
		}),
		fx.Provide(fx.Annotate(noop.NewMeterProvider, fx.As(new(metric.MeterProvider)))),
		http.Module(viper.GetString(serviceHttpAddrFlag)),

		grpc.Module(
			viper.GetString(serviceGrpcAddrFlag),
			viper.GetString(authIssuerURLFlag),
			viper.GetInt(maxRetriesJWKSFetchingFlag),
		),
	)

	return append(options, userOptions...)
}
