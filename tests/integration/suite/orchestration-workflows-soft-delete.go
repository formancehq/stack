package suite

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
)

var _ = WithModules([]*Module{modules.Auth, modules.Orchestration}, func() {
	When("populating 1 workflow", func() {
		var (
			workflow *shared.Workflow
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

			workflow = &response.CreateWorkflowResponse.Data
		})
		It("should be ok", func() {
			Expect(workflow.ID).NotTo(BeEmpty())
		})
		It("should delete the workflow", func() {
			response, err := Client().Orchestration.DeleteWorkflow(
				TestContext(),
				operations.DeleteWorkflowRequest{
					FlowID: workflow.ID,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))

		})
	})

	When("deleting a non-uuid workflow", func() {
		It("should return 400 with unknown", func() {
			response, err := Client().Orchestration.DeleteWorkflow(
				TestContext(),
				operations.DeleteWorkflowRequest{
					FlowID: "unknown",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(400))
		})
		It("should return 400, with empty spaces'      '", func() {
			response, err := Client().Orchestration.DeleteWorkflow(
				TestContext(),
				operations.DeleteWorkflowRequest{
					FlowID: "      ",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(400))
		})
	})
	When("deleting a non-existing workflow", func() {
		It("should return 404", func() {
			response, err := Client().Orchestration.DeleteWorkflow(
				TestContext(),
				operations.DeleteWorkflowRequest{
					FlowID: uuid.New(),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(404))
		})
	})
})
