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
		ScheduleOption: models.OPTIONS_RUN_INDEFINITELY,
		Duration:       c.cfg.PollingPeriod.Duration,
		// No need to restart this task, since the connector is not existing or
		// was uninstalled previously, the task does not exists in the database
		Restart: false,
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

var _ integration.Connector = &Connector{}

func newConnector(logger logging.Logger, cfg Config) *Connector {
	return &Connector{
		logger: logger.WithFields(map[string]any{
			"component": "connector",
		}),
		cfg: cfg,
	}
}
