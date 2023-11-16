package suite

import (
	"fmt"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Ledger}, func() {
	When("trying to list transactions of a non existent ledger", func() {
		var response *operations.ListTransactionsResponse
		BeforeEach(func() {
			var err error
			response, err = Client().Ledger.ListTransactions(TestContext(), operations.ListTransactionsRequest{
				Ledger: "default",
			})
			Expect(err).To(BeNil())

		})
		It("Should fail with a 404", func() {
			Expect(response.StatusCode).To(Equal(http.StatusNotFound))
		})
	})
})

var _ = WithModules([]*Module{modules.Ledger}, func() {
	BeforeEach(func() {
		createLedgerResponse, err := Client().Ledger.CreateLedger(TestContext(), operations.CreateLedgerRequest{
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
			transactions []shared.ExpandedTransaction
		)
		BeforeEach(func() {
			for i := 0; i < int(txCount); i++ {
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

				ret := response.CreateTransactionResponse
				transactions = append([]shared.ExpandedTransaction{
					{
						Timestamp: ret.Data.Timestamp,
						Postings:  ret.Data.Postings,
						Reference: ret.Data.Reference,
						Metadata:  ret.Data.Metadata,
						ID:        ret.Data.ID,
						PreCommitVolumes: map[string]map[string]shared.Volume{
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
						PostCommitVolumes: map[string]map[string]shared.Volume{
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
				rsp *shared.TransactionsCursorResponse
			)
			BeforeEach(func() {
				response, err := Client().Ledger.ListTransactions(
					TestContext(),
					operations.ListTransactionsRequest{
						Ledger:   "default",
						PageSize: ptr(pageSize),
						Expand:   pointer.For("volumes"),
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				rsp = response.TransactionsCursorResponse
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
					_, err := Client().Ledger.CreateTransaction(
						TestContext(),
						operations.CreateTransactionRequest{
							PostTransaction: shared.PostTransaction{
								Metadata: map[string]string{},
								Postings: []shared.Posting{
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

					response, err := Client().Ledger.ListTransactions(
						TestContext(),
						operations.ListTransactionsRequest{
							Cursor: rsp.Cursor.Next,
							Ledger: "default",
							Expand: pointer.For("volumes"),
						},
					)
					Expect(err).ToNot(HaveOccurred())
					Expect(response.StatusCode).To(Equal(200))

					rsp = response.TransactionsCursorResponse
				})
				It("should return next page", func() {
					Expect(rsp.Cursor.PageSize).To(Equal(pageSize))
					Expect(rsp.Cursor.Data).To(Equal(transactions[pageSize : 2*pageSize]))
					Expect(rsp.Cursor.Next).To(BeNil())
				})
				Then("following previous cursor", func() {
					BeforeEach(func() {
						response, err := Client().Ledger.ListTransactions(
							TestContext(),
							operations.ListTransactionsRequest{
								Cursor: rsp.Cursor.Previous,
								Ledger: "default",
								Expand: pointer.For("volumes"),
							},
						)
						Expect(err).ToNot(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))

						rsp = response.TransactionsCursorResponse
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
				rsp *operations.ListTransactionsResponse
			)
			BeforeEach(func() {
				rsp, err = Client().Ledger.ListTransactions(
					TestContext(),
					operations.ListTransactionsRequest{
						RequestBody: map[string]interface{}{
							"$match": map[string]any{
								"invalid-key": 0,
							},
						},
						Ledger:   "default",
						PageSize: ptr(pageSize),
					},
				)
				Expect(err).ToNot(HaveOccurred())
			})
			It("Should fail with "+string(shared.ErrorsEnumValidation)+" error code", func() {
				Expect(rsp.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(rsp.ErrorResponse.ErrorCode).To(Equal(shared.ErrorsEnumValidation))
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
		t1 shared.ExpandedTransaction
		t2 shared.ExpandedTransaction
		t3 shared.ExpandedTransaction
	)
	When("creating transactions", func() {
		BeforeEach(func() {
			response, err := Client().Ledger.CreateTransaction(
				TestContext(),
				operations.CreateTransactionRequest{
					PostTransaction: shared.PostTransaction{
						Metadata: m1,
						Postings: []shared.Posting{
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

			ret := response.CreateTransactionResponse
			t1 = shared.ExpandedTransaction{
				Timestamp: ret.Data.Timestamp,
				Postings:  ret.Data.Postings,
				Reference: ret.Data.Reference,
				Metadata:  ret.Data.Metadata,
				ID:        ret.Data.ID,
				PreCommitVolumes: map[string]map[string]shared.Volume{
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
				PostCommitVolumes: map[string]map[string]shared.Volume{
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

			response, err = Client().Ledger.CreateTransaction(
				TestContext(),
				operations.CreateTransactionRequest{
					PostTransaction: shared.PostTransaction{
						Metadata: m1,
						Postings: []shared.Posting{
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

			ret = response.CreateTransactionResponse
			t2 = shared.ExpandedTransaction{
				Timestamp: ret.Data.Timestamp,
				Postings:  ret.Data.Postings,
				Reference: ret.Data.Reference,
				Metadata:  ret.Data.Metadata,
				ID:        ret.Data.ID,
				PreCommitVolumes: map[string]map[string]shared.Volume{
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
				PostCommitVolumes: map[string]map[string]shared.Volume{
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

			ret = response.CreateTransactionResponse
			t3 = shared.ExpandedTransaction{
				Timestamp: ret.Data.Timestamp,
				Postings:  ret.Data.Postings,
				Reference: ret.Data.Reference,
				Metadata:  ret.Data.Metadata,
				ID:        ret.Data.ID,
				PreCommitVolumes: map[string]map[string]shared.Volume{
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
				PostCommitVolumes: map[string]map[string]shared.Volume{
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
			response, err := Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("3"))

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
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

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
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

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
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

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
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

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
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

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
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

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
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

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
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

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
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

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
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

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
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

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
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

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
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

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
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
			response, err := Client().Ledger.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
					Ledger: "default",
					Expand: pointer.For("volumes"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse := response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[2]).Should(Equal(t1))

			response, err = Client().Ledger.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
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
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[2]).Should(Equal(t1))

			response, err = Client().Ledger.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
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
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(0))

			response, err = Client().Ledger.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
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
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[2]).Should(Equal(t1))

			response, err = Client().Ledger.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
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
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(0))

			response, err = Client().Ledger.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
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
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(0))

			response, err = Client().Ledger.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
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
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[2]).Should(Equal(t1))

			response, err = Client().Ledger.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
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
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(2))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t1))

			response, err = Client().Ledger.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
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
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(0))

			response, err = Client().Ledger.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
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
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(2))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))

			response, err = Client().Ledger.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
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
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(1))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))

			response, err = Client().Ledger.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
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
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(0))

			response, err = Client().Ledger.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
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
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(2))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t1))

			response, err = Client().Ledger.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
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
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(1))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t1))

			response, err = Client().Ledger.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
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
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(0))
		})
		It("should be gettable on api", func() {
			response, err := Client().Ledger.GetTransaction(
				TestContext(),
				operations.GetTransactionRequest{
					Ledger: "default",
					ID:     t1.ID,
					Expand: pointer.For("volumes"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.GetTransactionResponse.Data).Should(Equal(t1))

			response, err = Client().Ledger.GetTransaction(
				TestContext(),
				operations.GetTransactionRequest{
					Ledger: "default",
					ID:     t2.ID,
					Expand: pointer.For("volumes"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.GetTransactionResponse.Data).Should(Equal(t2))

			response, err = Client().Ledger.GetTransaction(
				TestContext(),
				operations.GetTransactionRequest{
					Ledger: "default",
					ID:     t3.ID,
					Expand: pointer.For("volumes"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.GetTransactionResponse.Data).Should(Equal(t3))

			response, err = Client().Ledger.GetTransaction(
				TestContext(),
				operations.GetTransactionRequest{
					Ledger: "default",
					ID:     big.NewInt(666),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(404))
		})
	})

	When("counting and listing transactions empty", func() {
		It("should be countable on api even if empty", func() {
			response, err := Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("0"))
		})
		It("should be listed on api even if empty", func() {
			response, err := Client().Ledger.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.TransactionsCursorResponse.Cursor.Data).To(HaveLen(0))
		})
	})
})
