package api

import (
	"context"
	"net/http"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"

	"github.com/pkg/errors"

	"github.com/formancehq/payments/internal/app/metrics"
	"github.com/formancehq/payments/internal/app/models"

	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/integration"
	"github.com/formancehq/payments/internal/app/storage"
	"github.com/formancehq/payments/internal/app/task"
	"go.uber.org/dig"
	"go.uber.org/fx"
)

type connectorHandler struct {
	Handler  http.Handler
	Provider models.ConnectorProvider

	// TODO(polo): refactor to remove this ugly hack to access the connector manager
	initiatePayment paymentHandler
}

func addConnector[ConnectorConfig models.ConnectorConfigObject](loader integration.Loader[ConnectorConfig],
) fx.Option {
	return fx.Options(
		fx.Provide(func(store *storage.Storage,
			publisher message.Publisher,
			metricsRegistry metrics.MetricsRegistry,
		) *integration.ConnectorManager[ConnectorConfig] {
			schedulerFactory := integration.TaskSchedulerFactoryFn(func(
				resolver task.Resolver, maxTasks int,
			) *task.DefaultTaskScheduler {
				return task.NewDefaultScheduler(loader.Name(), store, func(ctx context.Context,
					descriptor models.TaskDescriptor,
					taskID uuid.UUID,
				) (*dig.Container, error) {
					container := dig.New()

					if err := container.Provide(func() ingestion.Ingester {
						return ingestion.NewDefaultIngester(loader.Name(), descriptor, store, publisher)
					}); err != nil {
						return nil, err
					}

					if err := container.Provide(func() storage.Reader {
						return store
					}); err != nil {
						return nil, err
					}

					return container, nil
				}, resolver, metricsRegistry, maxTasks)
			})

			return integration.NewConnectorManager(
				store, loader, schedulerFactory, publisher)
		}),
		fx.Provide(fx.Annotate(func(cm *integration.ConnectorManager[ConnectorConfig],
		) connectorHandler {
			return connectorHandler{
				Handler:         connectorRouter(loader.Name(), cm),
				Provider:        loader.Name(),
				initiatePayment: cm.InitiatePayment,
			}
		}, fx.ResultTags(`group:"connectorHandlers"`))),
		fx.Invoke(func(lc fx.Lifecycle, cm *integration.ConnectorManager[ConnectorConfig]) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := cm.Restore(ctx)
					if err != nil && !errors.Is(err, integration.ErrNotInstalled) {
						return err
					}

					return nil
				},
				OnStop: cm.Close,
			})
		}),
	)
}
