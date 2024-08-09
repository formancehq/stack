package fawry

import (
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

type Connector struct {
	logger logging.Logger
	cfg    Config

	taskMemoryState *TaskMemoryState
}

func (c *Connector) Install(ctx task.ConnectorContext) error {
	desc, err := models.EncodeTaskDescriptor(TaskDescriptor{
		Name: TaskMain,
	})
	if err != nil {
		return err
	}

	err = ctx.Scheduler().Schedule(ctx.Context(), desc, models.TaskSchedulerOptions{
		ScheduleOption: models.OPTIONS_RUN_NOW,
		RestartOption:  models.OPTIONS_RESTART_ALWAYS,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Connector) Uninstall(ctx task.ConnectorContext) error {
	return nil
}

func (c *Connector) Resolve(descriptor models.TaskDescriptor) task.Task {
	taskDescriptor, err := models.DecodeTaskDescriptor[TaskDescriptor](descriptor)
	if err != nil {
		panic(err)
	}

	return Resolve(c.logger, c.cfg, c.taskMemoryState)(taskDescriptor)
}

func (c *Connector) UpdateConfig(ctx task.ConnectorContext, config models.ConnectorConfigObject) error {
	panic("not supported")
}

func (c *Connector) InitiatePayment(ctx task.ConnectorContext, transfer *models.TransferInitiation) error {
	panic("not supported")
}

// ReversePayment
func (c *Connector) ReversePayment(ctx task.ConnectorContext, reversal *models.TransferReversal) error {
	panic("not supported")
}

// CreateExternalBankAccount
func (c *Connector) CreateExternalBankAccount(ctx task.ConnectorContext, account *models.BankAccount) error {
	panic("not supported")
}

// SupportedCurrenciesAndDecimals
func (c *Connector) SupportedCurrenciesAndDecimals(ctx task.ConnectorContext) map[string]int {
	return supportedCurrenciesWithDecimal
}

func newConnector() *Connector {
	return &Connector{}
}
