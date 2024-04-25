package triggers

import (
	"context"
	"database/sql"
	"time"

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
	Filter    *FilterEvaluationResult             `json:"filter"`
	Variables map[string]VariableEvaluationResult `json:"variables"`
}

type TriggerManager struct {
	db                  *bun.DB
	expressionEvaluator *expressionEvaluator
}

func (m *TriggerManager) ListTriggers(ctx context.Context, paramsQuery ListTriggersQuery) (*bunpaginate.Cursor[Trigger], error) {
	q := m.db.NewSelect()

	return bunpaginate.UsingOffset[ListTriggerParams, Trigger](ctx, q, bunpaginate.OffsetPaginatedQuery[ListTriggerParams](paramsQuery),
		func(query *bun.SelectQuery) *bun.SelectQuery {

			if paramsQuery.Options.Name != "" {
				query = query.Where("Name ILIKE '%?%';", paramsQuery.Options.Name)
			}
			return query.Where("deleted_at is null")
		})
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
		ret.Filter.Match, err = m.expressionEvaluator.evalFilter(event, *filter)
		if err != nil {
			ret.Filter.Error = err.Error()
		}
	}
	if (ret.Filter == nil || ret.Filter.Match) && len(trigger.Vars) > 0 {
		ret.Variables = map[string]VariableEvaluationResult{}
		for key, expr := range trigger.Vars {
			v := VariableEvaluationResult{}
			v.Value, err = m.expressionEvaluator.evalVariable(event, expr)
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

func (m *TriggerManager) ListTriggersOccurrences(ctx context.Context, q ListTriggersOccurrencesQuery) (*bunpaginate.Cursor[Occurrence], error) {
	query := m.db.NewSelect()

	return bunpaginate.UsingOffset[ListTriggersOccurrencesOptions, Occurrence](ctx, query, bunpaginate.OffsetPaginatedQuery[ListTriggersOccurrencesOptions](q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			query = query.Relation("WorkflowInstance")

			if q.Options.TriggerID != "" {
				query = query.Where("trigger_id = ?", q.Options.TriggerID)
			}

			return query
		})
}

func NewManager(db *bun.DB, expressionEvaluator *expressionEvaluator) *TriggerManager {
	return &TriggerManager{
		db:                  db,
		expressionEvaluator: expressionEvaluator,
	}
}

type ListTriggerParams struct {
	Name string
}
type ListTriggersQuery bunpaginate.OffsetPaginatedQuery[ListTriggerParams]

type ListTriggersOccurrencesOptions struct {
	TriggerID string
}

type ListTriggersOccurrencesQuery bunpaginate.OffsetPaginatedQuery[ListTriggersOccurrencesOptions]
