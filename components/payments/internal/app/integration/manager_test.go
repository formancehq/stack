package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"go.uber.org/dig"

	"github.com/formancehq/payments/internal/app/models"

	"github.com/formancehq/payments/internal/app/task"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func ChanClosed[T any](ch chan T) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}

type testContext[ConnectorConfig models.ConnectorConfigObject] struct {
	manager        *ConnectorManager[ConnectorConfig]
	taskStore      task.Repository
	connectorStore Repository
	loader         Loader[ConnectorConfig]
	provider       models.ConnectorProvider
}

func withManager[ConnectorConfig models.ConnectorConfigObject](builder *ConnectorBuilder,
	callback func(ctx *testContext[ConnectorConfig]),
) {
	l := logrus.New()
	if testing.Verbose() {
		l.SetLevel(logrus.DebugLevel)
	}

	DefaultContainerFactory := task.ContainerCreateFunc(func(ctx context.Context, descriptor models.TaskDescriptor, taskID uuid.UUID) (*dig.Container, error) {
		return dig.New(), nil
	})

	taskStore := task.NewInMemoryStore()
	managerStore := NewInMemoryStore()
	provider := models.ConnectorProvider(uuid.New().String())
	schedulerFactory := TaskSchedulerFactoryFn(func(resolver task.Resolver,
		maxTasks int,
	) *task.DefaultTaskScheduler {
		return task.NewDefaultScheduler(provider, taskStore,
			DefaultContainerFactory, resolver, maxTasks)
	})

	loader := NewLoaderBuilder[ConnectorConfig](provider).
		WithLoad(func(logger logging.Logger, config ConnectorConfig) Connector {
			return builder.Build()
		}).
		WithAllowedTasks(1).
		Build()
	manager := NewConnectorManager[ConnectorConfig](managerStore, loader, schedulerFactory, nil)

	defer func() {
		_ = manager.Uninstall(context.TODO())
	}()

	callback(&testContext[ConnectorConfig]{
		manager:        manager,
		taskStore:      taskStore,
		connectorStore: managerStore,
		loader:         loader,
		provider:       provider,
	})
}

func TestInstallConnector(t *testing.T) {
	t.Parallel()

	installed := make(chan struct{})
	builder := NewConnectorBuilder().
		WithInstall(func(ctx task.ConnectorContext) error {
			close(installed)

			return nil
		})
	withManager(builder, func(tc *testContext[models.EmptyConnectorConfig]) {
		err := tc.manager.Install(context.TODO(), models.EmptyConnectorConfig{})
		require.NoError(t, err)
		require.True(t, ChanClosed(installed))

		err = tc.manager.Install(context.TODO(), models.EmptyConnectorConfig{})
		require.Equal(t, ErrAlreadyInstalled, err)
	})
}

func TestUninstallConnector(t *testing.T) {
	t.Parallel()

	uninstalled := make(chan struct{})
	taskTerminated := make(chan struct{})
	taskStarted := make(chan struct{})
	builder := NewConnectorBuilder().
		WithResolve(func(name models.TaskDescriptor) task.Task {
			return func(ctx context.Context, stopChan task.StopChan) {
				close(taskStarted)
				defer close(taskTerminated)
				select {
				case flag := <-stopChan:
					flag <- struct{}{}
				case <-ctx.Done():
				}
			}
		}).
		WithInstall(func(ctx task.ConnectorContext) error {
			return ctx.Scheduler().Schedule(ctx.Context(), []byte(uuid.New().String()), models.TaskSchedulerOptions{
				ScheduleOption: models.OPTIONS_RUN_NOW,
				Restart:        false,
			})
		}).
		WithUninstall(func(ctx context.Context) error {
			close(uninstalled)

			return nil
		})
	withManager(builder, func(tc *testContext[models.EmptyConnectorConfig]) {
		err := tc.manager.Install(context.TODO(), models.EmptyConnectorConfig{})
		require.NoError(t, err)
		<-taskStarted
		require.NoError(t, tc.manager.Uninstall(context.TODO()))
		require.True(t, ChanClosed(uninstalled))
		// TODO: We need to give a chance to the connector to properly stop execution
		require.True(t, ChanClosed(taskTerminated))

		isInstalled, err := tc.manager.IsInstalled(context.TODO())
		require.NoError(t, err)
		require.False(t, isInstalled)
	})
}

func TestDisableConnector(t *testing.T) {
	t.Parallel()

	uninstalled := make(chan struct{})
	builder := NewConnectorBuilder().
		WithUninstall(func(ctx context.Context) error {
			close(uninstalled)

			return nil
		})
	withManager[models.EmptyConnectorConfig](builder, func(tc *testContext[models.EmptyConnectorConfig]) {
		err := tc.manager.Install(context.TODO(), models.EmptyConnectorConfig{})
		require.NoError(t, err)

		enabled, err := tc.manager.IsEnabled(context.TODO())
		require.NoError(t, err)
		require.True(t, enabled)

		require.NoError(t, tc.manager.Disable(context.TODO()))
		enabled, err = tc.manager.IsEnabled(context.TODO())
		require.NoError(t, err)
		require.False(t, enabled)
	})
}

func TestEnableConnector(t *testing.T) {
	t.Parallel()

	builder := NewConnectorBuilder()
	withManager[models.EmptyConnectorConfig](builder, func(tc *testContext[models.EmptyConnectorConfig]) {
		err := tc.connectorStore.Enable(context.TODO(), tc.loader.Name())
		require.NoError(t, err)

		err = tc.manager.Install(context.TODO(), models.EmptyConnectorConfig{})
		require.NoError(t, err)
	})
}

func TestRestoreEnabledConnector(t *testing.T) {
	t.Parallel()

	builder := NewConnectorBuilder()
	withManager(builder, func(tc *testContext[models.EmptyConnectorConfig]) {
		cfg, err := models.EmptyConnectorConfig{}.Marshal()
		require.NoError(t, err)

		err = tc.connectorStore.Install(context.TODO(), tc.loader.Name(), cfg)
		require.NoError(t, err)

		err = tc.manager.Restore(context.TODO())
		require.NoError(t, err)
		require.NotNil(t, tc.manager.connector)
	})
}

func TestRestoreNotInstalledConnector(t *testing.T) {
	t.Parallel()

	builder := NewConnectorBuilder()
	withManager(builder, func(tc *testContext[models.EmptyConnectorConfig]) {
		err := tc.manager.Restore(context.TODO())
		require.Equal(t, ErrNotInstalled, err)
	})
}
