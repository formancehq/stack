package suite

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/sdkerrors"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Webhooks}, func() {
	var (
		hook1 		shared.V2Hook
	)

	BeforeEach(func() {
	
		hookBodyParams := shared.V2HookBodyParams{
			Endpoint: "https://example.com",
			Name: ptr("Hook1"),
			Events: []string{
				"ledger.committed_transactions",
			},
		}

		resp1, err := Client().Webhooks.InsertHook(
			TestContext(),
			hookBodyParams,
		)
		Expect(err).NotTo(HaveOccurred())
		hook1 = resp1.V2HookResponse.Data

	})

	Context("Activate and Deactive Hook1", func(){
		It("should be activate", func(){
			response, err := Client().Webhooks.ActivateHook(
				TestContext(),
				operations.ActivateHookRequest{
					HookID: hook1.ID,
				},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(response.V2HookResponse.Data.Status).To(Equal(shared.V2HookStatusEnabled))
		})

		It("should be deactivate", func(){
			response, err := Client().Webhooks.DeactivateHook(
				TestContext(),
				operations.DeactivateHookRequest{
					HookID: hook1.ID,
				},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(response.V2HookResponse.Data.Status).To(Equal(shared.V2HookStatusDisabled))
		})

		It("Should return One hook deactivate", func(){
			response, err := Client().Webhooks.GetManyHooks(
				TestContext(),
				operations.GetManyHooksRequest{},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(response.V2HookCursorResponse.Cursor.Data[0].Status).To(Equal(shared.V2HookStatusDisabled))
		})

		
	})

	Context("Activate and Deactivate a bad ID Hook", func(){
		It("Activate should failed", func(){
			_, err := Client().Webhooks.ActivateHook(
				TestContext(),
				operations.ActivateHookRequest{
					HookID: "UNKNOWN",
				},
			)
			Expect(err).To(HaveOccurred())
			Expect(err.(*sdkerrors.WebhooksErrorResponse).ErrorCode).To(Equal(shared.WebhooksErrorsEnumNotFound))
		})

		It("Deactivate should failed", func(){
			_, err := Client().Webhooks.DeactivateHook(
				TestContext(),
				operations.DeactivateHookRequest{
					HookID: "UNKNOWN",
				},
			)
			Expect(err).To(HaveOccurred())
			Expect(err.(*sdkerrors.WebhooksErrorResponse).ErrorCode).To(Equal(shared.WebhooksErrorsEnumNotFound))
		})
	})

})
