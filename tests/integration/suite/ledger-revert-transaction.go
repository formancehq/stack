package suite

import (
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	ledgerevents "github.com/formancehq/ledger/pkg/events"
	"github.com/formancehq/stack/libs/events"
	. "github.com/formancehq/stack/tests/integration/internal"
)

var _ = WithModules([]*Module{modules.Ledger}, func() {
	BeforeEach(func() {
		createLedgerResponse, err := Client().Ledger.V2CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
			Ledger: "default",
		})
		Expect(err).To(BeNil())
		Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))
	})
	When("creating a transaction on a ledger", func() {
		var (
			msgs                      chan *nats.Msg
			cancelSubscription        func()
			timestamp                 = time.Now().Round(time.Second).UTC()
			createTransactionResponse *shared.V2CreateTransactionResponse
		)
		BeforeEach(func() {
			// Subscribe to nats subject
			cancelSubscription, msgs = SubscribeLedger()

			// Create a transaction
			response, err := Client().Ledger.V2CreateTransaction(
				TestContext(),
				operations.V2CreateTransactionRequest{
					V2PostTransaction: shared.V2PostTransaction{
						Metadata: map[string]string{},
						Postings: []shared.V2Posting{
							{
								Amount:      big.NewInt(100),
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

			createTransactionResponse = response.V2CreateTransactionResponse

			// Wait for created transaction event to drain events
			WaitOnChanWithTimeout(msgs, 5*time.Second)
		})
		AfterEach(func() {
			cancelSubscription()
		})
		Then("transferring funds from destination to another account", func() {
			BeforeEach(func() {
				response, err := Client().Ledger.V2CreateTransaction(
					TestContext(),
					operations.V2CreateTransactionRequest{
						V2PostTransaction: shared.V2PostTransaction{
							Metadata: map[string]string{},
							Postings: []shared.V2Posting{
								{
									Amount:      big.NewInt(100),
									Asset:       "USD",
									Source:      "alice",
									Destination: "foo",
								},
							},
							Timestamp: &timestamp,
						},
						Ledger: "default",
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))
			})
			Then("trying to revert the original transaction", func() {
				var (
					force    bool
					err      error
					response *operations.V2RevertTransactionResponse
				)
				revertTx := func() {
					response, err = Client().Ledger.V2RevertTransaction(
						TestContext(),
						operations.V2RevertTransactionRequest{
							Force:  pointer.For(force),
							ID:     big.NewInt(0),
							Ledger: "default",
						},
					)
				}
				JustBeforeEach(revertTx)
				It("Should fail", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(response.StatusCode).To(Equal(400))
					Expect(response.V2ErrorResponse.ErrorCode).To(Equal(shared.V2ErrorsEnumInsufficientFund))
				})
				Context("With forcing", func() {
					BeforeEach(func() {
						force = true
					})
					It("Should be ok", func() {
						Expect(err).ToNot(HaveOccurred())
						Expect(response.StatusCode).To(Equal(201))
					})
				})
			})
		})
		Then("reverting it", func() {
			BeforeEach(func() {
				response, err := Client().Ledger.V2RevertTransaction(
					TestContext(),
					operations.V2RevertTransactionRequest{
						Ledger: "default",
						ID:     createTransactionResponse.Data.ID,
					},
				)
				Expect(err).To(Succeed())
				Expect(response.StatusCode).To(Equal(201))
			})
			It("should trigger a new event", func() {
				// Wait for created transaction event
				msg := WaitOnChanWithTimeout(msgs, 5*time.Second)
				Expect(events.Check(msg.Data, "ledger", ledgerevents.EventTypeRevertedTransaction)).Should(Succeed())
			})
			It("should revert the original transaction", func() {
				response, err := Client().Ledger.V2GetTransaction(
					TestContext(),
					operations.V2GetTransactionRequest{
						Ledger: "default",
						ID:     createTransactionResponse.Data.ID,
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				Expect(response.V2GetTransactionResponse.Data.Reverted).To(BeTrue())
			})
			Then("trying to revert again", func() {
				It("should be rejected", func() {
					response, err := Client().Ledger.V2RevertTransaction(
						TestContext(),
						operations.V2RevertTransactionRequest{
							Ledger: "default",
							ID:     createTransactionResponse.Data.ID,
						},
					)
					Expect(err).To(BeNil())
					Expect(response.StatusCode).To(Equal(400))
					Expect(response.V2ErrorResponse.ErrorCode).To(Equal(shared.V2ErrorsEnumAlreadyRevert))
				})
			})
		})
	})
})
