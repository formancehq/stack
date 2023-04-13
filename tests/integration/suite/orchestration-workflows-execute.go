package suite

import (
	"time"

	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
)

var _ = Given("An empty environment", func() {
	When("creating a new workflow", func() {
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
										"id":     "world",
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
			Expect(err).ToNot(HaveOccurred())
		})
		It("should be ok", func() {
			Expect(createWorkflowResponse.Data.Id).NotTo(BeEmpty())
		})
		Then("executing it", func() {
			var runWorkflowResponse *formance.RunWorkflowResponse
			BeforeEach(func() {
				runWorkflowResponse, _, err = Client().OrchestrationApi.
					RunWorkflow(TestContext(), createWorkflowResponse.Data.Id).
					RequestBody(map[string]string{}).
					Execute()
				Expect(err).ToNot(HaveOccurred())
			})
			It("should be ok", func() {
				Expect(runWorkflowResponse.Data.Id).NotTo(BeEmpty())
			})
			Then("waiting for termination", func() {
				var instanceResponse *formance.GetWorkflowInstanceResponse
				BeforeEach(func() {
					Eventually(func() bool {
						instanceResponse, _, err = Client().OrchestrationApi.
							GetInstance(TestContext(), runWorkflowResponse.Data.Id).
							Execute()
						Expect(err).ToNot(HaveOccurred())
						return instanceResponse.Data.Terminated
					}).Should(BeTrue())
				})
				It("should be terminated successfully", func() {
					Expect(instanceResponse.Data.Id).NotTo(BeEmpty())
					if instanceResponse.Data.Error != nil {
						Expect(*instanceResponse.Data.Error).To(BeEmpty())
					}
					Expect(instanceResponse.Data.Terminated).To(BeTrue())
					Expect(instanceResponse.Data.CreatedAt).NotTo(BeZero())
					Expect(instanceResponse.Data.UpdatedAt).NotTo(BeZero())
					Expect(instanceResponse.Data.TerminatedAt).NotTo(BeZero())
					Expect(instanceResponse.Data.Status).To(HaveLen(1))
				})
				Then("checking ledger account balance", func() {
					var balancesCursorResponse *formance.BalancesCursorResponse
					BeforeEach(func() {
						balancesCursorResponse, _, err = Client().BalancesApi.
							GetBalances(TestContext(), "default").
							Address("bank").
							Execute()
						Expect(err).ToNot(HaveOccurred())
					})
					It("should return 100 USD/2 available", func() {
						Expect(balancesCursorResponse.Cursor.Data).To(HaveLen(1))
						Expect(balancesCursorResponse.Cursor.Data[0]).To(HaveLen(1))
						Expect(balancesCursorResponse.Cursor.Data[0]["bank"]).To(HaveLen(1))
						Expect(balancesCursorResponse.Cursor.Data[0]["bank"]["EUR/2"]).To(Equal(int64(100)))
					})
				})
				Then("reading history", func() {
					var getWorkflowInstanceHistoryResponse *formance.GetWorkflowInstanceHistoryResponse
					BeforeEach(func() {
						getWorkflowInstanceHistoryResponse, _, err = Client().OrchestrationApi.
							GetInstanceHistory(TestContext(), runWorkflowResponse.Data.Id).
							Execute()
					})
					It("should be ok", func() {
						Expect(getWorkflowInstanceHistoryResponse.Data).To(HaveLen(1))
						Expect(getWorkflowInstanceHistoryResponse.Data[0].Terminated).To(BeTrue())
						Expect(getWorkflowInstanceHistoryResponse.Data[0].TerminatedAt).NotTo(BeZero())
						Expect(getWorkflowInstanceHistoryResponse.Data[0].StartedAt).NotTo(BeZero())
						if getWorkflowInstanceHistoryResponse.Data[0].Error != nil {
							Expect(*getWorkflowInstanceHistoryResponse.Data[0].Error).To(BeEmpty())
						}
						Expect(getWorkflowInstanceHistoryResponse.Data[0].Input.StageSend).NotTo(BeNil())
						Expect(*getWorkflowInstanceHistoryResponse.Data[0].Input.StageSend).
							To(Equal(formance.StageSend{
								Amount: formance.NewMonetary("EUR/2", 100),
								Destination: &formance.StageSendDestination{
									Account: &formance.StageSendSourceAccount{
										Id:     "bank",
										Ledger: formance.PtrString("default"),
									},
								},
								Source: &formance.StageSendSource{
									Account: &formance.StageSendSourceAccount{
										Id:     "world",
										Ledger: formance.PtrString("default"),
									},
								},
							}))
					})
					Then("reading first stage history", func() {
						var getWorkflowInstanceHistoryStageResponse *formance.GetWorkflowInstanceHistoryStageResponse
						BeforeEach(func() {
							getWorkflowInstanceHistoryStageResponse, _, err = Client().OrchestrationApi.
								GetInstanceStageHistory(TestContext(), runWorkflowResponse.Data.Id, 0).
								Execute()
							Expect(err).ToNot(HaveOccurred())
						})
						It("should be properly terminated", func() {
							Expect(getWorkflowInstanceHistoryStageResponse.Data).To(HaveLen(1))
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].Error).To(BeNil())
							postings := []formance.Posting{{
								Amount:      100,
								Asset:       "EUR/2",
								Destination: "bank",
								Source:      "world",
							}}
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].Input).To(Equal(formance.WorkflowInstanceHistoryStageInput{
								CreateTransaction: &formance.ActivityCreateTransaction{
									Ledger: formance.PtrString("default"),
									Data: &formance.PostTransaction{
										Postings: postings,
										Metadata: metadata.Metadata{},
									},
								},
							}))
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].StartedAt).NotTo(BeZero())
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].TerminatedAt).NotTo(BeZero())
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].LastFailure).To(BeNil())
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].StartedAt).NotTo(BeZero())
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].Attempt).To(Equal(int32(1)))
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].NextExecution).To(BeNil())

							//TODO: fail here with zero timestamp
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].Output.CreateTransaction.Data.Timestamp).
								NotTo(BeZero())
							getWorkflowInstanceHistoryStageResponse.Data[0].Output.CreateTransaction.Data.Timestamp = time.Time{}
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].Output).To(Equal(&formance.WorkflowInstanceHistoryStageOutput{
								CreateTransaction: &formance.CreateTransactionResponse{
									Data: formance.Transaction{
										Postings:  postings,
										Reference: formance.PtrString(""),
										Metadata:  metadata.Metadata{},
									},
								},
							}))
						})
					})
				})
			})
		})
	})
})
