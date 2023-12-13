package internal

import (
	"encoding/json"
	"fmt"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/gomega"
)

func NatsClient() *nats.Conn {
	natsConn, err := nats.Connect(GetNatsAddress())
	Expect(err).ToNot(HaveOccurred())
	return natsConn
}

func SubscribeLedger() (func(), chan *nats.Msg) {
	msgs := make(chan *nats.Msg)
	subscription, err := NatsClient().Subscribe(fmt.Sprintf("%s-ledger", currentTest.id), func(msg *nats.Msg) {
		msgs <- msg
	})
	Expect(err).ToNot(HaveOccurred())
	return func() {
		Expect(subscription.Unsubscribe()).To(Succeed())
	}, msgs
}

func SubscribePayments() (func(), chan *nats.Msg) {
	msgs := make(chan *nats.Msg)
	subscription, err := NatsClient().Subscribe(fmt.Sprintf("%s-payments", currentTest.id), func(msg *nats.Msg) {
		msgs <- msg
	})
	Expect(err).ToNot(HaveOccurred())
	return func() {
		Expect(subscription.Unsubscribe()).To(Succeed())
	}, msgs
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
