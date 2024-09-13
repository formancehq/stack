package engine

import (
	"sync"

	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/go-libs/temporal"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	temporalworkflow "go.temporal.io/sdk/workflow"
)

const (
	defaultWorkerName = "default"
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
	workers := &Workers{
		logger:         logger,
		temporalClient: temporalClient,
		workers:        make(map[string]Worker),
		workflows:      workflows,
		activities:     activities,
		options:        options,
	}

	// For all operation outside of connectors handlers
	workers.AddWorker(defaultWorkerName)

	return workers
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
// A default one is instantiated when the workers struct is created
func (w *Workers) AddWorker(name string) error {
	w.rwMutex.Lock()
	defer w.rwMutex.Unlock()

	if _, ok := w.workers[name]; ok {
		return nil
	}

	worker := worker.New(w.temporalClient, name, w.options)

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

	w.workers[name] = Worker{
		worker: worker,
	}

	w.logger.Infof("worker for connector %s started", name)

	return nil
}

// Uninstalling a connector stops the worker
func (w *Workers) RemoveWorker(name string) error {
	w.rwMutex.Lock()
	defer w.rwMutex.Unlock()

	worker, ok := w.workers[name]
	if !ok {
		return nil
	}

	worker.worker.Stop()

	delete(w.workers, name)

	w.logger.Infof("worker for connector %s removed", name)

	return nil
}
