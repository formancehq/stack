package suite

import (
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"math/big"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
)

var _ = WithModules([]*Module{modules.Auth, modules.Orchestration, modules.Ledger}, func() {
	When("creating a new workflow which will fail with insufficient fund error", func() {
		var (
			createWorkflowResponse *shared.CreateWorkflowResponse
		)
		BeforeEach(func() {
			response, err := Client().Orchestration.CreateWorkflow(
				TestContext(),
				&shared.CreateWorkflowRequest{
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

			createWorkflowResponse = response.CreateWorkflowResponse
		})
		Then("executing it", func() {
			var runWorkflowResponse *shared.RunWorkflowResponse
			BeforeEach(func() {
				response, err := Client().Orchestration.RunWorkflow(
					TestContext(),
					operations.RunWorkflowRequest{
						RequestBody: map[string]string{},
						WorkflowID:  createWorkflowResponse.Data.ID,
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(201))

				runWorkflowResponse = response.RunWorkflowResponse
			})
			Then("waiting for first stage retried at least once", func() {
				var getWorkflowInstanceHistoryStageResponse *shared.GetWorkflowInstanceHistoryStageResponse
				BeforeEach(func() {
					Eventually(func(g Gomega) int64 {

						//TODO: error EOF response
						response, err := Client().Orchestration.GetInstanceStageHistory(
							TestContext(),
							operations.GetInstanceStageHistoryRequest{
								InstanceID: runWorkflowResponse.Data.ID,
								Number:     0,
							},
						)
						if err != nil {
							return 0
						}
						if response.StatusCode != 200 {
							return 0
						}

						getWorkflowInstanceHistoryStageResponse = response.GetWorkflowInstanceHistoryStageResponse
						g.Expect(getWorkflowInstanceHistoryStageResponse.Data).To(HaveLen(1))
						return getWorkflowInstanceHistoryStageResponse.Data[0].Attempt
					}).Should(BeNumerically(">", 2))
				})
				It("should be retried with insufficient fund error ", func() {
					Expect(getWorkflowInstanceHistoryStageResponse.Data[0].StartedAt).NotTo(BeZero())
					Expect(getWorkflowInstanceHistoryStageResponse.Data[0].NextExecution).NotTo(BeNil())
					Expect(getWorkflowInstanceHistoryStageResponse.Data[0].Attempt).To(BeNumerically(">", 2))
					Expect(getWorkflowInstanceHistoryStageResponse.Data[0]).To(Equal(shared.WorkflowInstanceHistoryStage{
						Name: "CreateTransaction",
						Input: shared.WorkflowInstanceHistoryStageInput{
							CreateTransaction: &shared.ActivityCreateTransaction{
								Ledger: ptr("default"),
								Data: &shared.PostTransaction{
									Postings: []shared.Posting{{
										Amount:      big.NewInt(100),
										Asset:       "EUR/2",
										Destination: "bank",
										Source:      "empty:account",
									}},
									Metadata: metadata.Metadata{},
								},
							},
						},
						LastFailure:   ptr("running numscript: script execution failed: no more fund to withdraw"),
						Attempt:       getWorkflowInstanceHistoryStageResponse.Data[0].Attempt,
						NextExecution: getWorkflowInstanceHistoryStageResponse.Data[0].NextExecution,
						StartedAt:     getWorkflowInstanceHistoryStageResponse.Data[0].StartedAt,
					}))
				})
			})
		})
	})
})
