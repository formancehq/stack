package system

import (
	"go.uber.org/fx"
)

func NewFXModule() fx.Option {
	return fx.Options(
		fx.Provide(NewController),
		fx.Provide(func(controller *DefaultController) Controller {
			return controller
		}),
	)
}
