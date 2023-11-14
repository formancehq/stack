package integration_test

import (
	"context"
	"testing"

	"github.com/formancehq/payments/cmd/connectors/internal/integration"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"
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
	manager        *integration.ConnectorsManager[ConnectorConfig]
	taskStore      task.Repository
	connectorStore integration.Repository
	loader         integration.Loader[ConnectorConfig]
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
	schedulerFactory := integration.TaskSchedulerFactoryFn(func(
		connectorID models.ConnectorID,
		resolver task.Resolver,
		maxTasks int,
	) *task.DefaultTaskScheduler {
		return task.NewDefaultScheduler(connectorID, taskStore,
			DefaultContainerFactory, resolver, metrics.NewNoOpMetricsRegistry(), maxTasks)
	})

	loader := integration.NewLoaderBuilder[ConnectorConfig](provider).
		WithLoad(func(logger logging.Logger, config ConnectorConfig) integration.Connector {
			return builder.Build()
		}).
		WithAllowedTasks(1).
		Build()
	manager := integration.NewConnectorManager[ConnectorConfig](provider, managerStore, loader, schedulerFactory, nil)

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
		_, err := tc.manager.Install(context.TODO(), "test1", models.EmptyConnectorConfig{
			Name: "test1",
		})
		require.NoError(t, err)
		require.True(t, ChanClosed(installed))

		_, err = tc.manager.Install(context.TODO(), "test1", models.EmptyConnectorConfig{
			Name: "test1",
		})
		require.Equal(t, integration.ErrAlreadyInstalled, err)

		connectors, err := tc.manager.FindAll(context.TODO())
		require.NoError(t, err)
		require.Len(t, connectors, 1)
		require.Equal(t, "test1", connectors[0].Name)

		isInstalled, err := tc.manager.IsInstalled(context.TODO(), connectors[0].ID)
		require.NoError(t, err)
		require.True(t, isInstalled)

		err = tc.manager.Uninstall(context.TODO(), connectors[0].ID)
		require.NoError(t, err)

		isInstalled, err = tc.manager.IsInstalled(context.TODO(), connectors[0].ID)
		require.NoError(t, err)
		require.False(t, isInstalled)
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
				RestartOption:  models.OPTIONS_RESTART_NEVER,
			})
		}).
		WithUninstall(func(ctx context.Context) error {
			close(uninstalled)

			return nil
		})
	withManager(builder, func(tc *testContext[models.EmptyConnectorConfig]) {
		_, err := tc.manager.Install(context.TODO(), "test1", models.EmptyConnectorConfig{
			Name: "test1",
		})
		require.NoError(t, err)
		<-taskStarted

		connectors, err := tc.manager.FindAll(context.TODO())
		require.NoError(t, err)
		require.Len(t, connectors, 1)
		require.Equal(t, "test1", connectors[0].Name)

		require.NoError(t, tc.manager.Uninstall(context.TODO(), connectors[0].ID))
		require.True(t, ChanClosed(uninstalled))
		// TODO: We need to give a chance to the connector to properly stop execution
		require.True(t, ChanClosed(taskTerminated))

		isInstalled, err := tc.manager.IsInstalled(context.TODO(), connectors[0].ID)
		require.NoError(t, err)
		require.False(t, isInstalled)
	})
}

func TestRestoreConnector(t *testing.T) {
	t.Parallel()

	builder := NewConnectorBuilder()
	withManager(builder, func(tc *testContext[models.EmptyConnectorConfig]) {
		cfg, err := models.EmptyConnectorConfig{
			Name: "test1",
		}.Marshal()
		require.NoError(t, err)

		connector := &models.Connector{
			ID: models.ConnectorID{
				Provider:  tc.provider,
				Reference: uuid.New(),
			},
			Name:     "test1",
			Provider: tc.provider,
		}

		err = tc.connectorStore.Install(context.TODO(), connector, cfg)
		require.NoError(t, err)

		err = tc.manager.Restore(context.TODO())
		require.NoError(t, err)
		require.Len(t, tc.manager.Connectors(), 1)

		require.NoError(t, tc.manager.Uninstall(context.TODO(), connector.ID))
	})
}

func TestRestoreNotInstalledConnector(t *testing.T) {
	t.Parallel()

	builder := NewConnectorBuilder()
	withManager(builder, func(tc *testContext[models.EmptyConnectorConfig]) {
		err := tc.manager.Restore(context.TODO())
		require.NoError(t, err)
		require.Len(t, tc.manager.Connectors(), 0)
	})
}
