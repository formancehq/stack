package suite

// Flag : WebhookAsyncCache
		// This test is commented because for Webhook V2, 
		// Worker and Runner have asynchrone cache. 
		// It needs a bit of time between the creation and activation
		// Of an Hook by the user and the moment where it's active in cache.

// import (
// 	"math/big"
// 	"net/http"
// 	"net/http/httptest"
// 	"time"

// 	"github.com/formancehq/stack/tests/integration/internal/modules"

// 	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
// 	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
// 	. "github.com/formancehq/stack/tests/integration/internal"
// 	webhooks "github.com/formancehq/webhooks/pkg/utils"
// 	. "github.com/onsi/ginkgo/v2"
// 	. "github.com/onsi/gomega"
// )

// var _ = WithModules([]*Module{modules.Ledger, modules.Webhooks}, func() {
// 	BeforeEach(func() {
// 		createLedgerResponse, err := Client().Ledger.V2CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
// 			Ledger: "default",
// 		})
// 		Expect(err).To(BeNil())
// 		Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))
// 	})
// 	var (
// 		httpServer *httptest.Server
// 		called     chan struct{}
// 		secret     = webhooks.NewSecret()
// 		hook1 shared.V2Hook
// 	)

// 	BeforeEach(func() {
// 		called = make(chan struct{})
// 		httpServer = httptest.NewServer(
// 			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				
// 				defer close(called)

				

// 			}))
// 		DeferCleanup(func() {
// 			httpServer.Close()
// 		})

		

// 		response, err := Client().Webhooks.InsertHook(
// 			TestContext(),
// 			shared.V2HookBodyParams{
// 				Endpoint: httpServer.URL,
// 				Secret:   &secret,
// 				Events: []string{
// 					"ledger.committed_transactions",
// 				},
// 			},
// 		)
// 		Expect(err).ToNot(HaveOccurred())
// 		Expect(response.StatusCode).To(Equal(http.StatusOK))
// 		hook1 = response.V2HookResponse.Data
// 		_, err = Client().Webhooks.ActivateHook(
// 			TestContext(),
// 			operations.ActivateHookRequest{
// 				HookID: hook1.ID,
// 			},
// 		)
// 		Expect(err).ToNot(HaveOccurred())

// 	})

// 	When("creating a transaction", func() {
// 		BeforeEach(func() {
// 			time.Sleep(1*time.Second)
// 			response, err := Client().Ledger.V2CreateTransaction(
// 				TestContext(),
// 				operations.V2CreateTransactionRequest{
// 					V2PostTransaction: shared.V2PostTransaction{
// 						Metadata: map[string]string{},
// 						Postings: []shared.V2Posting{
// 							{
// 								Amount:      big.NewInt(100),
// 								Asset:       "USD",
// 								Source:      "world",
// 								Destination: "alice",
// 							},
// 						},
// 					},
// 					Ledger: "default",
// 				},
// 			)
// 			Expect(err).ToNot(HaveOccurred())
// 			Expect(response.StatusCode).To(Equal(http.StatusOK))
// 		})

// 		It("should trigger a call to the webhook endpoint", func() {
// 			Eventually(ChanClosed(called)).Should(BeTrue())
// 		})
// 	})
// })
