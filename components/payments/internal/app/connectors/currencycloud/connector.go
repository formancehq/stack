package currencycloud

import (
	"context"

	"github.com/formancehq/payments/internal/app/models"

	"github.com/formancehq/payments/internal/app/integration"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

const Name = models.ConnectorProviderCurrencyCloud

type Connector struct {
	logger logging.Logger
	cfg    Config
}

func (c *Connector) InitiateTransfer(ctx task.ConnectorContext, transfer models.Transfer) error {
	// TODO implement me
	panic("implement me")
}

func (c *Connector) Install(ctx task.ConnectorContext) error {
	taskDescriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{Name: taskNameFetchTransactions})
	if err != nil {
		return err
	}

	return ctx.Scheduler().Schedule(ctx.Context(), taskDescriptor, models.TaskSchedulerOptions{
		ScheduleOption: models.OPTIONS_RUN_NOW,
		// No need to restart this task, since the connector is not existing or
		// was uninstalled previously, the task does not exists in the database
		Restart: false,
	})
}

func (c *Connector) Uninstall(ctx context.Context) error {
	return nil
}

func (c *Connector) Resolve(descriptor models.TaskDescriptor) task.Task {
	return resolveTasks(c.logger, c.cfg)
}

var _ integration.Connector = &Connector{}

func newConnector(logger logging.Logger, cfg Config) *Connector {
	return &Connector{
		logger: logger.WithFields(map[string]any{
			"component": "connector",
		}),
		cfg: cfg,
	}
}
