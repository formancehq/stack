package suite

import (
	"net/http"
	"net/http/httptest"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	webhooks "github.com/formancehq/webhooks/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("An environment configured with a webhook sent on created transaction", func() {
	var (
		httpServer *httptest.Server
		called     chan struct{}
		secret     = webhooks.NewSecret()
	)

	BeforeEach(func() {
		called = make(chan struct{})
		httpServer = httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer close(called)
			}))
		DeferCleanup(func() {
			httpServer.Close()
		})

		response, err := Client().Webhooks.InsertConfig(
			TestContext(),
			shared.ConfigUser{
				Endpoint: httpServer.URL,
				Secret:   &secret,
				EventTypes: []string{
					"ledger.committed_transactions",
				},
			},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusOK))
	})

	When("creating a transaction", func() {
		BeforeEach(func() {
			response, err := Client().Ledger.CreateTransaction(
				TestContext(),
				operations.CreateTransactionRequest{
					PostTransaction: shared.PostTransaction{
						Metadata: map[string]string{},
						Postings: []shared.Posting{
							{
								Amount:      100,
								Asset:       "USD",
								Source:      "world",
								Destination: "alice",
							},
						},
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
		})

		It("should trigger a call to the webhook endpoint", func() {
			Eventually(ChanClosed(called)).Should(BeTrue())
		})
	})
})
