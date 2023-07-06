package wise

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/formancehq/payments/internal/app/connectors/wise/client"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

const (
	taskNameMain           = "main"
	taskNameFetchTransfers = "fetch-transfers"
	taskNameFetchProfiles  = "fetch-profiles"
	taskNameTransfer       = "transfer"
)

// TaskDescriptor is the definition of a task.
type TaskDescriptor struct {
	Name      string   `json:"name" yaml:"name" bson:"name"`
	Key       string   `json:"key" yaml:"key" bson:"key"`
	ProfileID uint64   `json:"profileID" yaml:"profileID" bson:"profileID"`
	Transfer  Transfer `json:"transfer" yaml:"transfer" bson:"transfer"`
}

type Transfer struct {
	ID          uuid.UUID `json:"id" yaml:"id" bson:"id"`
	Source      string    `json:"source" yaml:"source" bson:"source"`
	Destination string    `json:"destination" yaml:"destination" bson:"destination"`
	Amount      int64     `json:"amount" yaml:"amount" bson:"amount"`
	Currency    string    `json:"currency" yaml:"currency" bson:"currency"`
}

func resolveTasks(logger logging.Logger, config Config) func(taskDefinition TaskDescriptor) task.Task {
	client := client.NewClient(config.APIKey)

	return func(taskDefinition TaskDescriptor) task.Task {
		switch taskDefinition.Key {
		case taskNameMain:
			return taskMain(logger)
		case taskNameFetchProfiles:
			return taskFetchProfiles(logger, client)
		case taskNameFetchTransfers:
			return taskFetchTransfers(logger, client, taskDefinition.ProfileID)
		case taskNameTransfer:
			return taskTransfer(logger, client, taskDefinition.Transfer)
		}

		// This should never happen.
		return func() error {
			return fmt.Errorf("key '%s': %w", taskDefinition.Key, ErrMissingTask)
		}
	}
}
