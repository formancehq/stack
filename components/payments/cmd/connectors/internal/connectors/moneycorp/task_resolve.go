package moneycorp

import (
	"fmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/moneycorp/client"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

const (
	taskNameMain                = "main"
	taskNameFetchAccounts       = "fetch-accounts"
	taskNameFetchRecipients     = "fetch-recipients"
	taskNameFetchTransactions   = "fetch-transactions"
	taskNameFetchBalances       = "fetch-balances"
	taskNameInitiatePayment     = "initiate-payment"
	taskNameUpdatePaymentStatus = "update-payment-status"
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

// clientID, apiKey, endpoint string, logger logging
func resolveTasks(logger logging.Logger, config Config) func(taskDefinition TaskDescriptor) task.Task {
	moneycorpClient, err := client.NewClient(
		config.ClientID,
		config.APIKey,
		config.Endpoint,
		logger,
	)
	if err != nil {
		logger.Error(err)

		return func(taskDescriptor TaskDescriptor) task.Task {
			return func() error {
				return fmt.Errorf("cannot build moneycorp client: %w", err)
			}
		}
	}

	return func(taskDescriptor TaskDescriptor) task.Task {
		switch taskDescriptor.Key {
		case taskNameMain:
			return taskMain()
		case taskNameFetchAccounts:
			return taskFetchAccounts(moneycorpClient)
		case taskNameFetchRecipients:
			return taskFetchRecipients(moneycorpClient, taskDescriptor.AccountID)
		case taskNameFetchTransactions:
			return taskFetchTransactions(moneycorpClient, taskDescriptor.AccountID)
		case taskNameFetchBalances:
			return taskFetchBalances(moneycorpClient, taskDescriptor.AccountID)
		case taskNameInitiatePayment:
			return taskInitiatePayment(moneycorpClient, taskDescriptor.TransferID)
		case taskNameUpdatePaymentStatus:
			return taskUpdatePaymentStatus(moneycorpClient, taskDescriptor.TransferID, taskDescriptor.PaymentID, taskDescriptor.Attempt)
		}

		// This should never happen.
		return func() error {
			return fmt.Errorf("key '%s': %w", taskDescriptor.Key, ErrMissingTask)
		}
	}
}
