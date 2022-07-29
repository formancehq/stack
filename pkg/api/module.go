package api

import (
	"context"
	"net/http"

	"go.uber.org/fx"
)

func Module(addr string) fx.Option {
	return fx.Options(
		fx.Provide(fx.Annotate(NewRouter, fx.As(new(http.Handler)))),
		fx.Invoke(func(lc fx.Lifecycle, handler http.Handler) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return StartServer(addr, handler)
				},
			})
		}),
	)
}
