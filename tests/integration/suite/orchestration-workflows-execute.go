package suite

import (
	"encoding/json"
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
)

var _ = Given("An empty environment", func() {
	When("creating a new workflow", func() {
		var (
			createWorkflowResponse *shared.CreateWorkflowResponse
		)
		BeforeEach(func() {
			response, err := Client().Orchestration.CreateWorkflow(
				TestContext(),
				shared.CreateWorkflowRequest{
					Name: ptr(uuid.New()),
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
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(201))

			createWorkflowResponse = response.CreateWorkflowResponse
		})
		It("should be ok", func() {
			Expect(createWorkflowResponse.Data.ID).NotTo(BeEmpty())
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
			It("should be ok", func() {
				Expect(runWorkflowResponse.Data.ID).NotTo(BeEmpty())
			})
			Then("waiting for termination", func() {
				var instanceResponse *shared.GetWorkflowInstanceResponse
				BeforeEach(func() {
					Eventually(func() bool {
						response, err := Client().Orchestration.GetInstance(
							TestContext(),
							operations.GetInstanceRequest{
								InstanceID: runWorkflowResponse.Data.ID,
							},
						)
						Expect(err).ToNot(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))

						instanceResponse = response.GetWorkflowInstanceResponse
						return response.GetWorkflowInstanceResponse.Data.Terminated
					}).Should(BeTrue())
				})
				It("should be terminated successfully", func() {
					Expect(instanceResponse.Data.ID).NotTo(BeEmpty())
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
					var balancesCursorResponse *shared.BalancesCursorResponse
					BeforeEach(func() {
						reponse, err := Client().Ledger.GetBalances(
							TestContext(),
							operations.GetBalancesRequest{
								Address: ptr("bank"),
								Ledger:  "default",
							},
						)
						Expect(err).ToNot(HaveOccurred())
						Expect(reponse.StatusCode).To(Equal(200))

						balancesCursorResponse = reponse.BalancesCursorResponse
					})
					It("should return 100 USD/2 available", func() {
						Expect(balancesCursorResponse.Cursor.Data).To(HaveLen(1))
						Expect(balancesCursorResponse.Cursor.Data[0]).To(HaveLen(1))
						Expect(balancesCursorResponse.Cursor.Data[0]["bank"]).To(HaveLen(1))
						Expect(balancesCursorResponse.Cursor.Data[0]["bank"]["EUR/2"]).To(Equal(int64(100)))
					})
				})
				Then("reading history", func() {
					var getWorkflowInstanceHistoryResponse *shared.GetWorkflowInstanceHistoryResponse
					BeforeEach(func() {
						response, err := Client().Orchestration.GetInstanceHistory(
							TestContext(),
							operations.GetInstanceHistoryRequest{
								InstanceID: runWorkflowResponse.Data.ID,
							},
						)
						Expect(err).ToNot(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))

						getWorkflowInstanceHistoryResponse = response.GetWorkflowInstanceHistoryResponse
					})
					It("should be ok", func() {
						Expect(getWorkflowInstanceHistoryResponse.Data).To(HaveLen(1))
						Expect(getWorkflowInstanceHistoryResponse.Data[0].Terminated).To(BeTrue())
						Expect(getWorkflowInstanceHistoryResponse.Data[0].TerminatedAt).NotTo(BeZero())
						Expect(getWorkflowInstanceHistoryResponse.Data[0].StartedAt).NotTo(BeZero())
						if getWorkflowInstanceHistoryResponse.Data[0].Error != nil {
							Expect(*getWorkflowInstanceHistoryResponse.Data[0].Error).To(BeEmpty())
						}
						Expect(getWorkflowInstanceHistoryResponse.Data[0].Input).NotTo(BeNil())
						var stageSend shared.StageSend
						buf, err := json.Marshal(getWorkflowInstanceHistoryResponse.Data[0].Input)
						Expect(err).ToNot(HaveOccurred())
						err = json.Unmarshal(buf, &stageSend)
						Expect(err).ToNot(HaveOccurred())
						Expect(stageSend).
							To(Equal(shared.StageSend{
								Amount: &shared.Monetary{
									Amount: 100,
									Asset:  "EUR/2",
								},
								Destination: &shared.StageSendDestination{
									Account: &shared.StageSendDestinationAccount{
										ID:     "bank",
										Ledger: ptr("default"),
									},
								},
								Source: &shared.StageSendSource{
									Account: &shared.StageSendSourceAccount{
										ID:     "world",
										Ledger: ptr("default"),
									},
								},
							}))
					})
					Then("reading first stage history", func() {
						var getWorkflowInstanceHistoryStageResponse *shared.GetWorkflowInstanceHistoryStageResponse
						BeforeEach(func() {
							response, err := Client().Orchestration.GetInstanceStageHistory(
								TestContext(),
								operations.GetInstanceStageHistoryRequest{
									InstanceID: runWorkflowResponse.Data.ID,
									Number:     0,
								},
							)
							Expect(err).ToNot(HaveOccurred())
							Expect(response.StatusCode).To(Equal(200))

							getWorkflowInstanceHistoryStageResponse = response.GetWorkflowInstanceHistoryStageResponse
						})
						It("should be properly terminated", func() {
							Expect(getWorkflowInstanceHistoryStageResponse.Data).To(HaveLen(1))
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].Error).To(BeNil())
							postings := []shared.Posting{{
								Amount:      100,
								Asset:       "EUR/2",
								Destination: "bank",
								Source:      "world",
							}}
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].Input).To(Equal(shared.WorkflowInstanceHistoryStageInput{
								CreateTransaction: &shared.ActivityCreateTransaction{
									Ledger: ptr("default"),
									Data: &shared.PostTransaction{
										Postings: postings,
										Metadata: metadata.Metadata{},
									},
								},
							}))
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].StartedAt).NotTo(BeZero())
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].TerminatedAt).NotTo(BeZero())
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].LastFailure).To(BeNil())
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].StartedAt).NotTo(BeZero())
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].Attempt).To(Equal(int64(1)))
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].NextExecution).To(BeNil())

							//TODO: fail here with zero timestamp
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].Output.CreateTransaction.Data.Timestamp).
								NotTo(BeZero())
							getWorkflowInstanceHistoryStageResponse.Data[0].Output.CreateTransaction.Data.Timestamp = time.Time{}
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].Output).To(Equal(&shared.WorkflowInstanceHistoryStageOutput{
								CreateTransaction: &shared.ActivityCreateTransactionOutput{
									Data: shared.Transaction{
										Postings:  postings,
										Reference: ptr(""),
										Metadata:  map[string]string{},
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
