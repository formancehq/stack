package storage

import (
	"context"
	"encoding/json"

	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type tasks struct {
	bun.BaseModel `bun:"table:tasks"`

	// Mandatory fields
	ConnectorID models.ConnectorID `bun:"connector_id,pk,type:character varying,notnull"`
	Tasks       json.RawMessage    `bun:"tasks,type:json,notnull"`
}

func (s *store) TasksUpsert(ctx context.Context, connectorID models.ConnectorID, ts models.Tasks) error {
	payload, err := json.Marshal(&ts)
	if err != nil {
		return errors.Wrap(err, "failed to marshal tasks")
	}

	tasks := tasks{
		ConnectorID: connectorID,
		Tasks:       payload,
	}

	_, err = s.db.NewInsert().
		Model(&tasks).
		On("CONFLICT (connector_id) DO UPDATE").
		Set("tasks = EXCLUDED.tasks").
		Exec(ctx)
	return e("failed to insert tasks", err)
}

func (s *store) TasksGet(ctx context.Context, connectorID models.ConnectorID) (*models.Tasks, error) {
	var ts tasks

	err := s.db.NewSelect().
		Model(&ts).
		Where("connector_id = ?", connectorID).
		Scan(ctx)
	if err != nil {
		return nil, e("failed to fetch tasks", err)
	}

	var tasks models.Tasks
	if err := json.Unmarshal(ts.Tasks, &tasks); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal tasks")
	}

	return &tasks, nil
}

func (s *store) TasksDeleteFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error {
	_, err := s.db.NewDelete().
		Model((*tasks)(nil)).
		Where("connector_id = ?", connectorID).
		Exec(ctx)

	return e("failed to delete tasks", err)
}
