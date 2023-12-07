package atlar

import (
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

const (
	taskNameFetchAccounts = "fetch_accounts"
	taskNameFetchPayments = "fetch_payments"
)

// TaskDescriptor is the definition of a task.
type TaskDescriptor struct {
	Name    string `json:"name" yaml:"name" bson:"name"`
	Key     string `json:"key" yaml:"key" bson:"key"`
	Main    bool   `json:"main,omitempty" yaml:"main" bson:"main"`
	Account string `json:"account,omitempty" yaml:"account" bson:"account"`
}

// clientID, apiKey, endpoint string, logger logging
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
			return FetchPaymentsTask(config, taskDescriptor.Account, client)
		default:
			return nil
		}
	}
}
