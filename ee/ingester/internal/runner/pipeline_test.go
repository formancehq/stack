package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	ingester "github.com/formancehq/stack/ee/ingester/internal"

	"github.com/formancehq/stack/ee/ingester/internal/drivers"
	"github.com/formancehq/stack/ee/ingester/internal/modules"

	"github.com/pkg/errors"

	"github.com/ThreeDotsLabs/watermill/message"
	api "github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func runPipeline(t *testing.T, ctx context.Context, pipeline ingester.Pipeline, subscriber modules.Module, connector drivers.Driver) (*PipelineHandler, <-chan ingester.State) {
	t.Helper()

	handler := NewPipelineHandler(
		pipeline,
		subscriber,
		connector,
		logging.Testing(),
		WithStateRetryInterval(50*time.Millisecond),
	)
	stateListener, cancelStateListener := handler.GetActiveState().Listen()
	t.Cleanup(cancelStateListener)

	go handler.Run(ctx)
	t.Cleanup(func() {
		require.NoError(t, handler.Shutdown(ctx))
	})

	return handler, stateListener
}

func TestPipelineReady(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	ctrl := gomock.NewController(t)
	module := modules.NewMockModule(ctrl)
	connector := drivers.NewMockDriver(ctrl)
	log := ingester.Log{
		ID:      "xxx",
		Payload: json.RawMessage(`{}`),
	}

	subscription := make(chan *message.Message, 1)
	module.EXPECT().
		Subscribe(gomock.Any()).
		DoAndReturn(newFakeSubscriberFactory(subscription))

	connector.EXPECT().
		Accept(gomock.Any(), ingester.NewLogWithModule("testing", log)).
		Return([]error{nil}, nil)

	pipelineConfiguration := ingester.NewPipelineConfiguration("testing", "testing")
	pipeline := ingester.NewPipeline(pipelineConfiguration, ingester.NewReadyState())

	_, stateListener := runPipeline(t, ctx, pipeline, module, connector)

	ShouldReceive(t, ingester.NewReadyState(), stateListener)

	marshalledLog, err := json.Marshal(log)
	require.NoError(t, err)

	subscription <- message.NewMessage(uuid.NewString(), marshalledLog)
	require.Eventually(t, ctrl.Satisfied, time.Second, 10*time.Millisecond)
}

func TestPipelineRetryFailingState(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	ctrl := gomock.NewController(t)
	module := modules.NewMockModule(ctrl)
	connector := drivers.NewMockDriver(ctrl)

	module.EXPECT().
		Subscribe(gomock.Any()).
		Return(nil, errors.New("failed to init module"))

	module.EXPECT().
		Subscribe(gomock.Any()).
		DoAndReturn(newFakeSubscriberFactory(make(chan *message.Message)))

	pipelineConfiguration := ingester.NewPipelineConfiguration("testing", "testing")
	pipeline := ingester.NewPipeline(pipelineConfiguration, ingester.NewReadyState())
	_, stateListener := runPipeline(t, ctx, pipeline, module, connector)

	ShouldReceive(t, ingester.NewReadyState(), stateListener)
}

func TestPipelinePause(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	ctrl := gomock.NewController(t)
	module := modules.NewMockModule(ctrl)
	connector := drivers.NewMockDriver(ctrl)

	state := ingester.NewPauseState(ingester.NewReadyState())
	pipelineConfiguration := ingester.NewPipelineConfiguration("testing", "testing")
	pipeline := ingester.NewPipeline(pipelineConfiguration, state)

	_, stateListener := runPipeline(t, ctx, pipeline, module, connector)

	ShouldReceive(t, state, stateListener)
}

func TestPipelineInit(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	ctrl := gomock.NewController(t)
	module := modules.NewMockModule(ctrl)
	connector := drivers.NewMockDriver(ctrl)

	const initialMessages = 10
	messages := make([]ingester.Log, initialMessages)
	for i := 0; i < initialMessages; i++ {
		messages[i] = ingester.Log{
			ID: fmt.Sprint(i),
		}
	}

	module.EXPECT().
		Pull(gomock.Any(), "").
		Return(&api.Cursor[ingester.Log]{
			Data: messages,
		}, nil)

	for i := 0; i < initialMessages; i++ {
		connector.EXPECT().
			Accept(gomock.Any(), ingester.NewLogWithModule("testing", messages[i])).
			Return([]error{nil}, nil)
	}

	// Once the data are pulled, the pipeline should switch to ready state and subscribe to the module
	module.EXPECT().
		Subscribe(gomock.Any()).
		DoAndReturn(AlwaysEmptySubscription)

	pipelineConfiguration := ingester.NewPipelineConfiguration("testing", "testing")
	pipeline := ingester.NewPipeline(pipelineConfiguration, ingester.NewInitState())
	_, stateListener := runPipeline(t, ctx, pipeline, module, connector)

	ShouldReceive(t, ingester.NewInitState(), stateListener)
	ShouldReceive(t, ingester.NewReadyState(), stateListener)
}
