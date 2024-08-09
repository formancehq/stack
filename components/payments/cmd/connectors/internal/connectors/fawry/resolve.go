package fawry

import (
	"fmt"
	"sync"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/fawry/client"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

const (
	TaskMain          = "main"
	TaskIngestPayment = "ingest-payment"
)

// internal state not pushed in the database
type TaskMemoryState struct {
	// We want to fetch the transactions once per service start.
	fetchTransactionsOnce map[string]*sync.Once
}

type TaskDescriptor struct {
	Name    string `json:"name" yaml:"name" bson:"name"`
	Key     string `json:"key" yaml:"key" bson:"key"`
	Payload string `json:"payload" yaml:"payload" bson:"payload"`
}

func Resolve(
	logger logging.Logger,
	config Config,
	taskMemoryState *TaskMemoryState,
) func(taskDefinition TaskDescriptor) task.Task {
	client := client.NewClient()
	err := client.Init()
	if err != nil {
		logger.Error(err)
		return func(taskDescriptor TaskDescriptor) task.Task {
			return func() error {
				return fmt.Errorf("cannot build fawry client: %w", err)
			}
		}
	}

	return func(taskDescriptor TaskDescriptor) task.Task {
		switch taskDescriptor.Key {
		case TaskMain:
			return func() error {
				return nil
			}
		case TaskIngestPayment:
			return taskIngestPayment
		default:
			return func() error {
				return fmt.Errorf("unknown task key: %s", taskDescriptor.Key)
			}
		}
	}
}
