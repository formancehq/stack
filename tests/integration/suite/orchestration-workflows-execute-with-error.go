package suite

import (
	"github.com/formancehq/formance-sdk-go"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
)

var _ = Given("An empty environment", func() {
	When("creating a new workflow which will fail with insufficient fund error", func() {
		var (
			createWorkflowResponse *formance.CreateWorkflowResponse
			err                    error
		)
		BeforeEach(func() {
			createWorkflowResponse, _, err = Client().OrchestrationApi.
				CreateWorkflow(TestContext()).
				Body(formance.WorkflowConfig{
					Name: formance.PtrString(uuid.New()),
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
				}).
				Execute()
			Expect(err).To(BeNil())
		})
		Then("executing it", func() {
			var (
				runWorkflowResponse *formance.RunWorkflowResponse
			)
			BeforeEach(func() {
				runWorkflowResponse, _, err = Client().OrchestrationApi.
					RunWorkflow(TestContext(), createWorkflowResponse.Data.Id).
					Execute()
				Expect(err).To(BeNil())
			})
			Then("waiting for first stage retried at least once", func() {
				var (
					getWorkflowInstanceHistoryStageResponse *formance.GetWorkflowInstanceHistoryStageResponse
				)
				BeforeEach(func() {
					Eventually(func(g Gomega) int32 {
						getWorkflowInstanceHistoryStageResponse, _, err = Client().OrchestrationApi.
							GetInstanceStageHistory(TestContext(), runWorkflowResponse.Data.Id, 0).
							Execute()
						g.Expect(err).To(BeNil())
						g.Expect(getWorkflowInstanceHistoryStageResponse.Data).To(HaveLen(1))
						return getWorkflowInstanceHistoryStageResponse.Data[0].Attempt
					}).Should(BeNumerically(">", 2))
				})
				It("should be retried with insufficient fund error ", func() {
					Expect(getWorkflowInstanceHistoryStageResponse.Data[0].StartedAt).NotTo(BeZero())
					Expect(getWorkflowInstanceHistoryStageResponse.Data[0].NextExecution).NotTo(BeNil())
					Expect(getWorkflowInstanceHistoryStageResponse.Data[0].Attempt).To(BeNumerically(">", 2))
					Expect(getWorkflowInstanceHistoryStageResponse.Data[0]).To(Equal(formance.WorkflowInstanceHistoryStage{
						Name: "CreateTransaction",
						Input: formance.WorkflowInstanceHistoryStageInput{
							CreateTransaction: &formance.ActivityCreateTransaction{
								Ledger: formance.PtrString("default"),
								Data: &formance.PostTransaction{
									Postings: []formance.Posting{{
										Amount:      100,
										Asset:       "EUR/2",
										Destination: "bank",
										Source:      "empty:account",
									}},
								},
							},
						},
						LastFailure:   formance.PtrString("[INSUFFICIENT_FUND] account had insufficient funds"),
						Attempt:       getWorkflowInstanceHistoryStageResponse.Data[0].Attempt,
						NextExecution: getWorkflowInstanceHistoryStageResponse.Data[0].NextExecution,
						StartedAt:     getWorkflowInstanceHistoryStageResponse.Data[0].StartedAt,
					}))
				})
			})
		})
	})
})
