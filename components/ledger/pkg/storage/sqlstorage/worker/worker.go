package worker

import (
	"context"
)

type Job[MODEL any] func(context.Context, []*MODEL) error

type modelHolder[MODEL any] struct {
	models  *MODEL
	errChan chan error
}

type Worker[MODEL any] struct {
	pending      []modelHolder[MODEL]
	writeChannel chan modelHolder[MODEL]
	jobs         chan []modelHolder[MODEL]
	releasedJob  chan struct{}
	workerJob    Job[MODEL]
	stopChan     chan chan struct{}
}

func NewWorker[MODEL any](workerJob Job[MODEL]) *Worker[MODEL] {
	return &Worker[MODEL]{
		workerJob:    workerJob,
		pending:      make([]modelHolder[MODEL], 0), // TODO(gfyrag): we need to limit the worker capacity
		jobs:         make(chan []modelHolder[MODEL]),
		writeChannel: make(chan modelHolder[MODEL], 1024), // TODO(gfyrag): Make configurable
		releasedJob:  make(chan struct{}, 1),
		stopChan:     make(chan chan struct{}, 1),
	}
}

func (w *Worker[MODEL]) writeLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case w.releasedJob <- struct{}{}:
		case modelHolders := <-w.jobs:
			models := make([]*MODEL, len(modelHolders))
			for i, holder := range modelHolders {
				models[i] = holder.models
			}
			err := w.workerJob(ctx, models)
			go func() {
				for _, holder := range modelHolders {
					select {
					case <-ctx.Done():
						return
					case holder.errChan <- err:
						close(holder.errChan)
					}
				}
			}()
		}
	}
}

// Run should be called in a goroutine
func (w *Worker[MODEL]) Run(ctx context.Context) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go w.writeLoop(ctx)

l:
	for {
		select {
		case <-ctx.Done():
			return
		case ch := <-w.stopChan:
			close(ch)
			return
		case mh := <-w.writeChannel:
			w.pending = append(w.pending, mh)
		case <-w.releasedJob:
			if len(w.pending) > 0 {
				for {
					select {
					case <-ctx.Done():
						return
					case ch := <-w.stopChan:
						close(ch)
						return
					case w.jobs <- w.pending:
						w.pending = make([]modelHolder[MODEL], 0)
						continue l
					}
				}
			}
			select {
			case <-ctx.Done():
				return
			case ch := <-w.stopChan:
				close(ch)
				return
			case mh := <-w.writeChannel:
				select {
				case <-ctx.Done():
					return
				case ch := <-w.stopChan:
					close(ch)
					return
				case w.jobs <- []modelHolder[MODEL]{mh}:
				}
			}
		}
	}
}

func (w *Worker[MODEL]) Stop(ctx context.Context) error {
	ch := make(chan struct{})
	select {
	case <-ctx.Done():
		return ctx.Err()
	case w.stopChan <- ch:
		select {
		case <-ch:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (w *Worker[MODEL]) WriteModels(ctx context.Context, model *MODEL) <-chan error {
	errChan := make(chan error, 1)

	select {
	case <-ctx.Done():
		errChan <- ctx.Err()
		close(errChan)
	case w.writeChannel <- modelHolder[MODEL]{
		models:  model,
		errChan: errChan,
	}:
	}

	return errChan
}
