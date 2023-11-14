package task

import (
	"context"

	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	UpdateTaskStatus(ctx context.Context, connectorID models.ConnectorID, descriptor models.TaskDescriptor, status models.TaskStatus, err string) error
	FindAndUpsertTask(ctx context.Context, connectorID models.ConnectorID, descriptor models.TaskDescriptor, status models.TaskStatus, schedulerOptions models.TaskSchedulerOptions, err string) (*models.Task, error)
	ListTasksByStatus(ctx context.Context, connectorID models.ConnectorID, status models.TaskStatus) ([]models.Task, error)
	ListTasks(ctx context.Context, connectorID models.ConnectorID, pagination storage.PaginatorQuery) ([]models.Task, storage.PaginationDetails, error)
	ReadOldestPendingTask(ctx context.Context, connectorID models.ConnectorID) (*models.Task, error)
	GetTask(ctx context.Context, taskID uuid.UUID) (*models.Task, error)
	GetTaskByDescriptor(ctx context.Context, connectorID models.ConnectorID, descriptor models.TaskDescriptor) (*models.Task, error)
}
