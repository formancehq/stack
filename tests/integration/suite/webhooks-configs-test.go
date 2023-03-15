package suite

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/formancehq/formance-sdk-go"
	. "github.com/formancehq/stack/tests/integration/internal"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/security"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("empty environment for webhooks configs", func() {
	When("testing configs", func() {
		Context("inserting a config with an endpoint to a success handler", func() {
			var (
				httpServer *httptest.Server
				insertResp *formance.ConfigResponse
				secret     = webhooks.NewSecret()
				err        error
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

				cfg := formance.ConfigUser{
					Endpoint: httpServer.URL,
					Secret:   &secret,
					EventTypes: []string{
						"ledger.committed_transactions",
					},
				}
				insertResp, _, err = Client().WebhooksApi.
					InsertConfig(TestContext()).ConfigUser(cfg).Execute()
				Expect(err).ToNot(HaveOccurred())
				DeferCleanup(func() {
					httpServer.Close()
				})
			})

			Context("testing the inserted one", func() {
				It("should return a successful attempt", func() {
					attemptResp, httpResp, err := Client().WebhooksApi.
						TestConfig(TestContext(), insertResp.Data.Id).Execute()
					Expect(err).ToNot(HaveOccurred())
					Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
					Expect(attemptResp.Data.Config).To(Equal(insertResp.Data))
					Expect(attemptResp.Data.Payload).To(Equal(`{"data":"test"}`))
					Expect(int(attemptResp.Data.StatusCode)).To(Equal(http.StatusOK))
					Expect(attemptResp.Data.Status).To(Equal("success"))
				})
			})
		})

		Context("inserting a config with an endpoint to a fail handler", func() {
			var insertResp *formance.ConfigResponse

			BeforeEach(func() {
				httpServer := httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, _ *http.Request) {
						http.Error(w,
							"WEBHOOKS RECEIVED: MOCK ERROR RESPONSE", http.StatusNotFound)
					}))

				cfg := formance.ConfigUser{
					Endpoint: httpServer.URL,
					EventTypes: []string{
						"ledger.committed_transactions",
					},
				}
				var err error
				insertResp, _, err = Client().WebhooksApi.
					InsertConfig(TestContext()).ConfigUser(cfg).Execute()
				Expect(err).ToNot(HaveOccurred())
				DeferCleanup(func() {
					httpServer.Close()
				})
			})

			Context("testing the inserted one", func() {
				It("should return a failed attempt", func() {
					attemptResp, httpResp, err := Client().WebhooksApi.
						TestConfig(TestContext(), insertResp.Data.Id).Execute()
					Expect(err).ToNot(HaveOccurred())
					Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
					Expect(attemptResp.Data.Config.Id).To(Equal(insertResp.Data.Id))
					Expect(attemptResp.Data.Payload).To(Equal(`{"data":"test"}`))
					Expect(int(attemptResp.Data.StatusCode)).To(Equal(http.StatusNotFound))
					Expect(attemptResp.Data.Status).To(Equal("failed"))
				})
			})
		})

		Context("testing an unknown ID", func() {
			It("should fail", func() {
				attemptResp, httpResp, err := Client().WebhooksApi.
					TestConfig(TestContext(), "unknown").Execute()
				Expect(err).To(HaveOccurred())
				Expect(attemptResp).To(BeNil())
				Expect(httpResp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})
})
