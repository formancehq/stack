package suite

import (
	"encoding/json"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"math/big"
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
)

var _ = WithModules([]*Module{modules.Orchestration, modules.Auth, modules.Ledger}, func() {
	When("creating a new workflow", func() {
		var (
			createWorkflowResponse *shared.V2CreateWorkflowResponse
		)
		BeforeEach(func() {
			response, err := Client().Orchestration.V2CreateWorkflow(
				TestContext(),
				&shared.V2CreateWorkflowRequest{
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
								"metadata": map[string]any{
									"foo": "${userID}",
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
		It("should be ok", func() {
			Expect(createWorkflowResponse.Data.ID).NotTo(BeEmpty())
		})
		Then("executing it", func() {
			var runWorkflowResponse *shared.V2RunWorkflowResponse
			BeforeEach(func() {
				response, err := Client().Orchestration.V2RunWorkflow(
					TestContext(),
					operations.V2RunWorkflowRequest{
						RequestBody: map[string]string{
							"userID": "bar",
						},
						WorkflowID: createWorkflowResponse.Data.ID,
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(201))

				runWorkflowResponse = response.V2RunWorkflowResponse
			})
			It("should be ok", func() {
				Expect(runWorkflowResponse.Data.ID).NotTo(BeEmpty())
			})
			Then("waiting for termination", func() {
				var instanceResponse *shared.V2GetWorkflowInstanceResponse
				BeforeEach(func() {
					Eventually(func() bool {
						response, err := Client().Orchestration.V2GetInstance(
							TestContext(),
							operations.V2GetInstanceRequest{
								InstanceID: runWorkflowResponse.Data.ID,
							},
						)
						Expect(err).ToNot(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))

						instanceResponse = response.V2GetWorkflowInstanceResponse

						return response.V2GetWorkflowInstanceResponse.Data.Terminated
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
				//Then("checking ledger account balance", func() {
				//	var balancesCursorResponse *shared.BalancesCursorResponse
				//	BeforeEach(func() {
				//		reponse, err := Client().Ledger.V2GetBalances(
				//			TestContext(),
				//			operations.GetBalancesRequest{
				//				Address: ptr("bank"),
				//				Ledger:  "default",
				//			},
				//		)
				//		Expect(err).ToNot(HaveOccurred())
				//		Expect(reponse.StatusCode).To(Equal(200))
				//
				//		balancesCursorResponse = reponse.BalancesCursorResponse
				//	})
				//	It("should return 100 USD/2 available", func() {
				//		Expect(balancesCursorResponse.Cursor.Data).To(HaveLen(1))
				//		Expect(balancesCursorResponse.Cursor.Data[0]).To(HaveLen(1))
				//		Expect(balancesCursorResponse.Cursor.Data[0]["bank"]).To(HaveLen(1))
				//		Expect(balancesCursorResponse.Cursor.Data[0]["bank"]["EUR/2"]).To(Equal(big.NewInt(100)))
				//	})
				//})
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
									Amount: big.NewInt(100),
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
								Metadata: map[string]string{
									"foo": "bar",
								},
							}))
					})
					Then("reading first stage history", func() {
						var getWorkflowInstanceHistoryStageResponse *shared.V2GetWorkflowInstanceHistoryStageResponse
						BeforeEach(func() {
							response, err := Client().Orchestration.V2GetInstanceStageHistory(
								TestContext(),
								operations.V2GetInstanceStageHistoryRequest{
									InstanceID: runWorkflowResponse.Data.ID,
									Number:     0,
								},
							)
							Expect(err).ToNot(HaveOccurred())
							Expect(response.StatusCode).To(Equal(200))

							getWorkflowInstanceHistoryStageResponse = response.V2GetWorkflowInstanceHistoryStageResponse
						})
						It("should be properly terminated", func() {
							Expect(getWorkflowInstanceHistoryStageResponse.Data).To(HaveLen(1))
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].Error).To(BeNil())
							postings := []shared.V2Posting{{
								Amount:      big.NewInt(100),
								Asset:       "EUR/2",
								Destination: "bank",
								Source:      "world",
							}}
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].Input).To(Equal(shared.V2WorkflowInstanceHistoryStageInput{
								CreateTransaction: &shared.V2ActivityCreateTransaction{
									Ledger: ptr("default"),
									Data: &shared.V2PostTransaction{
										Postings: postings,
										Metadata: metadata.Metadata{
											"foo": "bar",
										},
									},
								},
							}))
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].StartedAt).NotTo(BeZero())
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].TerminatedAt).NotTo(BeZero())
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].LastFailure).To(BeNil())
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].StartedAt).NotTo(BeZero())
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].Attempt).To(Equal(int64(1)))
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].NextExecution).To(BeNil())
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].Output.CreateTransaction.Data[0].Timestamp).
								NotTo(BeZero())
							getWorkflowInstanceHistoryStageResponse.Data[0].Output.CreateTransaction.Data[0].Timestamp = time.Time{}
							Expect(getWorkflowInstanceHistoryStageResponse.Data[0].Output).To(Equal(&shared.V2WorkflowInstanceHistoryStageOutput{
								CreateTransaction: &shared.V2ActivityCreateTransactionOutput{
									Data: []shared.OrchestrationV2Transaction{{
										Txid:     big.NewInt(0),
										Postings: postings,
										Metadata: map[string]string{
											"foo": "bar",
										},
									}},
								},
							}))
						})
					})
				})
			})
		})
	})
})
