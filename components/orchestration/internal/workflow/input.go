package workflow

import (
	"github.com/uptrace/bun"
	"go.temporal.io/sdk/workflow"
)

type Input struct {
	Workflow  Workflow          `json:"workflow"`
	Instance  Instance          `json:"instance"`
	Variables map[string]string `json:"variables"`
}

func (i Input) run(ctx workflow.Context, db *bun.DB) error {
	return i.Workflow.Config.run(ctx, db, i.Instance, i.Variables)
}
