package currencycloud

import (
	"fmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currencycloud/client"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
)

const (
	taskNameMain                = "main"
	taskNameFetchTransactions   = "fetch-transactions"
	taskNameFetchAccounts       = "fetch-accounts"
	taskNameFetchBeneficiaries  = "fetch-beneficiaries"
	taskNameFetchBalances       = "fetch-balances"
	taskNameInitiatePayment     = "initiate-payment"
	taskNameUpdatePaymentStatus = "update-payment-status"
)

// TaskDescriptor is the definition of a task.
type TaskDescriptor struct {
	Name       string `json:"name" yaml:"name" bson:"name"`
	Key        string `json:"key" yaml:"key" bson:"key"`
	TransferID string `json:"transferID" yaml:"transferID" bson:"transferID"`
	PaymentID  string `json:"paymentID" yaml:"paymentID" bson:"paymentID"`
	Attempt    int    `json:"attempt" yaml:"attempt" bson:"attempt"`
}

func resolveTasks(config Config) func(taskDefinition TaskDescriptor) task.Task {
	currencyCloudClient, err := client.NewClient(config.LoginID, config.APIKey, config.Endpoint)
	if err != nil {
		return func(taskDefinition TaskDescriptor) task.Task {
			return func() error {
				return fmt.Errorf("failed to initiate client: %w", err)
			}
		}
	}

	return func(taskDescriptor TaskDescriptor) task.Task {
		if taskDescriptor.Key == "" {
			// Keep the compatibility with previous version if the connector.
			// If the key is empty, use the name as the key.
			taskDescriptor.Key = taskDescriptor.Name
		}

		switch taskDescriptor.Key {
		case taskNameMain:
			return taskMain()
		case taskNameFetchAccounts:
			return taskFetchAccounts(currencyCloudClient)
		case taskNameFetchBeneficiaries:
			return taskFetchBeneficiaries(currencyCloudClient)
		case taskNameFetchTransactions:
			return taskFetchTransactions(currencyCloudClient, config)
		case taskNameFetchBalances:
			return taskFetchBalances(currencyCloudClient)
		case taskNameInitiatePayment:
			return taskInitiatePayment(currencyCloudClient, taskDescriptor.TransferID)
		case taskNameUpdatePaymentStatus:
			return taskUpdatePaymentStatus(currencyCloudClient, taskDescriptor.TransferID, taskDescriptor.PaymentID, taskDescriptor.Attempt)
		}

		// This should never happen.
		return func() error {
			return fmt.Errorf("key '%s': %w", taskDescriptor.Name, ErrMissingTask)
		}
	}
}
