package suite

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/sdkerrors"
	"github.com/formancehq/stack/tests/integration/internal/modules"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	webhookSecurity "github.com/formancehq/webhooks/pkg/security"
	webhooksUtils "github.com/formancehq/webhooks/pkg/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Webhooks}, func() {
	When("testing Hooks", func() {
		Context("inserting a HOok with an endpoint to a success handler", func() {
			var (
				httpServer *httptest.Server
				hook shared.V2Hook
				secret     = webhooksUtils.NewSecret()
				payload = "{\"Data\":\"payload_test\"}"
			)

			BeforeEach(func() {
				httpServer = httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						id := r.Header.Get("formance-webhook-id")
						ts := r.Header.Get("formance-webhook-timestamp")
						signatures := r.Header.Get("formance-webhook-signature")

						timeInt, err := strconv.ParseInt(ts, 10, 64)
						
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}

						payload, err := io.ReadAll(r.Body)
						
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
						
						ok, err := webhookSecurity.Verify(signatures, id, timeInt, secret, payload)
						if err != nil {
							
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
						if !ok {
							
							http.Error(w, "WEBHOOKS SIGNATURE VERIFICATION NOK", http.StatusBadRequest)
							return
						}
					}))

				hookBodyParam := shared.V2HookBodyParams{
					Endpoint: httpServer.URL,
					Secret:   &secret,
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
				hook = response.V2HookResponse.Data
				DeferCleanup(func() {
					httpServer.Close()
				})
			})

			Context("testing the inserted one", func() {
				It("should return a successful attempt", func() {
					response, err := Client().Webhooks.TestHook(
						TestContext(),
						operations.TestHookRequest{
							HookID: hook.ID,
							RequestBody: operations.TestHookRequestBody{
								Payload: &payload,
							},
							
						},
					)
					Expect(err).ToNot(HaveOccurred())
					Expect(response.StatusCode).To(Equal(http.StatusOK))

					attemptResp := response.V2AttemptResponse
					Expect(attemptResp.Data.HookID).To(Equal(hook.ID))
					Expect(attemptResp.Data.Payload).To(Equal(payload))
					Expect(int(attemptResp.Data.StatusCode)).To(Equal(http.StatusOK))
					Expect(attemptResp.Data.Status).To(Equal(shared.V2AttemptStatusSuccess))
				})
			})
		})

		Context("inserting a hook with an endpoint to a fail handler", func() {
			var hook2 *shared.V2Hook
			var payload = "{\"Data\":\"payload_test\"}"
			BeforeEach(func() {
				httpServer := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, _ *http.Request) {
						http.Error(w,
							"WEBHOOKS RECEIVED: MOCK ERROR RESPONSE", http.StatusNotFound)
					}))

				hookBodyParam := shared.V2HookBodyParams{
					Endpoint: httpServer.URL,
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
				hook2 = &response.V2HookResponse.Data
				DeferCleanup(func() {
					httpServer.Close()
				})
			})

			Context("testing the inserted one", func() {
				It("should return a failed attempt", func() {
					response, err := Client().Webhooks.TestHook(
						TestContext(),
						operations.TestHookRequest{
							HookID: hook2.ID,
							RequestBody: operations.TestHookRequestBody{
								Payload: ptr(payload),
							},
						},
					)
					Expect(err).ToNot(HaveOccurred())
					Expect(response.StatusCode).To(Equal(http.StatusOK))

					attemptResp := response.V2AttemptResponse
					Expect(attemptResp.Data.HookID).To(Equal(hook2.ID))
					Expect(attemptResp.Data.Payload).To(Equal(payload))
					Expect(int(attemptResp.Data.StatusCode)).To(Equal(http.StatusNotFound))
					Expect(attemptResp.Data.Status).To(Equal(shared.V2AttemptStatusAbort))
				})
			})
		})

		Context("testing an unknown ID", func() {
			It("should fail", func() {
				_, err := Client().Webhooks.TestHook(
					TestContext(),
					operations.TestHookRequest{
						HookID: "unknown",
					},
				)
				Expect(err).To(HaveOccurred())
				Expect(err.(*sdkerrors.WebhooksErrorResponse).ErrorCode).To(Equal(shared.WebhooksErrorsEnumNotFound))
			})
		})
	})
})
