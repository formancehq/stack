package stripe_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/formancehq/payments/internal/connectors/plugins"
	"github.com/formancehq/payments/internal/connectors/plugins/public/stripe"
	"github.com/formancehq/payments/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPlugin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stripe Plugin Suite")
}

var _ = Describe("Stripe Plugin", func() {
	var (
		plg *stripe.Plugin
	)

	BeforeEach(func() {
		plg = &stripe.Plugin{}
	})

	Context("install", func() {
		It("reports validation errors in the config", func(ctx SpecContext) {
			req := models.InstallRequest{Config: json.RawMessage(`{}`)}
			_, err := plg.Install(context.Background(), req)
			Expect(err).To(MatchError(ContainSubstring("config")))
		})
		It("returns valid install response", func(ctx SpecContext) {
			req := models.InstallRequest{Config: json.RawMessage(`{"apiKey":"dummy"}`)}
			res, err := plg.Install(context.Background(), req)
			Expect(err).To(BeNil())
			Expect(len(res.Capabilities) > 0).To(BeTrue())
			Expect(len(res.Workflow) > 0).To(BeTrue())
			Expect(res.Workflow[0].Name).To(Equal("fetch_accounts"))
		})
	})

	Context("uninstall", func() {
		It("returns valid uninstall response", func(ctx SpecContext) {
			req := models.UninstallRequest{ConnectorID: "dummyID"}
			_, err := plg.Uninstall(context.Background(), req)
			Expect(err).To(BeNil())
		})
	})

	Context("calling functions on uninstalled plugins", func() {
		It("fails when fetch next accounts is called before install", func(ctx SpecContext) {
			req := models.FetchNextAccountsRequest{
				State: json.RawMessage(`{}`),
			}
			_, err := plg.FetchNextAccounts(context.Background(), req)
			Expect(err).To(MatchError(plugins.ErrNotYetInstalled))
		})
	})
})
