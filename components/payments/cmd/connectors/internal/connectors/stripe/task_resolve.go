package stripe

import (
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/stripe/client"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

const (
	taskNameFetchAccounts            = "fetch_accounts"
	taskNameFetchPaymentsForAccounts = "fetch_transactions"
	taskNameFetchPayments            = "fetch_payments"
	taskNameFetchBalances            = "fetch_balance"
	taskNameFetchExternalAccounts    = "fetch_external_accounts"
	taskNameInitiatePayment          = "initiate-payment"
	taskNameUpdatePaymentStatus      = "update-payment-status"
)

// TaskDescriptor is the definition of a task.
type TaskDescriptor struct {
	Name       string `json:"name" yaml:"name" bson:"name"`
	Key        string `json:"key" yaml:"key" bson:"key"`
	Main       bool   `json:"main,omitempty" yaml:"main" bson:"main"`
	Account    string `json:"account,omitempty" yaml:"account" bson:"account"`
	TransferID string `json:"transferID" yaml:"transferID" bson:"transferID"`
	PaymentID  string `json:"paymentID" yaml:"paymentID" bson:"paymentID"`
	Attempt    int    `json:"attempt" yaml:"attempt" bson:"attempt"`
}

// clientID, apiKey, endpoint string, logger logging
func resolveTasks(logger logging.Logger, config Config) func(taskDefinition TaskDescriptor) task.Task {
	client := client.NewDefaultClient(config.APIKey)

	return func(taskDescriptor TaskDescriptor) task.Task {
		if taskDescriptor.Main {
			return MainTask(logger)
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
			return BalancesTask(taskDescriptor.Account, client)
		case taskNameInitiatePayment:
			return InitiatePaymentTask(logger, taskDescriptor.TransferID, client)
		case taskNameUpdatePaymentStatus:
			return UpdatePaymentStatusTask(logger, taskDescriptor.TransferID, taskDescriptor.PaymentID, taskDescriptor.Attempt, client)
		default:
			// For compatibility with old tasks
			return ConnectedAccountTask(config, taskDescriptor.Account, client)
		}
	}
}
