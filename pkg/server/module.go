package server

import (
	"os"

	"github.com/formancehq/go-libs/sharedlogging"
	"github.com/formancehq/go-libs/sharedotlp/pkg/sharedotlptraces"
	"github.com/formancehq/webhooks/pkg/httpserver"
	"github.com/formancehq/webhooks/pkg/storage/mongo"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func StartModule(addr string) fx.Option {
	var options []fx.Option

	options = append(options, sharedotlptraces.CLITracesModule(viper.GetViper()))

	options = append(options, fx.Provide(
		func() string { return addr },
		httpserver.NewMuxServer,
		mongo.NewStore,
		newServerHandler,
	))
	options = append(options, fx.Invoke(httpserver.RegisterHandler))
	options = append(options, fx.Invoke(httpserver.Run))

	sharedlogging.Debugf("starting server with env:")
	for _, e := range os.Environ() {
		sharedlogging.Debugf("%s", e)
	}

	return fx.Module("webhooks server", options...)
}
