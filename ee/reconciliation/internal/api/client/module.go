package client

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Options(
		fx.Provide(fx.Annotate(NewSDKFormance, fx.As(new(SDKFormance)))),
	)
}
