package triggers

import (
	"time"

	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/uptrace/bun"
	temporalworkflow "go.temporal.io/sdk/workflow"
)

type ProcessEventRequest struct {
	MessageID string               `json:"messageID"`
	Event     publish.EventMessage `json:"ledger"`
}

type triggerWorkflow struct {
	taskQueue string
	db        *bun.DB
}

func (w triggerWorkflow) RunTrigger(ctx temporalworkflow.Context, req ProcessEventRequest) error {

	triggers := make([]Trigger, 0)
	err := temporalworkflow.ExecuteActivity(
		temporalworkflow.WithActivityOptions(ctx, temporalworkflow.ActivityOptions{
			StartToCloseTimeout: 10 * time.Second,
		}),
		ListTriggersActivity,
		req,
	).Get(ctx, &triggers)
	if err != nil {
		return err
	}

	for _, trigger := range triggers {
		err := temporalworkflow.ExecuteActivity(
			temporalworkflow.WithActivityOptions(ctx, temporalworkflow.ActivityOptions{
				StartToCloseTimeout: 10 * time.Second,
			}),
			ProcessEventActivity,
			trigger,
			req,
		).Get(ctx, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewWorkflow(db *bun.DB, taskQueue string) *triggerWorkflow {
	return &triggerWorkflow{
		db:        db,
		taskQueue: taskQueue,
	}
}

var RunTrigger = triggerWorkflow{}.RunTrigger
