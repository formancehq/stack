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

var _ = Given("some empty environment", func() {
	When("setting metadata on a unknown account", func() {
		var (
			msgs               chan *nats.Msg
			cancelSubscription func()
			metadata           = map[string]interface{}{
				"clientType": "gold",
			}
		)
		BeforeEach(func() {
			// Subscribe to nats subject
			cancelSubscription, msgs = SubscribeLedger()

			// Create a transaction
			_, err := Client().AccountsApi.
				AddMetadataToAccount(TestContext(), "default", "foo").
				RequestBody(metadata).
				Execute()
			Expect(err).To(BeNil())
		})
		AfterEach(func() {
			cancelSubscription()
		})

		It("should trigger a new event", func() {
			// Wait for created transaction event
			// TODO: Check events content against schema
			msg := WaitOnChanWithTimeout(msgs, 5*time.Second)
			event := &bus.EventMessage{}
			Expect(json.Unmarshal(msg.Data, event)).To(BeNil())
		})
		It("should pop an account with the correct metadata on search service", func() {
			Eventually(func(g Gomega) bool {
				res, _, err := Client().SearchApi.Search(TestContext()).Query(formance.Query{
					Target: formance.PtrString("ACCOUNT"),
				}).Execute()
				g.Expect(err).To(BeNil())
				g.Expect(res.Cursor.Data).To(HaveLen(1))
				g.Expect(res.Cursor.Data[0]).To(Equal(map[string]any{
					"ledger":   "default",
					"metadata": metadata,
					"address":  "foo",
				}))

				return true
			}).Should(BeTrue())
		})
	})
})
