package server

import (
	"net/http"

	"github.com/numary/webhooks/internal/mux"
	"github.com/numary/webhooks/internal/storage/mongo"
	"github.com/numary/webhooks/internal/svix"
	"go.uber.org/fx"
)

func StartModule(httpClient *http.Client) fx.Option {
	return fx.Module("webhooks server",
		fx.Provide(
			func() *http.Client { return httpClient },
			mongo.NewConfigStore,
			svix.New,
			newServerHandler,
			mux.NewServer,
		),
		fx.Invoke(registerHandler),
	)
}

func registerHandler(mux *http.ServeMux, h http.Handler) {
	mux.Handle("/", h)
}
