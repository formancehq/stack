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

func (a Activities) SendWorkflowStartedEvent(ctx context.Context, instance Instance) error {
	return a.publisher.Publish(events.TopicOrchestration,
		events.NewMessage(ctx, events.StartedWorkflow, events.StartedWorkflowPayload{
			ID:         instance.WorkflowID,
			InstanceID: instance.ID,
		}))
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

func (a Activities) SendWorkflowStageStartedEvent(ctx context.Context, instance Instance, stage Stage) error {
	return a.publisher.Publish(events.TopicOrchestration,
		events.NewMessage(ctx, events.StartedWorkflowStage, events.StartedWorkflowStagePayload{
			ID:         instance.WorkflowID,
			InstanceID: instance.ID,
			Number:     stage.Number,
		}))
}

func (a Activities) SendWorkflowStageTerminationEvent(ctx context.Context, instance Instance, stage Stage) error {
	if stage.Error == nil {
		return a.publisher.Publish(events.TopicOrchestration,
			events.NewMessage(ctx, events.SucceededWorkflowStage, events.SucceededWorkflowStagePayload{
				ID:         instance.WorkflowID,
				InstanceID: instance.ID,
				Number:     stage.Number,
			}))
	} else {
		return a.publisher.Publish(events.TopicOrchestration,
			events.NewMessage(ctx, events.FailedWorkflowStage, events.FailedWorkflowStagePayload{
				ID:         instance.WorkflowID,
				InstanceID: instance.ID,
				Number:     stage.Number,
				Error:      *stage.Error,
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
var SendWorkflowStartedEventActivity = Activities{}.SendWorkflowStartedEvent
var SendWorkflowStageStartedEventActivity = Activities{}.SendWorkflowStageStartedEvent
var SendWorkflowStageTerminationEventActivity = Activities{}.SendWorkflowStageTerminationEvent
var InsertNewInstanceActivity = Activities{}.InsertNewInstance
var UpdateInstanceActivity = Activities{}.UpdateInstance
var InsertNewStageActivity = Activities{}.InsertNewStage
var UpdateStageActivity = Activities{}.UpdateStage

func NewActivities(publisher message.Publisher, db *bun.DB) Activities {
	return Activities{
		publisher: publisher,
		db:        db,
	}
}
