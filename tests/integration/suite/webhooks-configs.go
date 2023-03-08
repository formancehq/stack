package suite

import (
	"context"
	"net/http"

	"github.com/formancehq/formance-sdk-go"
	. "github.com/formancehq/stack/tests/integration/internal"
	webhooks "github.com/formancehq/webhooks/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("empty environment for webhooks configs", func() {
	When("inserting 0 config", func() {
		Context("getting all configs without filters", func() {
			It("should return 0 config", func() {
				resp, httpResp, err := Client().WebhooksApi.
					GetManyConfigs(context.Background()).Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Cursor.HasMore).To(BeFalse())
				Expect(resp.Cursor.Data).To(BeEmpty())
				Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})

	When("inserting 1 config", func() {
		var insertResp *formance.ConfigResponse
		var httpResp *http.Response
		var err error
		BeforeEach(func() {
			insertResp, httpResp, err = Client().WebhooksApi.
				InsertConfig(TestContext()).ConfigUser(formance.ConfigUser{
				Endpoint: "https://example.com",
				EventTypes: []string{
					"ledger.committed_transactions",
				},
			}).Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
		})

		Context("getting all configs without filters", func() {
			It("should return 1 config", func() {
				resp, _, err := Client().WebhooksApi.
					GetManyConfigs(context.Background()).Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Cursor.HasMore).To(BeFalse())
				Expect(resp.Cursor.Data).To(HaveLen(1))
				Expect(*resp.Cursor.Data[0].Endpoint).To(Equal("https://example.com"))
				Expect(*resp.Cursor.Data[0].Active).To(BeTrue())
			})
		})

		Context("getting all configs with known endpoint filter", func() {
			It("should return 1 config with the same endpoint", func() {
				resp, _, err := Client().WebhooksApi.
					GetManyConfigs(context.Background()).
					Endpoint("https://example.com").Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Cursor.HasMore).To(BeFalse())
				Expect(resp.Cursor.Data).To(HaveLen(1))
				Expect(*resp.Cursor.Data[0].Endpoint).To(Equal("https://example.com"))
			})
		})

		Context("getting all configs with unknown endpoint filter", func() {
			It("should return 0 config", func() {
				resp, _, err := Client().WebhooksApi.
					GetManyConfigs(context.Background()).
					Endpoint("https://unknown.com").Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Cursor.HasMore).To(BeFalse())
				Expect(resp.Cursor.Data).To(BeEmpty())
			})
		})

		Context("deactivating the inserted one", func() {
			BeforeEach(func() {
				_, _, err = Client().WebhooksApi.DeactivateConfig(TestContext(),
					*insertResp.Data.Id).Execute()
				Expect(err).NotTo(HaveOccurred())
			})

			Context("getting all configs", func() {
				It("should return 1 deactivated config", func() {
					resp, _, err := Client().WebhooksApi.
						GetManyConfigs(context.Background()).Execute()
					Expect(err).NotTo(HaveOccurred())
					Expect(resp.Cursor.HasMore).To(BeFalse())
					Expect(resp.Cursor.Data).To(HaveLen(1))
					Expect(*resp.Cursor.Data[0].Active).To(BeFalse())
				})
			})
		})

		Context("deactivating the inserted one, then reactivating it", func() {
			BeforeEach(func() {
				_, _, err = Client().WebhooksApi.DeactivateConfig(TestContext(),
					*insertResp.Data.Id).Execute()
				Expect(err).NotTo(HaveOccurred())

				_, _, err = Client().WebhooksApi.ActivateConfig(TestContext(),
					*insertResp.Data.Id).
					Execute()
				Expect(err).NotTo(HaveOccurred())
			})

			Context("getting all configs", func() {
				It("should return 1 activated config", func() {
					resp, _, err := Client().WebhooksApi.
						GetManyConfigs(context.Background()).Execute()
					Expect(err).NotTo(HaveOccurred())
					Expect(resp.Cursor.HasMore).To(BeFalse())
					Expect(resp.Cursor.Data).To(HaveLen(1))
					Expect(*resp.Cursor.Data[0].Active).To(BeTrue())
				})
			})
		})

		Context("changing the secret of the inserted one", func() {
			Context("without passing a secret", func() {
				BeforeEach(func() {
					resp, _, err := Client().WebhooksApi.
						ChangeConfigSecret(TestContext(), *insertResp.Data.Id).Execute()
					Expect(err).NotTo(HaveOccurred())
					Expect(*resp.Data.Secret).To(Not(Equal(*insertResp.Data.Secret)))
				})

				Context("getting all configs", func() {
					It("should return 1 config with a different secret", func() {
						resp, _, err := Client().WebhooksApi.
							GetManyConfigs(context.Background()).Execute()
						Expect(err).NotTo(HaveOccurred())
						Expect(resp.Cursor.HasMore).To(BeFalse())
						Expect(resp.Cursor.Data).To(HaveLen(1))
						Expect(*resp.Cursor.Data[0].Secret).To(Not(BeNil()))
						Expect(*resp.Cursor.Data[0].Secret).To(Not(Equal(*insertResp.Data.Secret)))
					})
				})
			})

			Context("bringing our own secret", func() {
				validSecret := webhooks.NewSecret()
				BeforeEach(func() {
					resp, _, err := Client().WebhooksApi.
						ChangeConfigSecret(TestContext(), *insertResp.Data.Id).
						ConfigChangeSecret(formance.ConfigChangeSecret{
							Secret: &validSecret,
						}).Execute()
					Expect(err).NotTo(HaveOccurred())
					Expect(*resp.Data.Secret).To(Equal(validSecret))
				})

				Context("getting all configs", func() {
					It("should return 1 config with the passed secret", func() {
						resp, _, err := Client().WebhooksApi.
							GetManyConfigs(context.Background()).Execute()
						Expect(err).NotTo(HaveOccurred())
						Expect(resp.Cursor.HasMore).To(BeFalse())
						Expect(resp.Cursor.Data).To(HaveLen(1))
						Expect(*resp.Cursor.Data[0].Secret).To(Equal(validSecret))
					})
				})
			})
		})
	})
})
