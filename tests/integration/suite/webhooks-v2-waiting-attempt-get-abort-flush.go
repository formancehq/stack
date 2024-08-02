package suite


// Flag : WebhookAsyncCache
		// This test is commented because for Webhook V2, 
		// Worker and Runner have asynchronized cache. 
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

				
// 		Describe("Try to manage the waiting attempts ", Ordered, func(){

// 			var (
// 				httpBadServer *httptest.Server
// 				httpGoodServer * httptest.Server 
// 				called     chan struct{}
// 				secret     = webhooks.NewSecret()
// 				hook1 shared.V2Hook
// 				waitinAttempt shared.V2Attempt
				
// 			)

// 			BeforeAll(func(){
// 				// CREATE LEDGER 
// 				createLedgerResponse, err := Client().Ledger.V2CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
// 					Ledger: "default",
// 				})
// 				Expect(err).To(BeNil())
// 				Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))

// 				// CREATE FAKE WEB SERVER FOR ENDPOINT
// 				httpBadServer = httptest.NewServer(
// 					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 							// ALWAYS SEND BAD RESPONSE
// 							w.WriteHeader(http.StatusUnauthorized)
// 							w.Write([]byte("401 Unauthorized"))
						
// 					}))
// 				DeferCleanup(func() {
// 					httpBadServer.Close()
// 				})

// 				called = make(chan struct{})
// 				httpGoodServer = httptest.NewServer(
// 					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 							// ALWAYS SEND BAD RESPONSE
// 							w.WriteHeader(http.StatusOK)
// 							w.Write([]byte("Success"))
// 							defer close(called)
// 					}))
// 				DeferCleanup(func() {
// 					httpGoodServer.Close()
// 				})
		
// 				// CREATE HOOK	
// 				response, err := Client().Webhooks.InsertHook(
// 					TestContext(),
// 					shared.V2HookBodyParams{
// 						Endpoint: httpBadServer.URL,
// 						Secret:   &secret,
// 						Events: []string{
// 							"ledger.committed_transactions",
// 						},
// 					},
// 				)
// 				Expect(err).ToNot(HaveOccurred())
// 				Expect(response.StatusCode).To(Equal(http.StatusOK))
// 				hook1 = response.V2HookResponse.Data

// 				// ACTIVATE HOOK
// 				_, err = Client().Webhooks.ActivateHook(
// 					TestContext(),
// 					operations.ActivateHookRequest{
// 						HookID: hook1.ID,
// 					},
// 				)
// 				Expect(err).ToNot(HaveOccurred())

// 				// NEED THIS TO LET CACHES REFRESH INSIDE WEBHOOKS_WORKER & WEBHOOKS_COLLECTOR
// 				time.Sleep(1*time.Second)
				
				
// 				// CREATE A TRANSACTION INSIDE THE LEDGER
// 				resp2, err := Client().Ledger.V2CreateTransaction(
// 					TestContext(),
// 					operations.V2CreateTransactionRequest{
// 						V2PostTransaction: shared.V2PostTransaction{
// 							Metadata: map[string]string{},
// 							Postings: []shared.V2Posting{
// 								{
// 									Amount:      big.NewInt(100),
// 									Asset:       "USD",
// 									Source:      "world",
// 									Destination: "alice",
// 								},
// 							},
// 						},
// 						Ledger: "default",
// 					},
// 				)
// 				Expect(err).ToNot(HaveOccurred())
// 				Expect(resp2.StatusCode).To(Equal(http.StatusOK))
				
// 				// RIGHT KNOW, A waiting attempt should be handle by the Collector because
// 				// Worker didn't successfully reach the endpoint (HttpBadServer)
// 			})
			
// 			It("should have a waiting attempt", func() {
				
// 				time.Sleep(1*time.Second)
				
// 				response, err := Client().Webhooks.GetWaitingAttempts(
// 					TestContext(),
// 					operations.GetWaitingAttemptsRequest{},
// 				)
// 				Expect(err).ToNot(HaveOccurred())
// 				Expect(response.V2AttemptCursorResponse.Cursor.HasMore).To(BeFalse())
// 				Expect(response.V2AttemptCursorResponse.Cursor.Data).To(HaveLen(1))
// 				waitinAttempt = response.V2AttemptCursorResponse.Cursor.Data[0]
// 				time.Sleep(1*time.Second)

// 				// Abort the waiting attempt
// 				resp, err := Client().Webhooks.AbortWaitingAttempt(
// 					TestContext(),		
// 					operations.AbortWaitingAttemptRequest{
// 						AttemptID: waitinAttempt.ID,
// 					},
// 				)
// 				Expect(err).ToNot(HaveOccurred())
// 				Expect(resp.V2AttemptResponse.Data.Status).To(Equal(shared.V2AttemptStatusAbort))
// 				// Check if no Waiting Attempts anymore.
// 				resp2, err := Client().Webhooks.GetWaitingAttempts(
// 					TestContext(),
// 					operations.GetWaitingAttemptsRequest{},
// 				)
// 				Expect(err).ToNot(HaveOccurred())
// 				Expect(resp2.V2AttemptCursorResponse.Cursor.HasMore).To(BeFalse())
// 				Expect(resp2.V2AttemptCursorResponse.Cursor.Data).To(HaveLen(0))


// 				// But We should have one in AbortedAttempts

// 				resp4, err := Client().Webhooks.GetAbortedAttempts(
// 					TestContext(),
// 					operations.GetAbortedAttemptsRequest{},
// 				)

// 				Expect(err).ToNot(HaveOccurred())
// 				Expect(resp4.V2AttemptCursorResponse.Cursor.HasMore).To(BeFalse())
// 				Expect(resp4.V2AttemptCursorResponse.Cursor.Data).To(HaveLen(1))

				

// 				// Change the endpoint of the Hook 

// 				_ , err = Client().Webhooks.UpdateEndpointHook(
// 					TestContext(),
// 					operations.UpdateEndpointHookRequest{
// 						HookID: hook1.ID,
// 						RequestBody: operations.UpdateEndpointHookRequestBody{
// 							Endpoint : &httpGoodServer.URL,
// 						},
// 					},
// 				)

// 				Expect(err).ToNot(HaveOccurred())

// 				//Wait for Cache refresh...
// 				time.Sleep(2*time.Second)

// 				// Chan should be still open because no Waiting Attempt in Cache
// 				Eventually(ChanClosed(called)).Should(BeFalse())
// 			})
	
			
// 		})
		

// 	})

