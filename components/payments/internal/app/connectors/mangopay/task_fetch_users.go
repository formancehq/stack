package mangopay

import (
	"context"
	"errors"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/mangopay/client"
	"github.com/formancehq/payments/internal/app/metrics"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	usersAttrs = append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "users"))
)

func taskFetchUsers(logger logging.Logger, client *client.Client) task.Task {
	return func(
		ctx context.Context,
		scheduler task.Scheduler,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		logger.Info(taskNameFetchUsers)

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), metric.WithAttributes(usersAttrs...))
		}()

		users, err := client.GetAllUsers(ctx)
		if err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, metric.WithAttributes(usersAttrs...))
			return err
		}

		for _, user := range users {
			logger.Infof("scheduling fetch-transactions: %s", user.ID)

			walletsTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
				Name:   "Fetch wallets from client by user",
				Key:    taskNameFetchWallets,
				UserID: user.ID,
			})
			if err != nil {
				return err
			}

			err = scheduler.Schedule(ctx, walletsTask, models.TaskSchedulerOptions{
				ScheduleOption: models.OPTIONS_RUN_NOW,
				Restart:        true,
			})
			if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
				return err
			}

			bankAccountsTask, err := models.EncodeTaskDescriptor(TaskDescriptor{
				Name:   "Fetch bank accounts from client by user",
				Key:    taskNameFetchBankAccounts,
				UserID: user.ID,
			})
			if err != nil {
				return err
			}

			err = scheduler.Schedule(ctx, bankAccountsTask, models.TaskSchedulerOptions{
				ScheduleOption: models.OPTIONS_RUN_NOW,
				Restart:        true,
			})
			if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
				return err
			}
		}
		metricsRegistry.ConnectorObjects().Add(ctx, int64(len(users)), metric.WithAttributes(usersAttrs...))

		return nil
	}
}
