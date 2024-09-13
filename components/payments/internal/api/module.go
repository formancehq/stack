package api

import (
	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/api/services"
	"github.com/formancehq/payments/internal/connectors/engine"
	"github.com/formancehq/payments/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

func TagVersion() fx.Annotation {
	return fx.ResultTags(`group:"apiVersions"`)
}

func NewModule(bind string) fx.Option {
	return fx.Options(
		fx.Invoke(func(m *chi.Mux, lc fx.Lifecycle) {
			lc.Append(httpserver.NewHook(m, httpserver.WithAddress(bind)))
		}),
		fx.Provide(fx.Annotate(NewRouter, fx.ParamTags(``, ``, ``, ``, `group:"apiVersions"`))),
		fx.Provide(func(storage storage.Storage, engine engine.Engine) backend.Backend {
			return services.New(storage, engine)
		}),
	)
}
