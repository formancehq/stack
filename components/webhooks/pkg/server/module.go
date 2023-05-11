package server

import (
	"net/http"
	"os"

	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func StartModule(addr string) fx.Option {
	var options []fx.Option

	options = append(options, otlptraces.CLITracesModule(viper.GetViper()))

	options = append(options, fx.Provide(
		newServerHandler,
	), fx.Invoke(func(lc fx.Lifecycle, handler http.Handler) {
		lc.Append(httpserver.NewHook(addr, handler))
	}))

	logging.Debugf("starting server with env:")
	for _, e := range os.Environ() {
		logging.Debugf("%s", e)
	}

	return fx.Module("webhooks server", options...)
}
