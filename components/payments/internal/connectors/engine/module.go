package engine

import (
	"context"

	"github.com/formancehq/payments/internal/connectors/engine/activities"
	"github.com/formancehq/payments/internal/connectors/engine/plugins"
	"github.com/formancehq/payments/internal/connectors/engine/workflow"
	"github.com/formancehq/payments/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/temporal"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/fx"
)

func Module(pluginPath map[string]string) fx.Option {
	ret := []fx.Option{
		fx.Supply(worker.Options{}),
		fx.Provide(New),
		fx.Provide(func() plugins.Plugins {
			return plugins.New(pluginPath)
		}),
		fx.Provide(func(temporalClient client.Client, plugins plugins.Plugins) workflow.Workflow {
			return workflow.New(temporalClient, plugins)
		}),
		fx.Provide(func(storage storage.Storage, plugins plugins.Plugins) activities.Activities {
			return activities.New(storage, plugins)
		}),
		fx.Provide(
			fx.Annotate(func(logger logging.Logger, temporalClient client.Client, workflows, activities []temporal.DefinitionSet, options worker.Options) *Workers {
				return NewWorkers(logger, temporalClient, workflows, activities, options)
			}, fx.ParamTags(``, ``, `group:"workflows"`, `group:"activities"`, ``)),
		),
		fx.Invoke(func(lc fx.Lifecycle, engine Engine, workers *Workers) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return engine.OnStart(ctx)
				},
				OnStop: func(ctx context.Context) error {
					workers.Close()
					return nil
				},
			})
		}),
		fx.Provide(fx.Annotate(func(a activities.Activities) temporal.DefinitionSet {
			return a.DefinitionSet()
		}, fx.ResultTags(`group:"activities"`))),
		fx.Provide(fx.Annotate(func(workflow workflow.Workflow) temporal.DefinitionSet {
			return workflow.DefinitionSet()
		}, fx.ResultTags(`group:"workflows"`))),
	}

	return fx.Options(ret...)
}
