package adyen

import (
	"fmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/adyen/client"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

const (
	taskNameMain = "main"
)

type TaskDescriptor struct {
	Name          string `json:"name" yaml:"name" bson:"name"`
	Key           string `json:"key" yaml:"key" bson:"key"`
	PollingPeriod int    `json:"pollingPeriod" yaml:"pollingPeriod" bson:"pollingPeriod"`
}

func resolveTasks(logger logging.Logger, config Config) func(taskDefinition TaskDescriptor) task.Task {
	adyenClient, err := client.NewClient(
		config.APIKey,
		config.LiveEndpointPrefix,
		logger,
	)
	if err != nil {
		logger.Error(err)

		return func(taskDescriptor TaskDescriptor) task.Task {
			return func() error {
				return fmt.Errorf("cannot build adyen client: %w", err)
			}
		}
	}

	_ = adyenClient

	return func(taskDescriptor TaskDescriptor) task.Task {
		switch taskDescriptor.Key {
		case taskNameMain:
			return taskMain()
		}

		// This should never happen.
		return func() error {
			return fmt.Errorf("key '%s': %w", taskDescriptor.Key, ErrMissingTask)
		}
	}
}
