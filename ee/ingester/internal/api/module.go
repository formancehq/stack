package api

import (
	"github.com/formancehq/stack/ee/ingester/internal/controller"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"go.uber.org/fx"
)

func NewModule(bind string) fx.Option {
	return fx.Options(
		fx.Provide(NewAPI),
		fx.Provide(func(ctrl *controller.Controller) Backend {
			return ctrl
		}),
		fx.Invoke(func(api *API, lc fx.Lifecycle) {
			lc.Append(httpserver.NewHook(api.Router(), httpserver.WithAddress(bind)))
		}),
		health.Module(),
	)
}
