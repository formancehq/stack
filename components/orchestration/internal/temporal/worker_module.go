package temporal

import (
	"context"
	"reflect"

	sdk "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/orchestration/internal/workflow/activities"
	"github.com/formancehq/orchestration/internal/workflow/stages"
	"github.com/uptrace/bun"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/fx"
)

func NewWorker(c client.Client, db *bun.DB, apiClient *sdk.APIClient, taskQueue string) worker.Worker {
	w := worker.New(c, taskQueue, worker.Options{})

	workflow := workflow.NewWorkflows(db)
	activities := activities.New(apiClient)

	valueOfActivities := reflect.ValueOf(activities)
	for i := 0; i < valueOfActivities.NumMethod(); i++ {
		w.RegisterActivityWithOptions(valueOfActivities.Method(i).Interface(), activity.RegisterOptions{
			Name: reflect.TypeOf(activities).Method(i).Name,
		})
	}
	RegisterWorkflows(workflow, w)

	return w
}

func RegisterWorkflows(workflows *workflow.Workflows, w interface {
	RegisterWorkflow(any)
}) {
	w.RegisterWorkflow(workflows.Run)
	for _, schema := range stages.All() {
		w.RegisterWorkflow(schema.GetWorkflow())
	}
}

func NewWorkerModule(taskQueue string) fx.Option {
	return fx.Options(
		fx.Provide(func(c client.Client, db *bun.DB, apiClient *sdk.APIClient) worker.Worker {
			return NewWorker(c, db, apiClient, taskQueue)
		}),
		fx.Invoke(func(lc fx.Lifecycle, w worker.Worker) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						err := w.Run(worker.InterruptCh())
						if err != nil {
							panic(err)
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					w.Stop()
					return nil
				},
			})
		}),
	)
}
