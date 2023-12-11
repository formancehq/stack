package atlar

import (
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

const (
	taskNameFetchAccounts             = "fetch_accounts"
	taskNameFetchPayments             = "fetch_payments"
	taskNameCreateExternalBankAccount = "create_external_bank_account"
	taskNameInitiatePayment           = "initiate_payment"
)

// TaskDescriptor is the definition of a task.
type TaskDescriptor struct {
	Name        string              `json:"name" yaml:"name" bson:"name"`
	Key         string              `json:"key" yaml:"key" bson:"key"`
	Main        bool                `json:"main,omitempty" yaml:"main" bson:"main"`
	Account     string              `json:"account,omitempty" yaml:"account" bson:"account"`
	BankAccount *models.BankAccount `json:"bankAccount,omitempty" yaml:"bankAccount" bson:"bankAccount"`
	TransferID  string              `json:"transferId,omitempty" yaml:"transferId" bson:"transferId"`
}

func resolveTasks(logger logging.Logger, config Config) func(taskDefinition TaskDescriptor) task.Task {
	client := createAtlarClient(&config)

	return func(taskDescriptor TaskDescriptor) task.Task {
		if taskDescriptor.Main {
			return MainTask(logger)
		}

		switch taskDescriptor.Key {
		case taskNameFetchAccounts:
			return FetchAccountsTask(config, client)
		case taskNameFetchPayments:
			return FetchPaymentsTask(config, client, taskDescriptor.Account)
		case taskNameCreateExternalBankAccount:
			return CreateExternalBankAccountTask(config, client, taskDescriptor.BankAccount)
		case taskNameInitiatePayment:
			return InitiatePaymentTask(config, client, taskDescriptor.TransferID)
		default:
			return nil
		}
	}
}
