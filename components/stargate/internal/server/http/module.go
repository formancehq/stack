package http

import (
	"github.com/formancehq/stack/components/stargate/internal/server/http/controllers"
	"github.com/formancehq/stack/components/stargate/internal/server/http/middlewares"
	"github.com/formancehq/stack/components/stargate/internal/server/http/opentelemetry"
	"github.com/formancehq/stack/components/stargate/internal/server/http/routes"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/formancehq/stack/libs/go-libs/logging"
	app "github.com/formancehq/stack/libs/go-libs/service"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func Module(bind string) fx.Option {
	return fx.Options(
		fx.Provide(opentelemetry.RegisterMetricsRegistry),
		fx.Provide(routes.NewRouter),
		fx.Provide(controllers.NewStargateController),
		health.Module(),

		fx.Invoke(func(lc fx.Lifecycle, h chi.Router, l logging.Logger) {
			if viper.GetBool(app.DebugFlag) {
				wrappedRouter := chi.NewRouter()
				wrappedRouter.Use(middlewares.Log())
				wrappedRouter.Mount("/", h)
				h = wrappedRouter
			}

			l.Infof("HTTP server listening on %s", bind)
			lc.Append(httpserver.NewHook(h, httpserver.WithAddress(bind)))
		}),
	)
}
