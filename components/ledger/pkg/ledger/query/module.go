package query

import (
	"context"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Supply(WorkerConfig{
			// TODO(gfyrag): Probably need to be configurable
			ChanSize: 1024,
		}),
		fx.Provide(NewWorker),
		fx.Provide(fx.Annotate(NewNoOpMonitor, fx.As(new(Monitor)))),
		fx.Invoke(func(worker *Worker, lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						if err := worker.Run(logging.ContextWithLogger(
							context.Background(),
							logging.FromContext(ctx),
						)); err != nil {
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
