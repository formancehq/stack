package runner

import (
	"context"
	"testing"
	"time"

	ingester "github.com/formancehq/stack/ee/ingester/internal"

	"github.com/formancehq/stack/ee/ingester/internal/drivers"
	"github.com/formancehq/stack/ee/ingester/internal/modules"

	"github.com/ThreeDotsLabs/watermill/message"
	api "github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func startRunner(t *testing.T, ctx context.Context, moduleFactory modules.Factory, store Store, connectorFactory drivers.Factory) *Runner {
	t.Helper()

	runner := NewRunner(
		store,
		moduleFactory,
		connectorFactory,
		logging.Testing(),
	)
	go func() {
		require.NoError(t, runner.Run(ctx))
	}()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		require.NoError(t, runner.Stop(ctx))
	})
	<-runner.Ready()

	return runner
}

func TestRunner(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	ctrl := gomock.NewController(t)
	moduleFactory := modules.NewMockFactory(ctrl)
	connectorFactory := drivers.NewMockFactory(ctrl)
	connector := drivers.NewMockDriver(ctrl)
	store := NewMockStore(ctrl)

	// Once the data are pulled, the pipeline should switch to ready state and subscribe to the module
	subscription := make(chan *message.Message, 1)

	pipelineConfiguration := ingester.NewPipelineConfiguration("module1", "connector")
	pipeline := ingester.NewPipeline(pipelineConfiguration, ingester.NewInitState())

	connectorFactory.EXPECT().
		Create(gomock.Any(), pipelineConfiguration.ConnectorID).
		Return(connector, nil, nil)
	connector.EXPECT().Start(gomock.Any()).Return(nil)

	moduleFactory.EXPECT().
		Create(pipelineConfiguration.Module).
		DoAndReturn(func(name string) (modules.Module, error) {
			module := modules.NewMockModule(ctrl)
			module.EXPECT().
				Pull(gomock.Any(), "").
				Return(&api.Cursor[ingester.Log]{
					Data: []ingester.Log{},
				}, nil)

			module.EXPECT().
				Subscribe(gomock.Any()).
				DoAndReturn(newFakeSubscriberFactory(subscription))

			return module, nil
		})

	store.EXPECT().
		StoreState(gomock.Any(), pipeline.ID, ingester.NewInitState()).
		Return(nil)

	store.EXPECT().
		StoreState(gomock.Any(), pipeline.ID, ingester.NewReadyState()).
		Return(nil)

	runner := startRunner(t, ctx, moduleFactory, store, connectorFactory)
	_, err := runner.StartPipeline(ctx, pipeline)
	require.NoError(t, err)

	event := `{"id": "xxx"}`
	connector.EXPECT().
		Accept(gomock.Any(), ingester.NewLogWithModule(pipelineConfiguration.Module, ingester.Log{
			ID: "xxx",
		})).
		Return([]error{nil}, nil)

	require.Eventually(t, func() bool {
		return runner.GetConnector("connector") != nil
	}, time.Second, 10*time.Millisecond)

	select {
	case <-runner.GetConnector("connector").Ready():
	case <-time.After(time.Second):
		require.Fail(t, "connector should be ready")
	}

	select {
	case subscription <- message.NewMessage(uuid.NewString(), []byte(event)):
	case <-time.After(time.Second):
		require.Fail(t, "message should have been handled")
	}

	require.Eventually(t, ctrl.Satisfied, 2*time.Second, 10*time.Millisecond)

	// notes(gfyrag): add this expectation AFTER the previous Eventually.
	// If configured before the Eventually, it will never finish as the stop call is made in a t.Cleanup defined earlier
	connector.EXPECT().Stop(gomock.Any()).Return(nil)
}
