package moneycorp

import (
	"context"
	"errors"

	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// taskMain is the main task of the connector. It launches the other tasks.
func taskMain() task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		scheduler task.Scheduler,
	) error {
		span := trace.SpanFromContext(ctx)
		span.SetName("moneycorp.taskMain")
		span.SetAttributes(
			attribute.String("connectorID", connectorID.String()),
		)

		taskAccounts, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name: "Fetch accounts from client",
			Key:  taskNameFetchAccounts,
		})
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		err = scheduler.Schedule(ctx, taskAccounts, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_NOW,
			RestartOption:  models.OPTIONS_RESTART_IF_NOT_ACTIVE,
		})
		if err != nil && !errors.Is(err, task.ErrAlreadyScheduled) {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}
