package task

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"
)

//nolint:gochecknoglobals // allow in tests
var DefaultContainerFactory = ContainerCreateFunc(func(ctx context.Context, descriptor models.TaskDescriptor, taskID uuid.UUID) (*dig.Container, error) {
	return dig.New(), nil
})

func newDescriptor() models.TaskDescriptor {
	return []byte(uuid.New().String())
}

func TaskTerminatedWithStatus(
	store *InMemoryStore,
	connectorID models.ConnectorID,
	descriptor models.TaskDescriptor,
	expectedStatus models.TaskStatus,
	errString string,
) func() bool {
	return func() bool {
		status, resultErr, ok := store.Result(connectorID, descriptor)
		if !ok {
			return false
		}

		if resultErr != errString {
			return false
		}

		return status == expectedStatus
	}
}

func TaskTerminated(store *InMemoryStore, connectorID models.ConnectorID, descriptor models.TaskDescriptor) func() bool {
	return TaskTerminatedWithStatus(store, connectorID, descriptor, models.TaskStatusTerminated, "")
}

func TaskFailed(store *InMemoryStore, connectorID models.ConnectorID, descriptor models.TaskDescriptor, errStr string) func() bool {
	return TaskTerminatedWithStatus(store, connectorID, descriptor, models.TaskStatusFailed, errStr)
}

func TaskPending(store *InMemoryStore, connectorID models.ConnectorID, descriptor models.TaskDescriptor) func() bool {
	return TaskTerminatedWithStatus(store, connectorID, descriptor, models.TaskStatusPending, "")
}

func TaskActive(store *InMemoryStore, connectorID models.ConnectorID, descriptor models.TaskDescriptor) func() bool {
	return TaskTerminatedWithStatus(store, connectorID, descriptor, models.TaskStatusActive, "")
}

func TestTaskScheduler(t *testing.T) {
	t.Parallel()

	l := logrus.New()
	if testing.Verbose() {
		l.SetLevel(logrus.DebugLevel)
	}

	t.Run("Nominal", func(t *testing.T) {
		t.Parallel()

		store := NewInMemoryStore()
		connectorID := models.ConnectorID{
			Reference: uuid.New(),
			Provider:  models.ConnectorProviderDummyPay,
		}
		done := make(chan struct{})
		scheduler := NewDefaultScheduler(connectorID, store,
			DefaultContainerFactory, ResolverFn(func(descriptor models.TaskDescriptor) Task {
				return func(ctx context.Context) error {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case <-done:
						return nil
					}
				}
			}), metrics.NewNoOpMetricsRegistry(), 1)

		descriptor := newDescriptor()
		err := scheduler.Schedule(context.TODO(), descriptor, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_NOW,
			RestartOption:  models.OPTIONS_RESTART_NEVER,
		})
		require.NoError(t, err)

		require.Eventually(t, TaskActive(store, connectorID, descriptor), time.Second, 100*time.Millisecond)
		close(done)
		require.Eventually(t, TaskTerminated(store, connectorID, descriptor), time.Second, 100*time.Millisecond)
	})

	t.Run("Duplicate task", func(t *testing.T) {
		t.Parallel()

		store := NewInMemoryStore()
		connectorID := models.ConnectorID{
			Reference: uuid.New(),
			Provider:  models.ConnectorProviderDummyPay,
		}
		scheduler := NewDefaultScheduler(connectorID, store, DefaultContainerFactory,
			ResolverFn(func(descriptor models.TaskDescriptor) Task {
				return func(ctx context.Context) error {
					<-ctx.Done()

					return ctx.Err()
				}
			}), metrics.NewNoOpMetricsRegistry(), 1)

		descriptor := newDescriptor()
		err := scheduler.Schedule(context.TODO(), descriptor, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_NOW,
			RestartOption:  models.OPTIONS_RESTART_NEVER,
		})
		require.NoError(t, err)
		require.Eventually(t, TaskActive(store, connectorID, descriptor), time.Second, 100*time.Millisecond)

		err = scheduler.Schedule(context.TODO(), descriptor, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_NOW,
			RestartOption:  models.OPTIONS_RESTART_NEVER,
		})
		require.Equal(t, ErrAlreadyScheduled, err)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		connectorID := models.ConnectorID{
			Reference: uuid.New(),
			Provider:  models.ConnectorProviderDummyPay,
		}
		store := NewInMemoryStore()
		scheduler := NewDefaultScheduler(connectorID, store, DefaultContainerFactory,
			ResolverFn(func(descriptor models.TaskDescriptor) Task {
				return func() error {
					return errors.New("test")
				}
			}), metrics.NewNoOpMetricsRegistry(), 1)

		descriptor := newDescriptor()
		err := scheduler.Schedule(context.TODO(), descriptor, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_NOW,
			RestartOption:  models.OPTIONS_RESTART_NEVER,
		})
		require.NoError(t, err)
		require.Eventually(t, TaskFailed(store, connectorID, descriptor, "test"), time.Second,
			100*time.Millisecond)
	})

	t.Run("Pending", func(t *testing.T) {
		t.Parallel()

		connectorID := models.ConnectorID{
			Reference: uuid.New(),
			Provider:  models.ConnectorProviderDummyPay,
		}
		store := NewInMemoryStore()
		descriptor1 := newDescriptor()
		descriptor2 := newDescriptor()

		task1Terminated := make(chan struct{})
		task2Terminated := make(chan struct{})

		scheduler := NewDefaultScheduler(connectorID, store, DefaultContainerFactory,
			ResolverFn(func(descriptor models.TaskDescriptor) Task {
				switch string(descriptor) {
				case string(descriptor1):
					return func(ctx context.Context) error {
						select {
						case <-task1Terminated:
							return nil
						case <-ctx.Done():
							return ctx.Err()
						}
					}
				case string(descriptor2):
					return func(ctx context.Context) error {
						select {
						case <-task2Terminated:
							return nil
						case <-ctx.Done():
							return ctx.Err()
						}
					}
				}

				panic("unknown descriptor")
			}), metrics.NewNoOpMetricsRegistry(), 1)

		require.NoError(t, scheduler.Schedule(context.TODO(), descriptor1, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_NOW,
			RestartOption:  models.OPTIONS_RESTART_NEVER,
		}))
		require.NoError(t, scheduler.Schedule(context.TODO(), descriptor2, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_NOW,
			RestartOption:  models.OPTIONS_RESTART_NEVER,
		}))
		require.Eventually(t, TaskActive(store, connectorID, descriptor1), time.Second, 100*time.Millisecond)
		require.Eventually(t, TaskPending(store, connectorID, descriptor2), time.Second, 100*time.Millisecond)
		close(task1Terminated)
		require.Eventually(t, TaskTerminated(store, connectorID, descriptor1), time.Second, 100*time.Millisecond)
		require.Eventually(t, TaskActive(store, connectorID, descriptor2), time.Second, 100*time.Millisecond)
		close(task2Terminated)
		require.Eventually(t, TaskTerminated(store, connectorID, descriptor2), time.Second, 100*time.Millisecond)
	})

	t.Run("Stop scheduler", func(t *testing.T) {
		t.Parallel()

		connectorID := models.ConnectorID{
			Reference: uuid.New(),
			Provider:  models.ConnectorProviderDummyPay,
		}
		store := NewInMemoryStore()
		mainDescriptor := newDescriptor()
		workerDescriptor := newDescriptor()

		scheduler := NewDefaultScheduler(connectorID, store, DefaultContainerFactory,
			ResolverFn(func(descriptor models.TaskDescriptor) Task {
				switch string(descriptor) {
				case string(mainDescriptor):
					return func(ctx context.Context, scheduler Scheduler) {
						<-ctx.Done()
						require.NoError(t, scheduler.Schedule(ctx, workerDescriptor, models.TaskSchedulerOptions{
							ScheduleOption: models.OPTIONS_RUN_NOW,
							RestartOption:  models.OPTIONS_RESTART_NEVER,
						}))
					}
				default:
					panic("should not be called")
				}
			}), metrics.NewNoOpMetricsRegistry(), 1)

		require.NoError(t, scheduler.Schedule(context.TODO(), mainDescriptor, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_NOW,
			RestartOption:  models.OPTIONS_RESTART_NEVER,
		}))
		require.Eventually(t, TaskActive(store, connectorID, mainDescriptor), time.Second, 100*time.Millisecond)
		require.NoError(t, scheduler.Shutdown(context.TODO()))
		require.Eventually(t, TaskTerminated(store, connectorID, mainDescriptor), time.Second, 100*time.Millisecond)
		require.Eventually(t, TaskPending(store, connectorID, workerDescriptor), time.Second, 100*time.Millisecond)
	})
}
