package task

import (
	"context"

	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	UpdateTaskStatus(ctx context.Context, provider models.ConnectorProvider, descriptor models.TaskDescriptor, status models.TaskStatus, err string) error
	FindAndUpsertTask(ctx context.Context, provider models.ConnectorProvider, descriptor models.TaskDescriptor, status models.TaskStatus, schedulerOptions models.TaskSchedulerOptions, err string) (*models.Task, error)
	ListTasksByStatus(ctx context.Context, provider models.ConnectorProvider, status models.TaskStatus) ([]models.Task, error)
	ListTasks(ctx context.Context, provider models.ConnectorProvider, pagination storage.PaginatorQuery) ([]models.Task, storage.PaginationDetails, error)
	ReadOldestPendingTask(ctx context.Context, provider models.ConnectorProvider) (*models.Task, error)
	GetTask(ctx context.Context, id uuid.UUID) (*models.Task, error)
	GetTaskByDescriptor(ctx context.Context, provider models.ConnectorProvider, descriptor models.TaskDescriptor) (*models.Task, error)
}
