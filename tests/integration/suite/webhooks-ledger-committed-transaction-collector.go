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
// 		count int
// 	)

// 	BeforeEach(func() {
// 		called = make(chan struct{})
// 		count = 0
// 		httpServer = httptest.NewServer(
// 			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 				count += 1
			
// 				if count == 1 {
// 					// FOR THE WORKER
// 					w.WriteHeader(http.StatusUnauthorized)
// 					w.Write([]byte("401 Unauthorized"))
// 				} else if count == 2 {
// 					// FOR THE COLLECTOR
// 					w.WriteHeader(http.StatusOK)
// 					w.Write([]byte("200 OK"))
// 					defer close(called)
// 				}
				
// 			}))
// 		DeferCleanup(func() {
// 			httpServer.Close()
// 		})

// 		response, err := Client().Webhooks.InsertConfig(
// 			TestContext(),
// 			shared.ConfigUser{
// 				Endpoint: httpServer.URL,
// 				Secret:   &secret,
// 				EventTypes: []string{
// 					"ledger.committed_transactions",
// 				},
// 			},
// 		)
// 		Expect(err).ToNot(HaveOccurred())
// 		Expect(response.StatusCode).To(Equal(http.StatusOK))
		
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
