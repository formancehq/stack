package api

import (
	"context"

	"github.com/bombsimon/logrusr/v3"
	"github.com/formancehq/payments/cmd/api/internal/api"
	"github.com/formancehq/payments/cmd/api/internal/storage"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlpmetrics"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.uber.org/fx"
)

const (
	stackURLFlag            = "stack-url"
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: runServer(version),
	}
}

func runServer(version string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		setLogger()

		databaseOptions, err := prepareDatabaseOptions(cmd.Context())
		if err != nil {
			return err
		}

		options := make([]fx.Option, 0)

		options = append(options, databaseOptions)
		options = append(options,
			otlptraces.CLITracesModule(),
			otlpmetrics.CLIMetricsModule(),
			auth.CLIAuthModule(),
			fx.Provide(fx.Annotate(noop.NewMeterProvider, fx.As(new(metric.MeterProvider)))),
		)
		options = append(options, publish.CLIPublisherModule(serviceName))
		options = append(options, api.HTTPModule(sharedapi.ServiceInfo{
			Version: version,
		}, viper.GetString(listenFlag), viper.GetString(stackURLFlag)))

		return service.New(cmd.OutOrStdout(), options...).Run(cmd.Context())
	}
}

func setLogger() {
	// Add a dedicated logger for opentelemetry in case of error
	otel.SetLogger(logrusr.New(logrus.New().WithField("component", "otlp")))
}

func prepareDatabaseOptions(ctx context.Context) (fx.Option, error) {
	configEncryptionKey := viper.GetString(configEncryptionKeyFlag)
	if configEncryptionKey == "" {
		return nil, errors.New("missing config encryption key")
	}

	connectionOptions, err := bunconnect.ConnectionOptionsFromFlags(ctx)
	if err != nil {
		return nil, err
	}

	return storage.Module(*connectionOptions, configEncryptionKey), nil
}
