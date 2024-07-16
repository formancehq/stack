package suite

import (
	"net/http"

	"github.com/formancehq/stack/tests/integration/internal/modules"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/sdkerrors"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	webhooks "github.com/formancehq/webhooks/pkg/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Webhooks}, func() {
	var (
		secret1     = webhooks.NewSecret()
		secret2 	= webhooks.NewSecret()
		hook1 shared.V2Hook
	)

	BeforeEach(func() {
		hookBodyParam := shared.V2HookBodyParams{
			Endpoint: "https://example.com",
			Secret:   &secret1,
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

	Context("changing the secret of the inserted one", func() {
		
			It("should work with no secret provided",func() {
				response, err := Client().Webhooks.UpdateSecretHook(
					TestContext(),
					operations.UpdateSecretHookRequest{
						RequestBody: operations.UpdateSecretHookRequestBody{
						},
						HookID: hook1.ID,
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				Expect(response.V2HookResponse.Data.Secret).To(Not(Equal(hook1.Secret)))
			})
			
			
			It("should work with a secret provided",func() {
				response, err := Client().Webhooks.UpdateSecretHook(
					TestContext(),
					operations.UpdateSecretHookRequest{
						RequestBody: operations.UpdateSecretHookRequestBody{
							Secret: &secret2,
						},
						HookID: hook1.ID,
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				Expect(response.V2HookResponse.Data.Secret).To(Equal(secret2))
			})
			
			It("should not work with an invalid secret provided",func() {
				_, err := Client().Webhooks.UpdateSecretHook(
					TestContext(),
					operations.UpdateSecretHookRequest{
						RequestBody: operations.UpdateSecretHookRequestBody{
							Secret: ptr("invalid_secret"),
						},
						HookID: hook1.ID,
					},
				)
				Expect(err).To(HaveOccurred())
				Expect(err.(*sdkerrors.WebhooksErrorResponse).ErrorCode).To(Equal(shared.WebhooksErrorsEnumValidationType))
				
			})
		

	})
})
