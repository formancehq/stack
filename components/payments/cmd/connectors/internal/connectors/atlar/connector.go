package atlar

import (
	"context"
	"errors"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

const Name = models.ConnectorProviderAtlar

var (
	mainTaskDescriptor = TaskDescriptor{
		Name: "Main task to periodically fetch transactions",
		Main: true,
	}
)

type Connector struct {
	logger logging.Logger
	cfg    Config
}

func (c *Connector) UpdateConfig(ctx task.ConnectorContext, config models.ConnectorConfigObject) error {
	cfg, ok := config.(Config)
	if !ok {
		return connectors.ErrInvalidConfig
	}

	restartTask := c.cfg.PollingPeriod.Duration != cfg.PollingPeriod.Duration

	c.cfg = cfg

	if restartTask {
		descriptor, err := models.EncodeTaskDescriptor(mainTaskDescriptor)
		if err != nil {
			return err
		}

		return ctx.Scheduler().Schedule(ctx.Context(), descriptor, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_PERIODICALLY,
			Duration:       c.cfg.PollingPeriod.Duration,
			// No need to restart this task, since the connector is not existing or
			// was uninstalled previously, the task does not exists in the database
			RestartOption: models.OPTIONS_STOP_AND_RESTART,
		})
	}

	return nil
}

func (c *Connector) Install(ctx task.ConnectorContext) error {
	descriptor, err := models.EncodeTaskDescriptor(mainTaskDescriptor)
	if err != nil {
		return err
	}

	return ctx.Scheduler().Schedule(ctx.Context(), descriptor, models.TaskSchedulerOptions{
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

func (c *Connector) SupportedCurrenciesAndDecimals() map[string]int {
	return supportedCurrenciesWithDecimal
}

func (c *Connector) Resolve(descriptor models.TaskDescriptor) task.Task {
	taskDescriptor, err := models.DecodeTaskDescriptor[TaskDescriptor](descriptor)
	if err != nil {
		panic(err)
	}

	return resolveTasks(c.logger, c.cfg)(taskDescriptor)
}

func (c *Connector) CreateExternalBankAccount(ctx task.ConnectorContext, bankAccount *models.BankAccount) error {
	descriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
		Name:        "Create external bank account",
		Key:         taskNameCreateExternalBankAccount,
		BankAccount: bankAccount,
	})
	if err != nil {
		return err
	}
	if err := ctx.Scheduler().Schedule(ctx.Context(), descriptor, models.TaskSchedulerOptions{
		ScheduleOption: models.OPTIONS_RUN_NOW_SYNC,
	}); err != nil {
		return err
	}

	// TODO: it might make sense to return the external account ID so the client can use it for initiating a payment
	return nil
}

func (c *Connector) InitiatePayment(ctx task.ConnectorContext, transfer *models.TransferInitiation) error {
	err := ValidateTransferInitiation(transfer)
	if err != nil {
		return err
	}

	descriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
		Name:       "Initiate payment",
		Key:        taskNameInitiatePayment,
		TransferID: transfer.ID.String(),
	})
	if err != nil {
		return err
	}

	scheduleOption := models.OPTIONS_RUN_NOW_SYNC
	scheduledAt := transfer.ScheduledAt
	if !scheduledAt.IsZero() {
		scheduleOption = models.OPTIONS_RUN_SCHEDULED_AT
	}

	if err := ctx.Scheduler().Schedule(ctx.Context(), descriptor, models.TaskSchedulerOptions{
		ScheduleOption: scheduleOption,
		ScheduleAt:     scheduledAt,
		RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
	}); err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
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
