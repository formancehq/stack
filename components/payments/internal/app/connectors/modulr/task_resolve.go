package modulr

import (
	"fmt"

	"github.com/formancehq/payments/internal/app/connectors/modulr/client"
	"github.com/formancehq/payments/internal/app/task"

	"github.com/formancehq/stack/libs/go-libs/logging"
)

const (
	taskNameMain                = "main"
	taskNameFetchTransactions   = "fetch-transactions"
	taskNameFetchAccounts       = "fetch-accounts"
	taskNameFetchBeneficiaries  = "fetch-beneficiaries"
	taskNameInitiatePayment     = "initiate-payment"
	taskNameUpdatePaymentStatus = "update-payment-status"
)

// TaskDescriptor is the definition of a task.
type TaskDescriptor struct {
	Name       string `json:"name" yaml:"name" bson:"name"`
	Key        string `json:"key" yaml:"key" bson:"key"`
	AccountID  string `json:"accountID" yaml:"accountID" bson:"accountID"`
	TransferID string `json:"transferID" yaml:"transferID" bson:"transferID"`
	Attempt    int    `json:"attempt" yaml:"attempt" bson:"attempt"`
}

func resolveTasks(logger logging.Logger, config Config) func(taskDefinition TaskDescriptor) task.Task {
	modulrClient, err := client.NewClient(config.APIKey, config.APISecret, config.Endpoint)
	if err != nil {
		return func(taskDefinition TaskDescriptor) task.Task {
			return func() error {
				return fmt.Errorf("key '%s': %w", taskDefinition.Key, ErrMissingTask)
			}
		}
	}

	return func(taskDefinition TaskDescriptor) task.Task {
		switch taskDefinition.Key {
		case taskNameMain:
			return taskMain(logger)
		case taskNameFetchAccounts:
			return taskFetchAccounts(logger, modulrClient)
		case taskNameFetchBeneficiaries:
			return taskFetchBeneficiaries(logger, modulrClient)
		case taskNameInitiatePayment:
			return taskInitiatePayment(logger, modulrClient, taskDefinition.TransferID)
		case taskNameUpdatePaymentStatus:
			return taskUpdatePaymentStatus(logger, modulrClient, taskDefinition.TransferID, taskDefinition.Attempt)
		case taskNameFetchTransactions:
			return taskFetchTransactions(logger, modulrClient, taskDefinition.AccountID)
		}

		// This should never happen.
		return func() error {
			return fmt.Errorf("key '%s': %w", taskDefinition.Key, ErrMissingTask)
		}
	}
}
