package workflow

import (
	"go.temporal.io/sdk/workflow"
	"time"
)

type Input struct {
	Workflow  Workflow          `json:"workflow"`
	Variables map[string]string `json:"variables"`
}

func (i Input) run(ctx workflow.Context) (*Instance, error) {

	instance := &Instance{}
	err := workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}), InsertNewInstance, i.Workflow.ID).Get(ctx, instance)
	if err != nil {
		return nil, err
	}

	err = i.Workflow.Config.run(ctx, *instance, i.Variables)
	if err != nil {
		instance.SetTerminatedWithError(workflow.Now(ctx), err)
	} else {
		instance.SetTerminated(workflow.Now(ctx))
	}

	err = workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}), UpdateInstance, instance).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflow.ExecuteActivity(workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}), SendWorkflowTerminationEventActivity, instance).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	return instance, nil
}
