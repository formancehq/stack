package cmd

import (
	"context"

	"github.com/bombsimon/logrusr/v3"
	"github.com/formancehq/payments/internal/app/api"
	"github.com/formancehq/payments/internal/app/storage"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/logging/logginglogrus"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"go.opentelemetry.io/otel"
	"go.uber.org/fx"
)

//nolint:gosec // false positive
const (
	postgresURIFlag                 = "postgres-uri"
	configEncryptionKeyFlag         = "config-encryption-key"
	otelTracesFlag                  = "otel-traces"
	envFlag                         = "env"
	authBasicEnabledFlag            = "auth-basic-enabled"
	authBasicCredentialsFlag        = "auth-basic-credentials"
	authBearerEnabledFlag           = "auth-bearer-enabled"
	authBearerIntrospectURLFlag     = "auth-bearer-introspect-url"
	authBearerAudienceFlag          = "auth-bearer-audience"
	authBearerAudiencesWildcardFlag = "auth-bearer-audiences-wildcard"
	authBearerUseScopesFlag         = "auth-bearer-use-scopes"
	listenFlag                      = "listen"
	autoMigrateFlag                 = "auto-migrate"

	serviceName = "Payments"
)

func newServer() *cobra.Command {
	return &cobra.Command{
		Use:          "server",
		Short:        "Launch server",
		SilenceUsage: true,
		RunE:         runServer,
	}
}

func runServer(cmd *cobra.Command, args []string) error {
	setLogger()

	if viper.GetBool(autoMigrateFlag) {
		if err := runMigrate(cmd, []string{"up"}); err != nil {
			return err
		}
	}

	databaseOptions, err := prepareDatabaseOptions()
	if err != nil {
		return err
	}

	options := make([]fx.Option, 0)

	if !viper.GetBool(debugFlag) {
		options = append(options, fx.NopLogger)
	}

	options = append(options, databaseOptions)
	options = append(options, otlptraces.CLITracesModule(viper.GetViper()))
	options = append(options, publish.CLIPublisherModule(viper.GetViper(), serviceName))
	options = append(options, api.HTTPModule(sharedapi.ServiceInfo{
		Version: Version,
	}, viper.GetString(listenFlag)))

	app := fx.New(options...)
	err = app.Start(cmd.Context())
	if err != nil {
		return err
	}

	select {
	case <-cmd.Context().Done():
		return app.Stop(context.Background())
	case <-app.Done():
		return app.Err()
	}
}

func setLogger() {
	log := logrus.New()

	if viper.GetBool(debugFlag) {
		log.SetLevel(logrus.DebugLevel)
	}

	if viper.GetBool(otelTracesFlag) {
		log.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		)))
		log.SetFormatter(&logrus.JSONFormatter{})
	}

	logging.SetFactory(logging.StaticLoggerFactory(logginglogrus.New(log)))

	// Add a dedicated logger for opentelemetry in case of error
	otel.SetLogger(logrusr.New(logrus.New().WithField("component", "otlp")))
}

func prepareDatabaseOptions() (fx.Option, error) {
	postgresURI := viper.GetString(postgresURIFlag)
	if postgresURI == "" {
		return nil, errors.New("missing postgres uri")
	}

	configEncryptionKey := viper.GetString(configEncryptionKeyFlag)
	if configEncryptionKey == "" {
		return nil, errors.New("missing config encryption key")
	}

	return storage.Module(postgresURI, configEncryptionKey), nil
}
