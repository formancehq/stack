package triggers

import (
	"fmt"
	"time"

	"github.com/formancehq/stack/libs/go-libs/publish"
	temporalworkflow "go.temporal.io/sdk/workflow"
)

type ProcessEventRequest struct {
	MessageID string               `json:"messageID"`
	Event     publish.EventMessage `json:"ledger"`
}

type triggerWorkflow struct {
	taskQueue string
}

func (w triggerWorkflow) RunTrigger(ctx temporalworkflow.Context, req ProcessEventRequest) error {

	fmt.Println("will list triggers")
	fmt.Println("will list triggers")
	fmt.Println("will list triggers")
	fmt.Println("will list triggers")
	fmt.Println("will list triggers")

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
		occurrence := &Occurrence{}
		err := temporalworkflow.ExecuteActivity(
			temporalworkflow.WithActivityOptions(ctx, temporalworkflow.ActivityOptions{
				StartToCloseTimeout: 10 * time.Second,
			}),
			ProcessEventActivity,
			trigger,
			req,
		).Get(ctx, occurrence)
		if err != nil {
			return err
		}

		err = temporalworkflow.ExecuteActivity(
			temporalworkflow.WithActivityOptions(ctx, temporalworkflow.ActivityOptions{
				StartToCloseTimeout: 10 * time.Second,
			}),
			UpdateTriggerOccurrence,
			occurrence,
		).Get(ctx, occurrence)
		if err != nil {
			return err
		}

		err = temporalworkflow.ExecuteActivity(
			temporalworkflow.WithActivityOptions(ctx, temporalworkflow.ActivityOptions{
				StartToCloseTimeout: 10 * time.Second,
			}),
			SendEventForTriggerTermination,
			occurrence,
		).Get(ctx, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewWorkflow(taskQueue string) *triggerWorkflow {
	return &triggerWorkflow{
		taskQueue: taskQueue,
	}
}

var RunTrigger = triggerWorkflow{}.RunTrigger
