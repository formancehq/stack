package suite

import (
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"net/http"
	"reflect"
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	ledgerevents "github.com/formancehq/ledger/pkg/events"
	"github.com/formancehq/stack/libs/events"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Search, modules.Ledger}, func() {
	BeforeEach(func() {
		createLedgerResponse, err := Client().Ledger.V2CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
			Ledger: "default",
		})
		Expect(err).To(BeNil())
		Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))
	})
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

			response, err := Client().Ledger.V2AddMetadataToAccount(
				TestContext(),
				operations.V2AddMetadataToAccountRequest{
					RequestBody: metadata,
					Address:     "foo",
					Ledger:      "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
		})
		AfterEach(func() {
			cancelSubscription()
		})
		It("should be available on api", func() {
			response, err := Client().Ledger.V2GetAccount(
				TestContext(),
				operations.V2GetAccountRequest{
					Address: "foo",
					Ledger:  "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			Expect(response.V2AccountResponse.Data).Should(Equal(shared.V2Account{
				Address:  "foo",
				Metadata: metadata,
			}))
		})
		It("should trigger a new event", func() {
			msg := WaitOnChanWithTimeout(msgs, 5*time.Second)
			Expect(events.Check(msg.Data, "ledger", ledgerevents.EventTypeSavedMetadata)).Should(Succeed())
		})
		It("should pop an account with the correct metadata on search service", func() {
			Eventually(func() bool {
				response, err := Client().Search.Search(
					TestContext(),
					shared.Query{
						Target: ptr("ACCOUNT"),
					},
				)
				if err != nil {
					return false
				}
				if response.StatusCode != 200 {
					return false
				}
				if len(response.Response.Cursor.Data) != 1 {
					return false
				}
				return reflect.DeepEqual(response.Response.Cursor.Data[0], map[string]any{
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
