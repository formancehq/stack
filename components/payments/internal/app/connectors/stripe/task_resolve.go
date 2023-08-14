package stripe

import (
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
)

const (
	taskNameFetchAccounts            = "fetch_accounts"
	taskNameFetchPaymentsForAccounts = "fetch_transactions"
	taskNameFetchPayments            = "fetch_payments"
	taskNameFetchBalances            = "fetch_balance"
	taskNameFetchExternalAccounts    = "fetch_external_accounts"
)

// TaskDescriptor is the definition of a task.
type TaskDescriptor struct {
	Name       string    `json:"name" yaml:"name" bson:"name"`
	Key        string    `json:"key" yaml:"key" bson:"key"`
	Main       bool      `json:"main,omitempty" yaml:"main" bson:"main"`
	Account    string    `json:"account,omitempty" yaml:"account" bson:"account"`
	TransferID uuid.UUID `json:"transferID,omitempty" yaml:"transferID" bson:"transferID"`
}

// clientID, apiKey, endpoint string, logger logging
func resolveTasks(logger logging.Logger, config Config) func(taskDefinition TaskDescriptor) task.Task {
	client := NewDefaultClient(config.APIKey)

	return func(taskDescriptor TaskDescriptor) task.Task {
		if taskDescriptor.Main {
			return MainTask(logger)
		}

		if taskDescriptor.TransferID != uuid.Nil {
			return TransferTask(config, taskDescriptor.TransferID)
		}

		switch taskDescriptor.Key {
		case taskNameFetchPayments:
			return FetchPaymentsTask(config, client)
		case taskNameFetchAccounts:
			return FetchAccountsTask(config, client)
		case taskNameFetchExternalAccounts:
			return FetchExternalAccountsTask(config, taskDescriptor.Account, client)
		case taskNameFetchPaymentsForAccounts:
			return ConnectedAccountTask(config, taskDescriptor.Account, client)
		case taskNameFetchBalances:
			return BalancesTask(config, taskDescriptor.Account, client)
		default:
			// For compatibility with old tasks
			return ConnectedAccountTask(config, taskDescriptor.Account, client)
		}
	}
}
