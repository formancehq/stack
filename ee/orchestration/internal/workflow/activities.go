package workflow

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/orchestration/pkg/events"
	"github.com/uptrace/bun"
	"go.temporal.io/sdk/activity"
)

type Activities struct {
	publisher message.Publisher
	db        *bun.DB
}

func (a Activities) SendWorkflowTerminationEvent(ctx context.Context, instance Instance) error {
	if instance.Error == "" {
		return a.publisher.Publish(events.TopicOrchestration,
			events.NewMessage(ctx, events.SucceededWorkflow, events.SucceededWorkflowPayload{
				ID:         instance.WorkflowID,
				InstanceID: instance.ID,
			}))
	} else {
		return a.publisher.Publish(events.TopicOrchestration,
			events.NewMessage(ctx, events.FailedWorkflow, events.FailedWorkflowPayload{
				ID:         instance.WorkflowID,
				InstanceID: instance.ID,
				Error:      instance.Error,
			}))
	}
}

func (a Activities) InsertNewInstance(ctx context.Context, workflowID string) (*Instance, error) {
	instance := NewInstance(activity.GetInfo(ctx).WorkflowExecution.ID, workflowID)
	if _, err := a.db.
		NewInsert().
		Model(&instance).
		Exec(ctx); err != nil {
		return nil, err
	}

	return &instance, nil
}

func (a Activities) UpdateInstance(ctx context.Context, instance *Instance) error {
	_, dbErr := a.db.NewUpdate().
		Model(instance).
		WherePK().
		Exec(ctx)
	return dbErr
}

func (a Activities) InsertNewStage(ctx context.Context, instance Instance, ind int) (*Stage, error) {
	stage := NewStage(instance.ID, activity.GetInfo(ctx).WorkflowExecution.RunID, ind)
	if _, err := a.db.NewInsert().
		Model(&stage).
		Exec(ctx); err != nil {
		return nil, err
	}

	return &stage, nil
}

func (a Activities) UpdateStage(ctx context.Context, stage Stage) error {
	_, err := a.db.NewUpdate().
		Model(&stage).
		WherePK().
		Exec(ctx)
	return err
}

var SendWorkflowTerminationEventActivity = Activities{}.SendWorkflowTerminationEvent
var InsertNewInstance = Activities{}.InsertNewInstance
var UpdateInstance = Activities{}.UpdateInstance
var InsertNewStage = Activities{}.InsertNewStage
var UpdateStage = Activities{}.UpdateStage

func NewActivities(publisher message.Publisher, db *bun.DB) Activities {
	return Activities{
		publisher: publisher,
		db:        db,
	}
}
