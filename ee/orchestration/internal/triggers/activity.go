package triggers

import (
	"context"

	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/uptrace/bun"
)

type Activities struct {
	db      *bun.DB
	manager *workflow.WorkflowManager
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

	for _, trigger := range triggers {
		ok := true
		var err error
		if trigger.Filter != nil && *trigger.Filter != "" {
			ok, err = evalFilter(request.Event.Payload, *trigger.Filter)
			if err != nil {
				logging.WithFields(map[string]any{
					"filter": *trigger.Filter,
				}).Errorf("unable to eval filter: %s", err)
			}
			continue
		}

		logging.FromContext(ctx).
			WithField("trigger-id", trigger.ID).
			Debugf("Checking expr '%s': %v", trigger.Filter, ok)

		if ok {
			ret = append(ret, trigger)
		}
	}

	return ret, nil
}

func (a Activities) ProcessTrigger(ctx context.Context, trigger Trigger, request ProcessEventRequest) error {
	var (
		evaluated map[string]string
		err       error
	)
	if trigger.Vars != nil {
		evaluated, err = evalVariables(request.Event.Payload, trigger.Vars)
		if err != nil {
			return err
		}
	}

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

func NewActivities(db *bun.DB, manager *workflow.WorkflowManager) Activities {
	return Activities{
		db:      db,
		manager: manager,
	}
}

var ProcessEventActivity = Activities{}.ProcessTrigger
var ListTriggersActivity = Activities{}.ListTriggers
