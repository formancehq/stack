package adyen

import (
	"context"
	"errors"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/contextutil"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
)

const Name = models.ConnectorProviderAdyen

var (
	connectorAttrs = []attribute.KeyValue{
		attribute.String("connector", Name.String()),
	}
)

type Connector struct {
	logger logging.Logger
	cfg    Config
}

func (c *Connector) UpdateConfig(ctx context.Context, config models.ConnectorConfigObject) error {
	cfg, ok := config.(Config)
	if !ok {
		return connectors.ErrInvalidConfig
	}

	c.cfg = cfg

	return nil
}

func (c *Connector) InitiatePayment(ctx task.ConnectorContext, transfer *models.TransferInitiation) error {
	return connectors.ErrNotImplemented
}

func (c *Connector) CreateExternalBankAccount(ctx task.ConnectorContext, bankAccount *models.BankAccount) error {
	return connectors.ErrNotImplemented
}

func (c *Connector) SupportedCurrenciesAndDecimals() map[string]int {
	return supportedCurrenciesWithDecimal
}

func (c *Connector) Install(ctx task.ConnectorContext) error {
	taskDescriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
		Name: "Main task to periodically fetch users and transactions",
		Key:  taskNameMain,
	})
	if err != nil {
		return err
	}

	return ctx.Scheduler().Schedule(ctx.Context(), taskDescriptor, models.TaskSchedulerOptions{
		// We want to polling every c.cfg.PollingPeriod.Duration seconds the users
		// and their transactions.
		ScheduleOption: models.OPTIONS_RUN_PERIODICALLY,
		Duration:       c.cfg.PollingPeriod.Duration,
		// No need to restart this task, since the connector is not existing or
		// was uninstalled previously, the task does not exists in the database
		RestartOption: models.OPTIONS_RESTART_NEVER,
	})
}

func (c *Connector) Uninstall(ctx context.Context) error {
	return nil
}

func (c *Connector) Resolve(descriptor models.TaskDescriptor) task.Task {
	taskDescriptor, err := models.DecodeTaskDescriptor[TaskDescriptor](descriptor)
	if err != nil {
		panic(err)
	}

	return resolveTasks(c.logger, c.cfg)(taskDescriptor)
}

func (c *Connector) HandleWebhook(ctx task.ConnectorContext, webhook *models.Webhook) error {
	// Detach the context since we're launching an async task and we're mostly
	// coming from a HTTP request.
	detachedCtx, _ := contextutil.Detached(ctx.Context())
	taskDescriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
		Name:      "handle webhook",
		Key:       taskNameHandleWebhook,
		WebhookID: webhook.ID,
	})
	if err != nil {
		return err
	}

	err = ctx.Scheduler().Schedule(detachedCtx, taskDescriptor, models.TaskSchedulerOptions{
		ScheduleOption: models.OPTIONS_RUN_NOW,
		RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
	})
	if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
		return err
	}

	return nil
}

var _ connectors.Connector = &Connector{}

func newConnector(logger logging.Logger, cfg Config) *Connector {
	return &Connector{
		logger: logger.WithFields(map[string]any{
			"component": "connector",
		}),
		cfg: cfg,
	}
}
