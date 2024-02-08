package modulr

import (
	"context"
	"errors"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

// taskMain is the main task of the connector. It launches the other tasks.
func taskMain(config Config) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		scheduler task.Scheduler,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"modulr.taskMain",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
		)
		defer span.End()

		taskAccounts, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name: "Fetch accounts from client",
			Key:  taskNameFetchAccounts,
		})
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		err = scheduler.Schedule(ctx, taskAccounts, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_PERIODICALLY,
			Duration:       config.PollingPeriod.Duration,
			RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			otel.RecordError(span, err)
			return err
		}

		taskBeneficiaries, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name: "Fetch beneficiaries from client",
			Key:  taskNameFetchBeneficiaries,
		})
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		err = scheduler.Schedule(ctx, taskBeneficiaries, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_PERIODICALLY,
			Duration:       config.PollingPeriod.Duration,
			RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}
