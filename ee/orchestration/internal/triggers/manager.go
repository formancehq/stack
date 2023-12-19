package triggers

import (
	"context"
	"database/sql"
	"time"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

var ErrWorkflowNotExists = errors.New("workflow does not exists")

type VariableEvaluationResult struct {
	Value string `json:"value,omitempty"`
	Error string `json:"error,omitempty"`
}

type FilterEvaluationResult struct {
	Match bool   `json:"match"`
	Error string `json:"error,omitempty"`
}

type TestTriggerResult struct {
	Filter    *FilterEvaluationResult
	Variables map[string]VariableEvaluationResult `json:"variables"`
}

type TriggerManager struct {
	db *bun.DB
}

func (m *TriggerManager) ListTriggers(ctx context.Context, query ListTriggersQuery) (*sharedapi.Cursor[Trigger], error) {
	ret := make([]Trigger, 0)
	q := m.db.NewSelect().
		Model(&ret).
		Where("deleted_at is null")

	return bunpaginate.UsingOffset[any, Trigger](ctx, q, bunpaginate.OffsetPaginatedQuery[any](query))
}

func (m *TriggerManager) TestTrigger(ctx context.Context, triggerID string, event map[string]any) (*TestTriggerResult, error) {
	trigger := &Trigger{}
	err := m.db.NewSelect().
		Model(trigger).
		Where("deleted_at is null").
		Where("id = ?", triggerID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	ret := TestTriggerResult{}
	if filter := trigger.Filter; filter != nil {
		ret.Filter = &FilterEvaluationResult{}
		ret.Filter.Match, err = evalFilter(event, *filter)
		if err != nil {
			ret.Filter.Error = err.Error()
		}
	}
	if (ret.Filter == nil || ret.Filter.Match) && len(trigger.Vars) > 0 {
		ret.Variables = map[string]VariableEvaluationResult{}
		for key, expr := range trigger.Vars {
			v := VariableEvaluationResult{}
			v.Value, err = evalVariable(event, expr)
			if err != nil {
				v.Error = err.Error()
			}
			ret.Variables[key] = v
		}
	}

	return &ret, nil
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

func (m *TriggerManager) ListTriggersOccurrences(ctx context.Context, query ListTriggersOccurrencesQuery) (*sharedapi.Cursor[Occurrence], error) {
	ret := make([]Occurrence, 0)
	q := m.db.NewSelect().
		Model(&ret)

	if query.Options.TriggerID != "" {
		q = q.Where("trigger_id = ?", query.Options.TriggerID)
	}

	return bunpaginate.UsingOffset[ListTriggersOccurrencesOptions, Occurrence](ctx, q, bunpaginate.OffsetPaginatedQuery[ListTriggersOccurrencesOptions](query))
}

func NewManager(db *bun.DB) *TriggerManager {
	return &TriggerManager{
		db: db,
	}
}

type ListTriggersQuery bunpaginate.OffsetPaginatedQuery[any]

type ListTriggersOccurrencesOptions struct {
	TriggerID string
}

type ListTriggersOccurrencesQuery bunpaginate.OffsetPaginatedQuery[ListTriggersOccurrencesOptions]
