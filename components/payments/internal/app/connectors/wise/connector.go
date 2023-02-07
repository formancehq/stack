package wise

import (
	"context"

	"github.com/formancehq/payments/internal/app/models"

	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/payments/internal/app/integration"
	"github.com/formancehq/payments/internal/app/task"
)

const Name = models.ConnectorProviderWise

type Connector struct {
	logger logging.Logger
	cfg    Config
}

func (c *Connector) InitiateTransfer(ctx task.ConnectorContext, transfer models.Transfer) error {
	descriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
		Name: "Initiate transfer",
		Key:  taskNameTransfer,
		Transfer: Transfer{
			ID:          transfer.ID,
			Source:      transfer.Source,
			Destination: transfer.Destination,
			Amount:      transfer.Amount,
			Currency:    transfer.Currency,
		},
	})
	if err != nil {
		return err
	}

	return ctx.Scheduler().Schedule(descriptor, true)
}

func (c *Connector) Install(ctx task.ConnectorContext) error {
	descriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
		Name: "Fetch profiles from client",
		Key:  taskNameFetchProfiles,
	})
	if err != nil {
		return err
	}

	return ctx.Scheduler().Schedule(descriptor, true)
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
