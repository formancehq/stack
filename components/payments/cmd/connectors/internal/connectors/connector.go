package connectors

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
	// UpdateConfig is used to update the configuration of the connector.
	UpdateConfig(ctx task.ConnectorContext, config models.ConnectorConfigObject) error
	// Resolve is used to recover state of a failed or restarted task
	Resolve(descriptor models.TaskDescriptor) task.Task
	// InitiateTransfer is used to initiate a transfer from the connector to a bank account.
	InitiatePayment(ctx task.ConnectorContext, transfer *models.TransferInitiation) error
	// ReverssePayment is used to reverse a transfer from the connector.
	ReversePayment(ctx task.ConnectorContext, transferReversal *models.TransferReversal) error
	// CreateExternalBankAccount is used to create a bank account on the connector side.
	CreateExternalBankAccount(ctx task.ConnectorContext, bankAccount *models.BankAccount) error
	// GetSupportedCurrenciesAndDecimals returns a map of supported currencies and their decimals.
	SupportedCurrenciesAndDecimals() map[string]int
}
