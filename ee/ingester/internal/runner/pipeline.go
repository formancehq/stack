package runner

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	ingester "github.com/formancehq/stack/ee/ingester/internal"

	"github.com/formancehq/stack/ee/ingester/internal/drivers"
	"github.com/formancehq/stack/ee/ingester/internal/modules"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

var (
	DefaultPullRetryPeriod    = 10 * time.Second
	DefaultPushRetryPeriod    = 10 * time.Second
	DefaultStateRetryInterval = 5 * time.Second
)

type PipelineHandlerConfig struct {
	ModulePullRetryPeriod    time.Duration
	ConnectorPushRetryPeriod time.Duration
	StateRetryInterval       time.Duration
}

type PipelineOption func(config *PipelineHandlerConfig)

func WithModulePullPeriod(v time.Duration) PipelineOption {
	return func(config *PipelineHandlerConfig) {
		config.ModulePullRetryPeriod = v
	}
}

func WithConnectorPushRetryPeriod(v time.Duration) PipelineOption {
	return func(config *PipelineHandlerConfig) {
		config.ConnectorPushRetryPeriod = v
	}
}

func WithStateRetryInterval(v time.Duration) PipelineOption {
	return func(config *PipelineHandlerConfig) {
		config.StateRetryInterval = v
	}
}

var (
	defaultPipelineOptions = []PipelineOption{
		WithModulePullPeriod(DefaultPullRetryPeriod),
		WithStateRetryInterval(DefaultStateRetryInterval),
		WithConnectorPushRetryPeriod(DefaultPushRetryPeriod),
	}
)

type PipelineHandler struct {
	mu sync.Mutex

	pipeline       ingester.Pipeline
	stopChannel    chan chan error
	module         modules.Module
	connector      drivers.Driver
	expectedState  *Signal[ingester.State]
	activeState    *Signal[ingester.State]
	pipelineConfig PipelineHandlerConfig
	stateHandler   *StateHandler
	logger         logging.Logger
}

// Pause can return following errors:
// * ErrAlreadyPaused
func (p *PipelineHandler) Pause() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	actualExpectedState := p.expectedState.Actual()
	if actualExpectedState.Label == ingester.StateLabelPause {
		return NewErrInvalidStateSwitch(p.pipeline.ID, actualExpectedState.Label, ingester.StateLabelStop)
	}
	p.expectedState.Signal(ingester.NewPauseState(
		*p.activeState.Actual(),
	))

	return nil
}

// Resume can return following errors:
// * ErrNotInPauseState
func (p *PipelineHandler) Resume() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	actualExpectedState := p.expectedState.Actual()
	if actualExpectedState.Label != ingester.StateLabelPause {
		return NewErrInvalidStateSwitch(p.pipeline.ID, actualExpectedState.Label, ingester.StateLabelPause)
	}
	p.expectedState.Signal(*actualExpectedState.PreviousState)

	return nil
}

func (p *PipelineHandler) Reset() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.expectedState.Signal(ingester.NewInitState())

	return nil
}

// Stop try to stop the pipeline
// It is asynchronous and can controller by watching the active state of the pipeline
// see GetActiveState
func (p *PipelineHandler) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	actualExpectedState := p.expectedState.Actual()
	if actualExpectedState.Label == ingester.StateLabelStop {
		return NewErrInvalidStateSwitch(p.pipeline.ID, actualExpectedState.Label, ingester.StateLabelStop)
	}

	p.expectedState.Signal(ingester.NewStopState(*actualExpectedState))

	return nil
}

func (p *PipelineHandler) switchToState(ctx context.Context, newState ingester.State) bool {
	if p.activeState.Actual() != nil &&
		newState.Label == p.activeState.Actual().Label {
		return true
	}
	p.logger.Infof("Switching to state '%s'", newState.Label)

	var fn func(ctx context.Context, readyChan chan struct{}) error
	switch newState.Label {
	case ingester.StateLabelInit:
		fn = p.init
	case ingester.StateLabelReady:
		fn = p.handleFlow
	case ingester.StateLabelPause:
		fn = p.pause
	case ingester.StateLabelStop:
		fn = p.stop
	}
	if err := p.stateHandler.Switch(ctx, fn); err != nil {
		p.logger.Errorf("Error switching to state '%s': %s", newState.Label, err)
		return true
	}
	p.activeState.Signal(newState)
	p.logger.Infof("Switched to state '%s'", newState.Label)

	return newState.Label != ingester.StateLabelStop
}

func (p *PipelineHandler) Run(ctx context.Context) {
	defer p.activeState.Close()

	stateChangedListener, cancelExpectedStateListener := p.expectedState.Listen()
	defer cancelExpectedStateListener()

	p.stateHandler = NewEmptyStateHandler(p.logger)
	p.stateHandler.Run(ctx, make(chan struct{}))

	for {
		select {
		case newState := <-stateChangedListener:
			if !p.switchToState(ctx, newState) {
				return
			}
		case errorChannel := <-p.stopChannel:
			p.logger.Debugf("stopping pipeline signal received...")
			err := p.stateHandler.Cancel(ctx)
			errorChannel <- err
			if err != nil {
				p.logger.Errorf("pipeline stopped with error: %s", err)
			} else {
				p.logger.Infof("pipeline stopped")
			}
			return
		case <-time.After(p.pipelineConfig.StateRetryInterval):
			if !p.switchToState(ctx, *p.expectedState.Actual()) {
				return
			}
		}
	}
}

func (p *PipelineHandler) handleFlow(ctx context.Context, ready chan struct{}) error {
	var messages <-chan *message.Message

	subscriptionContext, cancelSubscriptionContext := context.WithCancel(context.Background())
	defer func() {
		p.logger.Infof("cancel subscription")
		cancelSubscriptionContext()
	}()
	subscriptionContext = logging.ContextWithLogger(subscriptionContext, p.logger)

	p.logger.Infof("subscribing to messages")
	var err error
	messages, err = p.module.Subscribe(subscriptionContext)
	if err != nil {
		return err
	}

	p.logger.Infof("Pipeline ready")
	close(ready)

	processingMessages := sync.WaitGroup{}

	for {
		p.logger.Debugf("Wait for next message")
		select {
		case msg := <-messages:
			p.logger.Debugf("Got new message '%s': %s", msg.UUID, string(msg.Payload))
			processingMessages.Add(1)
			go func() {
				defer processingMessages.Done()

				logger := p.logger.WithField("msg", string(msg.Payload))
				log := ingester.Log{}
				if err := json.Unmarshal(msg.Payload, &log); err != nil {
					logger.Errorf("unable to unmarshal log: %s", err)
					msg.Nack()
					return
				}

				itemsErrors, err := p.connector.Accept(ctx, ingester.LogWithModule{
					Module: p.pipeline.Module,
					Log:    log,
				})
				if err == nil {
					err = itemsErrors[0]
				}
				if err != nil {
					logger.Errorf("Unable to send msg to connector: %s", err)
					msg.Nack()
					return
				}

				logger.Debug("Message sent to connector")
				msg.Ack()
			}()
		case <-ctx.Done():
			p.logger.Infof("waiting all messages completion")
			processingMessages.Wait()
			p.logger.Infof("all messages processed")
			return nil
		}
	}
}

func (p *PipelineHandler) init(ctx context.Context, ready chan struct{}) error {
	close(ready)

	wg := sync.WaitGroup{}
	cursor := p.expectedState.Actual().Cursor
	for {
		cursor, err := p.module.Pull(ctx, cursor)
		if err != nil {
			p.logger.Errorf("Error pulling module: %s", err)
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(p.pipelineConfig.ModulePullRetryPeriod):
				continue
			}
		}

		wg.Add(len(cursor.Data))
		for _, log := range cursor.Data {
			go func() {
				defer wg.Done()
				for {
					itemsErrors, err := p.connector.Accept(ctx, ingester.LogWithModule{
						Log:    log,
						Module: p.pipeline.Module,
					})
					if err == nil {
						err = itemsErrors[0]
					}
					if err != nil {
						p.logger.Errorf("Error pushing data on connector: %s", err)
						select {
						case <-ctx.Done():
							return
						case <-time.After(p.pipelineConfig.ConnectorPushRetryPeriod):
							continue
						}
					}
					break
				}
			}()
		}

		wg.Wait()

		if !cursor.HasMore {
			break
		}

		p.expectedState.Signal(ingester.NewInitStateWithCursor(cursor.Next))
	}

	p.expectedState.Signal(ingester.NewReadyState())
	return nil
}

func (p *PipelineHandler) pause(ctx context.Context, ready chan struct{}) error {
	close(ready)

	select {
	case <-ctx.Done():
		return nil
	}
}

func (p *PipelineHandler) stop(ctx context.Context, ready chan struct{}) error {
	close(ready)

	select {
	case <-ctx.Done():
		return nil
	}
}

func (p *PipelineHandler) Shutdown(ctx context.Context) error {
	p.logger.Infof("shutdowning pipeline")
	errorChannel := make(chan error, 1)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.stopChannel <- errorChannel:
		p.logger.Debugf("shutdowning pipeline signal sent")
		select {
		case err := <-errorChannel:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (p *PipelineHandler) GetActiveState() *Signal[ingester.State] {
	if p == nil {
		return nil
	}
	return p.activeState
}

func NewPipelineHandler(
	pipeline ingester.Pipeline,
	module modules.Module,
	connector drivers.Driver,
	logger logging.Logger,
	opts ...PipelineOption,
) *PipelineHandler {
	config := PipelineHandlerConfig{}
	for _, opt := range append(defaultPipelineOptions, opts...) {
		opt(&config)
	}

	return &PipelineHandler{
		pipeline:       pipeline,
		stopChannel:    make(chan chan error, 1),
		module:         module,
		connector:      connector,
		expectedState:  NewSignal(&pipeline.State),
		activeState:    NewSignal[ingester.State](nil),
		pipelineConfig: config,
		logger: logger.
			WithField("component", "pipeline").
			WithField("module", pipeline.Module).
			WithField("connector", pipeline.ConnectorID),
	}
}
