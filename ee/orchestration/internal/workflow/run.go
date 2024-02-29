package workflow

import (
	"go.temporal.io/sdk/workflow"
)

type Workflows struct{}

func (r Workflows) Run(ctx workflow.Context, input Input) (*Instance, error) {
	return input.run(ctx)
}

var Run = Workflows{}.Run

func NewWorkflows() *Workflows {
	return &Workflows{}
}
