package worker

import (
	"context"
	"time"
)

type WorkerJob[MODEL any] func(context.Context, []MODEL) error

type Model[MODEL any] struct {
	model           MODEL
	isErrChanClosed *bool
	errChan         chan error
}

type Worker[MODEL any] struct {
	ctx       context.Context
	batchSize int
	batchTime time.Duration

	models     []Model[MODEL]
	modelsChan chan []Model[MODEL]

	workerJob WorkerJob[MODEL]
}

func NewWorker[MODEL any](batchSize int, batchTime time.Duration, workerJob WorkerJob[MODEL]) *Worker[MODEL] {
	return &Worker[MODEL]{
		batchSize: batchSize,
		batchTime: batchTime,
		workerJob: workerJob,

		models:     nil,
		modelsChan: make(chan []Model[MODEL], 1024),
	}
}

// Run should be called in a goroutine
func (w *Worker[MODEL]) Run(ctx context.Context) {
	w.ctx = ctx
	ticker := time.NewTicker(w.batchTime)

	writeFn := func(context.Context) {
		if err := w.write(ctx); err != nil {
			// TODO(polo): should we stop the writer in case of an error ?
			// return err
			return
		}
	}

	for {
		select {
		case <-ctx.Done():
			return

		case ms := <-w.modelsChan:
			w.models = append(w.models, ms...)
			if len(w.models) >= w.batchSize {
				writeFn(ctx)
			}

		case <-ticker.C:
			writeFn(ctx)
		}
	}
}

func (w *Worker[MODEL]) write(ctx context.Context) error {
	if len(w.models) == 0 {
		return nil
	}

	models := make([]MODEL, len(w.models))
	for i, model := range w.models {
		models[i] = model.model
	}

	err := w.workerJob(ctx, models)

	for _, model := range w.models {
		if model.isErrChanClosed != nil && !*model.isErrChanClosed {
			model.errChan <- err
			close(model.errChan)
			*model.isErrChanClosed = true
		}
	}

	// Release the slice in order to release the underlying memory to the
	// garbage collector.
	w.models = nil

	return err
}

func (w *Worker[MODEL]) WriteModels(models []MODEL) <-chan error {
	errChan := make(chan error, 1)
	isClosed := false

	ms := make([]Model[MODEL], 0, len(models))
	for _, model := range models {
		ms = append(ms, Model[MODEL]{
			model:           model,
			isErrChanClosed: &isClosed,
			errChan:         errChan,
		})
	}

	// TODO(polo): check with antoine and max to remove batch transactions
	// TODO(polo): context + err if full

	w.modelsChan <- ms

	return errChan
}
