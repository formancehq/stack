package suite

import (
	"fmt"
	"math/big"
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("some empty environment", func() {
	const (
		pageSize = int64(10)
		txCount  = 2 * pageSize
	)
	When(fmt.Sprintf("creating %d transactions", txCount), func() {
		var (
			timestamp    = time.Now().Round(time.Second).UTC()
			transactions []shared.Transaction
		)
		BeforeEach(func() {
			for i := 0; i < int(txCount); i++ {
				response, err := Client().Transactions.CreateTransaction(
					TestContext(),
					operations.CreateTransactionRequest{
						IdempotencyKey: new(string),
						PostTransaction: shared.PostTransaction{
							Metadata: map[string]any{},
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

				ret := response.TransactionsResponse
				transactions = append([]shared.Transaction{
					{
						Timestamp: ret.Data[0].Timestamp,
						Postings:  ret.Data[0].Postings,
						Reference: ret.Data[0].Reference,
						Metadata:  ret.Data[0].Metadata,
						Txid:      ret.Data[0].Txid,
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
				response, err := Client().Transactions.ListTransactions(
					TestContext(),
					operations.ListTransactionsRequest{
						Ledger:   "default",
						PageSize: ptr(pageSize),
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
					response, err := Client().Transactions.ListTransactions(
						TestContext(),
						operations.ListTransactionsRequest{
							Cursor: rsp.Cursor.Next,
							Ledger: "default",
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
						response, err := Client().Transactions.ListTransactions(
							TestContext(),
							operations.ListTransactionsRequest{
								Cursor: rsp.Cursor.Previous,
								Ledger: "default",
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
	})
})

var _ = Given("some empty environment", func() {
	var (
		timestamp1 = time.Date(2023, 4, 10, 10, 0, 0, 0, time.UTC)
		timestamp2 = time.Date(2023, 4, 11, 10, 0, 0, 0, time.UTC)
		timestamp3 = time.Date(2023, 4, 12, 10, 0, 0, 0, time.UTC)

		m1 = map[string]any{
			"foo": "bar",
		}
	)

	var (
		t1 shared.Transaction
		t2 shared.Transaction
		t3 shared.Transaction
	)
	When("creating transactions", func() {
		BeforeEach(func() {
			response, err := Client().Transactions.CreateTransaction(
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

			ret := response.TransactionsResponse
			t1 = shared.Transaction{
				Timestamp: ret.Data[0].Timestamp,
				Postings:  ret.Data[0].Postings,
				Reference: ret.Data[0].Reference,
				Metadata:  ret.Data[0].Metadata,
				Txid:      ret.Data[0].Txid,
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

			response, err = Client().Transactions.CreateTransaction(
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

			ret = response.TransactionsResponse
			t2 = shared.Transaction{
				Timestamp: ret.Data[0].Timestamp,
				Postings:  ret.Data[0].Postings,
				Reference: ret.Data[0].Reference,
				Metadata:  ret.Data[0].Metadata,
				Txid:      ret.Data[0].Txid,
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

			response, err = Client().Transactions.CreateTransaction(
				TestContext(),
				operations.CreateTransactionRequest{
					PostTransaction: shared.PostTransaction{
						Metadata: map[string]any{},
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

			ret = response.TransactionsResponse
			t3 = shared.Transaction{
				Timestamp: ret.Data[0].Timestamp,
				Postings:  ret.Data[0].Postings,
				Reference: ret.Data[0].Reference,
				Metadata:  ret.Data[0].Metadata,
				Txid:      ret.Data[0].Txid,
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
			response, err := Client().Transactions.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("3"))

			response, err = Client().Transactions.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger:  "default",
					Account: ptr("foo:*"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("3"))

			response, err = Client().Transactions.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger:  "default",
					Account: ptr("not_existing"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("0"))

			response, err = Client().Transactions.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger:      "default",
					Destination: ptr("*:baz"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("1"))

			response, err = Client().Transactions.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger:      "default",
					Destination: ptr("not_existing"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("0"))

			response, err = Client().Transactions.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger: "default",
					Source: ptr("foo:*"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("0"))

			response, err = Client().Transactions.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger: "default",
					Source: ptr("world"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("3"))

			response, err = Client().Transactions.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger: "default",
					Metadata: map[string]any{
						"foo": "bar",
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("2"))

			response, err = Client().Transactions.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger: "default",
					Metadata: map[string]any{
						"foo": "not_existing",
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("0"))

			response, err = Client().Transactions.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger:    "default",
					StartTime: &timestamp2,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("2"))

			response, err = Client().Transactions.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger:    "default",
					StartTime: &timestamp3,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("1"))

			response, err = Client().Transactions.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger:    "default",
					StartTime: ptr(time.Now().UTC()),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("0"))

			response, err = Client().Transactions.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger:  "default",
					EndTime: &timestamp3,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("2"))

			response, err = Client().Transactions.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger:  "default",
					EndTime: &timestamp2,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("1"))

			response, err = Client().Transactions.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger:  "default",
					EndTime: ptr(time.Date(2023, 4, 9, 10, 0, 0, 0, time.UTC)),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("0"))
		})
		It("should be listed on api", func() {
			response, err := Client().Transactions.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse := response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[2]).Should(Equal(t1))

			response, err = Client().Transactions.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
					Account: ptr("foo:*"),
					Ledger:  "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[2]).Should(Equal(t1))

			response, err = Client().Transactions.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
					Account: ptr("not_existing"),
					Ledger:  "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(0))

			response, err = Client().Transactions.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
					Destination: ptr("foo:*"),
					Ledger:      "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[2]).Should(Equal(t1))

			response, err = Client().Transactions.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
					Destination: ptr("not_existing"),
					Ledger:      "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(0))

			response, err = Client().Transactions.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
					Source: ptr("foo:"),
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(0))

			response, err = Client().Transactions.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
					Source: ptr("world"),
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[2]).Should(Equal(t1))

			response, err = Client().Transactions.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
					Ledger: "default",
					Metadata: map[string]any{
						"foo": "bar",
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(2))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t1))

			response, err = Client().Transactions.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
					Ledger: "default",
					Metadata: map[string]any{
						"foo": "not_existing",
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(0))

			response, err = Client().Transactions.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
					Ledger:    "default",
					StartTime: &timestamp2,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(2))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))

			response, err = Client().Transactions.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
					Ledger:    "default",
					StartTime: &timestamp3,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(1))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))

			response, err = Client().Transactions.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
					Ledger:    "default",
					StartTime: ptr(time.Now().UTC()),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(0))

			response, err = Client().Transactions.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
					Ledger:  "default",
					EndTime: &timestamp3,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(2))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t1))

			response, err = Client().Transactions.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
					Ledger:  "default",
					EndTime: &timestamp2,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(1))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t1))

			response, err = Client().Transactions.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
					Ledger:  "default",
					EndTime: ptr(time.Date(2023, 4, 9, 10, 0, 0, 0, time.UTC)),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			transactionCursorResponse = response.TransactionsCursorResponse
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(0))
		})
		It("should be getable on api", func() {
			response, err := Client().Transactions.GetTransaction(
				TestContext(),
				operations.GetTransactionRequest{
					Ledger: "default",
					Txid:   t1.Txid,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.TransactionResponse.Data).Should(Equal(t1))

			response, err = Client().Transactions.GetTransaction(
				TestContext(),
				operations.GetTransactionRequest{
					Ledger: "default",
					Txid:   t2.Txid,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.TransactionResponse.Data).Should(Equal(t2))

			response, err = Client().Transactions.GetTransaction(
				TestContext(),
				operations.GetTransactionRequest{
					Ledger: "default",
					Txid:   t3.Txid,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.TransactionResponse.Data).Should(Equal(t3))

			response, err = Client().Transactions.GetTransaction(
				TestContext(),
				operations.GetTransactionRequest{
					Ledger: "default",
					Txid:   666,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(404))
		})
	})
})

var _ = Given("some empty environment", func() {
	When("counting and listing transactions empty", func() {
		It("should be countable on api even if empty", func() {
			response, err := Client().Transactions.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("0"))
		})
		It("should be listed on api even if empty", func() {
			response, err := Client().Transactions.ListTransactions(
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
