package triggers

import (
	"context"
	"strings"

	"github.com/formancehq/orchestration/internal/temporalworker"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/orchestration/pkg/events"
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
	publisher           message.Publisher
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
		Relation("Workflow").
		Where("trigger.deleted_at is null").
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

func (a Activities) EvalTriggerVariables(ctx context.Context, trigger Trigger, request ProcessEventRequest) (map[string]string, error) {
	return a.expressionEvaluator.evalVariables(request.Event.Payload, trigger.Vars)
}

func (a Activities) InsertTriggerOccurrence(ctx context.Context, occurrence Occurrence) error {
	_, err := a.db.NewInsert().
		Model(pointer.For(occurrence)).
		Exec(ctx)
	return err
}

func (a Activities) SendEventForTriggerTermination(ctx context.Context, occurrence Occurrence) error {
	if occurrence.Error == nil || *occurrence.Error == "" {
		return a.publisher.Publish(events.TopicOrchestration,
			events.NewMessage(ctx, events.SucceededTrigger, events.SucceededTriggerPayload{
				ID:                 occurrence.ID,
				WorkflowInstanceID: *occurrence.WorkflowInstanceID,
				TriggerID:          occurrence.TriggerID,
			}))
	} else {
		return a.publisher.Publish(events.TopicOrchestration,
			events.NewMessage(ctx, events.FailedTrigger, events.FailedTriggerPayload{
				ID:        occurrence.ID,
				TriggerID: occurrence.TriggerID,
				Error:     *occurrence.Error,
			}))
	}
}

func (a Activities) DefinitionSet() temporalworker.DefinitionSet {
	return temporalworker.NewDefinitionSet().
		Append(temporalworker.Definition{
			Func: a.EvalTriggerVariables,
			Name: "EvalTriggerVariables",
		}).
		Append(temporalworker.Definition{
			Func: a.InsertTriggerOccurrence,
			Name: "InsertTriggerOccurrence",
		}).
		Append(temporalworker.Definition{
			Func: a.ListTriggers,
			Name: "ListTriggers",
		}).
		Append(temporalworker.Definition{
			Func: a.SendEventForTriggerTermination,
			Name: "SendEventForTriggerTermination",
		})
}

func NewActivities(db *bun.DB, manager *workflow.WorkflowManager,
	expressionEvaluator *expressionEvaluator, publisher message.Publisher) Activities {
	return Activities{
		db:                  db,
		manager:             manager,
		expressionEvaluator: expressionEvaluator,
		publisher:           publisher,
	}
}

var EvalTriggerVariables = Activities{}.EvalTriggerVariables
var SendEventForTriggerTermination = Activities{}.SendEventForTriggerTermination
var ListTriggersActivity = Activities{}.ListTriggers
var InsertTriggerOccurrence = Activities{}.InsertTriggerOccurrence
