package cmd

import (
	"github.com/formancehq/stack/components/stargate/internal/server/grpc"
	"github.com/formancehq/stack/components/stargate/internal/server/http"
	"github.com/formancehq/stack/components/stargate/internal/server/http/controllers"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlpmetrics"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/publish"
	app "github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"google.golang.org/grpc/keepalive"
)

const (
	serverServiceName = "stargate-server"

	serviceGrpcAddrFlag = "grpc-addr"
	serviceHttpAddrFlag = "http-addr"

	natsRequestTimeout         = "nats-request-timeout"
	authIssuerURLFlag          = "auth-issuer-url"
	maxRetriesJWKSFetchingFlag = "max-retries-jwks-fetching"

	KeepAlivePolicyMinTimeFlag                    = "keepalive-policy-min-time"
	KeepAlivePolicyPermitWithoutStreamFlag        = "keepalive-policy-permit-without-stream"
	KeepAliveServerParamMaxConnectionIdleFlag     = "keepalive-server-param-max-connection-idle"
	KeepAliveServerParamMaxConnectionAgeFlag      = "keepalive-server-param-max-connection-age"
	KeepAliveServerParamMaxConnectionAgeGraceFlag = "keepalive-server-param-max-connection-age-grace"
	KeepAliveServerParamTimeFlag                  = "keepalive-server-param-time"
	KeepAliveServerParamTimeoutFlag               = "keepalive-server-param-timeout"
)

func newServer() *cobra.Command {
	return &cobra.Command{
		Use:          "server",
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
		otlptraces.CLITracesModule(viper.GetViper()),
		otlpmetrics.CLIMetricsModule(viper.GetViper()),

		// HTTP and GRPC APIs need to be started with a NATS conn, so we need to
		// create a NATS conn first.
		publish.NatsModule(v.GetString(publish.PublisherNatsClientIDFlag), v.GetString(publish.PublisherNatsURLFlag), serverServiceName),

		fx.Provide(func() controllers.StargateControllerConfig {
			return controllers.NewStargateControllerConfig(
				Version,
				viper.GetDuration(natsRequestTimeout),
			)
		}),
		fx.Provide(fx.Annotate(metric.NewNoopMeterProvider, fx.As(new(metric.MeterProvider)))),
		http.Module(viper.GetString(serviceHttpAddrFlag)),

		fx.Provide(func() keepalive.EnforcementPolicy {
			return grpc.NewKeepAlivePolicy(
				viper.GetDuration(KeepAlivePolicyMinTimeFlag),
				viper.GetBool(KeepAlivePolicyPermitWithoutStreamFlag),
			)
		}),
		fx.Provide(func() keepalive.ServerParameters {
			return grpc.NewKeepAliveServerParams(
				viper.GetDuration(KeepAliveServerParamMaxConnectionIdleFlag),
				viper.GetDuration(KeepAliveServerParamMaxConnectionAgeFlag),
				viper.GetDuration(KeepAliveServerParamMaxConnectionAgeGraceFlag),
				viper.GetDuration(KeepAliveServerParamTimeFlag),
				viper.GetDuration(KeepAliveServerParamTimeoutFlag),
			)
		}),
		grpc.Module(
			viper.GetString(serviceGrpcAddrFlag),
			viper.GetString(authIssuerURLFlag),
			viper.GetInt(maxRetriesJWKSFetchingFlag),
		),
	)

	return append(options, userOptions...)
}
