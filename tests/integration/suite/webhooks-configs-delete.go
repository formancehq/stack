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
		Expect(err).ToNot(HaveOccurred())
		Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
	})

	Context("deleting the inserted one", func() {
		BeforeEach(func() {
			httpResp, err := Client().WebhooksApi.
				DeleteConfig(TestContext(), insertResp.Data.Id).Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
		})

		Context("getting all configs", func() {
			It("should return 0 config", func() {
				resp, _, err := Client().WebhooksApi.
					GetManyConfigs(TestContext()).Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Cursor.HasMore).To(BeFalse())
				Expect(resp.Cursor.Data).To(BeEmpty())
			})
		})

		AfterEach(func() {
			httpResp, err := Client().WebhooksApi.
				DeleteConfig(TestContext(), insertResp.Data.Id).Execute()
			Expect(err).To(HaveOccurred())
			Expect(httpResp.StatusCode).To(Equal(http.StatusNotFound))
		})
	})

	Context("trying to delete an unknown ID", func() {
		It("should fail", func() {
			httpResp, err := Client().WebhooksApi.
				DeleteConfig(TestContext(), "unknown").Execute()
			Expect(err).To(HaveOccurred())
			Expect(httpResp.StatusCode).To(Equal(http.StatusNotFound))
		})
	})
})
