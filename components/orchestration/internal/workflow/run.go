package workflow

import (
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/uptrace/bun"
	"go.temporal.io/sdk/workflow"
)

type Workflows struct {
	db *bun.DB
}

func (r Workflows) Run(ctx workflow.Context, input Input) error {
	err := input.run(ctx, r.db)
	if err != nil {
		logging.Errorf("error running workflow: %s", err)
	}
	return err
}

var Run = Workflows{}.Run

func NewWorkflows(db *bun.DB) *Workflows {
	return &Workflows{
		db: db,
	}
}
