package suite

import (
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/ledger/pkg/bus"
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
			msgs               chan *nats.Msg
			cancelSubscription func()
			timestamp          = time.Now().Round(time.Second).UTC()
			rsp                *shared.CreateTransactionResponse
		)
		BeforeEach(func() {
			// Subscribe to nats subject
			cancelSubscription, msgs = SubscribeLedger()

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

			rsp = response.CreateTransactionResponse
		})
		AfterEach(func() {
			cancelSubscription()
		})
		It("should be available on api", func() {
			response, err := Client().Ledger.GetTransaction(
				TestContext(),
				operations.GetTransactionRequest{
					Ledger: "default",
					Txid:   rsp.Data.Txid,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			transactionResponse := response.GetTransactionResponse
			Expect(transactionResponse.Data).To(Equal(shared.ExpandedTransaction{
				Timestamp: rsp.Data.Timestamp,
				Postings:  rsp.Data.Postings,
				Reference: rsp.Data.Reference,
				Metadata:  rsp.Data.Metadata,
				Txid:      rsp.Data.Txid,
				PreCommitVolumes: map[string]map[string]shared.Volume{
					"world": {
						"USD": {
							Input:   0,
							Output:  0,
							Balance: ptr(int64(0)),
						},
					},
					"alice": {
						"USD": {
							Input:   0,
							Output:  0,
							Balance: ptr(int64(0)),
						},
					},
				},
				PostCommitVolumes: map[string]map[string]shared.Volume{
					"world": {
						"USD": {
							Input:   0,
							Output:  100,
							Balance: ptr(int64(-100)),
						},
					},
					"alice": {
						"USD": {
							Input:   100,
							Output:  0,
							Balance: ptr(int64(100)),
						},
					},
				},
			}))

			accResponse, err := Client().Ledger.GetAccount(
				TestContext(),
				operations.GetAccountRequest{
					Address: "alice",
					Ledger:  "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(accResponse.StatusCode).To(Equal(200))

			accountResponse := accResponse.AccountResponse
			Expect(accountResponse.Data).Should(Equal(shared.AccountWithVolumesAndBalances{
				Address:  "alice",
				Metadata: metadata.Metadata{},
				Volumes: map[string]map[string]int64{
					"USD": {
						"input":   100,
						"output":  0,
						"balance": 100,
					},
				},
				Balances: map[string]int64{
					"USD": 100,
				},
			}))
		})
		It("should trigger a new event", func() {
			// Wait for created transaction event
			msg := WaitOnChanWithTimeout(msgs, 5*time.Second)
			Expect(events.Check(msg.Data, "ledger", bus.EventTypeCommittedTransactions)).Should(Succeed())
		})
		It("should pop a transaction, two accounts and two assets entries on search service", func() {
			expectedTx := map[string]any{
				"reference": "",
				"metadata":  map[string]any{},
				"postings": []any{
					map[string]any{
						"source":      "world",
						"asset":       "USD",
						"amount":      float64(100),
						"destination": "alice",
					},
				},
				"txid":      float64(0),
				"timestamp": timestamp.Format(time.RFC3339),
				"ledger":    "default",
			}
			Eventually(func(g Gomega) bool {
				response, err := Client().Search.Search(
					TestContext(),
					shared.Query{
						Target: ptr("TRANSACTION"),
					},
				)
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(response.StatusCode).To(Equal(200))

				res := response.Response
				g.Expect(res.Cursor.Data).To(HaveLen(1))
				g.Expect(res.Cursor.Data[0]).To(Equal(expectedTx))

				return true
			}).Should(BeTrue())

			Eventually(func(g Gomega) []map[string]any {
				response, err := Client().Search.Search(
					TestContext(),
					shared.Query{
						Target: ptr("TRANSACTION"),
						Terms:  []string{"alice"},
					},
				)
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(response.StatusCode).To(Equal(200))

				res := response.Response
				g.Expect(res.Cursor.Data[0]).To(Equal(expectedTx))
				return res.Cursor.Data
			}).Should(HaveLen(1))

			Eventually(func(g Gomega) bool {
				response, err := Client().Search.Search(
					TestContext(),
					shared.Query{
						Target: ptr("ACCOUNT"),
					},
				)
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(response.StatusCode).To(Equal(200))

				res := response.Response
				g.Expect(res.Cursor.Data).To(HaveLen(2))
				g.Expect(res.Cursor.Data).To(ContainElements(
					map[string]any{
						"address": "world",
						"ledger":  "default",
					},
					map[string]any{
						"address": "alice",
						"ledger":  "default",
					},
				))
				return true
			}).Should(BeTrue())

			Eventually(func(g Gomega) bool {
				response, err := Client().Search.Search(
					TestContext(),
					shared.Query{
						Target: ptr("ASSET"),
					},
				)
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(response.StatusCode).To(Equal(200))

				res := response.Response
				g.Expect(res.Cursor.Data).To(HaveLen(2))
				g.Expect(res.Cursor.Data).To(ContainElements(
					map[string]any{
						"account": "world",
						"ledger":  "default",
						"output":  float64(100),
						"input":   float64(0),
						"name":    "USD",
					},
					map[string]any{
						"account": "alice",
						"ledger":  "default",
						"output":  float64(0),
						"input":   float64(100),
						"name":    "USD",
					},
				))
				return true
			}).Should(BeTrue())
		})
	})
})

type GenericOpenAPIError interface {
	Model() any
}

var _ = Given("some empty environment", func() {
	When("creating a transaction on a ledger with insufficient funds", func() {
		It("should fail", func() {
			response, err := Client().Ledger.CreateTransaction(
				TestContext(),
				operations.CreateTransactionRequest{
					PostTransaction: shared.PostTransaction{
						Metadata: map[string]string{},
						Postings: []shared.Posting{
							{
								Amount:      100,
								Asset:       "USD",
								Source:      "bob",
								Destination: "alice",
							},
						},
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(400))
			Expect(response.CreateTransactionResponse).To(BeNil())

			details := "https://play.numscript.org/?payload=eyJlcnJvciI6ImFjY291bnQgaGFkIGluc3VmZmljaWVudCBmdW5kcyJ9"
			Expect(response.ErrorResponse).Should(Equal(&shared.ErrorResponse{
				ErrorCode:    shared.ErrorsEnumInsufficientFund,
				ErrorMessage: "account had insufficient funds",
				Details:      &details,
			}))
		})
	})
})

var _ = Given("some empty environment", func() {
	When("creating a transaction on a ledger with an idempotency key", func() {
		var (
			err      error
			response *operations.CreateTransactionResponse
		)
		createTransaction := func() {
			response, err = Client().Ledger.CreateTransaction(
				TestContext(),
				operations.CreateTransactionRequest{
					IdempotencyKey: ptr("testing"),
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
					},
					Ledger: "default",
				},
			)
		}
		BeforeEach(createTransaction)
		It("should be ok", func() {
			Expect(err).To(Succeed())
			Expect(response.CreateTransactionResponse.Data.Txid).To(Equal(int64(0)))
		})
		Then("replaying with the same IK", func() {
			BeforeEach(createTransaction)
			It("should respond with the same tx id", func() {
				Expect(err).To(Succeed())
				Expect(response.CreateTransactionResponse.Data.Txid).To(Equal(int64(0)))
			})
		})
	})
})
