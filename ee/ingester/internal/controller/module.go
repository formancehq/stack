package controller

import (
	"github.com/formancehq/stack/ee/ingester/internal/drivers"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Options(
		fx.Provide(New),
		fx.Provide(NewDefaultRunner),
		fx.Provide(func(registry *drivers.Registry) ConfigValidator {
			return registry
		}),
	)
}
