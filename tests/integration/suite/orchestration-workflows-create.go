package suite

import (
	"github.com/formancehq/formance-sdk-go"
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
			Expect(err).To(BeNil())
		})
		It("should be ok", func() {
			Expect(createWorkflowResponse.Data.Id).NotTo(BeEmpty())
		})
	})
})
