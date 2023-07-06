package suite

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
)

var _ = Given("An empty environment", func() {
	var (
		workflow shared.Workflow
		ret      *shared.ListWorkflowsResponse
	)
	When("first listing workflows", func() {
		BeforeEach(func() {
			response, err := Client().Orchestration.ListWorkflows(
				TestContext(),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			ret = response.ListWorkflowsResponse
		})
		It("should respond with an empty list", func() {
			Expect(ret.Data).To(BeEmpty())
		})
	})
	When("populating 1 workflow", func() {
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

			workflow = response.CreateWorkflowResponse.Data
		})
		It("should be ok", func() {
			Expect(workflow.ID).NotTo(BeEmpty())
		})
		It("Should retrieve workflows", func() {

			response, err := Client().Orchestration.ListWorkflows(
				TestContext(),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			ret = response.ListWorkflowsResponse
		})
		It("should respond with a list of 1 element", func() {
			Expect(ret.Data).ToNot(BeEmpty())
			Expect(ret.Data).Should(HaveLen(1))
		})
	})
	When("Deleting a workflow", func() {
		BeforeEach(func() {
			response, err := Client().Orchestration.DeleteWorkflow(
				TestContext(),
				operations.DeleteWorkflowRequest{
					FlowID: workflow.ID,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))

		})
		It("Should retrieve workflows", func() {
			response, err := Client().Orchestration.ListWorkflows(
				TestContext(),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			ret = response.ListWorkflowsResponse
		})
		It("should have a list of 0 element", func() {
			Expect(ret.Data).To(BeEmpty())
		})
	})
})
