package workflow

import (
	"context"

	"github.com/uptrace/bun"
	"go.temporal.io/sdk/workflow"
)

type Input struct {
	Workflow  Workflow          `json:"workflow"`
	Instance  Instance          `json:"instance"`
	Variables map[string]string `json:"variables"`
}

func (i Input) run(ctx workflow.Context, db *bun.DB) error {
	instance := i.Instance
	err := i.Workflow.Config.run(ctx, db, instance, i.Variables)
	if err != nil {
		instance.SetTerminatedWithError(workflow.Now(ctx), err)
	} else {
		instance.SetTerminated(workflow.Now(ctx))
	}
	if _, dbErr := db.NewUpdate().
		Model(&instance).
		WherePK().
		Exec(context.Background()); dbErr != nil {
		workflow.GetLogger(ctx).Error("error updating instance into database", "error", dbErr)
	}
	return err
}
