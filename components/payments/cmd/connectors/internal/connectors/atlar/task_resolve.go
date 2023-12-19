package atlar

import (
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/atlar/client"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

const (
	taskNameFetchAccounts             = "fetch_accounts"
	taskNameFetchTransactions         = "fetch_transactions"
	taskNameCreateExternalBankAccount = "create_external_bank_account"
	taskNameInitiatePayment           = "initiate_payment"
	taskNameUpdatePaymentStatus       = "update_payment_status"
)

// TaskDescriptor is the definition of a task.
type TaskDescriptor struct {
	Name        string              `json:"name" yaml:"name" bson:"name"`
	Key         string              `json:"key" yaml:"key" bson:"key"`
	Main        bool                `json:"main,omitempty" yaml:"main" bson:"main"`
	BankAccount *models.BankAccount `json:"bankAccount,omitempty" yaml:"bankAccount" bson:"bankAccount"`
	TransferID  string              `json:"transferId,omitempty" yaml:"transferId" bson:"transferId"`
	PaymentID   string              `json:"paymentId,omitempty" yaml:"paymentId" bson:"paymentId"`
	Attempt     int                 `json:"attempt,omitempty" yaml:"attempt" bson:"attempt"`
}

func resolveTasks(logger logging.Logger, config Config) func(taskDefinition TaskDescriptor) task.Task {
	client := client.NewClient(config.BaseUrl, config.AccessKey, config.Secret)

	return func(taskDescriptor TaskDescriptor) task.Task {
		if taskDescriptor.Main {
			return MainTask(logger)
		}

		switch taskDescriptor.Key {
		case taskNameFetchAccounts:
			return FetchAccountsTask(config, client)
		case taskNameFetchTransactions:
			return FetchTransactionsTask(config, client)
		case taskNameCreateExternalBankAccount:
			return CreateExternalBankAccountTask(config, client, taskDescriptor.BankAccount)
		case taskNameInitiatePayment:
			return InitiatePaymentTask(config, client, taskDescriptor.TransferID)
		case taskNameUpdatePaymentStatus:
			return UpdatePaymentStatusTask(config, client, taskDescriptor.TransferID, taskDescriptor.PaymentID, taskDescriptor.Attempt)
		default:
			return nil
		}
	}
}
