package stages

import (
	"go.temporal.io/sdk/workflow"
)

type NoOp struct {
}

func (n NoOp) GetWorkflow() any {
	return RunNoOp
}

func RunNoOp(ctx workflow.Context, p NoOp) error {
	return nil
}

func init() {
	Register("noop", NoOp{})
}
