package triggers

import (
	"context"
	"database/sql"
	"time"

	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

var ErrWorkflowNotExists = errors.New("workflow does not exists")

type TriggerManager struct {
	db *bun.DB
}

func (m *TriggerManager) ListTriggers(ctx context.Context) ([]Trigger, error) {
	ret := make([]Trigger, 0)
	err := m.db.NewSelect().
		Model(&ret).
		Where("deleted_at is null").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (m *TriggerManager) GetTrigger(ctx context.Context, triggerID string) (*Trigger, error) {
	ret := &Trigger{}
	err := m.db.NewSelect().
		Model(ret).
		Where("deleted_at is null").
		Where("id = ?", triggerID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (m *TriggerManager) DeleteTrigger(ctx context.Context, triggerID string) error {
	updated, err := m.db.NewUpdate().
		Model(&Trigger{}).
		Where("deleted_at is null").
		Where("id = ?", triggerID).
		Set("deleted_at = ?", time.Now()).
		Exec(ctx)
	rowsAffected, err := updated.RowsAffected()
	if err != nil {
		panic(err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (m *TriggerManager) CreateTrigger(ctx context.Context, data TriggerData) (*Trigger, error) {

	if err := data.Validate(); err != nil {
		return nil, errors.Wrap(err, "validating data")
	}

	exists, err := m.db.NewSelect().
		Model(&workflow.Workflow{}).
		Where("deleted_at is null").
		Where("id = ?", data.WorkflowID).
		Exists(ctx)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrWorkflowNotExists
	}

	trigger, err := NewTrigger(data)
	if err != nil {
		return nil, err
	}

	_, err = m.db.NewInsert().
		Model(trigger).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	return trigger, nil
}

func (m *TriggerManager) ListTriggersOccurrences(ctx context.Context, triggerID string) ([]Occurrence, error) {
	ret := make([]Occurrence, 0)
	err := m.db.NewSelect().
		Model(&ret).
		Where("trigger_id = ?", triggerID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func NewManager(db *bun.DB) *TriggerManager {
	return &TriggerManager{
		db: db,
	}
}
