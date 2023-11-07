package integration

import (
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
)

type TaskSchedulerFactory interface {
	Make(connectorID models.ConnectorID, resolver task.Resolver, maxTasks int) *task.DefaultTaskScheduler
}

type TaskSchedulerFactoryFn func(connectorID models.ConnectorID, resolver task.Resolver, maxProcesses int) *task.DefaultTaskScheduler

func (fn TaskSchedulerFactoryFn) Make(connectorID models.ConnectorID, resolver task.Resolver, maxTasks int) *task.DefaultTaskScheduler {
	return fn(connectorID, resolver, maxTasks)
}
