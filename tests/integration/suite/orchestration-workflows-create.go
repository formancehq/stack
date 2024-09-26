package suite

import (
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
)

var _ = WithModules([]*Module{modules.Auth, modules.Orchestration}, func() {
	When("creating a new workflow", func() {
		var (
			createWorkflowResponse *shared.V2CreateWorkflowResponse
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

			createWorkflowResponse = response.V2CreateWorkflowResponse
		})
		It("should be ok", func() {
			Expect(createWorkflowResponse.Data.ID).NotTo(BeEmpty())
		})
	})
})
