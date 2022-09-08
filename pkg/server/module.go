package server

import (
	"github.com/numary/webhooks/pkg/httpserver"
	"github.com/numary/webhooks/pkg/storage/mongo"
	"go.uber.org/fx"
)

func StartModule(addr string) fx.Option {
	return fx.Module("webhooks server",
		fx.Provide(
			func() string { return addr },
			httpserver.NewMuxServer,
			mongo.NewStore,
			newServerHandler,
		),
		fx.Invoke(httpserver.RegisterHandler),
		fx.Invoke(httpserver.Run),
	)
}
