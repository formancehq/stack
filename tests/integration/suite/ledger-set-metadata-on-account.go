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

			_, err := Client().AccountsApi.
				AddMetadataToAccount(TestContext(), "default", "foo").
				RequestBody(metadata).
				Execute()
			Expect(err).To(BeNil())
		})
		AfterEach(func() {
			cancelSubscription()
		})
		It("should eventually be available on api", func() {
			accountResponse, _, err := Client().AccountsApi.
				GetAccount(TestContext(), "default", "foo").
				Execute()
			Expect(err).To(BeNil())
			Expect(accountResponse.Data).Should(Equal(formance.AccountWithVolumesAndBalances{
				Address:  "foo",
				Metadata: metadata,
				Volumes:  ptr(map[string]map[string]int64{}),
				Balances: ptr(map[string]int64{}),
			}))
		})
		It("should trigger a new event", func() {
			msg := WaitOnChanWithTimeout(msgs, 5*time.Second)
			Expect(events.Check(msg.Data, "ledger", bus.EventTypeSavedMetadata)).Should(BeNil())
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
