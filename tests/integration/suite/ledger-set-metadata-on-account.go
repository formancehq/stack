package suite

import (
	"reflect"
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
			metadata           = map[string]string{
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
			Expect(err).ToNot(HaveOccurred())
		})
		AfterEach(func() {
			cancelSubscription()
		})
		It("should be available on api", func() {
			accountResponse, _, err := Client().AccountsApi.
				GetAccount(TestContext(), "default", "foo").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(accountResponse.Data).Should(Equal(formance.AccountWithVolumesAndBalances{
				Address:  "foo",
				Metadata: metadata,
				Volumes:  map[string]map[string]int64{},
				Balances: map[string]int64{},
			}))
		})
		It("should trigger a new event", func() {
			msg := WaitOnChanWithTimeout(msgs, 5*time.Second)
			Expect(events.Check(msg.Data, "ledger", bus.EventTypeSavedMetadata)).Should(Succeed())
		})
		It("should pop an account with the correct metadata on search service", func() {
			Eventually(func() bool {
				res, _, err := Client().SearchApi.Search(TestContext()).Query(formance.Query{
					Target: formance.PtrString("ACCOUNT"),
				}).Execute()
				if err != nil {
					return false
				}
				if len(res.Cursor.Data) != 1 {
					return false
				}
				return reflect.DeepEqual(res.Cursor.Data[0], map[string]any{
					"ledger": "default",
					"metadata": map[string]any{
						"clientType": "gold",
					},
					"address": "foo",
				})
			}).Should(BeTrue())
		})
	})
})
