package suite

import (
	orchestrationevents "github.com/formancehq/orchestration/pkg/events"
	"github.com/formancehq/stack/libs/events"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"github.com/nats-io/nats.go"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
)

var _ = WithModules([]*Module{modules.Auth, modules.Orchestration, modules.Ledger}, func() {
	BeforeEach(func() {
		createLedgerResponse, err := Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
			Ledger: "default",
		})
		Expect(err).To(BeNil())
		Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))
	})
	When("creating a new workflow which will fail with insufficient fund error", func() {
		var (
			createWorkflowResponse *shared.V2CreateWorkflowResponse
		)
		BeforeEach(func() {
			response, err := Client().Orchestration.V2.CreateWorkflow(
				TestContext(),
				&shared.V2CreateWorkflowRequest{
					Name: ptr(uuid.New()),
					Stages: []map[string]interface{}{
						{
							"send": map[string]any{
								"source": map[string]any{
									"account": map[string]any{
										"id":     "empty:account",
										"ledger": "default",
									},
								},
								"destination": map[string]any{
									"account": map[string]any{
										"id":     "bank",
										"ledger": "default",
									},
								},
								"amount": map[string]any{
									"amount": 100,
									"asset":  "EUR/2",
								},
							},
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(201))

			createWorkflowResponse = response.V2CreateWorkflowResponse
		})
		Then("executing it", func() {
			var runWorkflowResponse *shared.V2RunWorkflowResponse
			BeforeEach(func() {
				response, err := Client().Orchestration.V2.RunWorkflow(
					TestContext(),
					operations.V2RunWorkflowRequest{
						RequestBody: map[string]string{},
						WorkflowID:  createWorkflowResponse.Data.ID,
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(201))

				runWorkflowResponse = response.V2RunWorkflowResponse
			})
			Then("waiting for first stage retried at least once", func() {
				var getWorkflowInstanceHistoryStageResponse *shared.V2GetWorkflowInstanceHistoryStageResponse
				BeforeEach(func() {
					Eventually(func(g Gomega) int64 {
						response, err := Client().Orchestration.V2.GetInstanceStageHistory(
							TestContext(),
							operations.V2GetInstanceStageHistoryRequest{
								InstanceID: runWorkflowResponse.Data.ID,
								Number:     0,
							},
						)
						g.Expect(err).To(BeNil())
						g.Expect(response.StatusCode).To(Equal(200))

						getWorkflowInstanceHistoryStageResponse = response.V2GetWorkflowInstanceHistoryStageResponse
						g.Expect(getWorkflowInstanceHistoryStageResponse.Data).To(HaveLen(1))
						return getWorkflowInstanceHistoryStageResponse.Data[0].Attempt
					}).Should(BeNumerically(">", 2))
				})
				It("should be retried with insufficient fund error", func() {
					Expect(getWorkflowInstanceHistoryStageResponse.Data[0].StartedAt).NotTo(BeZero())
					Expect(getWorkflowInstanceHistoryStageResponse.Data[0].NextExecution).NotTo(BeNil())
					Expect(getWorkflowInstanceHistoryStageResponse.Data[0].Attempt).To(BeNumerically(">", 2))
					Expect(getWorkflowInstanceHistoryStageResponse.Data[0]).To(Equal(shared.V2WorkflowInstanceHistoryStage{
						Name: "CreateTransaction",
						Input: shared.V2WorkflowInstanceHistoryStageInput{
							CreateTransaction: &shared.V2ActivityCreateTransaction{
								Ledger: ptr("default"),
								Data: &shared.V2PostTransaction{
									Postings: []shared.V2Posting{{
										Amount:      big.NewInt(100),
										Asset:       "EUR/2",
										Destination: "bank",
										Source:      "empty:account",
									}},
								},
							},
						},
						LastFailure:   ptr("running numscript: script execution failed: account(s) @empty:account had/have insufficient funds"),
						Attempt:       getWorkflowInstanceHistoryStageResponse.Data[0].Attempt,
						NextExecution: getWorkflowInstanceHistoryStageResponse.Data[0].NextExecution,
						StartedAt:     getWorkflowInstanceHistoryStageResponse.Data[0].StartedAt,
					}))
				})
			})
		})
	})
})

var _ = WithModules([]*Module{modules.Auth, modules.Orchestration, modules.Ledger}, func() {
	BeforeEach(func() {
		createLedgerResponse, err := Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
			Ledger: "default",
		})
		Expect(err).To(BeNil())
		Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))
	})
	When("creating a new workflow which will fail with invalid request", func() {
		var (
			createWorkflowResponse *shared.V2CreateWorkflowResponse
		)
		BeforeEach(func() {
			response, err := Client().Orchestration.V2.CreateWorkflow(
				TestContext(),
				&shared.V2CreateWorkflowRequest{
					Name: ptr(uuid.New()),
					Stages: []map[string]interface{}{
						{
							"send": map[string]any{
								"source": map[string]any{
									"account": map[string]any{
										"id":     "empty:account",
										"ledger": "default",
									},
								},
								"destination": map[string]any{
									"account": map[string]any{
										"id":     "bank",
										"ledger": "default",
									},
								},
								"amount": map[string]any{
									"amount": -1, // Invalid amount
									"asset":  "EUR/2",
								},
							},
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(201))

			createWorkflowResponse = response.V2CreateWorkflowResponse
		})
		Then("executing it", func() {
			var (
				runWorkflowResponse *shared.V2RunWorkflowResponse
				msgs                chan *nats.Msg
			)
			BeforeEach(func() {
				var closeSubscription func()
				closeSubscription, msgs = SubscribeOrchestration()
				DeferCleanup(func() {
					closeSubscription()
				})

				response, err := Client().Orchestration.V2.RunWorkflow(
					TestContext(),
					operations.V2RunWorkflowRequest{
						RequestBody: map[string]string{},
						WorkflowID:  createWorkflowResponse.Data.ID,
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(201))

				runWorkflowResponse = response.V2RunWorkflowResponse
			})
			It("should declare the workflow run instance as errored", func() {
				Eventually(func(g Gomega) string {
					response, err := Client().Orchestration.V2.GetInstanceStageHistory(
						TestContext(),
						operations.V2GetInstanceStageHistoryRequest{
							InstanceID: runWorkflowResponse.Data.ID,
							Number:     0,
						},
					)
					g.Expect(err).To(BeNil())
					g.Expect(response.StatusCode).To(Equal(200))
					g.Expect(response.V2GetWorkflowInstanceHistoryStageResponse.Data).To(HaveLen(1))
					g.Expect(response.V2GetWorkflowInstanceHistoryStageResponse.Data[0].Error).NotTo(BeNil())

					return *response.V2GetWorkflowInstanceHistoryStageResponse.Data[0].Error
				}).ShouldNot(BeEmpty())

				By("Should trigger appropriate events", func() {
					msg := WaitOnChanWithTimeout(msgs, 5*time.Second)
					Expect(events.Check(msg.Data, "orchestration", orchestrationevents.StartedWorkflow)).Should(Succeed())

					msg = WaitOnChanWithTimeout(msgs, 5*time.Second)
					Expect(events.Check(msg.Data, "orchestration", orchestrationevents.StartedWorkflowStage)).Should(Succeed())

					msg = WaitOnChanWithTimeout(msgs, 5*time.Second)
					Expect(events.Check(msg.Data, "orchestration", orchestrationevents.FailedWorkflowStage)).Should(Succeed())

					msg = WaitOnChanWithTimeout(msgs, 5*time.Second)
					Expect(events.Check(msg.Data, "orchestration", orchestrationevents.FailedWorkflow)).Should(Succeed())
				})
			})
		})
	})
})
