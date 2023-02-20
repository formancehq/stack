package app

import (
	"context"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
)

func DefaultLoggingContext(cmd *cobra.Command, debug bool) context.Context {
	l := logrus.New()
	l.SetOutput(cmd.OutOrStdout())
	if debug {
		l.Level = logrus.DebugLevel
	}
	if viper.GetBool(otlptraces.OtelTracesFlag) {
		l.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		)))
	}
	return logging.ContextWithLogger(cmd.Context(), logging.New(l))
}
