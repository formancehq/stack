package server

import (
	"github.com/formancehq/webhooks/pkg/httpserver"
	"github.com/formancehq/webhooks/pkg/storage/mongo"
	"github.com/numary/go-libs/sharedotlp/pkg/sharedotlptraces"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func StartModule(addr string) fx.Option {
	var options []fx.Option

	if mod := sharedotlptraces.CLITracesModule(viper.GetViper()); mod != nil {
		options = append(options, mod)
	}

	options = append(options, fx.Provide(
		func() string { return addr },
		httpserver.NewMuxServer,
		mongo.NewStore,
		newServerHandler,
	))
	options = append(options, fx.Invoke(httpserver.RegisterHandler))
	options = append(options, fx.Invoke(httpserver.Run))

	return fx.Module("webhooks server", options...)
}
