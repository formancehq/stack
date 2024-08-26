package api

import (
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
		)
		options = append(options, publish.FXModuleFromFlags(cmd, service.IsDebug(cmd)))
		listen, _ := cmd.Flags().GetString(listenFlag)
		stackURL, _ := cmd.Flags().GetString(stackURLFlag)
		otlpTraces, _ := cmd.Flags().GetBool(otlptraces.OtelTracesFlag)

		options = append(options, api.HTTPModule(sharedapi.ServiceInfo{
			Version: version,
			Debug:   service.IsDebug(cmd),
		}, listen, stackURL, otlpTraces))

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
