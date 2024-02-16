package storage

import (
	"context"
	"encoding/json"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func (s *Storage) UpdateTaskStatus(ctx context.Context, connectorID models.ConnectorID, descriptor models.TaskDescriptor, status models.TaskStatus, taskError string) error {
	_, err := s.db.NewUpdate().Model(&models.Task{}).
		Set("status = ?", status).
		Set("error = ?", taskError).
		Where("descriptor::TEXT = ?::TEXT", descriptor.ToMessage()).
		Where("connector_id = ?", connectorID).
		Exec(ctx)
	if err != nil {
		return e("failed to update task", err)
	}

	return nil
}

func (s *Storage) UpdateTaskState(ctx context.Context, connectorID models.ConnectorID, descriptor models.TaskDescriptor, state json.RawMessage) error {
	_, err := s.db.NewUpdate().Model(&models.Task{}).
		Set("state = ?", state).
		Where("descriptor::TEXT = ?::TEXT", descriptor.ToMessage()).
		Where("connector_id = ?", connectorID).
		Exec(ctx)
	if err != nil {
		return e("failed to update task", err)
	}

	return nil
}

func (s *Storage) FindAndUpsertTask(
	ctx context.Context,
	connectorID models.ConnectorID,
	descriptor models.TaskDescriptor,
	status models.TaskStatus,
	schedulerOptions models.TaskSchedulerOptions,
	taskErr string,
) (*models.Task, error) {
	_, err := s.GetTaskByDescriptor(ctx, connectorID, descriptor)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, e("failed to get task", err)
	}

	if err == nil {
		err = s.UpdateTaskStatus(ctx, connectorID, descriptor, status, taskErr)
		if err != nil {
			return nil, e("failed to update task", err)
		}
	} else {
		err = s.CreateTask(ctx, connectorID, descriptor, status, schedulerOptions)
		if err != nil {
			return nil, e("failed to upsert task", err)
		}
	}

	return s.GetTaskByDescriptor(ctx, connectorID, descriptor)
}

func (s *Storage) CreateTask(ctx context.Context, connectorID models.ConnectorID, descriptor models.TaskDescriptor, status models.TaskStatus, schedulerOptions models.TaskSchedulerOptions) error {
	_, err := s.db.NewInsert().Model(&models.Task{
		ConnectorID:      connectorID,
		Descriptor:       descriptor.ToMessage(),
		Status:           status,
		SchedulerOptions: schedulerOptions,
	}).Exec(ctx)
	if err != nil {
		return e("failed to create task", err)
	}

	return nil
}

func (s *Storage) ListTasksByStatus(ctx context.Context, connectorID models.ConnectorID, status models.TaskStatus) ([]*models.Task, error) {
	var tasks []*models.Task

	err := s.db.NewSelect().Model(&tasks).
		Where("connector_id = ?", connectorID).
		Where("status = ?", status).
		Scan(ctx)
	if err != nil {
		return nil, e("failed to get tasks", err)
	}

	return tasks, nil
}

type TaskQuery struct{}

type ListTasksQuery bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[TaskQuery]]

func NewListTasksQuery(opts PaginatedQueryOptions[TaskQuery]) ListTasksQuery {
	return ListTasksQuery{
		PageSize: opts.PageSize,
		Order:    bunpaginate.OrderAsc,
		Options:  opts,
	}
}

func (s *Storage) ListTasks(ctx context.Context, connectorID models.ConnectorID, q ListTasksQuery) (*api.Cursor[models.Task], error) {
	return PaginateWithOffset[PaginatedQueryOptions[TaskQuery], models.Task](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[TaskQuery]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			query = query.
				Where("connector_id = ?", connectorID).
				Order("created_at DESC")

			if q.Options.Sorter != nil {
				query = q.Options.Sorter.Apply(query)
			}

			return query
		},
	)
}

func (s *Storage) ReadOldestPendingTask(ctx context.Context, connectorID models.ConnectorID) (*models.Task, error) {
	var task models.Task
	err := s.db.NewSelect().Model(&task).
		Where("connector_id = ?", connectorID).
		Where("status = ?", models.TaskStatusPending).
		Order("created_at ASC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, e("failed to get task", err)
	}

	return &task, nil
}

func (s *Storage) GetTask(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	var task models.Task

	err := s.db.NewSelect().Model(&task).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, e("failed to get task", err)
	}

	return &task, nil
}

func (s *Storage) GetTaskByDescriptor(ctx context.Context, connectorID models.ConnectorID, descriptor models.TaskDescriptor) (*models.Task, error) {
	var task models.Task
	err := s.db.NewSelect().Model(&task).
		Where("connector_id = ?", connectorID).
		Where("descriptor::TEXT = ?::TEXT", descriptor.ToMessage()).
		Scan(ctx)
	if err != nil {
		return nil, e("failed to get task", err)
	}

	return &task, nil
}
