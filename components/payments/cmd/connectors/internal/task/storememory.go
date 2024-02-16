package task

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
)

type InMemoryStore struct {
	tasks    map[uuid.UUID]models.Task
	statuses map[string]models.TaskStatus
	created  map[string]time.Time
	errors   map[string]string
}

func (s *InMemoryStore) GetTask(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	task, ok := s.tasks[id]
	if !ok {
		return nil, storage.ErrNotFound
	}

	return &task, nil
}

func (s *InMemoryStore) GetTaskByDescriptor(
	ctx context.Context,
	connectorID models.ConnectorID,
	descriptor models.TaskDescriptor,
) (*models.Task, error) {
	id, err := descriptor.EncodeToString()
	if err != nil {
		return nil, err
	}

	status, ok := s.statuses[id]
	if !ok {
		return nil, storage.ErrNotFound
	}

	return &models.Task{
		Descriptor: descriptor.ToMessage(),
		Status:     status,
		Error:      s.errors[id],
		State:      nil,
		CreatedAt:  s.created[id],
	}, nil
}

func (s *InMemoryStore) ListTasks(ctx context.Context,
	connectorID models.ConnectorID,
	q storage.ListTasksQuery,
) (*api.Cursor[models.Task], error) {
	ret := make([]models.Task, 0)

	for id, status := range s.statuses {
		if !strings.HasPrefix(id, fmt.Sprintf("%s/", connectorID)) {
			continue
		}

		var descriptor models.TaskDescriptor

		ret = append(ret, models.Task{
			Descriptor: descriptor.ToMessage(),
			Status:     status,
			Error:      s.errors[id],
			State:      nil,
			CreatedAt:  s.created[id],
		})
	}

	return &api.Cursor[models.Task]{
		PageSize: 15,
		HasMore:  false,
		Previous: "",
		Next:     "",
		Data:     ret,
	}, nil
}

func (s *InMemoryStore) ReadOldestPendingTask(
	ctx context.Context,
	connectorID models.ConnectorID,
) (*models.Task, error) {
	var (
		oldestDate time.Time
		oldestID   string
	)

	for id, status := range s.statuses {
		if status != models.TaskStatusPending {
			continue
		}

		if oldestDate.IsZero() || s.created[id].Before(oldestDate) {
			oldestDate = s.created[id]
			oldestID = id
		}
	}

	if oldestDate.IsZero() {
		return nil, storage.ErrNotFound
	}

	descriptorStr := strings.Split(oldestID, "/")[1]

	var descriptor models.TaskDescriptor

	data, err := base64.StdEncoding.DecodeString(descriptorStr)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &descriptor)
	if err != nil {
		return nil, err
	}

	return &models.Task{
		Descriptor: descriptor.ToMessage(),
		Status:     models.TaskStatusPending,
		State:      nil,
		CreatedAt:  s.created[oldestID],
	}, nil
}

func (s *InMemoryStore) ListTasksByStatus(
	ctx context.Context,
	connectorID models.ConnectorID,
	taskStatus models.TaskStatus,
) ([]*models.Task, error) {
	cursor, err := s.ListTasks(ctx, connectorID, storage.NewListTasksQuery(storage.NewPaginatedQueryOptions(storage.TaskQuery{})))
	if err != nil {
		return nil, err
	}

	ret := make([]*models.Task, 0)

	for _, v := range cursor.Data {
		if v.Status != taskStatus {
			continue
		}

		ret = append(ret, &v)
	}

	return ret, nil
}

func (s *InMemoryStore) FindAndUpsertTask(
	ctx context.Context,
	connectorID models.ConnectorID,
	descriptor models.TaskDescriptor,
	status models.TaskStatus,
	options models.TaskSchedulerOptions,
	taskErr string,
) (*models.Task, error) {
	err := s.UpdateTaskStatus(ctx, connectorID, descriptor, status, taskErr)
	if err != nil {
		return nil, err
	}

	return &models.Task{
		Descriptor: descriptor.ToMessage(),
		Status:     status,
		Error:      taskErr,
		State:      nil,
	}, nil
}

func (s *InMemoryStore) UpdateTaskStatus(
	ctx context.Context,
	connectorID models.ConnectorID,
	descriptor models.TaskDescriptor,
	status models.TaskStatus,
	taskError string,
) error {
	taskID, err := descriptor.EncodeToString()
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s", connectorID, taskID)

	s.statuses[key] = status

	s.errors[key] = taskError
	if _, ok := s.created[key]; !ok {
		s.created[key] = time.Now()
	}

	return nil
}

func (s *InMemoryStore) Result(
	connectorID models.ConnectorID,
	descriptor models.TaskDescriptor,
) (models.TaskStatus, string, bool) {
	taskID, err := descriptor.EncodeToString()
	if err != nil {
		panic(err)
	}

	key := fmt.Sprintf("%s/%s", connectorID, taskID)

	status, ok := s.statuses[key]
	if !ok {
		return "", "", false
	}

	return status, s.errors[key], true
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		statuses: make(map[string]models.TaskStatus),
		errors:   make(map[string]string),
		created:  make(map[string]time.Time),
	}
}

var _ Repository = &InMemoryStore{}
