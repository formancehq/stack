package integration

import (
	"github.com/formancehq/payments/cmd/connectors/internal/task"
)

type TaskSchedulerFactory interface {
	Make(resolver task.Resolver, maxTasks int) *task.DefaultTaskScheduler
}

type TaskSchedulerFactoryFn func(resolver task.Resolver, maxProcesses int) *task.DefaultTaskScheduler

func (fn TaskSchedulerFactoryFn) Make(resolver task.Resolver, maxTasks int) *task.DefaultTaskScheduler {
	return fn(resolver, maxTasks)
}
