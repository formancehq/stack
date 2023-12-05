package internal

import (
	"encoding/json"
	"fmt"
	"github.com/formancehq/stack/libs/go-libs/publish"
	. "github.com/onsi/gomega"
)

func PublishPayments(message publish.EventMessage) {
	data, err := json.Marshal(message)
	Expect(err).WithOffset(1).To(BeNil())

	err = NatsClient().Publish(fmt.Sprintf("%s-payments", actualTestID), data)
	Expect(err).To(BeNil())
}
