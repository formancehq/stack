package query

import (
	"context"
	"time"

	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Supply(workerConfig{
			// TODO(gfyrag): Probably need to be configurable
			Interval: time.Second,
		}),
		fx.Provide(NewWorker),
		fx.Invoke(func(worker *Worker, lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						if err := worker.Run(ctx); err != nil {
							panic(err)
						}
					}()
					return nil
				},
				OnStop: worker.Stop,
			})
		}),
	)
}
