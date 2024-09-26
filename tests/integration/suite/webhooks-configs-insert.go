package suite

import (
	"net/http"
	"time"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/sdkerrors"
	"github.com/formancehq/stack/tests/integration/internal/modules"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Webhooks}, func() {
	It("inserting a valid config", func() {
		cfg := shared.ConfigUser{
			Endpoint: "https://example.com",
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

		insertResp := response.ConfigResponse
		Expect(insertResp.Data.Endpoint).To(Equal(cfg.Endpoint))
		Expect(insertResp.Data.EventTypes).To(Equal(cfg.EventTypes))
		Expect(insertResp.Data.Active).To(BeTrue())
		Expect(insertResp.Data.CreatedAt).NotTo(Equal(time.Time{}))
		Expect(insertResp.Data.UpdatedAt).NotTo(Equal(time.Time{}))
		_, err = uuid.Parse(insertResp.Data.ID)
		Expect(err).NotTo(HaveOccurred())
	})

	It("inserting an invalid config without event types", func() {
		cfg := shared.ConfigUser{
			Endpoint:   "https://example.com",
			EventTypes: []string{},
		}
		_, err := Client().Webhooks.V1.InsertConfig(
			TestContext(),
			cfg,
		)
		Expect(err).To(HaveOccurred())
	})

	It("inserting an invalid config without endpoint", func() {
		cfg := shared.ConfigUser{
			Endpoint: "",
			EventTypes: []string{
				"ledger.committed_transactions",
			},
		}
		_, err := Client().Webhooks.V1.InsertConfig(
			TestContext(),
			cfg,
		)
		Expect(err).To(HaveOccurred())
		Expect(err.(*sdkerrors.WebhooksErrorResponse).ErrorCode).To(Equal(shared.WebhooksErrorsEnumValidation))
	})

	It("inserting an invalid config with invalid secret", func() {
		secret := "invalid"
		cfg := shared.ConfigUser{
			Endpoint: "https://example.com",
			Secret:   &secret,
			EventTypes: []string{
				"ledger.committed_transactions",
			},
		}
		_, err := Client().Webhooks.V1.InsertConfig(
			TestContext(),
			cfg,
		)
		Expect(err).To(HaveOccurred())
	})
})
