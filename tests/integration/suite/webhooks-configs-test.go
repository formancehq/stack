package suite

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/sdkerrors"
	"github.com/formancehq/stack/tests/integration/internal/modules"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/security"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Webhooks}, func() {
	When("testing configs", func() {
		Context("inserting a config with an endpoint to a success handler", func() {
			var (
				httpServer *httptest.Server
				insertResp *shared.ConfigResponse
				secret     = webhooks.NewSecret()
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

						ok, err := security.Verify(signatures, id, timeInt, secret, payload)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
						if !ok {
							http.Error(w, "WEBHOOKS SIGNATURE VERIFICATION NOK", http.StatusBadRequest)
							return
						}
					}))

				cfg := shared.ConfigUser{
					Endpoint: httpServer.URL,
					Secret:   &secret,
					EventTypes: []string{
						"ledger.committed_transactions",
					},
				}
				response, err := Client().Webhooks.V1.InsertConfig(
					TestContext(),
					cfg,
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				insertResp = response.ConfigResponse
				DeferCleanup(func() {
					httpServer.Close()
				})
			})

			Context("testing the inserted one", func() {
				It("should return a successful attempt", func() {
					response, err := Client().Webhooks.V1.TestConfig(
						TestContext(),
						operations.TestConfigRequest{
							ID: insertResp.Data.ID,
						},
					)
					Expect(err).ToNot(HaveOccurred())
					Expect(response.StatusCode).To(Equal(http.StatusOK))

					attemptResp := response.AttemptResponse
					Expect(attemptResp.Data.Config.ID).To(Equal(insertResp.Data.ID))
					Expect(attemptResp.Data.Payload).To(Equal(`{"data":"test"}`))
					Expect(int(attemptResp.Data.StatusCode)).To(Equal(http.StatusOK))
					Expect(attemptResp.Data.Status).To(Equal("success"))
				})
			})
		})

		Context("inserting a config with an endpoint to a fail handler", func() {
			var insertResp *shared.ConfigResponse

			BeforeEach(func() {
				httpServer := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, _ *http.Request) {
						http.Error(w,
							"WEBHOOKS RECEIVED: MOCK ERROR RESPONSE", http.StatusNotFound)
					}))

				cfg := shared.ConfigUser{
					Endpoint: httpServer.URL,
					EventTypes: []string{
						"ledger.committed_transactions",
					},
				}
				response, err := Client().Webhooks.V1.InsertConfig(
					TestContext(),
					cfg,
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				insertResp = response.ConfigResponse
				DeferCleanup(func() {
					httpServer.Close()
				})
			})

			Context("testing the inserted one", func() {
				It("should return a failed attempt", func() {
					response, err := Client().Webhooks.V1.TestConfig(
						TestContext(),
						operations.TestConfigRequest{
							ID: insertResp.Data.ID,
						},
					)
					Expect(err).ToNot(HaveOccurred())
					Expect(response.StatusCode).To(Equal(http.StatusOK))

					attemptResp := response.AttemptResponse
					Expect(attemptResp.Data.Config.ID).To(Equal(insertResp.Data.ID))
					Expect(attemptResp.Data.Payload).To(Equal(`{"data":"test"}`))
					Expect(int(attemptResp.Data.StatusCode)).To(Equal(http.StatusNotFound))
					Expect(attemptResp.Data.Status).To(Equal("failed"))
				})
			})
		})

		Context("testing an unknown ID", func() {
			It("should fail", func() {
				_, err := Client().Webhooks.V1.TestConfig(
					TestContext(),
					operations.TestConfigRequest{
						ID: "unknown",
					},
				)
				Expect(err).To(HaveOccurred())
				Expect(err.(*sdkerrors.WebhooksErrorResponse).ErrorCode).To(Equal(shared.WebhooksErrorsEnumNotFound))
			})
		})
	})
})
