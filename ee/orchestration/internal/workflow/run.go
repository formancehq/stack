package workflow

import (
	"time"

	"github.com/formancehq/orchestration/internal/temporalworker"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

const SearchAttributeWorkflowID = "OrchestrationWorkflowID"

type Workflows struct {
	includeSearchAttributes bool
}

func (w Workflows) Initiate(ctx workflow.Context, input Input) (*Instance, error) {
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

	searchAttributes := map[string]any{}
	if w.includeSearchAttributes {
		searchAttributes = map[string]interface{}{
			SearchAttributeWorkflowID: input.Workflow.ID,
		}
	}

	if err := workflow.ExecuteChildWorkflow(
		workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
			ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
			WorkflowID:        workflow.GetInfo(ctx).WorkflowExecution.ID + "-main",
			SearchAttributes:  searchAttributes,
		}),
		Run,
		input,
		*instance,
	).GetChildWorkflowExecution().Get(ctx, nil); err != nil {
		return nil, err
	}

	return instance, nil
}

func (w Workflows) Run(ctx workflow.Context, i Input, instance Instance) error {
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

func (w Workflows) DefinitionSet() temporalworker.DefinitionSet {
	return temporalworker.NewDefinitionSet().
		Append(temporalworker.Definition{
			Func: w.Run,
			Name: "Run",
		}).Append(temporalworker.Definition{
		Func: w.Initiate,
		Name: "Initiate",
	})
}

var Initiate = Workflows{}.Initiate
var Run = Workflows{}.Run

func NewWorkflows(includeSearchAttributes bool) *Workflows {
	return &Workflows{
		includeSearchAttributes: includeSearchAttributes,
	}
}
