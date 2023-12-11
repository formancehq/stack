package suite

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	paymentsevents "github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/go-libs/publish"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/big"
	"net/http"
	"strings"
	"time"
)

var _ = WithModules([]*Module{modules.Auth, modules.Orchestration, modules.Ledger}, func() {
	BeforeEach(func() {
		createLedgerResponse, err := Client().Ledger.V2CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
			Ledger: "default",
		})
		Expect(err).To(BeNil())
		Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))
	})
	When("creating a new workflow and a trigger on payments creation", func() {
		var (
			createTriggerResponse *operations.CreateTriggerResponse
		)
		BeforeEach(func() {
			response, err := Client().Orchestration.CreateWorkflow(
				TestContext(),
				&shared.CreateWorkflowRequest{
					Name: ptr(uuid.NewString()),
					Stages: []map[string]interface{}{
						{
							"send": map[string]any{
								"source": map[string]any{
									"account": map[string]any{
										"id": "world",
									},
								},
								"destination": map[string]any{
									"account": map[string]any{
										"id": "${account}",
									},
								},
								"amount": map[string]any{
									"amount": "${amount}",
									"asset":  "${asset}",
								},
							},
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(201))

			createTriggerResponse, err = Client().Orchestration.CreateTrigger(
				TestContext(),
				&shared.TriggerData{
					Event:      paymentsevents.EventTypeSavedPayments,
					WorkflowID: response.CreateWorkflowResponse.Data.ID,
					Vars: map[string]any{
						"account": `replace(event.id, "-", "_")`,
						"amount":  "event.amount",
						"asset":   "event.asset",
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
		})
		It("should be ok and the trigger should be available on list", func() {
			Expect(createTriggerResponse.StatusCode).To(Equal(201))
			Expect(createTriggerResponse.CreateTriggerResponse.Data.ID).NotTo(BeEmpty())

			listTriggersResponse, err := Client().Orchestration.ListTriggers(TestContext())
			Expect(err).To(BeNil())
			Expect(listTriggersResponse.ListTriggersResponse.Data).Should(HaveLen(1))
		})
		Then("publishing a new payments in the event bus", func() {
			var payment map[string]any
			BeforeEach(func() {
				payment = map[string]any{
					"amount": 100,
					"asset":  "USD/2",
					"id":     uuid.NewString(),
				}
				PublishPayments(publish.EventMessage{
					Date:    time.Now(),
					App:     "payments",
					Type:    paymentsevents.EventTypeSavedPayments,
					Payload: payment,
				})
			})
			It("Should trigger the workflow", func() {
				var (
					listTriggersOccurrencesResponse *operations.ListTriggersOccurrencesResponse
					err                             error
				)
				Eventually(func(g Gomega) []shared.TriggerOccurrence {
					listTriggersOccurrencesResponse, err = Client().Orchestration.ListTriggersOccurrences(TestContext(), operations.ListTriggersOccurrencesRequest{
						TriggerID: createTriggerResponse.CreateTriggerResponse.Data.ID,
					})
					g.Expect(err).To(BeNil())
					return listTriggersOccurrencesResponse.ListTriggersOccurrencesResponse.Data
				}).ShouldNot(BeEmpty())

				var getInstanceResponse *operations.GetInstanceResponse
				Eventually(func() bool {
					getInstanceResponse, err = Client().Orchestration.GetInstance(TestContext(), operations.GetInstanceRequest{
						InstanceID: listTriggersOccurrencesResponse.ListTriggersOccurrencesResponse.Data[0].WorkflowInstanceID,
					})
					Expect(err).To(BeNil())

					return getInstanceResponse.GetWorkflowInstanceResponse.Data.Terminated
				}).Should(BeTrue())

				Expect(*getInstanceResponse.GetWorkflowInstanceResponse.Data.Error).To(BeEmpty())

				listTransactionsResponse, err := Client().Ledger.V2ListTransactions(TestContext(), operations.V2ListTransactionsRequest{
					Ledger: "default",
				})
				Expect(err).To(BeNil())
				Expect(listTransactionsResponse.V2TransactionsCursorResponse.Cursor.Data).To(HaveLen(1))
				Expect(listTransactionsResponse.V2TransactionsCursorResponse.Cursor.Data[0].Postings).To(HaveLen(1))
				Expect(listTransactionsResponse.V2TransactionsCursorResponse.Cursor.Data[0].Postings[0].Source).
					To(Equal("world"))
				Expect(listTransactionsResponse.V2TransactionsCursorResponse.Cursor.Data[0].Postings[0].Destination).
					To(Equal(strings.Replace(payment["id"].(string), "-", "_", -1)))
				Expect(listTransactionsResponse.V2TransactionsCursorResponse.Cursor.Data[0].Postings[0].Asset).
					To(Equal("USD/2"))
				Expect(listTransactionsResponse.V2TransactionsCursorResponse.Cursor.Data[0].Postings[0].Amount).
					To(Equal(big.NewInt(100)))
			})
		})
		Then("deleting the trigger", func() {
			BeforeEach(func() {
				_, err := Client().Orchestration.DeleteTrigger(TestContext(), operations.DeleteTriggerRequest{
					TriggerID: createTriggerResponse.CreateTriggerResponse.Data.ID,
				})
				Expect(err).To(BeNil())
			})
			It("should not appear on list", func() {
				listTriggersResponse, err := Client().Orchestration.ListTriggers(TestContext())
				Expect(err).To(BeNil())
				Expect(listTriggersResponse.ListTriggersResponse.Data).Should(HaveLen(0))
			})
		})
	})
})
