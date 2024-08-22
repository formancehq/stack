package server

import (
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/webhooks/pkg/storage"

	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"go.uber.org/fx"
)

func FXModuleFromFlags(cmd *cobra.Command, addr string, debug bool) fx.Option {
	var options []fx.Option

	options = append(options, otlptraces.FXModuleFromFlags(cmd))

	options = append(options, fx.Provide(
		func(
			store storage.Store,
			httpClient *http.Client,
			logger logging.Logger,
			info ServiceInfo,
			authenticator auth.Authenticator,
		) http.Handler {
			return newServerHandler(store, httpClient, logger, info, authenticator, debug)
		},
	), fx.Invoke(func(lc fx.Lifecycle, handler http.Handler) {
		lc.Append(httpserver.NewHook(handler, httpserver.WithAddress(addr)))
	}))

	logging.Debugf("starting server with env:")
	for _, e := range os.Environ() {
		logging.Debugf("%s", e)
	}

	return fx.Module("webhooks server", options...)
}
