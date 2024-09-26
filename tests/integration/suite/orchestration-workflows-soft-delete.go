package suite

import (
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/sdkerrors"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
)

var _ = WithModules([]*Module{modules.Auth, modules.Orchestration}, func() {
	When("populating 1 workflow", func() {
		var (
			workflow *shared.V2Workflow
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

			workflow = &response.V2CreateWorkflowResponse.Data
		})
		It("should be ok", func() {
			Expect(workflow.ID).NotTo(BeEmpty())
		})
		It("should delete the workflow", func() {
			response, err := Client().Orchestration.V2.DeleteWorkflow(
				TestContext(),
				operations.V2DeleteWorkflowRequest{
					FlowID: workflow.ID,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))

		})
	})

	When("deleting a non-uuid workflow", func() {
		It("should return 400 with unknown", func() {
			_, err := Client().Orchestration.V2.DeleteWorkflow(
				TestContext(),
				operations.V2DeleteWorkflowRequest{
					FlowID: "unknown",
				},
			)
			Expect(err).To(HaveOccurred())
			Expect(err.(*sdkerrors.V2Error).ErrorCode).To(Equal(sdkerrors.SchemasErrorCodeValidation))
		})
	})
	When("deleting a non-existing workflow", func() {
		It("should return 404", func() {
			_, err := Client().Orchestration.V2.DeleteWorkflow(
				TestContext(),
				operations.V2DeleteWorkflowRequest{
					FlowID: uuid.New(),
				},
			)
			Expect(err).To(HaveOccurred())
			Expect(err.(*sdkerrors.V2Error).ErrorCode).To(Equal(sdkerrors.SchemasErrorCodeNotFound))
		})
	})
})
