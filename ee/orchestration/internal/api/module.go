package api

import (
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

func TagVersion() fx.Annotation {
	return fx.ResultTags(`group:"apiVersions"`)
}

func NewModule(debug bool) fx.Option {
	return fx.Options(
		fx.Provide(fx.Annotate(func(
			backend Backend,
			info ServiceInfo,
			healthController *health.HealthController,
			authenticator auth.Authenticator,
			versions ...Version) *chi.Mux {
			return NewRouter(backend, info, healthController, authenticator, debug, versions...)
		}, fx.ParamTags(``, ``, ``, ``, `group:"apiVersions"`))),
		fx.Provide(NewDefaultBackend),
	)
}
