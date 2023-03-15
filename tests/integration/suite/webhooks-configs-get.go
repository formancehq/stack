package suite

import (
	"net/http"

	"github.com/formancehq/formance-sdk-go"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("empty environment for webhooks configs", func() {
	It("should return 0 config", func() {
		resp, httpResp, err := Client().WebhooksApi.
			GetManyConfigs(TestContext()).Execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Cursor.HasMore).To(BeFalse())
		Expect(resp.Cursor.Data).To(BeEmpty())
		Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
	})

	When("inserting 2 configs", func() {
		var (
			insertResp1 *formance.ConfigResponse
			insertResp2 *formance.ConfigResponse
		)

		BeforeEach(func() {
			var (
				httpResp *http.Response
				err      error
				cfg1     = formance.ConfigUser{
					Endpoint: "https://example1.com",
					EventTypes: []string{
						"ledger.committed_transactions",
					},
				}
				cfg2 = formance.ConfigUser{
					Endpoint: "https://example2.com",
					EventTypes: []string{
						"ledger.saved_metadata",
					},
				}
			)

			insertResp1, httpResp, err = Client().WebhooksApi.
				InsertConfig(TestContext()).ConfigUser(cfg1).Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(httpResp.StatusCode).To(Equal(http.StatusOK))

			insertResp2, httpResp, err = Client().WebhooksApi.
				InsertConfig(TestContext()).ConfigUser(cfg2).Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
		})

		Context("getting all configs without filters", func() {
			It("should return 2 configs", func() {
				resp, _, err := Client().WebhooksApi.
					GetManyConfigs(TestContext()).Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Cursor.HasMore).To(BeFalse())
				Expect(resp.Cursor.Data).To(HaveLen(2))
				Expect(resp.Cursor.Data[0].Endpoint).To(Equal(insertResp2.Data.Endpoint))
				Expect(resp.Cursor.Data[1].Endpoint).To(Equal(insertResp1.Data.Endpoint))
			})
		})

		Context("getting all configs with known endpoint filter", func() {
			It("should return 1 config with the same endpoint", func() {
				resp, _, err := Client().WebhooksApi.
					GetManyConfigs(TestContext()).
					Endpoint(insertResp1.Data.Endpoint).Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Cursor.HasMore).To(BeFalse())
				Expect(resp.Cursor.Data).To(HaveLen(1))
				Expect(resp.Cursor.Data[0].Endpoint).To(Equal(insertResp1.Data.Endpoint))
			})
		})

		Context("getting all configs with unknown endpoint filter", func() {
			It("should return 0 config", func() {
				resp, _, err := Client().WebhooksApi.
					GetManyConfigs(TestContext()).
					Endpoint("https://unknown.com").Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Cursor.HasMore).To(BeFalse())
				Expect(resp.Cursor.Data).To(BeEmpty())
			})
		})

		Context("getting all configs with known ID filter", func() {
			It("should return 1 config with the same ID", func() {
				resp, _, err := Client().WebhooksApi.
					GetManyConfigs(TestContext()).
					Id(insertResp1.Data.Id).Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Cursor.HasMore).To(BeFalse())
				Expect(resp.Cursor.Data).To(HaveLen(1))
				Expect(resp.Cursor.Data[0].Id).To(Equal(insertResp1.Data.Id))
			})
		})

		Context("getting all configs with unknown ID filter", func() {
			It("should return 0 config", func() {
				resp, _, err := Client().WebhooksApi.
					GetManyConfigs(TestContext()).
					Id("unknown").Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Cursor.HasMore).To(BeFalse())
				Expect(resp.Cursor.Data).To(BeEmpty())
			})
		})
	})
})
