package suite

import (
	"fmt"
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/metadata"
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
			transactions []shared.ExpandedTransaction
		)
		BeforeEach(func() {
			for i := 0; i < int(txCount); i++ {
				response, err := Client().Ledger.CreateTransaction(
					TestContext(),
					operations.CreateTransactionRequest{
						IdempotencyKey: new(string),
						PostTransaction: shared.PostTransaction{
							Metadata: map[string]string{},
							Postings: []shared.Posting{
								{
									Amount:      100,
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
						Txid:      ret.Data.Txid,
						PreCommitVolumes: map[string]map[string]shared.Volume{
							"world": {
								"USD": {
									Input:   0,
									Output:  int64(i * 100),
									Balance: ptr(int64(-i * 100)),
								},
							},
							fmt.Sprintf("account:%d", i): {
								"USD": {
									Input:   0,
									Output:  0,
									Balance: ptr(int64(0)),
								},
							},
						},
						PostCommitVolumes: map[string]map[string]shared.Volume{
							"world": {
								"USD": {
									Input:   0,
									Output:  int64((i + 1) * 100),
									Balance: ptr(int64(-(i + 1) * 100)),
								},
							},
							fmt.Sprintf("account:%d", i): {
								"USD": {
									Input:   100,
									Output:  0,
									Balance: ptr(int64(100)),
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
					response, err := Client().Ledger.ListTransactions(
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
						response, err := Client().Ledger.ListTransactions(
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
								Amount:      100,
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
				Txid:      ret.Data.Txid,
				PreCommitVolumes: map[string]map[string]shared.Volume{
					"world": {
						"USD": {
							Input:   0,
							Output:  0,
							Balance: ptr(int64(0)),
						},
					},
					"foo:foo": {
						"USD": {
							Input:   0,
							Output:  0,
							Balance: ptr(int64(0)),
						},
					},
				},
				PostCommitVolumes: map[string]map[string]shared.Volume{
					"world": {
						"USD": {
							Input:   0,
							Output:  100,
							Balance: ptr(int64(-100)),
						},
					},
					"foo:foo": {
						"USD": {
							Input:   100,
							Output:  0,
							Balance: ptr(int64(100)),
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
								Amount:      100,
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
				Txid:      ret.Data.Txid,
				PreCommitVolumes: map[string]map[string]shared.Volume{
					"world": {
						"USD": {
							Input:   0,
							Output:  100,
							Balance: ptr(int64(-100)),
						},
					},
					"foo:bar": {
						"USD": {
							Input:   0,
							Output:  0,
							Balance: ptr(int64(0)),
						},
					},
				},
				PostCommitVolumes: map[string]map[string]shared.Volume{
					"world": {
						"USD": {
							Input:   0,
							Output:  200,
							Balance: ptr(int64(-200)),
						},
					},
					"foo:bar": {
						"USD": {
							Input:   100,
							Output:  0,
							Balance: ptr(int64(100)),
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
								Amount:      100,
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
				Txid:      ret.Data.Txid,
				PreCommitVolumes: map[string]map[string]shared.Volume{
					"world": {
						"USD": {
							Input:   0,
							Output:  200,
							Balance: ptr(int64(-200)),
						},
					},
					"foo:baz": {
						"USD": {
							Input:   0,
							Output:  0,
							Balance: ptr(int64(0)),
						},
					},
				},
				PostCommitVolumes: map[string]map[string]shared.Volume{
					"world": {
						"USD": {
							Input:   0,
							Output:  300,
							Balance: ptr(int64(-300)),
						},
					},
					"foo:baz": {
						"USD": {
							Input:   100,
							Output:  0,
							Balance: ptr(int64(100)),
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
					Ledger:  "default",
					Account: ptr("foo:"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("3"))

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger:  "default",
					Account: ptr("not_existing"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("0"))

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger:      "default",
					Destination: ptr(":baz"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("1"))

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger:      "default",
					Destination: ptr("not_existing"),
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
					Source: ptr("foo:"),
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
					Source: ptr("world"),
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
					Metadata: map[string]string{
						"foo": "bar",
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
					Metadata: map[string]string{
						"foo": "not_existing",
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
					Ledger:    "default",
					StartTime: &timestamp2,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("2"))

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger:    "default",
					StartTime: &timestamp3,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("1"))

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger:    "default",
					StartTime: ptr(time.Now().UTC()),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("0"))

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger:  "default",
					EndTime: &timestamp3,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("2"))

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger:  "default",
					EndTime: &timestamp2,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("1"))

			response, err = Client().Ledger.CountTransactions(
				TestContext(),
				operations.CountTransactionsRequest{
					Ledger:  "default",
					EndTime: ptr(time.Date(2023, 4, 9, 10, 0, 0, 0, time.UTC)),
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
					Account: ptr("foo:"),
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

			response, err = Client().Ledger.ListTransactions(
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

			response, err = Client().Ledger.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
					Destination: ptr("foo:"),
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

			response, err = Client().Ledger.ListTransactions(
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

			response, err = Client().Ledger.ListTransactions(
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

			response, err = Client().Ledger.ListTransactions(
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

			response, err = Client().Ledger.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
					Ledger: "default",
					Metadata: map[string]string{
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

			response, err = Client().Ledger.ListTransactions(
				TestContext(),
				operations.ListTransactionsRequest{
					Ledger: "default",
					Metadata: map[string]string{
						"foo": "not_existing",
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

			response, err = Client().Ledger.ListTransactions(
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

			response, err = Client().Ledger.ListTransactions(
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

			response, err = Client().Ledger.ListTransactions(
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

			response, err = Client().Ledger.ListTransactions(
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

			response, err = Client().Ledger.ListTransactions(
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
			response, err := Client().Ledger.GetTransaction(
				TestContext(),
				operations.GetTransactionRequest{
					Ledger: "default",
					Txid:   t1.Txid,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.GetTransactionResponse.Data).Should(Equal(t1))

			response, err = Client().Ledger.GetTransaction(
				TestContext(),
				operations.GetTransactionRequest{
					Ledger: "default",
					Txid:   t2.Txid,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.GetTransactionResponse.Data).Should(Equal(t2))

			response, err = Client().Ledger.GetTransaction(
				TestContext(),
				operations.GetTransactionRequest{
					Ledger: "default",
					Txid:   t3.Txid,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.GetTransactionResponse.Data).Should(Equal(t3))

			response, err = Client().Ledger.GetTransaction(
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
