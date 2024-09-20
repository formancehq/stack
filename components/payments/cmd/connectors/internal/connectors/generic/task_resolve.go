package generic

import (
	"fmt"

	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/generic/client"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
)

const (
	taskNameMain               = "main"
	taskNameFetchAccounts      = "fetch-accounts"
	taskNameFetchBalances      = "fetch-balances"
	taskNameFetchBeneficiaries = "fetch-beneficiaries"
	taskNameFetchTransactions  = "fetch-transactions"
)

// TaskDescriptor is the definition of a task.
type TaskDescriptor struct {
	Name      string `json:"name" yaml:"name" bson:"name"`
	Key       string `json:"key" yaml:"key" bson:"key"`
	AccountID string `json:"account_id" yaml:"account_id" bson:"account_id"`
}

func resolveTasks(logger logging.Logger, config Config) func(taskDefinition TaskDescriptor) task.Task {
	genericClient := client.NewClient(
		config.APIKey,
		config.Endpoint,
		logger,
	)

	return func(taskDescriptor TaskDescriptor) task.Task {
		switch taskDescriptor.Key {
		case taskNameMain:
			return taskMain()
		case taskNameFetchAccounts:
			return taskFetchAccounts(genericClient, &config)
		case taskNameFetchBeneficiaries:
			return taskFetchBeneficiaries(genericClient, &config)
		case taskNameFetchBalances:
			return taskFetchBalances(genericClient, &config, taskDescriptor.AccountID)
		case taskNameFetchTransactions:
			return taskFetchTransactions(genericClient, &config)
		}

		// This should never happen.
		return func() error {
			return fmt.Errorf("key '%s': %w", taskDescriptor.Key, ErrMissingTask)
		}
	}
}
