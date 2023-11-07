package storage

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/pkg/errors"

	"github.com/formancehq/payments/internal/models"
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

func (s *Storage) ListTasksByStatus(ctx context.Context, connectorID models.ConnectorID, status models.TaskStatus) ([]models.Task, error) {
	var tasks []models.Task

	err := s.db.NewSelect().Model(&tasks).
		Where("connector_id = ?", connectorID).
		Where("status = ?", status).
		Scan(ctx)
	if err != nil {
		return nil, e("failed to get tasks", err)
	}

	return tasks, nil
}

func (s *Storage) ListTasks(ctx context.Context, connectorID models.ConnectorID, pagination PaginatorQuery) ([]models.Task, PaginationDetails, error) {
	var tasks []models.Task

	query := s.db.NewSelect().Model(&tasks).
		Where("connector_id = ?", connectorID)

	query = pagination.apply(query, "task.created_at")

	err := query.Scan(ctx)
	if err != nil {
		return nil, PaginationDetails{}, e("failed to get tasks", err)
	}

	var (
		hasMore                       = len(tasks) > pagination.pageSize
		hasPrevious                   bool
		firstReference, lastReference string
	)

	if hasMore {
		if pagination.cursor.Next || pagination.cursor.Reference == "" {
			tasks = tasks[:pagination.pageSize]
		} else {
			tasks = tasks[1:]
		}
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
	})

	if len(tasks) > 0 {
		firstReference = tasks[0].CreatedAt.Format(time.RFC3339Nano)
		lastReference = tasks[len(tasks)-1].CreatedAt.Format(time.RFC3339Nano)

		query = s.db.NewSelect().Model(&tasks).
			Where("connector_id = ?", connectorID)

		hasPrevious, err = pagination.hasPrevious(ctx, query, "task.created_at", firstReference)
		if err != nil {
			return nil, PaginationDetails{}, e("failed to check if there is a previous page", err)
		}
	}

	paginationDetails, err := pagination.paginationDetails(hasMore, hasPrevious, firstReference, lastReference)
	if err != nil {
		return nil, PaginationDetails{}, e("failed to get pagination details", err)
	}

	return tasks, paginationDetails, nil
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
