package suite

import (
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/sdkerrors"
	"github.com/formancehq/stack/tests/integration/internal/modules"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Ledger}, func() {
	When("trying to list transactions of a non existent ledger", func() {
		BeforeEach(func() {
			_, err := Client().Ledger.V2ListTransactions(TestContext(), operations.V2ListTransactionsRequest{
				Ledger: "default",
			})
			Expect(err).NotTo(BeNil())
			Expect(err.(*sdkerrors.V2ErrorResponse).ErrorCode).To(Equal(shared.V2ErrorsEnumLedgerNotFound))
		})
		It("Should fail with a 404", func() {
		})
	})
})

var _ = WithModules([]*Module{modules.Ledger}, func() {
	BeforeEach(func() {
		createLedgerResponse, err := Client().Ledger.V2CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
			Ledger: "default",
		})
		Expect(err).To(BeNil())
		Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))
	})
	const (
		pageSize = int64(10)
		txCount  = 2 * pageSize
	)
	When(fmt.Sprintf("creating %d transactions", txCount), func() {
		var (
			timestamp    = time.Now().Round(time.Second).UTC()
			transactions []shared.V2ExpandedTransaction
		)
		BeforeEach(func() {
			for i := 0; i < int(txCount); i++ {
				response, err := Client().Ledger.V2CreateTransaction(
					TestContext(),
					operations.V2CreateTransactionRequest{
						V2PostTransaction: shared.V2PostTransaction{
							Metadata: map[string]string{},
							Postings: []shared.V2Posting{
								{
									Amount:      big.NewInt(100),
									Asset:       "USD",
									Source:      "world",
									Destination: fmt.Sprintf("account:%d", i),
								},
							},
							Timestamp: &timestamp,
						},
						Ledger: "default",
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				ret := response.V2CreateTransactionResponse
				transactions = append([]shared.V2ExpandedTransaction{
					{
						Timestamp: ret.Data.Timestamp,
						Postings:  ret.Data.Postings,
						Reference: ret.Data.Reference,
						Metadata:  ret.Data.Metadata,
						ID:        ret.Data.ID,
						PreCommitVolumes: map[string]map[string]shared.V2Volume{
							"world": {
								"USD": {
									Input:   big.NewInt(0),
									Output:  big.NewInt(int64(i * 100)),
									Balance: big.NewInt(int64(-i * 100)),
								},
							},
							fmt.Sprintf("account:%d", i): {
								"USD": {
									Input:   big.NewInt(0),
									Output:  big.NewInt(0),
									Balance: big.NewInt(0),
								},
							},
						},
						PostCommitVolumes: map[string]map[string]shared.V2Volume{
							"world": {
								"USD": {
									Input:   big.NewInt(0),
									Output:  big.NewInt(int64((i + 1) * 100)),
									Balance: big.NewInt(int64(-(i + 1) * 100)),
								},
							},
							fmt.Sprintf("account:%d", i): {
								"USD": {
									Input:   big.NewInt(100),
									Output:  big.NewInt(0),
									Balance: big.NewInt(int64(100)),
								},
							},
						},
					},
				}, transactions...)
			}
		})
		AfterEach(func() {
			transactions = nil
		})
		Then(fmt.Sprintf("listing transactions using page size of %d", pageSize), func() {
			var (
				rsp *shared.V2TransactionsCursorResponse
			)
			BeforeEach(func() {
				response, err := Client().Ledger.V2ListTransactions(
					TestContext(),
					operations.V2ListTransactionsRequest{
						Ledger:   "default",
						PageSize: ptr(pageSize),
						Expand:   pointer.For("volumes"),
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				rsp = response.V2TransactionsCursorResponse
				Expect(rsp.Cursor.HasMore).To(BeTrue())
				Expect(rsp.Cursor.Previous).To(BeNil())
				Expect(rsp.Cursor.Next).NotTo(BeNil())
			})
			It("Should be ok", func() {
				Expect(rsp.Cursor.PageSize).To(Equal(pageSize))
				Expect(rsp.Cursor.Data).To(Equal(transactions[:pageSize]))
			})
			Then("following next cursor", func() {
				BeforeEach(func() {

					// Create a new transaction to ensure cursor is stable
					_, err := Client().Ledger.V2CreateTransaction(
						TestContext(),
						operations.V2CreateTransactionRequest{
							V2PostTransaction: shared.V2PostTransaction{
								Metadata: map[string]string{},
								Postings: []shared.V2Posting{
									{
										Amount:      big.NewInt(100),
										Asset:       "USD",
										Source:      "world",
										Destination: "account:0",
									},
								},
								Timestamp: pointer.For(time.Now()),
							},
							Ledger: "default",
						},
					)
					Expect(err).ToNot(HaveOccurred())

					response, err := Client().Ledger.V2ListTransactions(
						TestContext(),
						operations.V2ListTransactionsRequest{
							Cursor: rsp.Cursor.Next,
							Ledger: "default",
							Expand: pointer.For("volumes"),
						},
					)
					Expect(err).ToNot(HaveOccurred())
					Expect(response.StatusCode).To(Equal(200))

					rsp = response.V2TransactionsCursorResponse
				})
				It("should return next page", func() {
					Expect(rsp.Cursor.PageSize).To(Equal(pageSize))
					Expect(rsp.Cursor.Data).To(Equal(transactions[pageSize : 2*pageSize]))
					Expect(rsp.Cursor.Next).To(BeNil())
				})
				Then("following previous cursor", func() {
					BeforeEach(func() {
						response, err := Client().Ledger.V2ListTransactions(
							TestContext(),
							operations.V2ListTransactionsRequest{
								Cursor: rsp.Cursor.Previous,
								Ledger: "default",
								Expand: pointer.For("volumes"),
							},
						)
						Expect(err).ToNot(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))

						rsp = response.V2TransactionsCursorResponse
					})
					It("should return first page", func() {
						Expect(rsp.Cursor.PageSize).To(Equal(pageSize))
						Expect(rsp.Cursor.Data).To(Equal(transactions[:pageSize]))
						Expect(rsp.Cursor.Previous).To(BeNil())
					})
				})
			})
		})

		Then("listing transactions using invalid filter", func() {
			var (
				err error
			)
			BeforeEach(func() {
				_, err = Client().Ledger.V2ListTransactions(
					TestContext(),
					operations.V2ListTransactionsRequest{
						RequestBody: map[string]interface{}{
							"$match": map[string]any{
								"invalid-key": 0,
							},
						},
						Ledger:   "default",
						PageSize: ptr(pageSize),
					},
				)
				Expect(err).To(HaveOccurred())
			})
			It("Should fail with "+string(shared.V2ErrorsEnumValidation)+" error code", func() {
				Expect(err.(*sdkerrors.V2ErrorResponse).ErrorCode).To(Equal(shared.V2ErrorsEnumValidation))
			})
		})
	})
	var (
		timestamp1 = time.Date(2023, 4, 10, 10, 0, 0, 0, time.UTC)
		timestamp2 = time.Date(2023, 4, 11, 10, 0, 0, 0, time.UTC)
		timestamp3 = time.Date(2023, 4, 12, 10, 0, 0, 0, time.UTC)

		m1 = metadata.Metadata{
			"foo": "bar",
		}
	)

	var (
		t1 shared.V2ExpandedTransaction
		t2 shared.V2ExpandedTransaction
		t3 shared.V2ExpandedTransaction
	)
	When("creating transactions", func() {
		BeforeEach(func() {
			response, err := Client().Ledger.V2CreateTransaction(
				TestContext(),
				operations.V2CreateTransactionRequest{
					V2PostTransaction: shared.V2PostTransaction{
						Metadata: m1,
						Postings: []shared.V2Posting{
							{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Source:      "world",
								Destination: "foo:foo",
							},
						},
						Timestamp: &timestamp1,
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			ret := response.V2CreateTransactionResponse
			t1 = shared.V2ExpandedTransaction{
				Timestamp: ret.Data.Timestamp,
				Postings:  ret.Data.Postings,
				Reference: ret.Data.Reference,
				Metadata:  ret.Data.Metadata,
				ID:        ret.Data.ID,
				PreCommitVolumes: map[string]map[string]shared.V2Volume{
					"world": {
						"USD": {
							Input:   big.NewInt(0),
							Output:  big.NewInt(0),
							Balance: big.NewInt(int64(0)),
						},
					},
					"foo:foo": {
						"USD": {
							Input:   big.NewInt(0),
							Output:  big.NewInt(0),
							Balance: big.NewInt(int64(0)),
						},
					},
				},
				PostCommitVolumes: map[string]map[string]shared.V2Volume{
					"world": {
						"USD": {
							Input:   big.NewInt(0),
							Output:  big.NewInt(100),
							Balance: big.NewInt(int64(-100)),
						},
					},
					"foo:foo": {
						"USD": {
							Input:   big.NewInt(100),
							Output:  big.NewInt(0),
							Balance: big.NewInt(int64(100)),
						},
					},
				},
			}

			response, err = Client().Ledger.V2CreateTransaction(
				TestContext(),
				operations.V2CreateTransactionRequest{
					V2PostTransaction: shared.V2PostTransaction{
						Metadata: m1,
						Postings: []shared.V2Posting{
							{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Source:      "world",
								Destination: "foo:bar",
							},
						},
						Timestamp: &timestamp2,
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			ret = response.V2CreateTransactionResponse
			t2 = shared.V2ExpandedTransaction{
				Timestamp: ret.Data.Timestamp,
				Postings:  ret.Data.Postings,
				Reference: ret.Data.Reference,
				Metadata:  ret.Data.Metadata,
				ID:        ret.Data.ID,
				PreCommitVolumes: map[string]map[string]shared.V2Volume{
					"world": {
						"USD": {
							Input:   big.NewInt(0),
							Output:  big.NewInt(100),
							Balance: big.NewInt(int64(-100)),
						},
					},
					"foo:bar": {
						"USD": {
							Input:   big.NewInt(0),
							Output:  big.NewInt(0),
							Balance: big.NewInt(int64(0)),
						},
					},
				},
				PostCommitVolumes: map[string]map[string]shared.V2Volume{
					"world": {
						"USD": {
							Input:   big.NewInt(0),
							Output:  big.NewInt(200),
							Balance: big.NewInt(int64(-200)),
						},
					},
					"foo:bar": {
						"USD": {
							Input:   big.NewInt(100),
							Output:  big.NewInt(0),
							Balance: big.NewInt(int64(100)),
						},
					},
				},
			}

			response, err = Client().Ledger.V2CreateTransaction(
				TestContext(),
				operations.V2CreateTransactionRequest{
					V2PostTransaction: shared.V2PostTransaction{
						Metadata: map[string]string{},
						Postings: []shared.V2Posting{
							{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Source:      "world",
								Destination: "foo:baz",
							},
						},
						Timestamp: &timestamp3,
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			ret = response.V2CreateTransactionResponse
			t3 = shared.V2ExpandedTransaction{
				Timestamp: ret.Data.Timestamp,
				Postings:  ret.Data.Postings,
				Reference: ret.Data.Reference,
				Metadata:  ret.Data.Metadata,
				ID:        ret.Data.ID,
				PreCommitVolumes: map[string]map[string]shared.V2Volume{
					"world": {
						"USD": {
							Input:   big.NewInt(0),
							Output:  big.NewInt(200),
							Balance: big.NewInt(int64(-200)),
						},
					},
					"foo:baz": {
						"USD": {
							Input:   big.NewInt(0),
							Output:  big.NewInt(0),
							Balance: big.NewInt(int64(0)),
						},
					},
				},
				PostCommitVolumes: map[string]map[string]shared.V2Volume{
					"world": {
						"USD": {
							Input:   big.NewInt(0),
							Output:  big.NewInt(300),
							Balance: big.NewInt(int64(-300)),
						},
					},
					"foo:baz": {
						"USD": {
							Input:   big.NewInt(100),
							Output:  big.NewInt(0),
							Balance: big.NewInt(int64(100)),
						},
					},
				},
			}
		})
		It("should be countable on api", func() {
			response, err := Client().Ledger.V2CountTransactions(
				TestContext(),
				operations.V2CountTransactionsRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("3"))

			response, err = Client().Ledger.V2CountTransactions(
				TestContext(),
				operations.V2CountTransactionsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"account": "foo:",
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("3"))

			response, err = Client().Ledger.V2CountTransactions(
				TestContext(),
				operations.V2CountTransactionsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"account": "not_existing",
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("0"))

			response, err = Client().Ledger.V2CountTransactions(
				TestContext(),
				operations.V2CountTransactionsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"destination": ":baz",
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("1"))

			response, err = Client().Ledger.V2CountTransactions(
				TestContext(),
				operations.V2CountTransactionsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"destination": "not_existing",
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("0"))

			response, err = Client().Ledger.V2CountTransactions(
				TestContext(),
				operations.V2CountTransactionsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"source": "foo:",
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("0"))

			response, err = Client().Ledger.V2CountTransactions(
				TestContext(),
				operations.V2CountTransactionsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"source": "world",
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("3"))

			response, err = Client().Ledger.V2CountTransactions(
				TestContext(),
				operations.V2CountTransactionsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"metadata[foo]": "bar",
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("2"))

			response, err = Client().Ledger.V2CountTransactions(
				TestContext(),
				operations.V2CountTransactionsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"metadata[foo]": "not_existing",
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("0"))

			response, err = Client().Ledger.V2CountTransactions(
				TestContext(),
				operations.V2CountTransactionsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$gte": map[string]any{
							"timestamp": timestamp2.Format(time.RFC3339),
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("2"))

			response, err = Client().Ledger.V2CountTransactions(
				TestContext(),
				operations.V2CountTransactionsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$gte": map[string]any{
							"timestamp": timestamp3.Format(time.RFC3339),
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("1"))

			response, err = Client().Ledger.V2CountTransactions(
				TestContext(),
				operations.V2CountTransactionsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$gte": map[string]any{
							"timestamp": time.Now().UTC().Format(time.RFC3339),
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("0"))

			response, err = Client().Ledger.V2CountTransactions(
				TestContext(),
				operations.V2CountTransactionsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$lt": map[string]any{
							"timestamp": timestamp3.Format(time.RFC3339),
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("2"))

			response, err = Client().Ledger.V2CountTransactions(
				TestContext(),
				operations.V2CountTransactionsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$lt": map[string]any{
							"timestamp": timestamp2.Format(time.RFC3339),
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("1"))

			response, err = Client().Ledger.V2CountTransactions(
				TestContext(),
				operations.V2CountTransactionsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$lt": map[string]any{
							"timestamp": time.Date(2023, 4, 9, 10, 0, 0, 0, time.UTC).Format(time.RFC3339),
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("0"))
		})
		It("should be listed on api", func() {
			response, err := Client().Ledger.V2ListTransactions(
				TestContext(),
				operations.V2ListTransactionsRequest{
					Ledger: "default",
					Expand: pointer.For("volumes"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse := response.V2TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[2]).Should(Equal(t1))

			response, err = Client().Ledger.V2ListTransactions(
				TestContext(),
				operations.V2ListTransactionsRequest{
					Ledger: "default",
					Expand: pointer.For("volumes"),
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"account": "foo:",
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.V2TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[2]).Should(Equal(t1))

			response, err = Client().Ledger.V2ListTransactions(
				TestContext(),
				operations.V2ListTransactionsRequest{
					Ledger: "default",
					Expand: pointer.For("volumes"),
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"account": "not_existing",
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.V2TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(0))

			response, err = Client().Ledger.V2ListTransactions(
				TestContext(),
				operations.V2ListTransactionsRequest{
					Ledger: "default",
					Expand: pointer.For("volumes"),
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"destination": "foo:",
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.V2TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[2]).Should(Equal(t1))

			response, err = Client().Ledger.V2ListTransactions(
				TestContext(),
				operations.V2ListTransactionsRequest{
					Ledger: "default",
					Expand: pointer.For("volumes"),
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"destination": "not_existing",
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.V2TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(0))

			response, err = Client().Ledger.V2ListTransactions(
				TestContext(),
				operations.V2ListTransactionsRequest{
					Ledger: "default",
					Expand: pointer.For("volumes"),
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"source": "foo:",
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.V2TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(0))

			response, err = Client().Ledger.V2ListTransactions(
				TestContext(),
				operations.V2ListTransactionsRequest{
					Ledger: "default",
					Expand: pointer.For("volumes"),
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"source": "world",
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.V2TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[2]).Should(Equal(t1))

			response, err = Client().Ledger.V2ListTransactions(
				TestContext(),
				operations.V2ListTransactionsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"metadata[foo]": "bar",
						},
					},
					Expand: pointer.For("volumes"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.V2TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(2))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t1))

			response, err = Client().Ledger.V2ListTransactions(
				TestContext(),
				operations.V2ListTransactionsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"metadata[foo]": "not_existing",
						},
					},
					Expand: pointer.For("volumes"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.V2TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(0))

			response, err = Client().Ledger.V2ListTransactions(
				TestContext(),
				operations.V2ListTransactionsRequest{
					Ledger: "default",
					Expand: pointer.For("volumes"),
					RequestBody: map[string]interface{}{
						"$gte": map[string]any{
							"timestamp": timestamp2.Format(time.RFC3339),
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.V2TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(2))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))

			response, err = Client().Ledger.V2ListTransactions(
				TestContext(),
				operations.V2ListTransactionsRequest{
					Ledger: "default",
					Expand: pointer.For("volumes"),
					RequestBody: map[string]interface{}{
						"$gte": map[string]any{
							"timestamp": timestamp3.Format(time.RFC3339),
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.V2TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(1))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))

			response, err = Client().Ledger.V2ListTransactions(
				TestContext(),
				operations.V2ListTransactionsRequest{
					Ledger: "default",
					Expand: pointer.For("volumes"),
					RequestBody: map[string]interface{}{
						"$gte": map[string]any{
							"timestamp": time.Now().UTC().Format(time.RFC3339),
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.V2TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(0))

			response, err = Client().Ledger.V2ListTransactions(
				TestContext(),
				operations.V2ListTransactionsRequest{
					Ledger: "default",
					Expand: pointer.For("volumes"),
					RequestBody: map[string]interface{}{
						"$lt": map[string]any{
							"timestamp": timestamp3.Format(time.RFC3339),
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.V2TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(2))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t1))

			response, err = Client().Ledger.V2ListTransactions(
				TestContext(),
				operations.V2ListTransactionsRequest{
					Ledger: "default",
					Expand: pointer.For("volumes"),
					RequestBody: map[string]interface{}{
						"$lt": map[string]any{
							"timestamp": timestamp2.Format(time.RFC3339),
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.V2TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(1))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t1))

			response, err = Client().Ledger.V2ListTransactions(
				TestContext(),
				operations.V2ListTransactionsRequest{
					Ledger: "default",
					Expand: pointer.For("volumes"),
					RequestBody: map[string]interface{}{
						"$lt": map[string]any{
							"timestamp": time.Date(2023, 4, 9, 10, 0, 0, 0, time.UTC).Format(time.RFC3339),
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.V2TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(0))
		})
		It("should be gettable on api", func() {
			response, err := Client().Ledger.V2GetTransaction(
				TestContext(),
				operations.V2GetTransactionRequest{
					Ledger: "default",
					ID:     t1.ID,
					Expand: pointer.For("volumes"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.V2GetTransactionResponse.Data).Should(Equal(t1))

			response, err = Client().Ledger.V2GetTransaction(
				TestContext(),
				operations.V2GetTransactionRequest{
					Ledger: "default",
					ID:     t2.ID,
					Expand: pointer.For("volumes"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.V2GetTransactionResponse.Data).Should(Equal(t2))

			response, err = Client().Ledger.V2GetTransaction(
				TestContext(),
				operations.V2GetTransactionRequest{
					Ledger: "default",
					ID:     t3.ID,
					Expand: pointer.For("volumes"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.V2GetTransactionResponse.Data).Should(Equal(t3))

			response, err = Client().Ledger.V2GetTransaction(
				TestContext(),
				operations.V2GetTransactionRequest{
					Ledger: "default",
					ID:     big.NewInt(666),
				},
			)
			Expect(err).To(HaveOccurred())
			Expect(err.(*sdkerrors.V2ErrorResponse).ErrorCode).To(Equal(shared.V2ErrorsEnumNotFound))
		})
	})

	When("counting and listing transactions empty", func() {
		It("should be countable on api even if empty", func() {
			response, err := Client().Ledger.V2CountTransactions(
				TestContext(),
				operations.V2CountTransactionsRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("0"))
		})
		It("should be listed on api even if empty", func() {
			response, err := Client().Ledger.V2ListTransactions(
				TestContext(),
				operations.V2ListTransactionsRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.V2TransactionsCursorResponse.Cursor.Data).To(HaveLen(0))
		})
	})
})
