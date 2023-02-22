package suite

import (
	"encoding/json"
	"time"

	"github.com/formancehq/formance-sdk-go"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/nats-io/nats.go"
	"github.com/numary/ledger/pkg/bus"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func waitOnChanWithTimeout[T any](ch chan T, timeout time.Duration) T {
	select {
	case t := <-ch:
		return t
	case <-time.After(timeout):
		Fail("should have received a created transaction event")
	}
	panic("cannot happen")
}

var _ = Given("some empty environment", func() {
	When("creating a transaction on a ledger", func() {
		var (
			msgs               chan *nats.Msg
			cancelSubscription func()
			timestamp          = time.Now().Round(time.Second).UTC()
		)
		BeforeEach(func() {
			// Subscribe to nats subject
			cancelSubscription, msgs = SubscribeLedger()

			// Create a transaction
			_, _, err := Client().TransactionsApi.
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
		It("should trigger a new event", func() {
			// Wait for created transaction event
			msg := waitOnChanWithTimeout(msgs, 5*time.Second)
			event := &bus.EventMessage{}
			Expect(json.Unmarshal(msg.Data, event)).To(BeNil())
		})
		It("should pop a transaction on search service", func() {
			Eventually(func(g Gomega) bool {
				res, _, err := Client().SearchApi.Search(TestContext()).Query(formance.Query{
					Target: formance.PtrString("TRANSACTION"),
				}).Execute()
				g.Expect(err).To(BeNil())
				g.Expect(res.Cursor.Data).To(HaveLen(1))
				g.Expect(res.Cursor.Data[0]).To(Equal(map[string]any{
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
				}))

				return true
			}).Should(BeTrue())

			Eventually(func(g Gomega) []any {
				res, _, err := Client().SearchApi.Search(TestContext()).Query(formance.Query{
					Target: formance.PtrString("TRANSACTION"),
					Terms:  []string{"alice"},
				}).Execute()
				g.Expect(err).To(BeNil())
				g.Expect(res.Cursor.Data[0]).To(Equal(map[string]any{
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
				}))
				return res.Cursor.Data
			}).Should(HaveLen(1))
		})
	})
})
