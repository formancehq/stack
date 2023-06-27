package mangopay

import (
	"fmt"

	"github.com/formancehq/payments/internal/app/connectors/mangopay/client"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

const (
	taskNameFetchUsers        = "fetch-users"
	taskNameFetchTransactions = "fetch-transactions"
)

// TaskDescriptor is the definition of a task.
type TaskDescriptor struct {
	Name   string `json:"name" yaml:"name" bson:"name"`
	Key    string `json:"key" yaml:"key" bson:"key"`
	UserID string `json:"userID" yaml:"userID" bson:"userID"`
}

// clientID, apiKey, endpoint string, logger logging
func resolveTasks(logger logging.Logger, config Config) func(taskDefinition models.TaskDescriptor) task.Task {
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

	return func(taskDefinition models.TaskDescriptor) task.Task {
		taskDescriptor, err := models.DecodeTaskDescriptor[TaskDescriptor](taskDefinition)
		if err != nil {
			logger.Error(err)

			return nil
		}

		switch taskDescriptor.Key {
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
