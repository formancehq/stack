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

	Context("changing the secret of the inserted one", func() {
		Context("without passing a secret", func() {
			BeforeEach(func() {
				resp, httpResp, err := Client().WebhooksApi.
					ChangeConfigSecret(TestContext(), insertResp.Data.Id).
					ConfigChangeSecret(formance.ConfigChangeSecret{
						Secret: "",
					}).Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Data.Secret).To(Not(Equal(insertResp.Data.Secret)))
			})

			Context("getting all configs", func() {
				It("should return 1 config with a different secret", func() {
					resp, _, err := Client().WebhooksApi.
						GetManyConfigs(TestContext()).Execute()
					Expect(err).NotTo(HaveOccurred())
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
				resp, httpResp, err := Client().WebhooksApi.
					ChangeConfigSecret(TestContext(), insertResp.Data.Id).
					ConfigChangeSecret(formance.ConfigChangeSecret{
						Secret: newSecret,
					}).Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
				Expect(resp.Data.Secret).To(Equal(newSecret))
			})

			Context("getting all configs", func() {
				It("should return 1 config with the passed secret", func() {
					resp, _, err := Client().WebhooksApi.
						GetManyConfigs(TestContext()).Execute()
					Expect(err).NotTo(HaveOccurred())
					Expect(resp.Cursor.HasMore).To(BeFalse())
					Expect(resp.Cursor.Data).To(HaveLen(1))
					Expect(resp.Cursor.Data[0].Secret).To(Equal(newSecret))
				})
			})
		})

		Context("bringing our own invalid secret", func() {
			invalidSecret := "invalid"
			It("should return a bad request error", func() {
				resp, httpResp, err := Client().WebhooksApi.
					ChangeConfigSecret(TestContext(), insertResp.Data.Id).
					ConfigChangeSecret(formance.ConfigChangeSecret{
						Secret: invalidSecret,
					}).Execute()
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(httpResp.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})
	})
})
