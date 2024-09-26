package suite

import (
	"math/big"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/formancehq/go-libs/api"
	"github.com/formancehq/go-libs/publish"
	orchestrationevents "github.com/formancehq/orchestration/pkg/events"
	paymentsevents "github.com/formancehq/payments/pkg/events"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Auth, modules.Orchestration, modules.Ledger}, func() {
	When("creating a new workflow and a trigger on payments creation", func() {
		var (
			createTriggerResponse *operations.CreateTriggerResponse
			srv                   *httptest.Server
		)
		AfterEach(func() {
			srv.Close()
		})
		BeforeEach(func() {
			createLedgerResponse, err := Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
				Ledger: "default",
			})
			Expect(err).To(BeNil())
			Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))

			srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") == "" {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				_, _ = w.Write([]byte(`{"data": {"name": "foo"}}`))
			}))
			response, err := Client().Orchestration.V1.CreateWorkflow(
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

			createTriggerResponse, err = Client().Orchestration.V1.CreateTrigger(
				TestContext(),
				&shared.TriggerData{
					Event:      paymentsevents.EventTypeSavedPayments,
					WorkflowID: response.CreateWorkflowResponse.Data.ID,
					Vars: map[string]any{
						"account": `link(event, "destination_account").name`,
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

			listTriggersResponse, err := Client().Orchestration.V2.ListTriggers(TestContext(), operations.V2ListTriggersRequest{})
			Expect(err).To(BeNil())
			Expect(listTriggersResponse.V2ListTriggersResponse.Cursor.Data).Should(HaveLen(1))
		})
		Then("publishing a new payments in the event bus", func() {
			var (
				payment map[string]any
				msgs    chan *nats.Msg
			)
			BeforeEach(func() {
				payment = map[string]any{
					"amount": 1000000,
					"asset":  "USD/2",
					"id":     uuid.NewString(),
					"links": []api.Link{
						{
							Name: "destination_account",
							URI:  srv.URL,
						},
					},
				}
				PublishPayments(publish.EventMessage{
					Date:    time.Now(),
					App:     "payments",
					Type:    paymentsevents.EventTypeSavedPayments,
					Payload: payment,
				})
				var closeSubscription func()
				closeSubscription, msgs = SubscribeOrchestration()
				DeferCleanup(func() {
					closeSubscription()
				})
			})
			It("Should trigger the workflow", func() {
				var (
					listTriggersOccurrencesResponse *operations.V2ListTriggersOccurrencesResponse
					err                             error
				)
				Eventually(func(g Gomega) bool {
					listTriggersOccurrencesResponse, err = Client().Orchestration.V2.ListTriggersOccurrences(TestContext(), operations.V2ListTriggersOccurrencesRequest{
						TriggerID: createTriggerResponse.CreateTriggerResponse.Data.ID,
					})
					g.Expect(err).To(BeNil())
					g.Expect(listTriggersOccurrencesResponse.V2ListTriggersOccurrencesResponse.Cursor.Data).NotTo(BeEmpty())
					occurrence := listTriggersOccurrencesResponse.V2ListTriggersOccurrencesResponse.Cursor.Data[0]
					g.Expect(occurrence.WorkflowInstance.Terminated).To(BeTrue())
					g.Expect(occurrence.WorkflowInstance.TerminatedAt).ShouldNot(BeNil())
					g.Expect(occurrence.WorkflowInstance.Error).Should(BeNil())
					return true
				}).Should(BeTrue())

				var getInstanceResponse *operations.V2GetInstanceResponse
				Eventually(func() bool {
					getInstanceResponse, err = Client().Orchestration.V2.GetInstance(TestContext(), operations.V2GetInstanceRequest{
						InstanceID: *listTriggersOccurrencesResponse.V2ListTriggersOccurrencesResponse.Cursor.Data[0].WorkflowInstanceID,
					})
					Expect(err).To(BeNil())

					return getInstanceResponse.V2GetWorkflowInstanceResponse.Data.Terminated
				}).Should(BeTrue())

				Expect(getInstanceResponse.V2GetWorkflowInstanceResponse.Data.Error).To(BeNil())

				listTransactionsResponse, err := Client().Ledger.V2.ListTransactions(TestContext(), operations.V2ListTransactionsRequest{
					Ledger: "default",
				})
				Expect(err).To(BeNil())
				Expect(listTransactionsResponse.V2TransactionsCursorResponse.Cursor.Data).To(HaveLen(1))
				Expect(listTransactionsResponse.V2TransactionsCursorResponse.Cursor.Data[0].Postings).To(HaveLen(1))
				Expect(listTransactionsResponse.V2TransactionsCursorResponse.Cursor.Data[0].Postings[0].Source).
					To(Equal("world"))
				Expect(listTransactionsResponse.V2TransactionsCursorResponse.Cursor.Data[0].Postings[0].Destination).
					To(Equal("foo"))
				Expect(listTransactionsResponse.V2TransactionsCursorResponse.Cursor.Data[0].Postings[0].Asset).
					To(Equal("USD/2"))
				Expect(listTransactionsResponse.V2TransactionsCursorResponse.Cursor.Data[0].Postings[0].Amount).
					To(Equal(big.NewInt(1000000)))

				By("And trigger a new succeeded workflow event", func() {
					Eventually(msgs).Should(ReceiveEvent("orchestration", orchestrationevents.SucceededTrigger))
				})
			})
		})
		Then("deleting the trigger", func() {
			BeforeEach(func() {
				_, err := Client().Orchestration.V2.DeleteTrigger(TestContext(), operations.V2DeleteTriggerRequest{
					TriggerID: createTriggerResponse.CreateTriggerResponse.Data.ID,
				})
				Expect(err).To(BeNil())
			})
			It("should not appear on list", func() {
				listTriggersResponse, err := Client().Orchestration.V2.ListTriggers(TestContext(), operations.V2ListTriggersRequest{})
				Expect(err).To(BeNil())
				Expect(listTriggersResponse.V2ListTriggersResponse.Cursor.Data).Should(HaveLen(0))
			})
		})
	})
})

var _ = WithModules([]*Module{modules.Auth, modules.Orchestration}, func() {
	When("creating a new workflow with a delay of 5s and a trigger on payments creation", func() {
		var (
			createTriggerResponse *operations.CreateTriggerResponse
		)
		BeforeEach(func() {
			response, err := Client().Orchestration.V1.CreateWorkflow(
				TestContext(),
				&shared.CreateWorkflowRequest{
					Name:   ptr(uuid.NewString()),
					Stages: []map[string]interface{}{{"delay": map[string]any{"duration": "5s"}}},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(201))

			createTriggerResponse, err = Client().Orchestration.V1.CreateTrigger(
				TestContext(),
				&shared.TriggerData{
					Event:      paymentsevents.EventTypeSavedPayments,
					WorkflowID: response.CreateWorkflowResponse.Data.ID,
					Vars: map[string]any{
						"fail": `link(event, "unknown").name`,
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
		})
		Then("publishing a new empty payment", func() {
			var (
				payment map[string]any
				msgs    chan *nats.Msg
			)
			BeforeEach(func() {
				payment = map[string]any{
					"id": uuid.NewString(),
				}
				PublishPayments(publish.EventMessage{
					Date:    time.Now(),
					App:     "payments",
					Type:    paymentsevents.EventTypeSavedPayments,
					Payload: payment,
				})
				var closeSubscription func()
				closeSubscription, msgs = SubscribeOrchestration()
				DeferCleanup(func() {
					closeSubscription()
				})
			})
			It("Should create a trigger workflow with an error", func() {
				var (
					listTriggersOccurrencesResponse *operations.V2ListTriggersOccurrencesResponse
					err                             error
				)
				Eventually(func(g Gomega) bool {
					listTriggersOccurrencesResponse, err = Client().Orchestration.V2.ListTriggersOccurrences(TestContext(), operations.V2ListTriggersOccurrencesRequest{
						TriggerID: createTriggerResponse.CreateTriggerResponse.Data.ID,
					})
					g.Expect(err).To(BeNil())
					g.Expect(listTriggersOccurrencesResponse.V2ListTriggersOccurrencesResponse.Cursor.Data).NotTo(BeEmpty())
					occurrence := listTriggersOccurrencesResponse.V2ListTriggersOccurrencesResponse.Cursor.Data[0]
					g.Expect(occurrence.WorkflowInstanceID).To(BeNil())
					g.Expect(occurrence.WorkflowInstance).To(BeNil())
					return true
				}).Should(BeTrue())

				By("should also trigger a new failed workflow event", func() {
					Eventually(msgs).Should(ReceiveEvent("orchestration", orchestrationevents.FailedTrigger))
				})
			})
		})
	})
})
