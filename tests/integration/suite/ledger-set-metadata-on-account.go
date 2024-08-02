package suite

import (
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"math/big"
	"net/http"
	"reflect"
	"time"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	ledgerevents "github.com/formancehq/ledger/pkg/events"
	"github.com/formancehq/stack/libs/events"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Search, modules.Ledger}, func() {
	BeforeEach(func() {
		createLedgerResponse, err := Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
			Ledger: "default",
		})
		Expect(err).To(BeNil())
		Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))
	})
	When("creating a new transaction", func() {
		BeforeEach(func() {
			// Create a transaction
			response, err := Client().Ledger.V2.CreateTransaction(
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
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			Eventually(func() bool {
				response, err := Client().Search.V1.Search(
					TestContext(),
					shared.Query{
						Target: ptr("ACCOUNT"),
						Terms:  []string{"address=alice"},
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
					"ledger":   "default",
					"address":  "alice",
					"metadata": map[string]any{},
				})
			}).Should(BeTrue())
		})
		Then("setting a metadata on the destination account", func() {
			BeforeEach(func() {
				response, err := Client().Ledger.V2.AddMetadataToAccount(
					TestContext(),
					operations.V2AddMetadataToAccountRequest{
						RequestBody: map[string]string{
							"foo": "bar",
						},
						Address: "alice",
						Ledger:  "default",
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(204))
			})
			It("should be ok", func() {
				Eventually(func() bool {
					response, err := Client().Search.V1.Search(
						TestContext(),
						shared.Query{
							Target: ptr("ACCOUNT"),
							Terms:  []string{"address=alice"},
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
							"foo": "bar",
						},
						"address": "alice",
					})
				}).Should(BeTrue())
			})
		})
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

			response, err := Client().Ledger.V2.AddMetadataToAccount(
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
			response, err := Client().Ledger.V2.GetAccount(
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
				response, err := Client().Search.V1.Search(
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
