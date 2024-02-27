package workflow

import (
	"github.com/formancehq/orchestration/internal/workflow/activities"
	"github.com/formancehq/orchestration/internal/workflow/stages"
	"github.com/uptrace/bun"
	"go.temporal.io/sdk/client"
	"go.uber.org/fx"
)

func NewModule(taskQueue string) fx.Option {
	ret := []fx.Option{
		fx.Provide(func(db *bun.DB, temporalClient client.Client) *WorkflowManager {
			return NewManager(db, temporalClient, taskQueue)
		}),
		fx.Provide(fx.Annotate(NewWorkflows, fx.ResultTags(`group:"workflows"`), fx.As(new(any)))),
		fx.Provide(fx.Annotate(activities.New, fx.ResultTags(`group:"activities"`), fx.As(new(any)))),
		fx.Provide(fx.Annotate(NewActivities, fx.ResultTags(`group:"activities"`), fx.As(new(any)))),
	}

	for _, schema := range stages.All() {
		ret = append(ret, fx.Supply(
			fx.Annotate(
				schema.GetWorkflow(),
				fx.ResultTags(`group:"workflows"`),
				fx.As(new(any)),
			),
		))
	}

	return fx.Options(ret...)
}
