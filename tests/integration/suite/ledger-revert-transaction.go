package suite

import (
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/ledger/pkg/bus"
	"github.com/formancehq/ledger/pkg/core"
	"github.com/formancehq/stack/libs/events"
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
			createTransactionResponse *shared.CreateTransactionResponse
		)
		BeforeEach(func() {
			// Subscribe to nats subject
			cancelSubscription, msgs = SubscribeLedger()
			_ = msgs

			// Create a transaction
			response, err := Client().Ledger.CreateTransaction(
				TestContext(),
				operations.CreateTransactionRequest{
					PostTransaction: shared.PostTransaction{
						Metadata: map[string]string{},
						Postings: []shared.Posting{
							{
								Amount:      100,
								Asset:       "USD",
								Source:      "world",
								Destination: "alice",
							},
						},
						Timestamp: &timestamp,
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			createTransactionResponse = response.CreateTransactionResponse

			// Wait for created transaction event to drain events
			WaitOnChanWithTimeout(msgs, 5*time.Second)
		})
		AfterEach(func() {
			cancelSubscription()
		})
		Then("reverting it", func() {
			BeforeEach(func() {
				response, err := Client().Ledger.RevertTransaction(
					TestContext(),
					operations.RevertTransactionRequest{
						Ledger: "default",
						Txid:   createTransactionResponse.Data.Txid,
					},
				)
				Expect(err).To(Succeed())
				Expect(response.StatusCode).To(Equal(201))
			})
			It("should trigger a new event", func() {
				// Wait for created transaction event
				msg := WaitOnChanWithTimeout(msgs, 5*time.Second)
				Expect(events.Check(msg.Data, "ledger", bus.EventTypeRevertedTransaction)).Should(Succeed())
			})
			It("should set a metadata on the original transaction", func() {
				response, err := Client().Ledger.GetTransaction(
					TestContext(),
					operations.GetTransactionRequest{
						Ledger: "default",
						Txid:   createTransactionResponse.Data.Txid,
					},
				)
				Expect(err).To(Succeed())
				Expect(response.StatusCode).To(Equal(200))

				Expect(core.IsReverted(response.GetTransactionResponse.Data.Metadata)).To(BeTrue())
			})
			Then("trying to revert again", func() {
				It("should be rejected", func() {
					response, err := Client().Ledger.RevertTransaction(
						TestContext(),
						operations.RevertTransactionRequest{
							Ledger: "default",
							Txid:   createTransactionResponse.Data.Txid,
						},
					)
					Expect(err).To(BeNil())
					Expect(response.StatusCode).To(Equal(400))
				})
			})
		})
	})
})
