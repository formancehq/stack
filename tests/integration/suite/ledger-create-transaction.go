package suite

import (
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	ledgerevents "github.com/formancehq/ledger/pkg/events"
	"github.com/formancehq/stack/libs/events"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Search, modules.Ledger}, func() {
	When("creating a transaction on a ledger", func() {
		var (
			msgs               chan *nats.Msg
			cancelSubscription func()
			timestamp          = time.Now().Round(time.Second).UTC()
			rsp                *shared.CreateTransactionResponse
		)
		BeforeEach(func() {
			// Subscribe to nats subject
			cancelSubscription, msgs = SubscribeLedger()

			// Create a transaction
			response, err := Client().Ledger.CreateTransaction(
				TestContext(),
				operations.CreateTransactionRequest{
					PostTransaction: shared.PostTransaction{
						Metadata: map[string]string{},
						Postings: []shared.Posting{
							{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Source:      "world",
								Destination: "alice",
							},
						},
						Timestamp: &timestamp,
						Reference: pointer.For("foo"),
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			rsp = response.CreateTransactionResponse
		})
		AfterEach(func() {
			cancelSubscription()
		})
		It("should be available on api", func() {
			response, err := Client().Ledger.GetTransaction(
				TestContext(),
				operations.GetTransactionRequest{
					Ledger: "default",
					ID:     rsp.Data.ID,
					Expand: pointer.For("volumes"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			transactionResponse := response.GetTransactionResponse
			Expect(transactionResponse.Data).To(Equal(shared.ExpandedTransaction{
				Timestamp: rsp.Data.Timestamp,
				Postings:  rsp.Data.Postings,
				Reference: rsp.Data.Reference,
				Metadata:  rsp.Data.Metadata,
				ID:        rsp.Data.ID,
				PreCommitVolumes: map[string]map[string]shared.Volume{
					"world": {
						"USD": {
							Input:   big.NewInt(0),
							Output:  big.NewInt(0),
							Balance: big.NewInt(0),
						},
					},
					"alice": {
						"USD": {
							Input:   big.NewInt(0),
							Output:  big.NewInt(0),
							Balance: big.NewInt(0),
						},
					},
				},
				PostCommitVolumes: map[string]map[string]shared.Volume{
					"world": {
						"USD": {
							Input:   big.NewInt(0),
							Output:  big.NewInt(100),
							Balance: big.NewInt(-100),
						},
					},
					"alice": {
						"USD": {
							Input:   big.NewInt(100),
							Output:  big.NewInt(0),
							Balance: big.NewInt(100),
						},
					},
				},
			}))

			accResponse, err := Client().Ledger.GetAccount(
				TestContext(),
				operations.GetAccountRequest{
					Address: "alice",
					Ledger:  "default",
					Expand:  pointer.For("volumes"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(accResponse.StatusCode).To(Equal(200))

			accountResponse := accResponse.AccountResponse
			Expect(accountResponse.Data).Should(Equal(shared.Account{
				Address:  "alice",
				Metadata: metadata.Metadata{},
				Volumes: map[string]shared.Volume{
					"USD": {
						Input:   big.NewInt(100),
						Output:  big.NewInt(0),
						Balance: big.NewInt(100),
					},
				},
			}))
		})
		Then("trying to commit a new transaction with the same reference", func() {
			var (
				response *operations.CreateTransactionResponse
				err      error
			)
			BeforeEach(func() {
				// Create a transaction
				response, err = Client().Ledger.CreateTransaction(
					TestContext(),
					operations.CreateTransactionRequest{
						PostTransaction: shared.PostTransaction{
							Metadata: map[string]string{},
							Postings: []shared.Posting{
								{
									Amount:      big.NewInt(100),
									Asset:       "USD",
									Source:      "world",
									Destination: "alice",
								},
							},
							Timestamp: &timestamp,
							Reference: pointer.For("foo"),
						},
						Ledger: "default",
					},
				)
				Expect(err).ToNot(HaveOccurred())
			})
			It("Should fail with "+string(shared.ErrorsEnumConflict)+" error code", func() {
				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(response.ErrorResponse.ErrorCode).To(Equal(shared.ErrorsEnumConflict))
			})
		})
		It("should trigger a new event", func() {
			// Wait for created transaction event
			msg := WaitOnChanWithTimeout(msgs, 5*time.Second)
			Expect(events.Check(msg.Data, "ledger", ledgerevents.EventTypeCommittedTransactions)).Should(Succeed())
		})
		It("should pop a transaction, two accounts and two assets entries on search service", func() {
			expectedTx := map[string]any{
				"metadata":  map[string]any{},
				"reference": "foo",
				"postings": []any{
					map[string]any{
						"source":      "world",
						"asset":       "USD",
						"amount":      float64(100),
						"destination": "alice",
					},
				},
				"txid":      float64(0),
				"timestamp": timestamp.Format(time.RFC3339),
				"ledger":    "default",
			}
			Eventually(func(g Gomega) bool {
				response, err := Client().Search.Search(
					TestContext(),
					shared.Query{
						Target: ptr("TRANSACTION"),
					},
				)
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(response.StatusCode).To(Equal(200))

				res := response.Response
				g.Expect(res.Cursor.Data).To(HaveLen(1))
				g.Expect(res.Cursor.Data[0]).To(Equal(expectedTx))

				return true
			}).Should(BeTrue())

			Eventually(func(g Gomega) []map[string]any {
				response, err := Client().Search.Search(
					TestContext(),
					shared.Query{
						Target: ptr("TRANSACTION"),
						Terms:  []string{"alice"},
					},
				)
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(response.StatusCode).To(Equal(200))

				res := response.Response
				g.Expect(res.Cursor.Data[0]).To(Equal(expectedTx))
				return res.Cursor.Data
			}).Should(HaveLen(1))

			Eventually(func(g Gomega) bool {
				response, err := Client().Search.Search(
					TestContext(),
					shared.Query{
						Target: ptr("ACCOUNT"),
					},
				)
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(response.StatusCode).To(Equal(200))

				res := response.Response
				g.Expect(res.Cursor.Data).To(HaveLen(2))
				g.Expect(res.Cursor.Data).To(ContainElements(
					map[string]any{
						"address":  "world",
						"ledger":   "default",
						"metadata": map[string]any{},
					},
					map[string]any{
						"address":  "alice",
						"ledger":   "default",
						"metadata": map[string]any{},
					},
				))
				return true
			}).Should(BeTrue())
		})
	})

	When("creating a transaction on a ledger with insufficient funds", func() {
		It("should fail", func() {
			response, err := Client().Ledger.CreateTransaction(
				TestContext(),
				operations.CreateTransactionRequest{
					PostTransaction: shared.PostTransaction{
						Metadata: map[string]string{},
						Postings: []shared.Posting{
							{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Source:      "bob",
								Destination: "alice",
							},
						},
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(400))
			Expect(response.CreateTransactionResponse).To(BeNil())

			Expect(response.ErrorResponse).Should(Equal(&shared.ErrorResponse{
				ErrorCode:    shared.ErrorsEnumInsufficientFund,
				ErrorMessage: "running numscript: script execution failed: no more fund to withdraw",
			}))
		})
	})

	When("creating a transaction on a ledger with an idempotency key", func() {
		var (
			err      error
			response *operations.CreateTransactionResponse
		)
		createTransaction := func() {
			response, err = Client().Ledger.CreateTransaction(
				TestContext(),
				operations.CreateTransactionRequest{
					IdempotencyKey: ptr("testing"),
					PostTransaction: shared.PostTransaction{
						Metadata: map[string]string{},
						Postings: []shared.Posting{
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
		}
		BeforeEach(createTransaction)
		It("should be ok", func() {
			Expect(err).To(Succeed())
			Expect(response.CreateTransactionResponse.Data.ID).To(Equal(big.NewInt(0)))
		})
		Then("replaying with the same IK", func() {
			BeforeEach(createTransaction)
			It("should respond with the same tx id", func() {
				Expect(err).To(Succeed())
				Expect(response.CreateTransactionResponse.Data.ID).To(Equal(big.NewInt(0)))
			})
		})
	})
	// TODO(gfyrag): test negative amount with a variable
	When("creating a transaction on a ledger with a negative amount in the script", func() {
		var (
			err      error
			response *operations.CreateTransactionResponse
		)
		BeforeEach(func() {
			response, err = Client().Ledger.CreateTransaction(
				TestContext(),
				operations.CreateTransactionRequest{
					IdempotencyKey: ptr("testing"),
					PostTransaction: shared.PostTransaction{
						Metadata: map[string]string{},
						Script: &shared.PostTransactionScript{
							Plain: `send [COIN -100] (
								source = @world
								destination = @bob
							)`,
							Vars: map[string]interface{}{},
						},
					},
					Ledger: "default",
				},
			)
		})
		It("should fail with "+string(shared.ErrorsEnumCompilationFailed)+" code", func() {
			Expect(err).To(Succeed())
			Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
			Expect(response.ErrorResponse).NotTo(BeNil())
			Expect(response.ErrorResponse.ErrorCode).To(Equal(shared.ErrorsEnumCompilationFailed))
			Expect(response.ErrorResponse.Details).To(Equal(pointer.For("https://play.numscript.org/?payload=eyJlcnJvciI6Ilx1MDAxYlszMW0tLVx1MDAzZVx1MDAxYlswbSBlcnJvcjoxOjE1XHJcbiAgXHUwMDFiWzM0bXxcdTAwMWJbMG1cclxuXHUwMDFiWzMxbTEgfCBcdTAwMWJbMG1cdTAwMWJbOTBtc2VuZCBbQ09JTiAtMTAwXHUwMDFiWzBtXVx1MDAxYls5MG0gKFxyXG5cdTAwMWJbMG0gIFx1MDAxYlszNG18XHUwMDFiWzBtICAgICAgICAgICAgICAgIFx1MDAxYlszMW1eXHUwMDFiWzBtIG5vIHZpYWJsZSBhbHRlcm5hdGl2ZSBhdCBpbnB1dCAnW0NPSU4tMTAwXSdcclxuIn0=")))
		})
	})
	When("creating a transaction on a ledger with a negative amount in the script", func() {
		var (
			err      error
			response *operations.CreateTransactionResponse
		)
		BeforeEach(func() {
			response, err = Client().Ledger.CreateTransaction(
				TestContext(),
				operations.CreateTransactionRequest{
					IdempotencyKey: ptr("testing"),
					PostTransaction: shared.PostTransaction{
						Metadata: map[string]string{},
						Script: &shared.PostTransactionScript{
							Plain: `vars {
								monetary $amount
							}
							send $amount (
								source = @world
								destination = @bob
							)`,
							Vars: map[string]interface{}{
								"amount": "USD -100",
							},
						},
					},
					Ledger: "default",
				},
			)
		})
		It("should fail with "+string(shared.ErrorsEnumCompilationFailed)+" code", func() {
			Expect(err).To(Succeed())
			Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
			Expect(response.ErrorResponse).NotTo(BeNil())
			Expect(response.ErrorResponse.ErrorCode).To(Equal(shared.ErrorsEnumCompilationFailed))
			Expect(response.ErrorResponse.Details).To(Equal(pointer.For("https://play.numscript.org/?payload=eyJlcnJvciI6ImludmFsaWQgSlNPTiB2YWx1ZSBmb3IgdmFyaWFibGUgJGFtb3VudCBvZiB0eXBlIG1vbmV0YXJ5OiB2YWx1ZSBbVVNEIC0xMDBdOiBuZWdhdGl2ZSBhbW91bnQifQ==")))
		})
	})
	When("creating a transaction on a ledger with error on script", func() {
		var (
			err      error
			response *operations.CreateTransactionResponse
		)
		BeforeEach(func() {
			response, err = Client().Ledger.CreateTransaction(
				TestContext(),
				operations.CreateTransactionRequest{
					IdempotencyKey: ptr("testing"),
					PostTransaction: shared.PostTransaction{
						Metadata: map[string]string{},
						Script: &shared.PostTransactionScript{
							Plain: `XXX`,
							Vars:  map[string]interface{}{},
						},
					},
					Ledger: "default",
				},
			)
		})
		It("should fail with "+string(shared.ErrorsEnumCompilationFailed)+" code", func() {
			Expect(err).To(Succeed())
			Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
			Expect(response.ErrorResponse).NotTo(BeNil())
			Expect(response.ErrorResponse.ErrorCode).To(Equal(shared.ErrorsEnumCompilationFailed))
			Expect(response.ErrorResponse.Details).To(Equal(pointer.For("https://play.numscript.org/?payload=eyJlcnJvciI6Ilx1MDAxYlszMW0tLVx1MDAzZVx1MDAxYlswbSBlcnJvcjoxOjBcclxuICBcdTAwMWJbMzRtfFx1MDAxYlswbVxyXG5cdTAwMWJbMzFtMSB8IFx1MDAxYlswbVx1MDAxYls5MG1cdTAwMWJbMG1YWFhcdTAwMWJbOTBtXHJcblx1MDAxYlswbSAgXHUwMDFiWzM0bXxcdTAwMWJbMG0gXHUwMDFiWzMxbV5eXHUwMDFiWzBtIG1pc21hdGNoZWQgaW5wdXQgJ1hYWCcgZXhwZWN0aW5nIHtORVdMSU5FLCAndmFycycsICdzZXRfdHhfbWV0YScsICdzZXRfYWNjb3VudF9tZXRhJywgJ3ByaW50JywgJ2ZhaWwnLCAnc2VuZCcsICdzYXZlJ31cclxuIn0=")))
		})
	})
	When("creating a transaction with no postings", func() {
		var (
			err      error
			response *operations.CreateTransactionResponse
		)
		BeforeEach(func() {
			response, err = Client().Ledger.CreateTransaction(
				TestContext(),
				operations.CreateTransactionRequest{
					IdempotencyKey: ptr("testing"),
					PostTransaction: shared.PostTransaction{
						Metadata: map[string]string{},
						Script: &shared.PostTransactionScript{
							Plain: `vars {
								monetary $amount
							}
							set_tx_meta("foo", "bar")
							`,
							Vars: map[string]interface{}{
								"amount": "USD 100",
							},
						},
					},
					Ledger: "default",
				},
			)
		})
		It("should fail with "+string(shared.ErrorsEnumNoPostings)+" code", func() {
			Expect(err).To(Succeed())
			Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
			Expect(response.ErrorResponse).NotTo(BeNil())
			Expect(response.ErrorResponse.ErrorCode).To(Equal(shared.ErrorsEnumNoPostings))
		})
	})
	When("creating a transaction with metadata override", func() {
		var (
			err      error
			response *operations.CreateTransactionResponse
		)
		BeforeEach(func() {
			response, err = Client().Ledger.CreateTransaction(
				TestContext(),
				operations.CreateTransactionRequest{
					IdempotencyKey: ptr("testing"),
					PostTransaction: shared.PostTransaction{
						Metadata: map[string]string{
							"foo": "baz",
						},
						Script: &shared.PostTransactionScript{
							Plain: `send [COIN 100] (
								source = @world
								destination = @bob
							)
							set_tx_meta("foo", "bar")`,
							Vars: map[string]interface{}{},
						},
					},
					Ledger: "default",
				},
			)
		})
		It("should fail with "+string(shared.ErrorsEnumMetadataOverride)+" code", func() {
			Expect(err).To(Succeed())
			Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
			Expect(response.ErrorResponse).NotTo(BeNil())
			Expect(response.ErrorResponse.ErrorCode).To(Equal(shared.ErrorsEnumMetadataOverride))
		})
	})
})

type GenericOpenAPIError interface {
	Model() any
}
