package worker

import (
	"fmt"
	"runtime/debug"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/publish"
)

func handleMessage(msg *message.Message) error {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
			debug.PrintStack()
		}
	}()

	var event *publish.EventMessage
	span, event, err := publish.UnmarshalMessage(msg)
	if err != nil {
		logging.FromContext(msg.Context()).Error(err.Error())
		return err
	}
	_ = span

	// ctx, span := workflow.Tracer.Start(msg.Context(), "Trigger:HandleEvent",
	// 	trace.WithLinks(trace.Link{
	// 		SpanContext: span.SpanContext(),
	// 	}),
	// 	trace.WithAttributes(
	// 		attribute.String("event-id", msg.UUID),
	// 		attribute.Bool("duplicate", false),
	// 		attribute.String("event-type", event.Type),
	// 		attribute.String("event-payload", string(msg.Payload)),
	// 	),
	// )
	// defer span.End()
	// defer func() {
	// 	if err != nil {
	// 		span.RecordError(err)
	// 	}
	// }()

	switch event.Type {
	case EventTypeCommittedTransactions:
		return handleCommittedTransactions(event)
	case EventTypeRevertedTransaction:
		return handleRevertedTransactions(event)
	case EventTypeSavedPayments:
		return handleSavedPayments(event)
	default:
		return nil
	}
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
