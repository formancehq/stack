package suite

import (
	"net/http"
	"time"

	"github.com/formancehq/formance-sdk-go"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("empty environment for webhooks configs", func() {
	It("inserting a valid config", func() {
		cfg := formance.ConfigUser{
			Endpoint: "https://example.com",
			EventTypes: []string{
				"ledger.committed_transactions",
			},
		}
		insertResp, httpResp, err := Client().WebhooksApi.
			InsertConfig(TestContext()).ConfigUser(cfg).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(httpResp.StatusCode).To(Equal(http.StatusOK))
		Expect(insertResp.Data.Endpoint).To(Equal(cfg.Endpoint))
		Expect(insertResp.Data.EventTypes).To(Equal(cfg.EventTypes))
		Expect(insertResp.Data.Active).To(BeTrue())
		Expect(insertResp.Data.CreatedAt).NotTo(Equal(time.Time{}))
		Expect(insertResp.Data.UpdatedAt).NotTo(Equal(time.Time{}))
		_, err = uuid.Parse(insertResp.Data.Id)
		Expect(err).NotTo(HaveOccurred())
	})

	It("inserting an invalid config without event types", func() {
		cfg := formance.ConfigUser{
			Endpoint:   "https://example.com",
			EventTypes: []string{},
		}
		insertResp, httpResp, err := Client().WebhooksApi.
			InsertConfig(TestContext()).ConfigUser(cfg).Execute()
		Expect(err).To(HaveOccurred())
		Expect(insertResp).To(BeNil())
		Expect(httpResp.StatusCode).To(Equal(http.StatusBadRequest))
	})

	It("inserting an invalid config without endpoint", func() {
		cfg := formance.ConfigUser{
			Endpoint: "",
			EventTypes: []string{
				"ledger.committed_transactions",
			},
		}
		insertResp, httpResp, err := Client().WebhooksApi.
			InsertConfig(TestContext()).ConfigUser(cfg).Execute()
		Expect(err).To(HaveOccurred())
		Expect(insertResp).To(BeNil())
		Expect(httpResp.StatusCode).To(Equal(http.StatusBadRequest))
	})

	It("inserting an invalid config with invalid secret", func() {
		secret := "invalid"
		cfg := formance.ConfigUser{
			Endpoint: "https://example.com",
			Secret:   &secret,
			EventTypes: []string{
				"ledger.committed_transactions",
			},
		}
		insertResp, httpResp, err := Client().WebhooksApi.
			InsertConfig(TestContext()).ConfigUser(cfg).Execute()
		Expect(err).To(HaveOccurred())
		Expect(insertResp).To(BeNil())
		Expect(httpResp.StatusCode).To(Equal(http.StatusBadRequest))
	})
})
