package service

import (
	"io"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
)

func GetDefaultLogger(w io.Writer, debug bool) logging.Logger {
	l := logging.DefaultLogger(w, debug)

	if viper.GetBool(otlptraces.OtelTracesFlag) {
		l.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		)))
	}
	return logging.NewLogrus(l)
}
