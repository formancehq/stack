package server

import (
	"net/http"

	"github.com/numary/webhooks/pkg/httpserver"
	"github.com/numary/webhooks/pkg/storage/mongo"
	"github.com/numary/webhooks/pkg/svix"
	"go.uber.org/fx"
)

func StartModule(httpClient *http.Client, addr string) fx.Option {
	return fx.Module("webhooks server",
		fx.Provide(
			func() (*http.Client, string) { return httpClient, addr },
			httpserver.NewMuxServer,
			mongo.NewConfigStore,
			svix.New,
			newServerHandler,
		),
		fx.Invoke(httpserver.RegisterHandler),
		fx.Invoke(httpserver.Run),
	)
}
