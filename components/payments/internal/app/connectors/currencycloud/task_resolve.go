package currencycloud

import (
	"fmt"

	"github.com/formancehq/payments/internal/app/connectors/currencycloud/client"

	"github.com/formancehq/payments/internal/app/task"

	"github.com/formancehq/stack/libs/go-libs/logging"
)

const (
	taskNameMain              = "main"
	taskNameFetchTransactions = "fetch-transactions"
)

// TaskDescriptor is the definition of a task.
type TaskDescriptor struct {
	Name string `json:"name" yaml:"name" bson:"name"`
	Key  string `json:"key" yaml:"key" bson:"key"`
}

func resolveTasks(logger logging.Logger, config Config) func(taskDefinition TaskDescriptor) task.Task {
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
			return taskMain(logger)
		case taskNameFetchTransactions:
			return taskFetchTransactions(logger, currencyCloudClient, config)
		}

		// This should never happen.
		return func() error {
			return fmt.Errorf("key '%s': %w", taskDescriptor.Name, ErrMissingTask)
		}
	}
}
