package mangopay

import (
	"fmt"

	"github.com/formancehq/payments/internal/app/connectors"
	"github.com/formancehq/payments/internal/app/connectors/mangopay/client"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

const (
	taskNameMain              = "main"
	taskNameFetchUsers        = "fetch-users"
	taskNameFetchTransactions = "fetch-transactions"
)

// TaskDescriptor is the definition of a task.
type TaskDescriptor struct {
	Name          string              `json:"name" yaml:"name" bson:"name"`
	Key           string              `json:"key" yaml:"key" bson:"key"`
	UserID        string              `json:"userID" yaml:"userID" bson:"userID"`
	PollingPeriod connectors.Duration `json:"pollingPeriod" yaml:"pollingPeriod" bson:"pollingPeriod"`
}

// clientID, apiKey, endpoint string, logger logging
func resolveTasks(logger logging.Logger, config Config) func(taskDefinition TaskDescriptor) task.Task {
	mangopayClient, err := client.NewClient(
		config.ClientID,
		config.APIKey,
		config.Endpoint,
		logger,
	)
	if err != nil {
		logger.Error(err)

		return nil
	}

	return func(taskDescriptor TaskDescriptor) task.Task {
		switch taskDescriptor.Key {
		case taskNameMain:
			return taskMain(logger)
		case taskNameFetchUsers:
			return taskFetchUsers(logger, mangopayClient)
		case taskNameFetchTransactions:
			return taskFetchTransactions(logger, mangopayClient, taskDescriptor.UserID)
		}

		// This should never happen.
		return func() error {
			return fmt.Errorf("key '%s': %w", taskDescriptor.Key, ErrMissingTask)
		}
	}
}
