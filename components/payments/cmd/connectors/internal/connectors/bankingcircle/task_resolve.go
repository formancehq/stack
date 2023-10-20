package bankingcircle

import (
	"fmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/bankingcircle/client"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

const (
	taskNameMain                = "main"
	taskNameFetchPayments       = "fetch-payments"
	taskNameFetchAccounts       = "fetch-accounts"
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

func resolveTasks(logger logging.Logger, config Config) func(taskDefinition TaskDescriptor) task.Task {
	bankingCircleClient, err := client.NewClient(
		config.Username,
		config.Password,
		config.Endpoint,
		config.AuthorizationEndpoint,
		config.UserCertificate,
		config.UserCertificateKey,
		logger,
	)
	if err != nil {
		logger.Error(err)

		return func(taskDescriptor TaskDescriptor) task.Task {
			return func() error {
				return fmt.Errorf("cannot build banking circle client: %w", err)
			}
		}
	}

	return func(taskDescriptor TaskDescriptor) task.Task {
		switch taskDescriptor.Key {
		case taskNameMain:
			return taskMain(logger)
		case taskNameFetchPayments:
			return taskFetchPayments(logger, bankingCircleClient)
		case taskNameFetchAccounts:
			return taskFetchAccounts(logger, bankingCircleClient)
		case taskNameInitiatePayment:
			return taskInitiatePayment(logger, bankingCircleClient, taskDescriptor.TransferID)
		case taskNameUpdatePaymentStatus:
			return taskUpdatePaymentStatus(logger, bankingCircleClient, taskDescriptor.TransferID, taskDescriptor.PaymentID, taskDescriptor.Attempt)
		}

		// This should never happen.
		return func() error {
			return fmt.Errorf("key '%s': %w", taskDescriptor.Key, ErrMissingTask)
		}
	}
}
