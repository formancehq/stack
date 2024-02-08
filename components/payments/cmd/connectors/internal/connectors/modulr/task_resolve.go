package modulr

import (
	"fmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/modulr/client"
	"github.com/formancehq/payments/cmd/connectors/internal/task"

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

const (
	timeTemplate = "2006-01-02T15:04:05-0700"
)

// TaskDescriptor is the definition of a task.
type TaskDescriptor struct {
	Name       string `json:"name" yaml:"name" bson:"name"`
	Key        string `json:"key" yaml:"key" bson:"key"`
	AccountID  string `json:"accountID" yaml:"accountID" bson:"accountID"`
	TransferID string `json:"transferID" yaml:"transferID" bson:"transferID"`
	PaymentID  string `json:"paymentID" yaml:"paymentID" bson:"paymentID"`
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
			return taskMain(config)
		case taskNameFetchAccounts:
			return taskFetchAccounts(config, modulrClient)
		case taskNameFetchBeneficiaries:
			return taskFetchBeneficiaries(config, modulrClient)
		case taskNameInitiatePayment:
			return taskInitiatePayment(modulrClient, taskDefinition.TransferID)
		case taskNameUpdatePaymentStatus:
			return taskUpdatePaymentStatus(modulrClient, taskDefinition.TransferID, taskDefinition.PaymentID, taskDefinition.Attempt)
		case taskNameFetchTransactions:
			return taskFetchTransactions(config, modulrClient, taskDefinition.AccountID)
		}

		// This should never happen.
		return func() error {
			return fmt.Errorf("key '%s': %w", taskDefinition.Key, ErrMissingTask)
		}
	}
}
