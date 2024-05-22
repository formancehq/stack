package worker

import (
	"strings"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.uber.org/fx"
)

func Module(topics []string) fx.Option {
	return fx.Options(
		fx.Invoke(func(logger logging.Logger, r *message.Router, s message.Subscriber) {
			logger.Infof("Listening events from topics: %s", strings.Join(topics, ","))
			registerListener(logger, r, s, topics)
		}),
	)
}
