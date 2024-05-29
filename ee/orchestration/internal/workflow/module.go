package workflow

import (
	"github.com/formancehq/orchestration/internal/temporalworker"
	"github.com/formancehq/orchestration/internal/workflow/activities"
	"github.com/formancehq/orchestration/internal/workflow/stages"
	"github.com/iancoleman/strcase"
	"github.com/uptrace/bun"
	"go.temporal.io/sdk/client"
	"go.uber.org/fx"
)

func NewModule(taskQueue string) fx.Option {
	ret := []fx.Option{
		fx.Provide(func(db *bun.DB, temporalClient client.Client) *WorkflowManager {
			return NewManager(db, temporalClient, taskQueue, true)
		}),
		fx.Provide(func() *Workflows {
			return NewWorkflows(true)
		}),
		fx.Provide(activities.New),
		fx.Provide(NewActivities),
		fx.Provide(fx.Annotate(func(a activities.Activities) temporalworker.DefinitionSet {
			return a.DefinitionSet()
		}, fx.ResultTags(`group:"activities"`))),
		fx.Provide(fx.Annotate(func(a Activities) temporalworker.DefinitionSet {
			return a.DefinitionSet()
		}, fx.ResultTags(`group:"activities"`))),
		fx.Provide(fx.Annotate(func(workflow *Workflows) temporalworker.DefinitionSet {
			return workflow.DefinitionSet()
		}, fx.ResultTags(`group:"workflows"`))),
	}

	set := temporalworker.NewDefinitionSet()
	for name, schema := range stages.All() {
		set = set.Append(temporalworker.Definition{
			Name: "Run" + strcase.ToCamel(name),
			Func: schema.GetWorkflow(),
		})
	}
	ret = append(ret, fx.Supply(
		fx.Annotate(set, fx.ResultTags(`group:"workflows"`)),
	))

	return fx.Options(ret...)
}
