package suite

import (
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"net/http"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	webhooks "github.com/formancehq/webhooks/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Webhooks}, func() {
	var (
		secret     = webhooks.NewSecret()
		insertResp *shared.ConfigResponse
	)

	BeforeEach(func() {
		cfg := shared.ConfigUser{
			Endpoint: "https://example.com",
			Secret:   &secret,
			EventTypes: []string{
				"ledger.committed_transactions",
			},
		}
		response, err := Client().Webhooks.InsertConfig(
			TestContext(),
			cfg,
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusOK))

		insertResp = response.ConfigResponse
	})

	Context("deleting the inserted one", func() {
		BeforeEach(func() {
			response, err := Client().Webhooks.DeleteConfig(
				TestContext(),
				operations.DeleteConfigRequest{
					ID: insertResp.Data.ID,
				},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})

		Context("getting all configs", func() {
			It("should return 0 config", func() {
				response, err := Client().Webhooks.GetManyConfigs(
					TestContext(),
					operations.GetManyConfigsRequest{},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))

				Expect(response.ConfigsResponse.Cursor.HasMore).To(BeFalse())
				Expect(response.ConfigsResponse.Cursor.Data).To(BeEmpty())
			})
		})

		AfterEach(func() {
			response, err := Client().Webhooks.DeleteConfig(
				TestContext(),
				operations.DeleteConfigRequest{
					ID: insertResp.Data.ID,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusNotFound))
			Expect(response.WebhooksErrorResponse).ToNot(BeNil())
		})
	})

	Context("trying to delete an unknown ID", func() {
		It("should fail", func() {
			response, err := Client().Webhooks.DeleteConfig(
				TestContext(),
				operations.DeleteConfigRequest{
					ID: "unknown",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusNotFound))
			Expect(response.WebhooksErrorResponse).ToNot(BeNil())
		})
	})
})
