package suite

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
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
		Expect(err).NotTo(HaveOccurred())
		Expect(response.StatusCode).To(Equal(200))

		insertResp = response.ConfigResponse
	})

	Context("deactivating the inserted one", func() {
		BeforeEach(func() {
			response, err := Client().Webhooks.DeactivateConfig(
				TestContext(),
				operations.DeactivateConfigRequest{
					ID: insertResp.Data.ID,
				},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			Expect(response.ConfigResponse.Data.Active).To(BeFalse())
		})

		Context("getting all configs", func() {
			It("should return 1 deactivated config", func() {
				response, err := Client().Webhooks.GetManyConfigs(
					TestContext(),
					operations.GetManyConfigsRequest{},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				Expect(response.ConfigsResponse.Cursor.Data).To(HaveLen(1))
				Expect(response.ConfigsResponse.Cursor.Data[0].Active).To(BeFalse())
			})
		})
	})

	Context("deactivating the inserted one, then reactivating it", func() {
		BeforeEach(func() {
			response, err := Client().Webhooks.DeactivateConfig(
				TestContext(),
				operations.DeactivateConfigRequest{
					ID: insertResp.Data.ID,
				},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.ConfigResponse.Data.Active).To(BeFalse())

			activateConfigResponse, err := Client().Webhooks.ActivateConfig(
				TestContext(),
				operations.ActivateConfigRequest{
					ID: insertResp.Data.ID,
				},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(activateConfigResponse.StatusCode).To(Equal(200))
			Expect(activateConfigResponse.ConfigResponse.Data.Active).To(BeTrue())
		})

		Context("getting all configs", func() {
			It("should return 1 activated config", func() {
				response, err := Client().Webhooks.GetManyConfigs(
					TestContext(),
					operations.GetManyConfigsRequest{},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				Expect(response.ConfigsResponse.Cursor.Data).To(HaveLen(1))
				Expect(response.ConfigsResponse.Cursor.Data[0].Active).To(BeTrue())
			})
		})
	})

	Context("trying to deactivate an unknown ID", func() {
		It("should fail", func() {
			response, err := Client().Webhooks.DeactivateConfig(
				TestContext(),
				operations.DeactivateConfigRequest{
					ID: "unknown",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(404))
			Expect(response.ConfigResponse).To(BeNil())
			Expect(response.ErrorResponse).ToNot(BeNil())
		})
	})
})
