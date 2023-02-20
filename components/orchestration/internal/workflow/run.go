package workflow

import (
	"github.com/uptrace/bun"
	"go.temporal.io/sdk/workflow"
)

type Workflows struct {
	db *bun.DB
}

func (r Workflows) Run(ctx workflow.Context, input Input) error {
	return input.run(ctx, r.db)
}

var Run = Workflows{}.Run

func NewWorkflows(db *bun.DB) *Workflows {
	return &Workflows{
		db: db,
	}
}
