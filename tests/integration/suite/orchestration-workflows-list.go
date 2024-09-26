package suite

import (
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
)

var _ = WithModules([]*Module{modules.Auth, modules.Orchestration}, func() {
	var (
		workflow shared.V2Workflow
		list     []shared.V2Workflow
	)
	When("first listing workflows", func() {
		BeforeEach(func() {
			response, err := Client().Orchestration.V2.ListWorkflows(
				TestContext(),
				operations.V2ListWorkflowsRequest{},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			list = response.V2ListWorkflowsResponse.Cursor.Data
		})
		It("should respond with an empty list", func() {
			Expect(list).To(BeEmpty())
		})
	})
	When("populating 1 workflow", func() {
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

			workflow = response.V2CreateWorkflowResponse.Data
		})
		It("should be ok", func() {
			Expect(workflow.ID).NotTo(BeEmpty())
		})
		JustBeforeEach(func() {
			response, err := Client().Orchestration.V2.ListWorkflows(
				TestContext(),
				operations.V2ListWorkflowsRequest{},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			list = response.V2ListWorkflowsResponse.Cursor.Data
		})
		It("should respond with a list of 1 elements", func() {
			Expect(list).ToNot(BeEmpty())
			Expect(list).Should(HaveLen(1))
		})
		When("Deleting a workflow", func() {
			JustBeforeEach(func() {
				response, err := Client().Orchestration.V2.DeleteWorkflow(
					TestContext(),
					operations.V2DeleteWorkflowRequest{
						FlowID: workflow.ID,
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(204))

			})
			JustBeforeEach(func() {
				response, err := Client().Orchestration.V2.ListWorkflows(
					TestContext(),
					operations.V2ListWorkflowsRequest{},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				list = response.V2ListWorkflowsResponse.Cursor.Data

			})
			It("should have a list of 0 element", func() {
				Expect(list).To(BeEmpty())
			})
		})
	})
})
