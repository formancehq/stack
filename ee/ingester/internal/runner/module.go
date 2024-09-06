package runner

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/stack/ee/ingester/internal/drivers"
	"github.com/formancehq/stack/ee/ingester/internal/httpclient"
	"github.com/formancehq/stack/ee/ingester/internal/modules"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.uber.org/fx"
)

// NewModule create a new fx module
func NewModule(stack string, pullConfiguration modules.PullConfiguration) fx.Option {
	return fx.Options(
		fx.Provide(drivers.NewRegistry),
		fx.Provide(func(registry *drivers.Registry) drivers.Factory {
			return registry
		}),
		// decorate the original Factory (implemented by *Registry)
		// to abstract the fact we want to batch logs
		fx.Decorate(fx.Annotate(
			drivers.NewWithBatchingConnectorFactory,
			fx.As(new(drivers.Factory)),
		)),
		fx.Provide(func(
			subscriber message.Subscriber,
			httpClient *httpclient.StackAuthenticatedClient,
			logger logging.Logger,
		) modules.Factory {
			return modules.NewModuleFactory(
				subscriber,
				stack,
				httpClient,
				pullConfiguration,
				logger,
			)
		}),
		fx.Provide(NewRunner),
		fx.Invoke(func(lc fx.Lifecycle, runner *Runner) {
			lc.Append(fx.Hook{
				OnStart: runner.StartAsync,
				OnStop:  runner.Stop,
			})
		}),
		fx.Invoke(func(lc fx.Lifecycle, store Store, runner *Runner) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return RestorePipelines(ctx, store, runner)
				},
			})
		}),
	)
}

func As[T any](v T) T {
	return v
}
