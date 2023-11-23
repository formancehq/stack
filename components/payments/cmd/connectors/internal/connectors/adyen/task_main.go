package adyen

import (
	"context"

	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskMain() task.Task {
	return func(
		ctx context.Context,
		logger logging.Logger,
		scheduler task.Scheduler,
	) error {
		logger.Info(taskNameMain)

		return nil
	}
}
