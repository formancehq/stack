package v3

import (
	"github.com/formancehq/payments/internal/api"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Options(
		fx.Supply(fx.Annotate(api.Version{
			Version: 3,
			Builder: newRouter,
		}, api.TagVersion())),
	)
}
