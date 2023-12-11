package temporalworker

import (
	"context"
	"reflect"

	temporalworkflow "go.temporal.io/sdk/workflow"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/fx"
)

func registerWorkflow(worker worker.Worker, workflow any) {
	valueOfWorkflow := reflect.ValueOf(workflow)
	switch valueOfWorkflow.Kind() {
	case reflect.Func:
		worker.RegisterWorkflow(workflow)
	case reflect.Struct:
		for i := 0; i < valueOfWorkflow.NumMethod(); i++ {
			name := reflect.TypeOf(workflow).Method(i).Name
			worker.RegisterWorkflowWithOptions(valueOfWorkflow.Method(i).Interface(), temporalworkflow.RegisterOptions{
				Name: name,
			})
		}
	case reflect.Ptr:
		registerWorkflow(worker, valueOfWorkflow.Elem().Interface())
	}
}

func registerActivity(worker worker.Worker, act any) {
	valueOfActivities := reflect.ValueOf(act)
	switch valueOfActivities.Kind() {
	case reflect.Struct:
		for i := 0; i < valueOfActivities.NumMethod(); i++ {
			name := reflect.TypeOf(act).Method(i).Name
			worker.RegisterActivityWithOptions(valueOfActivities.Method(i).Interface(), activity.RegisterOptions{
				Name: name,
			})
		}
	case reflect.Func:
		worker.RegisterActivity(act)
	case reflect.Ptr:
		registerActivity(worker, valueOfActivities.Elem().Interface())
	}
}

func NewWorker(c client.Client, taskQueue string, workflows, activities []any) worker.Worker {
	worker := worker.New(c, taskQueue, worker.Options{})

	for _, workflow := range workflows {
		registerWorkflow(worker, workflow)
	}

	for _, act := range activities {
		registerActivity(worker, act)
	}

	return worker
}

func NewWorkerModule(taskQueue string) fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(func(c client.Client, workflows, activities []any) worker.Worker {
				return NewWorker(c, taskQueue, workflows, activities)
			}, fx.ParamTags(``, `group:"workflows"`, `group:"activities"`)),
		),
		fx.Invoke(func(lc fx.Lifecycle, w worker.Worker) {
			willStop := false
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						err := w.Run(worker.InterruptCh())
						if err != nil {
							if !willStop {
								panic(err)
							}
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					willStop = true
					w.Stop()
					return nil
				},
			})
		}),
	)
}
