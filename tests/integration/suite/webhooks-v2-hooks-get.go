package suite

import (
	"net/http"

	"github.com/formancehq/stack/tests/integration/internal/modules"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/sdkerrors"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Webhooks}, func() {
	
	It("should return 0 hook", func() {
		response, err := Client().Webhooks.GetManyHooks(
			TestContext(),
			operations.GetManyHooksRequest{},
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusOK))

		Expect(response.V2HookCursorResponse.Cursor.HasMore).To(BeFalse())
		Expect(response.V2HookCursorResponse.Cursor.Data).To(BeEmpty())
	})


	When("inserting  4 Hooks", func() {
		var (
			
			hook1 *shared.V2Hook
			hook2 *shared.V2Hook
			hook3 *shared.V2Hook
			hook4 *shared.V2Hook
			
		)

		BeforeEach(func() {
			var (
				err  error
				
				name1 = "Hook1"
				name2 = "Hook2"
				name3 = "Hook3"
				name4 = "Hook4"

				hookBody1 = shared.V2HookBodyParams{
					Name: &name1,
					Endpoint: "https://example1.com",
					Events: []string{
						"ledger.committed_transactions",
					},
				}
				hookBody2 = shared.V2HookBodyParams{
					Name: &name2,
					Endpoint: "https://example.hook2",
					Events: []string{
						"ledger.committed_transactions",
					},
				}
				hookBody3 = shared.V2HookBodyParams{
					Name: &name3,
					Endpoint: "https://example.hook2",
					Events: []string{
						"ledger.committed_transactions",
					},
				}
				hookBody4 = shared.V2HookBodyParams{
					Name: &name4,
					Endpoint: "https://example.hook4",
					Events: []string{
						"ledger.committed_transactions",
					},
				}
			)

			
			response1, err := Client().Webhooks.InsertHook(
				TestContext(),
				hookBody1,
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response1.StatusCode).To(Equal(http.StatusOK))
			hook1 = &response1.V2HookResponse.Data
			Expect(hook1.Name).To(Equal(*(hookBody1.Name)))

			response2, err := Client().Webhooks.InsertHook(
				TestContext(),
				hookBody2,
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response2.StatusCode).To(Equal(http.StatusOK))
			hook2 = &response2.V2HookResponse.Data
			Expect(hook2.Name).To(Equal(*(hookBody2.Name)))
			
			response3, err := Client().Webhooks.InsertHook(
				TestContext(),
				hookBody3,
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response3.StatusCode).To(Equal(http.StatusOK))
			hook3 = &response3.V2HookResponse.Data
			Expect(hook3.Name).To(Equal(*(hookBody3.Name)))

			response4, err := Client().Webhooks.InsertHook(
				TestContext(),
				hookBody4,
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response4.StatusCode).To(Equal(http.StatusOK))
			hook4 = &response4.V2HookResponse.Data
			Expect(hook4.Name).To(Equal(*(hookBody4.Name)))
		})

		Context("getting all  V2 hooks without filters", func() {
			
			It("should return 4 Hooks", func() {
				response, err := Client().Webhooks.GetManyHooks(
					TestContext(),
					operations.GetManyHooksRequest{},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))

				resp := response.V2HookCursorResponse
				Expect(resp.Cursor.HasMore).To(BeFalse())
				Expect(resp.Cursor.Data).To(HaveLen(4))

			})
		})

		Context("getting all chooks with known endpoint filter", func() {
			
			It("should return 2 hooks with the same endpoint", func() {
				response, err := Client().Webhooks.GetManyHooks(
					TestContext(),
					operations.GetManyHooksRequest{
						Endpoint: ptr(hook2.Endpoint),
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))

				cursor := response.V2HookCursorResponse.Cursor
				Expect(cursor.HasMore).To(BeFalse())
				Expect(cursor.Data).To(HaveLen(2))
				Expect(cursor.Data[0].Endpoint).To(Equal(hook2.Endpoint))
			})
		})

		Context("getting all hooks with unknown endpoint filter", func() {
			
			It("should return 0 hook", func() {
				response, err := Client().Webhooks.GetManyHooks(
					TestContext(),
					operations.GetManyHooksRequest{
						Endpoint: ptr("https://unknown.com"),
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))

				resp := response.V2HookCursorResponse.Cursor
				Expect(resp.HasMore).To(BeFalse())
				Expect(resp.Data).To(BeEmpty())
			})
		})


		Context("getting ONE hook", func() {
			It("should return 1 Hook with the same ID", func() {
				response, err := Client().Webhooks.GetHook(
					TestContext(),
					operations.GetHookRequest{
						HookID: hook1.ID,
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))

				resp := response.V2HookResponse
				Expect(resp.Data.ID).To(Equal(hook1.ID))
			})
			It("should return error because false ID", func(){
				_, err := Client().Webhooks.GetHook(
					TestContext(),
					operations.GetHookRequest{
						HookID: "BADID",
					},
				)
				Expect(err).To(HaveOccurred())
				Expect(err.(*sdkerrors.WebhooksErrorResponse).ErrorCode).To(Equal(shared.WebhooksErrorsEnumNotFound))
			})
		})

	})
})
