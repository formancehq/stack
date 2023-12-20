package triggers

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/uptrace/bun"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Activities struct {
	db                  *bun.DB
	manager             *workflow.WorkflowManager
	expressionEvaluator *expressionEvaluator
}

func (a Activities) processTrigger(ctx context.Context, request ProcessEventRequest, trigger Trigger) bool {
	ctx, span := workflow.Tracer.Start(ctx, "Triggers:CheckRequirements", trace.WithAttributes(
		attribute.String("trigger-id", trigger.ID),
	))
	defer span.End()

	if trigger.Filter != nil && *trigger.Filter != "" {

		ok, err := a.expressionEvaluator.evalFilter(request.Event.Payload, *trigger.Filter)
		if err != nil {
			span.SetAttributes(
				attribute.String("filter-error", err.Error()),
			)
		}
		span.SetAttributes(
			attribute.String("filter", *trigger.Filter),
			attribute.Bool("match", ok),
		)

		if !ok {
			return false
		}
	}

	return true
}

func (a Activities) ListTriggers(ctx context.Context, request ProcessEventRequest) ([]Trigger, error) {
	ret := make([]Trigger, 0)

	triggers := make([]Trigger, 0)
	if err := a.db.NewSelect().
		Model(&triggers).
		Where("deleted_at is null").
		Where("event = ?", request.Event.Type).
		Scan(ctx); err != nil {
		return nil, err
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("found-triggers", strings.Join(collectionutils.Map(triggers, Trigger.GetID), ", ")))

	for _, trigger := range triggers {
		if a.processTrigger(trace.ContextWithSpan(ctx, span), request, trigger) {
			ret = append(ret, trigger)
		}
	}

	return ret, nil
}

func (a Activities) ProcessTrigger(ctx context.Context, trigger Trigger, request ProcessEventRequest) error {

	span := trace.SpanFromContext(ctx)
	var (
		evaluated map[string]string
		err       error
	)
	if trigger.Vars != nil {
		evaluated, err = a.expressionEvaluator.evalVariables(request.Event.Payload, trigger.Vars)
		if err != nil {
			span.RecordError(err)
			return err
		}
	}

	data, err := json.Marshal(evaluated)
	if err != nil {
		panic(err)
	}

	span.SetAttributes(attribute.String("variables", string(data)))

	instance, err := a.manager.RunWorkflow(ctx, trigger.WorkflowID, evaluated)
	if err != nil {
		return err
	}

	_, err = a.db.NewInsert().
		Model(pointer.For(NewTriggerOccurrence(request.MessageID, trigger.ID, instance.ID, request.Event))).
		On("CONFLICT (trigger_id, event_id) DO NOTHING").
		Exec(ctx)

	return err
}

func NewActivities(db *bun.DB, manager *workflow.WorkflowManager, expressionEvaluator *expressionEvaluator) Activities {
	return Activities{
		db:                  db,
		manager:             manager,
		expressionEvaluator: expressionEvaluator,
	}
}

var ProcessEventActivity = Activities{}.ProcessTrigger
var ListTriggersActivity = Activities{}.ListTriggers
