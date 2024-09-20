package connectors

import (
	"github.com/bombsimon/logrusr/v3"
	sharedapi "github.com/formancehq/go-libs/api"
	"github.com/formancehq/go-libs/auth"
	"github.com/formancehq/go-libs/bun/bunconnect"
	"github.com/formancehq/go-libs/otlp/otlpmetrics"
	"github.com/formancehq/go-libs/otlp/otlptraces"
	"github.com/formancehq/go-libs/publish"
	"github.com/formancehq/go-libs/service"
	"github.com/formancehq/payments/cmd/connectors/internal/api"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.uber.org/fx"
)

const (
	stackURLFlag            = "stack-url"
	postgresURIFlag         = "postgres-uri"
	configEncryptionKeyFlag = "config-encryption-key"
	envFlag                 = "env"
	listenFlag              = "listen"

	serviceName = "Payments"
)

func newServer(version string) *cobra.Command {
	return &cobra.Command{
		Use:          "serve",
		Aliases:      []string{"server"},
		Short:        "Launch server",
		SilenceUsage: true,
		RunE:         runServer(version),
	}
}

func runServer(version string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		setLogger()

		databaseOptions, err := prepareDatabaseOptions(cmd, service.IsDebug(cmd))
		if err != nil {
			return err
		}

		options := make([]fx.Option, 0)

		options = append(options, databaseOptions)
		options = append(options,
			otlptraces.FXModuleFromFlags(cmd),
			otlpmetrics.FXModuleFromFlags(cmd),
			auth.FXModuleFromFlags(cmd),
			fx.Provide(fx.Annotate(noop.NewMeterProvider, fx.As(new(metric.MeterProvider)))),
			fx.Provide(metrics.RegisterMetricsRegistry),
		)
		options = append(options, publish.FXModuleFromFlags(cmd, service.IsDebug(cmd)))
		listen, _ := cmd.Flags().GetString(listenFlag)
		stackURL, _ := cmd.Flags().GetString(stackURLFlag)
		otelTraces, _ := cmd.Flags().GetBool(otlptraces.OtelTracesFlag)

		options = append(options, api.HTTPModule(sharedapi.ServiceInfo{
			Version: version,
			Debug:   service.IsDebug(cmd),
		}, listen, stackURL, otelTraces))

		return service.New(cmd.OutOrStdout(), options...).Run(cmd)
	}
}

func setLogger() {
	// Add a dedicated logger for opentelemetry in case of error
	otel.SetLogger(logrusr.New(logrus.New().WithField("component", "otlp")))
}

func prepareDatabaseOptions(cmd *cobra.Command, debug bool) (fx.Option, error) {
	configEncryptionKey, _ := cmd.Flags().GetString(configEncryptionKeyFlag)
	if configEncryptionKey == "" {
		return nil, errors.New("missing config encryption key")
	}

	connectionOptions, err := bunconnect.ConnectionOptionsFromFlags(cmd)
	if err != nil {
		return nil, err
	}

	return storage.Module(*connectionOptions, configEncryptionKey, debug), nil
}
