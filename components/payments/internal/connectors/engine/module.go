package engine

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/go-libs/temporal"
	"github.com/formancehq/payments/internal/connectors/engine/activities"
	"github.com/formancehq/payments/internal/connectors/engine/plugins"
	"github.com/formancehq/payments/internal/connectors/engine/webhooks"
	"github.com/formancehq/payments/internal/connectors/engine/workflow"
	"github.com/formancehq/payments/internal/events"
	"github.com/formancehq/payments/internal/storage"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/fx"
)

func Module(pluginPath map[string]string, stack, stackURL string) fx.Option {
	ret := []fx.Option{
		fx.Supply(worker.Options{}),
		fx.Provide(func(
			temporalClient client.Client,
			workers *Workers,
			plugins plugins.Plugins,
			storage storage.Storage,
			webhooks webhooks.Webhooks,
		) Engine {
			return New(temporalClient, workers, plugins, storage, webhooks, stack)
		}),
		fx.Provide(func(publisher message.Publisher) *events.Events {
			return events.New(publisher, stackURL)
		}),
		fx.Provide(func() plugins.Plugins {
			return plugins.New(pluginPath)
		}),
		fx.Provide(func() webhooks.Webhooks {
			return webhooks.New()
		}),
		fx.Provide(func(temporalClient client.Client, plugins plugins.Plugins, webhooks webhooks.Webhooks) workflow.Workflow {
			return workflow.New(temporalClient, plugins, webhooks, stack)
		}),
		fx.Provide(func(storage storage.Storage, events *events.Events, plugins plugins.Plugins) activities.Activities {
			return activities.New(storage, events, plugins)
		}),
		fx.Provide(
			fx.Annotate(func(
				logger logging.Logger,
				temporalClient client.Client,
				workflows,
				activities []temporal.DefinitionSet,
				options worker.Options,
			) *Workers {
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
