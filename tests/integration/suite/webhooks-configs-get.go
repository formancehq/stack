package suite

import (
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"net/http"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Webhooks}, func() {
	It("should return 0 config", func() {
		response, err := Client().Webhooks.V1.GetManyConfigs(
			TestContext(),
			operations.GetManyConfigsRequest{},
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusOK))

		Expect(response.ConfigsResponse.Cursor.HasMore).To(BeFalse())
		Expect(response.ConfigsResponse.Cursor.Data).To(BeEmpty())
	})

	When("inserting 2 configs", func() {
		var (
			insertResp1 *shared.ConfigResponse
			insertResp2 *shared.ConfigResponse
		)

		BeforeEach(func() {
			var (
				err  error
				cfg1 = shared.ConfigUser{
					Endpoint: "https://example1.com",
					EventTypes: []string{
						"ledger.committed_transactions",
					},
				}
				cfg2 = shared.ConfigUser{
					Endpoint: "https://example2.com",
					EventTypes: []string{
						"ledger.saved_metadata",
					},
				}
			)

			response, err := Client().Webhooks.V1.InsertConfig(
				TestContext(),
				cfg1,
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			insertResp1 = response.ConfigResponse

			response, err = Client().Webhooks.V1.InsertConfig(
				TestContext(),
				cfg2,
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			insertResp2 = response.ConfigResponse
		})

		Context("getting all configs without filters", func() {
			It("should return 2 configs", func() {
				response, err := Client().Webhooks.V1.GetManyConfigs(
					TestContext(),
					operations.GetManyConfigsRequest{},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))

				resp := response.ConfigsResponse
				Expect(resp.Cursor.HasMore).To(BeFalse())
				Expect(resp.Cursor.Data).To(HaveLen(2))
				Expect(resp.Cursor.Data[0].Endpoint).To(Equal(insertResp2.Data.Endpoint))
				Expect(resp.Cursor.Data[1].Endpoint).To(Equal(insertResp1.Data.Endpoint))
			})
		})

		Context("getting all configs with known endpoint filter", func() {
			It("should return 1 config with the same endpoint", func() {
				response, err := Client().Webhooks.V1.GetManyConfigs(
					TestContext(),
					operations.GetManyConfigsRequest{
						Endpoint: ptr(insertResp1.Data.Endpoint),
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))

				resp := response.ConfigsResponse
				Expect(resp.Cursor.HasMore).To(BeFalse())
				Expect(resp.Cursor.Data).To(HaveLen(1))
				Expect(resp.Cursor.Data[0].Endpoint).To(Equal(insertResp1.Data.Endpoint))
			})
		})

		Context("getting all configs with unknown endpoint filter", func() {
			It("should return 0 config", func() {
				response, err := Client().Webhooks.V1.GetManyConfigs(
					TestContext(),
					operations.GetManyConfigsRequest{
						Endpoint: ptr("https://unknown.com"),
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))

				resp := response.ConfigsResponse
				Expect(resp.Cursor.HasMore).To(BeFalse())
				Expect(resp.Cursor.Data).To(BeEmpty())
			})
		})

		Context("getting all configs with known ID filter", func() {
			It("should return 1 config with the same ID", func() {
				response, err := Client().Webhooks.V1.GetManyConfigs(
					TestContext(),
					operations.GetManyConfigsRequest{
						ID: ptr(insertResp1.Data.ID),
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))

				resp := response.ConfigsResponse
				Expect(resp.Cursor.HasMore).To(BeFalse())
				Expect(resp.Cursor.Data).To(HaveLen(1))
				Expect(resp.Cursor.Data[0].ID).To(Equal(insertResp1.Data.ID))
			})
		})

		Context("getting all configs with unknown ID filter", func() {
			It("should return 0 config", func() {
				response, err := Client().Webhooks.V1.GetManyConfigs(
					TestContext(),
					operations.GetManyConfigsRequest{
						ID: ptr("unknown"),
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))

				resp := response.ConfigsResponse
				Expect(resp.Cursor.HasMore).To(BeFalse())
				Expect(resp.Cursor.Data).To(BeEmpty())
			})
		})
	})
})
