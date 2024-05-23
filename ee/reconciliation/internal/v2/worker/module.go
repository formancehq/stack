package worker

import (
	"context"
	"strings"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.uber.org/fx"
)

func Module(
	topics []string,
	transactionBasedDelay time.Duration,
) fx.Option {
	return fx.Options(
		fx.Invoke(func(
			lc fx.Lifecycle,
			logger logging.Logger,
			l *Listener,
			r *message.Router,
			s message.Subscriber,
		) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error {
					logger.Infof("Listening events from topics: %s", strings.Join(topics, ","))
					l.registerListener(logger, r, s, topics)
					return nil
				},
			})
			logger.Infof("Listening events from topics: %s", strings.Join(topics, ","))
			l.registerListener(logger, r, s, topics)
		}),
	)
}
