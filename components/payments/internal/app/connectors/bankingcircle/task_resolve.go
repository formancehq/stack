package bankingcircle

import (
	"fmt"

	"github.com/formancehq/payments/internal/app/connectors/bankingcircle/client"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

const (
	taskNameMain          = "main"
	taskNameFetchPayments = "fetch-payments"
)

// TaskDescriptor is the definition of a task.
type TaskDescriptor struct {
	Name string `json:"name" yaml:"name" bson:"name"`
	Key  string `json:"key" yaml:"key" bson:"key"`
}

func resolveTasks(logger logging.Logger, config Config) func(taskDefinition TaskDescriptor) task.Task {
	bankingCircleClient, err := client.NewClient(
		config.Username,
		config.Password,
		config.Endpoint,
		config.AuthorizationEndpoint,
		config.UserCertificate,
		config.UserCertificateKey,
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
		case taskNameFetchPayments:
			return taskFetchPayments(logger, bankingCircleClient)
		}

		// This should never happen.
		return func() error {
			return fmt.Errorf("key '%s': %w", taskDescriptor.Key, ErrMissingTask)
		}
	}
}
