package worker

import (
	"fmt"
	"runtime/debug"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func handleMessage(msg *message.Message) error {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
			debug.PrintStack()
		}
	}()

	return nil
}

func registerListener(logger logging.Logger, r *message.Router, s message.Subscriber, topics []string) {
	for _, topic := range topics {
		r.AddNoPublisherHandler(fmt.Sprintf("reco-listen-%s-events", topic), topic, s, func(msg *message.Message) error {
			if err := handleMessage(msg); err != nil {
				logger.Errorf("error handling message: %s", err)
				return err
			}
			return nil
		})
	}
}
