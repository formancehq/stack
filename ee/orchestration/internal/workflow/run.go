package workflow

import (
	"time"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

const SearchAttributeWorkflowID = "OrchestrationWorkflowID"

type Workflows struct{}

func (r Workflows) Initiate(ctx workflow.Context, input Input) (*Instance, error) {
	instance := &Instance{}
	err := workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}), InsertNewInstanceActivity, input.Workflow.ID).Get(ctx, instance)
	if err != nil {
		return nil, err
	}

	err = workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}), SendWorkflowStartedEventActivity, instance).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	if err := workflow.ExecuteChildWorkflow(
		workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
			ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
			WorkflowID:        workflow.GetInfo(ctx).WorkflowExecution.ID + "-main",
			SearchAttributes: map[string]interface{}{
				SearchAttributeWorkflowID: input.Workflow.ID,
			},
		}),
		Run,
		input,
		*instance,
	).GetChildWorkflowExecution().Get(ctx, nil); err != nil {
		return nil, err
	}

	return instance, nil
}

func (r Workflows) Run(ctx workflow.Context, i Input, instance Instance) error {
	err := i.Workflow.Config.run(ctx, instance, i.Variables)
	if err != nil {
		instance.SetTerminatedWithError(workflow.Now(ctx), err)
	} else {
		instance.SetTerminated(workflow.Now(ctx))
	}

	err = workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}), UpdateInstanceActivity, instance).Get(ctx, nil)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}), SendWorkflowTerminationEventActivity, instance).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

var Initiate = Workflows{}.Initiate
var Run = Workflows{}.Run

func NewWorkflows() *Workflows {
	return &Workflows{}
}
