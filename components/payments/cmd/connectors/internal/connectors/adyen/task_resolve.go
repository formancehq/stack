package adyen

import (
	"fmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/adyen/client"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
)

const (
	taskNameMain          = "main"
	taskNameFetchAccounts = "fetch-accounts"
	taskNameHandleWebhook = "handle-webhook"
)

type TaskDescriptor struct {
	Name          string    `json:"name" yaml:"name" bson:"name"`
	Key           string    `json:"key" yaml:"key" bson:"key"`
	PollingPeriod int       `json:"pollingPeriod" yaml:"pollingPeriod" bson:"pollingPeriod"`
	WebhookID     uuid.UUID `json:"webhookId" yaml:"webhookId" bson:"webhookId"`
}

func resolveTasks(logger logging.Logger, config Config) func(taskDefinition TaskDescriptor) task.Task {
	adyenClient, err := client.NewClient(
		config.APIKey,
		config.HMACKey,
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

	return func(taskDescriptor TaskDescriptor) task.Task {
		switch taskDescriptor.Key {
		case taskNameMain:
			return taskMain()
		case taskNameFetchAccounts:
			return taskFetchAccounts(adyenClient)
		case taskNameHandleWebhook:
			return taskHandleStandardWebhooks(adyenClient, taskDescriptor.WebhookID)
		}

		// This should never happen.
		return func() error {
			return fmt.Errorf("key '%s': %w", taskDescriptor.Key, ErrMissingTask)
		}
	}
}
