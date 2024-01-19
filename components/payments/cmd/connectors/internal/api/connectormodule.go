package api

import (
	"context"
	"net/http"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/payments/cmd/connectors/internal/api/backend"
	manager "github.com/formancehq/payments/cmd/connectors/internal/api/connectors_manager"
	"github.com/formancehq/payments/cmd/connectors/internal/api/service"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/messages"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"go.uber.org/fx"
)

type connectorHandler struct {
	Handler        http.Handler
	WebhookHandler http.Handler
	Provider       models.ConnectorProvider

	// TODO(polo): refactor to remove this ugly hack to access the connector manager
	initiatePayment           service.InitiatePaymentHandler
	reversePayment            service.ReversePaymentHandler
	createExternalBankAccount service.BankAccountHandler
}

func addConnector[ConnectorConfig models.ConnectorConfigObject](loader manager.Loader[ConnectorConfig],
) fx.Option {
	return fx.Options(
		fx.Provide(func(store *storage.Storage,
			publisher message.Publisher,
			metricsRegistry metrics.MetricsRegistry,
			messages *messages.Messages,
		) *manager.ConnectorsManager[ConnectorConfig] {
			schedulerFactory := manager.TaskSchedulerFactoryFn(func(
				connectorID models.ConnectorID, resolver task.Resolver, maxTasks int,
			) *task.DefaultTaskScheduler {
				return task.NewDefaultScheduler(connectorID, store, func(ctx context.Context,
					descriptor models.TaskDescriptor,
					taskID uuid.UUID,
				) (*dig.Container, error) {
					container := dig.New()

					if err := container.Provide(func() ingestion.Ingester {
						return ingestion.NewDefaultIngester(loader.Name(), descriptor, store, publisher, messages)
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

			return manager.NewConnectorManager(
				loader.Name(), store, loader, schedulerFactory, publisher, messages)
		}),
		fx.Provide(func(cm *manager.ConnectorsManager[ConnectorConfig]) backend.ManagerBackend[ConnectorConfig] {
			return backend.NewDefaultManagerBackend[ConnectorConfig](cm)
		}),
		fx.Provide(fx.Annotate(func(
			b backend.ManagerBackend[ConnectorConfig],
			cm *manager.ConnectorsManager[ConnectorConfig],
		) connectorHandler {
			return connectorHandler{
				Handler:                   connectorRouter(loader.Name(), b),
				WebhookHandler:            webhookConnectorRouter(loader.Name(), loader.Router(), b),
				Provider:                  loader.Name(),
				initiatePayment:           cm.InitiatePayment,
				reversePayment:            cm.ReversePayment,
				createExternalBankAccount: cm.CreateExternalBankAccount,
			}
		}, fx.ResultTags(`group:"connectorHandlers"`))),
		fx.Invoke(func(lc fx.Lifecycle, cm *manager.ConnectorsManager[ConnectorConfig]) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					ctx, span := otel.Tracer().Start(ctx, "connectorsManager.Restore")
					defer span.End()

					err := cm.Restore(ctx)
					if err != nil && !errors.Is(err, manager.ErrNotInstalled) {
						return err
					}

					return nil
				},
				OnStop: cm.Close,
			})
		}),
	)
}
