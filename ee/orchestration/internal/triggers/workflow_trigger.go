package triggers

import (
	"time"

	"go.temporal.io/api/enums/v1"

	"github.com/formancehq/orchestration/internal/temporalworker"

	"github.com/formancehq/go-libs/pointer"
	"github.com/formancehq/orchestration/internal/workflow"

	"github.com/formancehq/go-libs/publish"
	temporalworkflow "go.temporal.io/sdk/workflow"
)

const SearchAttributeTriggerID = "OrchestrationTriggerID"

type ProcessEventRequest struct {
	Event publish.EventMessage `json:"ledger"`
}

type triggerWorkflow struct {
	taskQueue               string
	includeSearchAttributes bool
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
		searchAttributes := map[string]any{}
		if w.includeSearchAttributes {
			searchAttributes[SearchAttributeTriggerID] = trigger.ID
		}

		if err := temporalworkflow.ExecuteChildWorkflow(
			temporalworkflow.WithChildOptions(ctx, temporalworkflow.ChildWorkflowOptions{
				TaskQueue:         w.taskQueue,
				ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
				SearchAttributes:  searchAttributes,
			}),
			ExecuteTrigger,
			req,
			trigger).GetChildWorkflowExecution().Get(ctx, nil); err != nil {
			return err
		}
	}

	return nil
}

func (w triggerWorkflow) ExecuteTrigger(ctx temporalworkflow.Context, req ProcessEventRequest, trigger Trigger) error {

	vars := make(map[string]string)
	occurrence := NewTriggerOccurrence(trigger.ID, req.Event, temporalworkflow.Now(ctx))
	err := temporalworkflow.ExecuteActivity(
		temporalworkflow.WithActivityOptions(ctx, temporalworkflow.ActivityOptions{
			StartToCloseTimeout: 10 * time.Second,
		}),
		EvalTriggerVariables,
		trigger,
		req,
	).Get(ctx, &vars)
	if err != nil {
		occurrence.Error = pointer.For(err.Error())
	} else {
		instance := &workflow.Instance{}
		if err := temporalworkflow.ExecuteChildWorkflow(
			temporalworkflow.WithChildOptions(ctx, temporalworkflow.ChildWorkflowOptions{
				TaskQueue: w.taskQueue,
			}),
			workflow.Initiate,
			workflow.Input{
				Workflow:  *trigger.Workflow,
				Variables: vars,
			}).Get(ctx, instance); err != nil {
			return err
		}

		occurrence.WorkflowInstanceID = pointer.For(instance.ID)
	}

	err = temporalworkflow.ExecuteActivity(
		temporalworkflow.WithActivityOptions(ctx, temporalworkflow.ActivityOptions{
			StartToCloseTimeout: 10 * time.Second,
		}),
		InsertTriggerOccurrence,
		occurrence,
	).Get(ctx, nil)
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

	return nil
}

func (w triggerWorkflow) DefinitionSet() temporalworker.DefinitionSet {
	return temporalworker.NewDefinitionSet().
		Append(temporalworker.Definition{
			Func: w.RunTrigger,
			Name: "RunTrigger",
		}).
		Append(temporalworker.Definition{
			Func: w.ExecuteTrigger,
			Name: "ExecuteTrigger",
		})
}

func NewWorkflow(taskQueue string, includeSearchAttributes bool) *triggerWorkflow {
	return &triggerWorkflow{
		taskQueue:               taskQueue,
		includeSearchAttributes: includeSearchAttributes,
	}
}

var ExecuteTrigger = triggerWorkflow{}.ExecuteTrigger
var RunTrigger = triggerWorkflow{}.RunTrigger
