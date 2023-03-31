package internal

import (
	"fmt"

	"github.com/nats-io/nats.go"
	. "github.com/onsi/gomega"
)

func natsPort() string {
	return "4222" // TODO: Make configurable
}

func natsAddress() string {
	return "localhost:4222" // TODO: Make configurable
}

func NatsClient() *nats.Conn {
	natsConn, err := nats.Connect(natsAddress())
	Expect(err).ToNot(HaveOccurred())
	return natsConn
}

func SubscribeLedger() (func(), chan *nats.Msg) {
	msgs := make(chan *nats.Msg)
	subscription, err := NatsClient().Subscribe(fmt.Sprintf("%s-ledger", actualTestID), func(msg *nats.Msg) {
		msgs <- msg
	})
	Expect(err).ToNot(HaveOccurred())
	return func() {
		Expect(subscription.Unsubscribe()).To(Succeed())
	}, msgs
}

func SubscribePayments() (func(), chan *nats.Msg) {
	msgs := make(chan *nats.Msg)
	subscription, err := NatsClient().Subscribe(fmt.Sprintf("%s-payments", actualTestID), func(msg *nats.Msg) {
		msgs <- msg
	})
	Expect(err).ToNot(HaveOccurred())
	return func() {
		Expect(subscription.Unsubscribe()).To(Succeed())
	}, msgs
}
