package triggers

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/orchestration/pkg/events"
	sharedlogging "github.com/formancehq/stack/libs/go-libs/logging"

	"go.temporal.io/sdk/temporal"

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

func (a Activities) ProcessTrigger(ctx context.Context, trigger Trigger, request ProcessEventRequest) (*Occurrence, error) {

	span := trace.SpanFromContext(ctx)
	var (
		evaluated    map[string]string
		triggerError error
		occurrence   = NewTriggerOccurrence(request.MessageID, trigger.ID, request.Event)
	)
	if trigger.Vars != nil {
		evaluated, triggerError = a.expressionEvaluator.evalVariables(request.Event.Payload, trigger.Vars)
	}
	if triggerError == nil {
		data, triggerError := json.Marshal(evaluated)
		if triggerError != nil {
			panic(triggerError)
		}

		span.SetAttributes(attribute.String("variables", string(data)))

		instance, triggerError := a.manager.RunWorkflow(ctx, trigger.WorkflowID, evaluated)
		if triggerError != nil {
			return nil, triggerError
		}

		occurrence.WorkflowInstanceID = pointer.For(instance.ID)
	} else {
		triggerError = temporal.NewNonRetryableApplicationError("unable to eval variables", "VARIABLES_EVAL", triggerError)
		span.RecordError(triggerError)
		occurrence.Error = pointer.For(triggerError.Error())
	}

	_, err := a.db.NewInsert().
		Model(pointer.For(occurrence)).
		On("CONFLICT (trigger_id, event_id) DO NOTHING").
		Exec(ctx)
	if err != nil {
		sharedlogging.FromContext(ctx).Errorf("unable to save trigger occurrence: %s", err)
	}

	return &occurrence, nil
}

func (a Activities) SendEventForTriggerTermination(ctx context.Context, occurrence Occurrence) error {
	if occurrence.Error == nil || *occurrence.Error == "" {
		return a.publisher.Publish(events.SucceededTrigger,
			events.NewMessage(ctx, events.SucceededTrigger, events.SucceededTriggerPayload{
				ID:      occurrence.TriggerID,
				EventID: occurrence.EventID,
			}))
	} else {
		return a.publisher.Publish(events.FailedTrigger,
			events.NewMessage(ctx, events.FailedTrigger, events.FailedTriggerPayload{
				ID:      occurrence.TriggerID,
				Error:   *occurrence.Error,
				EventID: occurrence.EventID,
			}))
	}
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

var ProcessEventActivity = Activities{}.ProcessTrigger
var SendEventForTriggerTermination = Activities{}.SendEventForTriggerTermination
var ListTriggersActivity = Activities{}.ListTriggers
