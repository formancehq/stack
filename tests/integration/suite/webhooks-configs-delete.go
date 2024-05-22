package suite

import (
	"net/http"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/sdkerrors"
	"github.com/formancehq/stack/tests/integration/internal/modules"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	webhooks "github.com/formancehq/webhooks/pkg/utils"
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

	Context("Config: deleting the inserted one", func() {
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
				for _, data := range response.ConfigsResponse.Cursor.Data {
					Expect(data.ID).ToNot(Equal(insertResp.Data.ID))
				}
				
			})
		})
	})

	Context("Config trying to delete an unknown ID", func() {
		It("should fail", func() {
			_, err := Client().Webhooks.DeleteConfig(
				TestContext(),
				operations.DeleteConfigRequest{
					ID: "unknown",
				},
			)
			Expect(err).To(HaveOccurred())
			Expect(err.(*sdkerrors.WebhooksErrorResponse).ErrorCode).To(Equal(shared.WebhooksErrorsEnumNotFound))
		})
	})
})
