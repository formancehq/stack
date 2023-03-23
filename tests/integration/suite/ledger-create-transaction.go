package suite_test

import (
	"time"

	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/ledger/pkg/bus"
	"github.com/formancehq/stack/libs/events"
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
			err                error
			rsp                *formance.TransactionResponse
		)
		BeforeEach(func() {
			// Subscribe to nats subject
			cancelSubscription, msgs = SubscribeLedger()

			// Create a transaction
			rsp, _, err = Client().TransactionsApi.
				CreateTransaction(TestContext(), "default").
				PostTransaction(formance.PostTransaction{
					Timestamp: &timestamp,
					Postings: []formance.Posting{{
						Amount:      100,
						Asset:       "USD",
						Source:      "world",
						Destination: "alice",
					}},
				}).
				Execute()
			Expect(err).To(BeNil())
		})
		AfterEach(func() {
			cancelSubscription()
		})
		It("should eventually be available on api", func() {
			transactionResponse, _, err := Client().TransactionsApi.
				GetTransaction(TestContext(), "default", rsp.Data.Txid).
				Execute()
			Expect(err).To(BeNil())
			Expect(transactionResponse.Data).To(Equal(rsp.Data))

			accountResponse, _, err := Client().AccountsApi.
				GetAccount(TestContext(), "default", "alice").
				Execute()
			Expect(err).To(BeNil())
			Expect(accountResponse.Data).Should(Equal(formance.AccountWithVolumesAndBalances{
				Address:  "alice",
				Metadata: map[string]interface{}{},
				Volumes: ptr(map[string]map[string]int64{
					"USD": {
						"input":   100,
						"output":  0,
						"balance": 100,
					},
				}),
				Balances: ptr(map[string]int64{
					"USD": 100,
				}),
			}))
		})
		It("should trigger a new event", func() {
			// Wait for created transaction event
			msg := WaitOnChanWithTimeout(msgs, 5*time.Second)
			Expect(events.Check(msg.Data, "ledger", bus.EventTypeCommittedTransactions)).Should(BeNil())
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
				res, _, err := Client().SearchApi.Search(TestContext()).Query(formance.Query{
					Target: formance.PtrString("TRANSACTION"),
				}).Execute()
				g.Expect(err).To(BeNil())
				g.Expect(res.Cursor.Data).To(HaveLen(1))
				g.Expect(res.Cursor.Data[0]).To(Equal(expectedTx))

				return true
			}).Should(BeTrue())

			Eventually(func(g Gomega) []any {
				res, _, err := Client().SearchApi.Search(TestContext()).Query(formance.Query{
					Target: formance.PtrString("TRANSACTION"),
					Terms:  []string{"alice"},
				}).Execute()
				g.Expect(err).To(BeNil())
				g.Expect(res.Cursor.Data[0]).To(Equal(expectedTx))
				return res.Cursor.Data
			}).Should(HaveLen(1))

			Eventually(func(g Gomega) bool {
				res, _, err := Client().SearchApi.Search(TestContext()).Query(formance.Query{
					Target: formance.PtrString("ACCOUNT"),
				}).Execute()
				g.Expect(err).To(BeNil())
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
				res, _, err := Client().SearchApi.Search(TestContext()).Query(formance.Query{
					Target: formance.PtrString("ASSET"),
				}).Execute()
				g.Expect(err).To(BeNil())
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
			resp, httpResp, err := Client().TransactionsApi.
				CreateTransaction(TestContext(), "default").
				PostTransaction(formance.PostTransaction{
					Postings: []formance.Posting{{
						Amount:      100,
						Asset:       "USD",
						Source:      "bob",
						Destination: "alice",
					}},
				}).Execute()
			Expect(err).To(HaveOccurred())
			Expect(httpResp.StatusCode).To(Equal(400))
			Expect(resp).To(BeNil())

			apiErr, ok := err.(GenericOpenAPIError)
			Expect(ok).To(BeTrue())

			details := "https://play.numscript.org/?payload=eyJlcnJvciI6ImFjY291bnQgaGFkIGluc3VmZmljaWVudCBmdW5kcyJ9"
			Expect(apiErr.Model()).Should(Equal(formance.ErrorResponse{
				ErrorCode:    formance.INSUFFICIENT_FUND,
				ErrorMessage: "[INSUFFICIENT_FUND] account had insufficient funds",
				Details:      &details,
			}))
		})
	})
})
