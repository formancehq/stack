package suite

import (
	"net/http"

	"github.com/formancehq/formance-sdk-go"
	. "github.com/formancehq/stack/tests/integration/internal"
	webhooks "github.com/formancehq/webhooks/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("empty environment for webhooks configs", func() {
	var (
		secret     = webhooks.NewSecret()
		insertResp *formance.ConfigResponse
		httpResp   *http.Response
		err        error
	)

	BeforeEach(func() {
		cfg := formance.ConfigUser{
			Endpoint: "https://example.com",
			Secret:   &secret,
			EventTypes: []string{
				"ledger.committed_transactions",
			},
		}
		insertResp, httpResp, err = Client().WebhooksApi.
			InsertConfig(TestContext()).ConfigUser(cfg).Execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
	})

	Context("deactivating the inserted one", func() {
		BeforeEach(func() {
			resp, httpResp, err := Client().WebhooksApi.
				DeactivateConfig(TestContext(), insertResp.Data.Id).Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Data.Active).To(BeFalse())
		})

		Context("getting all configs", func() {
			It("should return 1 deactivated config", func() {
				resp, _, err := Client().WebhooksApi.
					GetManyConfigs(TestContext()).Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Cursor.Data).To(HaveLen(1))
				Expect(resp.Cursor.Data[0].Active).To(BeFalse())
			})
		})
	})

	Context("deactivating the inserted one, then reactivating it", func() {
		BeforeEach(func() {
			resp, httpResp, err := Client().WebhooksApi.
				DeactivateConfig(TestContext(), insertResp.Data.Id).Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Data.Active).To(BeFalse())

			resp, httpResp, err = Client().WebhooksApi.
				ActivateConfig(TestContext(), insertResp.Data.Id).Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Data.Active).To(BeTrue())
		})

		Context("getting all configs", func() {
			It("should return 1 activated config", func() {
				resp, _, err := Client().WebhooksApi.
					GetManyConfigs(TestContext()).Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Cursor.Data).To(HaveLen(1))
				Expect(resp.Cursor.Data[0].Active).To(BeTrue())
			})
		})
	})

	Context("trying to deactivate an unknown ID", func() {
		It("should fail", func() {
			resp, httpResp, err := Client().WebhooksApi.
				DeactivateConfig(TestContext(), "unknown").Execute()
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
			Expect(httpResp.StatusCode).To(Equal(http.StatusNotFound))
		})
	})
})
