package integration

import (
	"context"

	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
)

// Connector provide entry point to a payment provider.
type Connector interface {
	// Install is used to start the connector. The implementation if in charge of scheduling all required resources.
	Install(ctx task.ConnectorContext) error
	// Uninstall is used to uninstall the connector. It has to close all related resources opened by the connector.
	Uninstall(ctx context.Context) error
	// Resolve is used to recover state of a failed or restarted task
	Resolve(descriptor models.TaskDescriptor) task.Task
	// InitiateTransfer is used to initiate a transfer from the connector to a bank account.
	InitiatePayment(ctx task.ConnectorContext, transfer *models.TransferInitiation) error
}
