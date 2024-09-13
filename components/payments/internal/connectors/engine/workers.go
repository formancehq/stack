package engine

import (
	"sync"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/temporal"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	temporalworkflow "go.temporal.io/sdk/workflow"
)

type Workers struct {
	logger logging.Logger

	temporalClient client.Client

	workers map[string]Worker
	rwMutex sync.RWMutex

	workflows  []temporal.DefinitionSet
	activities []temporal.DefinitionSet

	options worker.Options
}

type Worker struct {
	worker worker.Worker
}

func NewWorkers(logger logging.Logger, temporalClient client.Client, workflows, activities []temporal.DefinitionSet, options worker.Options) *Workers {
	return &Workers{
		logger:         logger,
		temporalClient: temporalClient,
		workers:        make(map[string]Worker),
		workflows:      workflows,
		activities:     activities,
		options:        options,
	}
}

// Close is called when app is terminated
func (w *Workers) Close() {
	w.rwMutex.Lock()
	defer w.rwMutex.Unlock()

	for _, worker := range w.workers {
		worker.worker.Stop()
	}
}

// Installing a new connector lauches a new worker
func (w *Workers) AddWorker(connectorID models.ConnectorID) error {
	w.rwMutex.Lock()
	defer w.rwMutex.Unlock()

	if _, ok := w.workers[connectorID.String()]; ok {
		return nil
	}

	worker := worker.New(w.temporalClient, connectorID.Reference, w.options)

	for _, set := range w.workflows {
		for _, workflow := range set {
			worker.RegisterWorkflowWithOptions(workflow.Func, temporalworkflow.RegisterOptions{
				Name: workflow.Name,
			})
		}
	}

	for _, set := range w.activities {
		for _, act := range set {
			worker.RegisterActivityWithOptions(act.Func, activity.RegisterOptions{
				Name: act.Name,
			})
		}
	}

	go func() {
		err := worker.Run(nil)
		if err != nil {
			w.logger.Errorf("worker loop stopped: %v", err)
		}
	}()

	w.workers[connectorID.String()] = Worker{
		worker: worker,
	}

	w.logger.Infof("worker for connector %s started", connectorID.String())

	return nil
}

// Uninstalling a connector stops the worker
func (w *Workers) RemoveWorker(connectorID models.ConnectorID) error {
	w.rwMutex.Lock()
	defer w.rwMutex.Unlock()

	worker, ok := w.workers[connectorID.String()]
	if !ok {
		return nil
	}

	worker.worker.Stop()

	delete(w.workers, connectorID.String())

	w.logger.Infof("worker for connector %s removed", connectorID.String())

	return nil
}
