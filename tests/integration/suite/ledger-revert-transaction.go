package suite

import (
	"time"

	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/ledger/pkg/bus"
	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/stack/libs/events"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("some empty environment", func() {
	When("creating a transaction on a ledger", func() {
		var (
			msgs                      chan *nats.Msg
			cancelSubscription        func()
			timestamp                 = time.Now().Round(time.Second).UTC()
			err                       error
			createTransactionResponse *formance.TransactionResponse
		)
		BeforeEach(func() {
			// Subscribe to nats subject
			cancelSubscription, msgs = SubscribeLedger()
			_ = msgs

			// Create a transaction
			createTransactionResponse, _, err = Client().TransactionsApi.
				CreateTransaction(TestContext(), "default").
				PostTransaction(formance.PostTransaction{
					Timestamp: &timestamp,
					Postings: []formance.Posting{{
						Amount:      100,
						Asset:       "USD",
						Source:      "world",
						Destination: "alice",
					}},
					Metadata: metadata.Metadata{},
				}).
				Execute()
			Expect(err).ToNot(HaveOccurred())

			// Wait for created transaction event to drain events
			WaitOnChanWithTimeout(msgs, 5*time.Second)
		})
		AfterEach(func() {
			cancelSubscription()
		})
		Then("reverting it", func() {
			BeforeEach(func() {
				_, _, err = Client().TransactionsApi.
					RevertTransaction(TestContext(), "default", createTransactionResponse.Data.Txid).
					Execute()
				Expect(err).To(Succeed())
			})
			It("should trigger a new event", func() {
				// Wait for created transaction event
				msg := WaitOnChanWithTimeout(msgs, 5*time.Second)
				Expect(events.Check(msg.Data, "ledger", bus.EventTypeRevertedTransaction)).Should(Succeed())
			})
			It("should set a metadata on the original transaction", func() {
				rsp, _, err := Client().TransactionsApi.
					GetTransaction(TestContext(), "default", createTransactionResponse.Data.Txid).
					Execute()
				Expect(err).To(Succeed())
				Expect(core.IsReverted(rsp.Data.Metadata)).To(BeTrue())
			})
			Then("trying to revert again", func() {
				It("should be rejected", func() {
					_, _, err = Client().TransactionsApi.
						RevertTransaction(TestContext(), "default", createTransactionResponse.Data.Txid).
						Execute()
					Expect(err).NotTo(BeNil())
				})
			})
		})
	})
})
