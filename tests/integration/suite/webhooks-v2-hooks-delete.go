package suite

import (
	"net/http"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/sdkerrors"
	"github.com/formancehq/stack/tests/integration/internal/modules"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Webhooks}, func() {
	var (
		hook *shared.V2Hook
	)

	BeforeEach(func() {

		hookBodyParam := shared.V2HookBodyParams{
			Endpoint: "https://example1.com",
			Name: ptr("Test1"),
			Events: []string{
				"ledger.committed_transactions",
			},
		}

		resp, err := Client().Webhooks.InsertHook(
			TestContext(),
			hookBodyParam,
		)
		Expect(err).ToNot(HaveOccurred())
		hook = &resp.V2HookResponse.Data
		

	})

	Context("Hook: deleting the inserted one", func() {
		BeforeEach(func() {
			


			response, err := Client().Webhooks.DeleteHook(
				TestContext(),
				operations.DeleteHookRequest{
					HookID: hook.ID,
				},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(response.V2HookResponse.Data.Status).To(Equal(shared.V2HookStatusDeleted))
		})

		Context("Hook : getting all hooks", func() {
			It("None should be Hook1", func() {
				response, err := Client().Webhooks.GetManyHooks(
					TestContext(),
					operations.GetManyHooksRequest{},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))

				Expect(response.V2HookCursorResponse.Cursor.HasMore).To(BeFalse())
				for _, data := range response.V2HookCursorResponse.Cursor.Data {
					Expect(data.ID).ToNot(Equal(hook.ID))
				}
				
			})
		})
	})

	Context("Hook: trying to delete an unknown ID", func() {
		It("should fail", func() {
			_, err := Client().Webhooks.DeleteHook(
				TestContext(),
				operations.DeleteHookRequest{
					HookID: "unknown",
				},
			)
			Expect(err).To(HaveOccurred())
			Expect(err.(*sdkerrors.WebhooksErrorResponse).ErrorCode).To(Equal(shared.WebhooksErrorsEnumNotFound))
		})
	})
})
