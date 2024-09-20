package internal

import (
	"encoding/json"
	"fmt"
	"github.com/formancehq/stack/libs/events"
	"github.com/formancehq/go-libs/publish"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

func NatsClient() *nats.Conn {
	natsConn, err := nats.Connect(GetNatsAddress())
	Expect(err).ToNot(HaveOccurred())
	return natsConn
}

func Subscribe(service string) (func(), chan *nats.Msg) {
	msgs := make(chan *nats.Msg)
	subscription, err := NatsClient().Subscribe(fmt.Sprintf("%s-%s", currentTest.id, service), func(msg *nats.Msg) {
		msgs <- msg
	})
	Expect(err).ToNot(HaveOccurred())
	return func() {
		Expect(subscription.Unsubscribe()).To(Succeed())
	}, msgs
}

func SubscribeLedger() (func(), chan *nats.Msg) {
	return Subscribe("ledger")
}

func SubscribePayments() (func(), chan *nats.Msg) {
	return Subscribe("payments")
}

func SubscribeOrchestration() (func(), chan *nats.Msg) {
	return Subscribe("orchestration")
}

func PublishPayments(message publish.EventMessage) {
	data, err := json.Marshal(message)
	Expect(err).WithOffset(1).To(BeNil())

	err = NatsClient().Publish(fmt.Sprintf("%s-payments", currentTest.id), data)
	Expect(err).To(BeNil())
}

func Topic(service string) string {
	return fmt.Sprintf("%s-%s", currentTest.id, service)
}

func CreateTopic(service string) {
	js, err := NatsClient().JetStream()
	Expect(err).To(Succeed())

	_, err = js.AddStream(&nats.StreamConfig{
		Name:      Topic(service),
		Subjects:  []string{Topic(service)},
		Retention: nats.LimitsPolicy,
	})
	Expect(err).To(Succeed())
}

func PublishLedger(message publish.EventMessage) {
	data, err := json.Marshal(message)
	Expect(err).WithOffset(1).To(BeNil())

	err = NatsClient().Publish(Topic("ledger"), data)
	Expect(err).To(BeNil())
}

type receiveEventMatcher struct {
	eventName   string
	serviceName string
	err         error
}

func (r *receiveEventMatcher) Match(actual interface{}) (success bool, err error) {

	var data []byte
	// Handle different types of channels and extract data accordingly.
	switch v := actual.(type) {
	case chan *nats.Msg:
		select {
		case msg := <-v:
			data = msg.Data
		default:
			return false, nil
		}
	case chan []byte:
		select {
		case msg := <-v:
			data = msg
		default:
			return false, nil
		}
	default:
		return false, fmt.Errorf("expected chan *nats.Msg or chan []uint8, got %T", actual)
	}

	r.err = events.Check(data, r.serviceName, r.eventName)
	return r.err == nil, nil
}

func (r *receiveEventMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected event to match schema: \r\n%s", r.err)
}

func (r *receiveEventMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected event not to match schema: \n%s", r.err)
}

var _ types.GomegaMatcher = (*receiveEventMatcher)(nil)

func ReceiveEvent(serviceName, eventName string) *receiveEventMatcher {
	return &receiveEventMatcher{
		eventName:   eventName,
		serviceName: serviceName,
	}
}
