package suite

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/formancehq/formance-sdk-go"
	. "github.com/formancehq/stack/tests/integration/internal"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/security"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("empty environment for webhooks configs", func() {
	When("inserting 0 config", func() {
		Context("getting all configs without filters", func() {
			It("should return 0 config", func() {
				resp, httpResp, err := Client().WebhooksApi.
					GetManyConfigs(TestContext()).Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Cursor.HasMore).To(BeFalse())
				Expect(resp.Cursor.Data).To(BeEmpty())
				Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
			})
		})
	})

	When("inserting 1 config", func() {
		secret := webhooks.NewSecret()
		var insertResp *formance.ConfigResponse
		BeforeEach(func() {
			var err error
			var httpResp *http.Response
			cfg := formance.ConfigUser{
				Endpoint: "https://example.com",
				Secret:   &secret,
				EventTypes: []string{
					"ledger.committed_transactions",
				},
			}
			insertResp, httpResp, err = Client().WebhooksApi.
				InsertConfig(TestContext()).ConfigUser(cfg).Execute()
			Expect(err).To(Not(HaveOccurred()))
			Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
			Expect(insertResp.Data.Endpoint).To(Equal(cfg.Endpoint))
			Expect(insertResp.Data.EventTypes).To(Equal(cfg.EventTypes))
			Expect(insertResp.Data.Active).To(BeTrue())
			Expect(time.Until(insertResp.Data.CreatedAt)).To(BeNumerically("<", time.Second))
			Expect(time.Until(insertResp.Data.UpdatedAt)).To(BeNumerically("<", time.Second))
			_, err = uuid.Parse(insertResp.Data.Id)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("getting all configs without filters", func() {
			It("should return 1 config", func() {
				resp, _, err := Client().WebhooksApi.
					GetManyConfigs(TestContext()).Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Cursor.HasMore).To(BeFalse())
				Expect(resp.Cursor.Data).To(HaveLen(1))
				Expect(resp.Cursor.Data[0].Endpoint).To(Equal(insertResp.Data.Endpoint))
				Expect(resp.Cursor.Data[0].Active).To(BeTrue())
			})
		})

		Context("getting all configs with known endpoint filter", func() {
			It("should return 1 config with the same endpoint", func() {
				resp, _, err := Client().WebhooksApi.
					GetManyConfigs(TestContext()).
					Endpoint(insertResp.Data.Endpoint).Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Cursor.HasMore).To(BeFalse())
				Expect(resp.Cursor.Data).To(HaveLen(1))
				Expect(resp.Cursor.Data[0].Endpoint).To(Equal(insertResp.Data.Endpoint))
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
					Id(insertResp.Data.Id).Execute()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Cursor.HasMore).To(BeFalse())
				Expect(resp.Cursor.Data).To(HaveLen(1))
				Expect(resp.Cursor.Data[0].Id).To(Equal(insertResp.Data.Id))
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
					Expect(resp.Cursor.HasMore).To(BeFalse())
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
					Expect(resp.Cursor.HasMore).To(BeFalse())
					Expect(resp.Cursor.Data).To(HaveLen(1))
					Expect(resp.Cursor.Data[0].Active).To(BeTrue())
				})
			})
		})

		Context("changing the secret of the inserted one", func() {
			Context("without passing a secret", func() {
				BeforeEach(func() {
					resp, httpResp, err := Client().WebhooksApi.
						ChangeConfigSecret(TestContext(), insertResp.Data.Id).Execute()
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

			Context("bringing our own secret", func() {
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
	})

	When("testing configs", func() {
		Context("inserting a config with an endpoint to a success handler", func() {
			var httpServer *httptest.Server
			var insertResp *formance.ConfigResponse
			secret := webhooks.NewSecret()

			BeforeEach(func() {
				httpServer = httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						id := r.Header.Get("formance-webhook-id")
						ts := r.Header.Get("formance-webhook-timestamp")
						signatures := r.Header.Get("formance-webhook-signature")
						timeInt, err := strconv.ParseInt(ts, 10, 64)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}

						payload, err := io.ReadAll(r.Body)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}

						ok, err := security.Verify(signatures, id, timeInt, secret, payload)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
						if !ok {
							http.Error(w, "WEBHOOKS SIGNATURE VERIFICATION NOK", http.StatusBadRequest)
							return
						}
					}))

				cfg := formance.ConfigUser{
					Endpoint: httpServer.URL,
					Secret:   &secret,
					EventTypes: []string{
						"ledger.committed_transactions",
					},
				}
				var err error
				insertResp, _, err = Client().WebhooksApi.
					InsertConfig(TestContext()).ConfigUser(cfg).Execute()
				Expect(err).To(Not(HaveOccurred()))
				DeferCleanup(func() {
					httpServer.Close()
				})
			})

			Context("testing the inserted one", func() {
				It("should return a successful attempt", func() {
					attemptResp, httpResp, err := Client().WebhooksApi.
						TestConfig(TestContext(), insertResp.Data.Id).Execute()
					Expect(err).To(Not(HaveOccurred()))
					Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
					Expect(attemptResp.Data.Config.Id).To(Equal(insertResp.Data.Id))
					Expect(attemptResp.Data.Payload).To(Equal(`{"data":"test"}`))
					Expect(int(attemptResp.Data.StatusCode)).To(Equal(http.StatusOK))
					Expect(attemptResp.Data.Status).To(Equal("success"))
				})
			})
		})

		Context("inserting a config with an endpoint to a fail handler", func() {
			var httpServer *httptest.Server
			var insertResp *formance.ConfigResponse
			BeforeEach(func() {
				httpServer = httptest.NewServer(http.HandlerFunc(
					func(w http.ResponseWriter, _ *http.Request) {
						http.Error(w, "WEBHOOKS RECEIVED: MOCK ERROR RESPONSE", http.StatusNotFound)
					}))

				cfg := formance.ConfigUser{
					Endpoint: httpServer.URL,
					EventTypes: []string{
						"ledger.committed_transactions",
					},
				}
				var err error
				insertResp, _, err = Client().WebhooksApi.
					InsertConfig(TestContext()).ConfigUser(cfg).Execute()
				Expect(err).To(Not(HaveOccurred()))
				DeferCleanup(func() {
					httpServer.Close()
				})
			})

			Context("testing the inserted one", func() {
				It("should return a failed attempt", func() {
					attemptResp, httpResp, err := Client().WebhooksApi.
						TestConfig(TestContext(), insertResp.Data.Id).Execute()
					Expect(err).To(Not(HaveOccurred()))
					Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
					Expect(attemptResp.Data.Config.Id).To(Equal(insertResp.Data.Id))
					Expect(attemptResp.Data.Payload).To(Equal(`{"data":"test"}`))
					Expect(int(attemptResp.Data.StatusCode)).To(Equal(http.StatusNotFound))
					Expect(attemptResp.Data.Status).To(Equal("failed"))
				})
			})
		})
	})
})
