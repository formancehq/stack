package stripe

import (
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/stripe/client"
)

type Timeline struct {
	state                  TimelineState
	firstIDAfterStartingAt string
	startingAt             time.Time
	config                 TimelineConfig
	client                 client.Client
}

func NewTimeline(client client.Client, cfg TimelineConfig, state TimelineState, options ...TimelineOption) *Timeline {
	defaultOptions := make([]TimelineOption, 0)

	c := &Timeline{
		config: cfg,
		state:  state,
		client: client,
	}

	options = append(defaultOptions, append([]TimelineOption{
		WithStartingAt(time.Now()),
	}, options...)...)

	for _, opt := range options {
		opt.apply(c)
	}

	return c
}

type TimelineOption interface {
	apply(c *Timeline)
}
type TimelineOptionFn func(c *Timeline)

func (fn TimelineOptionFn) apply(c *Timeline) {
	fn(c)
}

func WithStartingAt(v time.Time) TimelineOptionFn {
	return func(c *Timeline) {
		c.startingAt = v
	}
}

func (tl *Timeline) State() TimelineState {
	return tl.state
}
