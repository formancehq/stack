package suite

import (
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"math/big"
	"net/http"
	"net/http/httptest"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	webhooks "github.com/formancehq/webhooks/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Ledger, modules.Webhooks}, func() {
	BeforeEach(func() {
		createLedgerResponse, err := Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
			Ledger: "default",
		})
		Expect(err).To(BeNil())
		Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))
	})
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

		response, err := Client().Webhooks.V1.InsertConfig(
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
			response, err := Client().Ledger.V2.CreateTransaction(
				TestContext(),
				operations.V2CreateTransactionRequest{
					V2PostTransaction: shared.V2PostTransaction{
						Metadata: map[string]string{},
						Postings: []shared.V2Posting{
							{
								Amount:      big.NewInt(100),
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
