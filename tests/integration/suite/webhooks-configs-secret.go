package suite

import (
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"net/http"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
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
		response, err := Client().Webhooks.V1.InsertConfig(
			TestContext(),
			cfg,
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusOK))

		insertResp = response.ConfigResponse
	})

	Context("changing the secret of the inserted one", func() {
		Context("without passing a secret", func() {
			BeforeEach(func() {
				response, err := Client().Webhooks.V1.ChangeConfigSecret(
					TestContext(),
					operations.ChangeConfigSecretRequest{
						ConfigChangeSecret: &shared.ConfigChangeSecret{
							Secret: "",
						},
						ID: insertResp.Data.ID,
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				Expect(response.ConfigResponse.Data.Secret).To(Not(Equal(insertResp.Data.Secret)))
			})

			Context("getting all configs", func() {
				It("should return 1 config with a different secret", func() {
					response, err := Client().Webhooks.V1.GetManyConfigs(
						TestContext(),
						operations.GetManyConfigsRequest{},
					)
					Expect(err).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(http.StatusOK))

					resp := response.ConfigsResponse
					Expect(resp.Cursor.HasMore).To(BeFalse())
					Expect(resp.Cursor.Data).To(HaveLen(1))
					Expect(resp.Cursor.Data[0].Secret).To(Not(BeNil()))
					Expect(resp.Cursor.Data[0].Secret).To(Not(Equal(insertResp.Data.Secret)))
				})
			})
		})

		Context("bringing our own valid secret", func() {
			newSecret := webhooks.NewSecret()
			BeforeEach(func() {
				response, err := Client().Webhooks.V1.ChangeConfigSecret(
					TestContext(),
					operations.ChangeConfigSecretRequest{
						ConfigChangeSecret: &shared.ConfigChangeSecret{
							Secret: newSecret,
						},
						ID: insertResp.Data.ID,
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))

				Expect(response.ConfigResponse.Data.Secret).To(Equal(newSecret))
			})

			Context("getting all configs", func() {
				It("should return 1 config with the passed secret", func() {
					response, err := Client().Webhooks.V1.GetManyConfigs(
						TestContext(),
						operations.GetManyConfigsRequest{},
					)
					Expect(err).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(http.StatusOK))

					resp := response.ConfigsResponse
					Expect(resp.Cursor.HasMore).To(BeFalse())
					Expect(resp.Cursor.Data).To(HaveLen(1))
					Expect(resp.Cursor.Data[0].Secret).To(Equal(newSecret))
				})
			})
		})

		Context("bringing our own invalid secret", func() {
			invalidSecret := "invalid"
			It("should return a bad request error", func() {
				_, err := Client().Webhooks.V1.ChangeConfigSecret(
					TestContext(),
					operations.ChangeConfigSecretRequest{
						ConfigChangeSecret: &shared.ConfigChangeSecret{
							Secret: invalidSecret,
						},
						ID: insertResp.Data.ID,
					},
				)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
