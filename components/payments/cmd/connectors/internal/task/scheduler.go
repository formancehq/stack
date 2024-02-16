package task

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/alitto/pond"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

var (
	ErrValidation       = errors.New("validation error")
	ErrAlreadyScheduled = errors.New("already scheduled")
	ErrUnableToResolve  = errors.New("unable to resolve task")
)

type Scheduler interface {
	Schedule(ctx context.Context, p models.TaskDescriptor, options models.TaskSchedulerOptions) error
}

type taskHolder struct {
	descriptor models.TaskDescriptor
	cancel     func()
	logger     logging.Logger
	stopChan   StopChan
}

type ContainerCreateFunc func(ctx context.Context, descriptor models.TaskDescriptor, taskID uuid.UUID) (*dig.Container, error)

type DefaultTaskScheduler struct {
	connectorID      models.ConnectorID
	store            Repository
	metricsRegistry  metrics.MetricsRegistry
	containerFactory ContainerCreateFunc
	tasks            map[string]*taskHolder
	mu               sync.Mutex
	resolver         Resolver
	stopped          bool
	workerPool       *pond.WorkerPool
}

func (s *DefaultTaskScheduler) ListTasks(ctx context.Context, q storage.ListTasksQuery) (*api.Cursor[models.Task], error) {
	return s.store.ListTasks(ctx, s.connectorID, q)
}

func (s *DefaultTaskScheduler) ReadTask(ctx context.Context, taskID uuid.UUID) (*models.Task, error) {
	return s.store.GetTask(ctx, taskID)
}

func (s *DefaultTaskScheduler) ReadTaskByDescriptor(ctx context.Context, descriptor models.TaskDescriptor) (*models.Task, error) {
	taskDescriptor, err := json.Marshal(descriptor)
	if err != nil {
		return nil, err
	}

	return s.store.GetTaskByDescriptor(ctx, s.connectorID, taskDescriptor)
}

// Schedule schedules a task to be executed.
// Schedule waits for:
//   - Context to be done
//   - Task creation if the scheduler option is not equal to OPTIONS_RUN_NOW_SYNC
//   - Task termination if the scheduler option is equal to OPTIONS_RUN_NOW_SYNC
func (s *DefaultTaskScheduler) Schedule(ctx context.Context, descriptor models.TaskDescriptor, options models.TaskSchedulerOptions) error {
	select {
	case err := <-s.schedule(ctx, descriptor, options):
		return err
	case <-ctx.Done():
		return nil
	}
}

// schedule schedules a task to be executed.
// It returns an error chan that will be closed when the task is terminated if
// the scheduler option is equal to OPTIONS_RUN_NOW_SYNC. Otherwise, it will
// return an error chan that will be closed immediately after task creation.
func (s *DefaultTaskScheduler) schedule(ctx context.Context, descriptor models.TaskDescriptor, options models.TaskSchedulerOptions) <-chan error {
	s.mu.Lock()
	defer s.mu.Unlock()

	returnErrorFunc := func(err error) <-chan error {
		errChan := make(chan error, 1)
		if err != nil {
			errChan <- err
		}
		close(errChan)
		return errChan
	}

	taskID, err := descriptor.EncodeToString()
	if err != nil {
		return returnErrorFunc(err)
	}

	if _, ok := s.tasks[taskID]; ok {
		switch options.RestartOption {
		case models.OPTIONS_STOP_AND_RESTART, models.OPTIONS_RESTART_ALWAYS:
			// We still want to restart the task
		default:
			return returnErrorFunc(ErrAlreadyScheduled)
		}
	}

	switch options.RestartOption {
	case models.OPTIONS_RESTART_NEVER:
		_, err := s.ReadTaskByDescriptor(ctx, descriptor)
		if err == nil {
			return returnErrorFunc(nil)
		}
	case models.OPTIONS_RESTART_IF_NOT_ACTIVE:
		task, err := s.ReadTaskByDescriptor(ctx, descriptor)
		if err == nil && task.Status == models.TaskStatusActive {
			return nil
		}
	case models.OPTIONS_STOP_AND_RESTART:
		err := s.stopTask(ctx, descriptor)
		if err != nil {
			return returnErrorFunc(err)
		}
	case models.OPTIONS_RESTART_ALWAYS:
		// Do nothing
	}

	errChan := s.startTask(ctx, descriptor, options)

	return errChan
}

func (s *DefaultTaskScheduler) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stopped = true

	s.logger(ctx).Infof("Stopping scheduler...")
	s.workerPool.Stop()

	for name, task := range s.tasks {
		task.logger.Debugf("Stopping task")

		if task.stopChan != nil {
			errCh := make(chan struct{})
			task.stopChan <- errCh
			select {
			case <-errCh:
			case <-time.After(time.Second): // TODO: Make configurable
				task.logger.Debugf("Stopping using stop chan timeout, canceling context")
				task.cancel()
			}
		} else {
			task.cancel()
		}

		delete(s.tasks, name)
	}

	return nil
}

func (s *DefaultTaskScheduler) Restore(ctx context.Context) error {
	tasks, err := s.store.ListTasksByStatus(ctx, s.connectorID, models.TaskStatusActive)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if task.SchedulerOptions.Restart {
			task.SchedulerOptions.RestartOption = models.OPTIONS_RESTART_ALWAYS
		}

		errChan := s.startTask(ctx, task.GetDescriptor(), task.SchedulerOptions)
		select {
		case err := <-errChan:
			if err != nil {
				s.logger(ctx).Errorf("Unable to restore task %s: %s", task.ID, err)
			}
		case <-ctx.Done():
		}
	}

	return nil
}

func (s *DefaultTaskScheduler) registerTaskError(ctx context.Context, holder *taskHolder, taskErr any) {
	var taskError string

	switch v := taskErr.(type) {
	case error:
		taskError = v.Error()
	default:
		taskError = fmt.Sprintf("%s", v)
	}

	holder.logger.Errorf("Task terminated with error: %s", taskErr)

	err := s.store.UpdateTaskStatus(ctx, s.connectorID, holder.descriptor, models.TaskStatusFailed, taskError)
	if err != nil {
		holder.logger.Errorf("Error updating task status: %s", taskError)
	}
}

func (s *DefaultTaskScheduler) deleteTask(ctx context.Context, holder *taskHolder) {
	s.mu.Lock()
	defer s.mu.Unlock()

	taskID, err := holder.descriptor.EncodeToString()
	if err != nil {
		holder.logger.Errorf("Error encoding task descriptor: %s", err)

		return
	}

	delete(s.tasks, taskID)

	if s.stopped {
		return
	}

	oldestPendingTask, err := s.store.ReadOldestPendingTask(ctx, s.connectorID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return
		}

		logging.FromContext(ctx).Error(err)

		return
	}

	p := s.resolver.Resolve(oldestPendingTask.GetDescriptor())
	if p == nil {
		logging.FromContext(ctx).Errorf("unable to resolve task")
		return
	}

	errChan := s.startTask(ctx, oldestPendingTask.GetDescriptor(), models.TaskSchedulerOptions{
		ScheduleOption: models.OPTIONS_RUN_NOW,
	})
	select {
	case err, ok := <-errChan:
		if !ok {
			return
		}
		if err != nil {
			logging.FromContext(ctx).Error(err)
		}
	case <-ctx.Done():
		return
	}
}

type StopChan chan chan struct{}

// Lock should be held when calling this function
func (s *DefaultTaskScheduler) stopTask(ctx context.Context, descriptor models.TaskDescriptor) error {
	taskID, err := descriptor.EncodeToString()
	if err != nil {
		return err
	}

	task, ok := s.tasks[taskID]
	if !ok {
		return nil
	}

	task.logger.Infof("Stopping task...")

	if task.stopChan != nil {
		errCh := make(chan struct{})
		task.stopChan <- errCh
		select {
		case <-errCh:
		case <-time.After(time.Second): // TODO: Make configurable
			task.logger.Debugf("Stopping using stop chan timeout, canceling context")
			task.cancel()
		}
	} else {
		task.cancel()
	}

	err = s.store.UpdateTaskStatus(ctx, s.connectorID, descriptor, models.TaskStatusStopped, "")
	if err != nil {
		task.logger.Errorf("Error updating task status: %s", err)
		return err
	}

	delete(s.tasks, taskID)

	return nil
}

func (s *DefaultTaskScheduler) startTask(ctx context.Context, descriptor models.TaskDescriptor, options models.TaskSchedulerOptions) <-chan error {
	errChan := make(chan error, 1)

	task, err := s.store.FindAndUpsertTask(ctx, s.connectorID, descriptor,
		models.TaskStatusActive, options, "")
	if err != nil {
		errChan <- errors.Wrap(err, "finding task and update")
		close(errChan)
		return errChan
	}

	logger := s.logger(ctx).WithFields(map[string]interface{}{
		"task-id": task.ID,
	})

	taskResolver := s.resolver.Resolve(task.GetDescriptor())
	if taskResolver == nil {
		errChan <- ErrUnableToResolve
		close(errChan)
		return errChan
	}

	ctx, cancel := context.WithCancel(ctx)

	holder := &taskHolder{
		cancel:     cancel,
		logger:     logger,
		descriptor: descriptor,
	}

	container, err := s.containerFactory(ctx, descriptor, task.ID)
	if err != nil {
		// TODO: Handle error
		panic(err)
	}

	err = container.Provide(func() context.Context {
		return ctx
	})
	if err != nil {
		panic(err)
	}

	err = container.Provide(func() Scheduler {
		return s
	})
	if err != nil {
		panic(err)
	}

	err = container.Provide(func() models.ConnectorID {
		return s.connectorID
	})
	if err != nil {
		panic(err)
	}

	err = container.Provide(func() models.TaskID {
		return models.TaskID(task.ID)
	})

	err = container.Provide(func() StopChan {
		s.mu.Lock()
		defer s.mu.Unlock()

		holder.stopChan = make(StopChan, 1)

		return holder.stopChan
	})
	if err != nil {
		panic(err)
	}

	err = container.Provide(func() logging.Logger {
		return logger
	})
	if err != nil {
		panic(err)
	}

	err = container.Provide(func() metrics.MetricsRegistry {
		return s.metricsRegistry
	})
	if err != nil {
		panic(err)
	}

	err = container.Provide(func() StateResolver {
		return StateResolverFn(func(ctx context.Context, v any) error {
			t, err := s.store.GetTask(ctx, task.ID)
			if err != nil {
				return err
			}

			if t.State == nil || len(t.State) == 0 {
				return nil
			}

			return json.Unmarshal(t.State, v)
		})
	})
	if err != nil {
		panic(err)
	}

	taskID, err := holder.descriptor.EncodeToString()
	if err != nil {
		errChan <- err
		close(errChan)
		return errChan
	}

	s.tasks[taskID] = holder

	sendError := false
	switch options.ScheduleOption {
	case models.OPTIONS_RUN_NOW_SYNC:
		sendError = true
		fallthrough
	case models.OPTIONS_RUN_NOW:
		options.Duration = 0
		fallthrough
	case models.OPTIONS_RUN_SCHEDULED_AT:
		if !options.ScheduleAt.IsZero() {
			options.Duration = time.Until(options.ScheduleAt)
			if options.Duration < 0 {
				options.Duration = 0
			}
		}
		fallthrough
	case models.OPTIONS_RUN_IN_DURATION:
		go func() {
			if options.Duration > 0 {
				logger.Infof("Waiting %s before starting task...", options.Duration)
				// todo(gfyrag): need to listen on stopChan if the application is stopped
				time.Sleep(options.Duration)
			}

			logger.Infof("Starting task...")

			defer func() {
				defer s.deleteTask(ctx, holder)

				if sendError {
					defer close(errChan)
				}

				if e := recover(); e != nil {
					switch v := e.(type) {
					case error:
						if errors.Is(v, pond.ErrSubmitOnStoppedPool) {
							// Pool is stopped and task is marked as active,
							// nothing to do as they will be restarted on
							// next startup
							return
						}
					}

					s.registerTaskError(ctx, holder, e)
					debug.PrintStack()

					if sendError {
						switch v := e.(type) {
						case error:
							errChan <- v
						default:
							errChan <- fmt.Errorf("%s", v)
						}
					}
				}
			}()

			done := make(chan struct{})
			s.workerPool.Submit(func() {
				defer close(done)
				err = container.Invoke(taskResolver)
			})
			select {
			case <-done:
			case <-ctx.Done():
				return
			}
			if err != nil {
				s.registerTaskError(ctx, holder, err)

				if sendError {
					errChan <- err
					return
				}

				return
			}

			logger.Infof("Task terminated with success")

			err = s.store.UpdateTaskStatus(ctx, s.connectorID, descriptor, models.TaskStatusTerminated, "")
			if err != nil {
				logger.Errorf("Error updating task status: %s", err)
				if sendError {
					errChan <- err
				}
			}
		}()
	case models.OPTIONS_RUN_PERIODICALLY:
		go func() {
			defer func() {
				defer s.deleteTask(ctx, holder)

				if e := recover(); e != nil {
					s.registerTaskError(ctx, holder, e)
					debug.PrintStack()

					return
				}
			}()

			processFunc := func() (bool, error) {
				done := make(chan struct{})
				s.workerPool.Submit(func() {
					defer close(done)
					err = container.Invoke(taskResolver)
				})
				select {
				case <-done:
				case <-ctx.Done():
					return true, nil
				case ch := <-holder.stopChan:
					logger.Infof("Stopping task...")
					close(ch)
					return true, nil
				}
				if err != nil {
					s.registerTaskError(ctx, holder, err)
					return false, err
				}

				return false, err
			}

			// launch it once before starting the ticker
			stopped, err := processFunc()
			if err != nil {
				// error is already registered
				return
			}

			if stopped {
				// Task is stopped or context is done
				return
			}

			logger.Infof("Starting task...")
			ticker := time.NewTicker(options.Duration)
			for {
				select {
				case ch := <-holder.stopChan:
					logger.Infof("Stopping task...")
					close(ch)
					return
				case <-ctx.Done():
					return
				case <-ticker.C:
					logger.Infof("Polling trigger, running task...")
					stop, err := processFunc()
					if err != nil {
						// error is already registered
						return
					}

					if stop {
						// Task is stopped or context is done
						return
					}
				}
			}

		}()
	}

	if !sendError {
		close(errChan)
	}

	return errChan
}

func (s *DefaultTaskScheduler) logger(ctx context.Context) logging.Logger {
	return logging.FromContext(ctx).WithFields(map[string]any{
		"component":   "scheduler",
		"connectorID": s.connectorID,
	})
}

var _ Scheduler = &DefaultTaskScheduler{}

func NewDefaultScheduler(
	connectorID models.ConnectorID,
	store Repository,
	containerFactory ContainerCreateFunc,
	resolver Resolver,
	metricsRegistry metrics.MetricsRegistry,
	maxTasks int,
) *DefaultTaskScheduler {
	return &DefaultTaskScheduler{
		connectorID:      connectorID,
		store:            store,
		metricsRegistry:  metricsRegistry,
		tasks:            map[string]*taskHolder{},
		containerFactory: containerFactory,
		resolver:         resolver,
		workerPool:       pond.New(maxTasks, maxTasks),
	}
}
