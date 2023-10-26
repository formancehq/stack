package api

import (
	"io"

	"github.com/bombsimon/logrusr/v3"
	"github.com/formancehq/payments/cmd/api/internal/api"
	"github.com/formancehq/payments/cmd/api/internal/storage"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: runServer(version),
	}
}

func runServer(version string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		setLogger()

		databaseOptions, err := prepareDatabaseOptions(cmd.OutOrStdout())
		if err != nil {
			return err
		}

		options := make([]fx.Option, 0)

		options = append(options, databaseOptions)
		options = append(options,
			otlptraces.CLITracesModule(viper.GetViper()),
			otlpmetrics.CLIMetricsModule(viper.GetViper()),
			fx.Provide(fx.Annotate(noop.NewMeterProvider, fx.As(new(metric.MeterProvider)))),
		)
		options = append(options, publish.CLIPublisherModule(viper.GetViper(), serviceName))
		options = append(options, api.HTTPModule(sharedapi.ServiceInfo{
			Version: version,
		}, viper.GetString(listenFlag)))

		return service.New(cmd.OutOrStdout(), options...).Run(cmd.Context())
	}
}

func setLogger() {
	// Add a dedicated logger for opentelemetry in case of error
	otel.SetLogger(logrusr.New(logrus.New().WithField("component", "otlp")))
}

func prepareDatabaseOptions(output io.Writer) (fx.Option, error) {
	postgresURI := viper.GetString(postgresURIFlag)
	if postgresURI == "" {
		return nil, errors.New("missing postgres uri")
	}

	configEncryptionKey := viper.GetString(configEncryptionKeyFlag)
	if configEncryptionKey == "" {
		return nil, errors.New("missing config encryption key")
	}

	return storage.Module(postgresURI, configEncryptionKey, viper.GetBool(service.DebugFlag), output), nil
}
