package suite

import (
	"net/http"

	"github.com/formancehq/stack/tests/integration/internal/modules"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Webhooks}, func() {
	var (
		endpoint2 = "https://example2.com"
		hook1 shared.V2Hook
	)

	BeforeEach(func() {
		hookBodyParam := shared.V2HookBodyParams{
			Endpoint: "https://example.com",
			Events: []string{
				"ledger.committed_transactions",
			},
		}
		response, err := Client().Webhooks.InsertHook(
			TestContext(),
			hookBodyParam,
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusCreated))

		hook1 = response.V2HookResponse.Data
	})

	Context("changing the endpoint of the inserted one", func() {
		
			It("should work ",func() {
				response, err := Client().Webhooks.UpdateEndpointHook(
					TestContext(),
					operations.UpdateEndpointHookRequest{
						RequestBody: operations.UpdateEndpointHookRequestBody{
							Endpoint: &endpoint2,
						},
						HookID: hook1.ID,
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				Expect(response.V2HookResponse.Data.Endpoint).To(Equal(endpoint2))
			})
			
		

	})
})
