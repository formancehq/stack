package suite

import (
	"fmt"
	"time"

	"github.com/formancehq/formance-sdk-go"
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
			transactions []formance.ExpandedTransaction
		)
		BeforeEach(func() {
			for i := 0; i < int(txCount); i++ {
				ret, _, err := Client().TransactionsApi.
					CreateTransaction(TestContext(), "default").
					PostTransaction(formance.PostTransaction{
						Timestamp: &timestamp,
						Postings: []formance.Posting{{
							Amount:      100,
							Asset:       "USD",
							Source:      "world",
							Destination: fmt.Sprintf("account:%d", i),
						}},
						Metadata: metadata.Metadata{},
					}).
					Execute()
				Expect(err).ToNot(HaveOccurred())
				transactions = append([]formance.ExpandedTransaction{
					{
						Timestamp: ret.Data.Timestamp,
						Postings:  ret.Data.Postings,
						Reference: ret.Data.Reference,
						Metadata:  ret.Data.Metadata,
						Txid:      ret.Data.Txid,
						PreCommitVolumes: map[string]map[string]formance.Volume{
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
						PostCommitVolumes: map[string]map[string]formance.Volume{
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
				rsp *formance.TransactionsCursorResponse
				err error
			)
			BeforeEach(func() {
				rsp, _, err = Client().TransactionsApi.
					ListTransactions(TestContext(), "default").
					PageSize(pageSize).
					Execute()
				Expect(err).ToNot(HaveOccurred())
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
					rsp, _, err = Client().TransactionsApi.
						ListTransactions(TestContext(), "default").
						Cursor(*rsp.Cursor.Next).
						Execute()
					Expect(err).ToNot(HaveOccurred())
				})
				It("should return next page", func() {
					Expect(rsp.Cursor.PageSize).To(Equal(pageSize))
					Expect(rsp.Cursor.Data).To(Equal(transactions[pageSize : 2*pageSize]))
					Expect(rsp.Cursor.Next).To(BeNil())
				})
				Then("following previous cursor", func() {
					BeforeEach(func() {
						rsp, _, err = Client().TransactionsApi.
							ListTransactions(TestContext(), "default").
							Cursor(*rsp.Cursor.Previous).
							Execute()
						Expect(err).ToNot(HaveOccurred())
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
		t1 formance.ExpandedTransaction
		t2 formance.ExpandedTransaction
		t3 formance.ExpandedTransaction
	)
	When("creating transactions", func() {
		BeforeEach(func() {
			ret, _, err := Client().TransactionsApi.
				CreateTransaction(TestContext(), "default").
				PostTransaction(formance.PostTransaction{
					Timestamp: &timestamp1,
					Postings: []formance.Posting{{
						Amount:      100,
						Asset:       "USD",
						Source:      "world",
						Destination: "foo:foo",
					}},
					Metadata: m1,
				}).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			t1 = formance.ExpandedTransaction{
				Timestamp: ret.Data.Timestamp,
				Postings:  ret.Data.Postings,
				Reference: ret.Data.Reference,
				Metadata:  ret.Data.Metadata,
				Txid:      ret.Data.Txid,
				PreCommitVolumes: map[string]map[string]formance.Volume{
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
				PostCommitVolumes: map[string]map[string]formance.Volume{
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

			ret, _, err = Client().TransactionsApi.
				CreateTransaction(TestContext(), "default").
				PostTransaction(formance.PostTransaction{
					Timestamp: &timestamp2,
					Postings: []formance.Posting{{
						Amount:      100,
						Asset:       "USD",
						Source:      "world",
						Destination: "foo:bar",
					}},
					Metadata: m1,
				}).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			t2 = formance.ExpandedTransaction{
				Timestamp: ret.Data.Timestamp,
				Postings:  ret.Data.Postings,
				Reference: ret.Data.Reference,
				Metadata:  ret.Data.Metadata,
				Txid:      ret.Data.Txid,
				PreCommitVolumes: map[string]map[string]formance.Volume{
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
				PostCommitVolumes: map[string]map[string]formance.Volume{
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

			ret, _, err = Client().TransactionsApi.
				CreateTransaction(TestContext(), "default").
				PostTransaction(formance.PostTransaction{
					Timestamp: &timestamp3,
					Postings: []formance.Posting{{
						Amount:      100,
						Asset:       "USD",
						Source:      "world",
						Destination: "foo:baz",
					}},
					Metadata: metadata.Metadata{},
				}).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			t3 = formance.ExpandedTransaction{
				Timestamp: ret.Data.Timestamp,
				Postings:  ret.Data.Postings,
				Reference: ret.Data.Reference,
				Metadata:  ret.Data.Metadata,
				Txid:      ret.Data.Txid,
				PreCommitVolumes: map[string]map[string]formance.Volume{
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
				PostCommitVolumes: map[string]map[string]formance.Volume{
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
			transactionResponse, err := Client().TransactionsApi.
				CountTransactions(TestContext(), "default").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionResponse.Header.Get("Count")).Should(Equal("3"))

			transactionResponse, err = Client().TransactionsApi.
				CountTransactions(TestContext(), "default").
				Account("foo:").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionResponse.Header.Get("Count")).Should(Equal("3"))

			transactionResponse, err = Client().TransactionsApi.
				CountTransactions(TestContext(), "default").
				Account("not_existing").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionResponse.Header.Get("Count")).Should(Equal("0"))

			transactionResponse, err = Client().TransactionsApi.
				CountTransactions(TestContext(), "default").
				Destination(":baz").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionResponse.Header.Get("Count")).Should(Equal("1"))

			transactionResponse, err = Client().TransactionsApi.
				CountTransactions(TestContext(), "default").
				Destination("not_existing").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionResponse.Header.Get("Count")).Should(Equal("0"))

			transactionResponse, err = Client().TransactionsApi.
				CountTransactions(TestContext(), "default").
				Source("foo:").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionResponse.Header.Get("Count")).Should(Equal("0"))

			transactionResponse, err = Client().TransactionsApi.
				CountTransactions(TestContext(), "default").
				Source("world").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionResponse.Header.Get("Count")).Should(Equal("3"))

			transactionResponse, err = Client().TransactionsApi.
				CountTransactions(TestContext(), "default").
				Metadata(map[string]string{
					"foo": "bar",
				}).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionResponse.Header.Get("Count")).Should(Equal("2"))

			transactionResponse, err = Client().TransactionsApi.
				CountTransactions(TestContext(), "default").
				Metadata(map[string]string{
					"foo": "not_existing",
				}).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionResponse.Header.Get("Count")).Should(Equal("0"))

			transactionResponse, err = Client().TransactionsApi.
				CountTransactions(TestContext(), "default").
				StartTime(timestamp2).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionResponse.Header.Get("Count")).Should(Equal("2"))

			transactionResponse, err = Client().TransactionsApi.
				CountTransactions(TestContext(), "default").
				StartTime(timestamp3).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionResponse.Header.Get("Count")).Should(Equal("1"))

			transactionResponse, err = Client().TransactionsApi.
				CountTransactions(TestContext(), "default").
				StartTime(time.Now().UTC()).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionResponse.Header.Get("Count")).Should(Equal("0"))

			transactionResponse, err = Client().TransactionsApi.
				CountTransactions(TestContext(), "default").
				EndTime(timestamp3).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionResponse.Header.Get("Count")).Should(Equal("2"))

			transactionResponse, err = Client().TransactionsApi.
				CountTransactions(TestContext(), "default").
				EndTime(timestamp2).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionResponse.Header.Get("Count")).Should(Equal("1"))

			transactionResponse, err = Client().TransactionsApi.
				CountTransactions(TestContext(), "default").
				EndTime(time.Date(2023, 4, 9, 10, 0, 0, 0, time.UTC)).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionResponse.Header.Get("Count")).Should(Equal("0"))
		})
		It("should be listed on api", func() {
			transactionCursorResponse, _, err := Client().TransactionsApi.
				ListTransactions(TestContext(), "default").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[2]).Should(Equal(t1))

			transactionCursorResponse, _, err = Client().TransactionsApi.
				ListTransactions(TestContext(), "default").
				Account("foo:").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[2]).Should(Equal(t1))

			transactionCursorResponse, _, err = Client().TransactionsApi.
				ListTransactions(TestContext(), "default").
				Account("foo:").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[2]).Should(Equal(t1))

			transactionCursorResponse, _, err = Client().TransactionsApi.
				ListTransactions(TestContext(), "default").
				Account("not_existing").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(0))

			transactionCursorResponse, _, err = Client().TransactionsApi.
				ListTransactions(TestContext(), "default").
				Destination("foo:").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[2]).Should(Equal(t1))

			transactionCursorResponse, _, err = Client().TransactionsApi.
				ListTransactions(TestContext(), "default").
				Destination("not_existing").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(0))

			transactionCursorResponse, _, err = Client().TransactionsApi.
				ListTransactions(TestContext(), "default").
				Source("foo:").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(0))

			transactionCursorResponse, _, err = Client().TransactionsApi.
				ListTransactions(TestContext(), "default").
				Source("world").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[2]).Should(Equal(t1))

			transactionCursorResponse, _, err = Client().TransactionsApi.
				ListTransactions(TestContext(), "default").
				Metadata(map[string]string{
					"foo": "bar",
				}).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(2))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t1))

			transactionCursorResponse, _, err = Client().TransactionsApi.
				ListTransactions(TestContext(), "default").
				Metadata(map[string]string{
					"foo": "not_existing",
				}).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(0))

			transactionCursorResponse, _, err = Client().TransactionsApi.
				ListTransactions(TestContext(), "default").
				StartTime(timestamp2).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(2))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t2))

			transactionCursorResponse, _, err = Client().TransactionsApi.
				ListTransactions(TestContext(), "default").
				StartTime(timestamp3).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(1))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t3))

			transactionCursorResponse, _, err = Client().TransactionsApi.
				ListTransactions(TestContext(), "default").
				StartTime(time.Now().UTC()).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(0))

			transactionCursorResponse, _, err = Client().TransactionsApi.
				ListTransactions(TestContext(), "default").
				EndTime(timestamp3).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(2))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t2))
			Expect(transactionCursorResponse.Cursor.Data[1]).Should(Equal(t1))

			transactionCursorResponse, _, err = Client().TransactionsApi.
				ListTransactions(TestContext(), "default").
				EndTime(timestamp2).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(1))
			Expect(transactionCursorResponse.Cursor.Data[0]).Should(Equal(t1))

			transactionCursorResponse, _, err = Client().TransactionsApi.
				ListTransactions(TestContext(), "default").
				EndTime(time.Date(2023, 4, 9, 10, 0, 0, 0, time.UTC)).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionCursorResponse.Cursor.Data).Should(HaveLen(0))
		})
		It("should be getable on api", func() {
			transactionResponse, _, err := Client().TransactionsApi.
				GetTransaction(TestContext(), "default", t1.Txid).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionResponse.Data).Should(Equal(t1))

			transactionResponse, _, err = Client().TransactionsApi.
				GetTransaction(TestContext(), "default", t2.Txid).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionResponse.Data).Should(Equal(t2))

			transactionResponse, _, err = Client().TransactionsApi.
				GetTransaction(TestContext(), "default", t3.Txid).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionResponse.Data).Should(Equal(t3))

			_, resp, err := Client().TransactionsApi.
				GetTransaction(TestContext(), "default", 666).
				Execute()
			Expect(err).To(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(404))
		})
	})
})

var _ = Given("some empty environment", func() {
	When("counting and listing transactions empty", func() {
		It("should be countable on api even if empty", func() {
			transactionResponse, err := Client().TransactionsApi.
				CountTransactions(TestContext(), "default").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionResponse.Header.Get("Count")).Should(Equal("0"))
		})
		It("should be listed on api even if empty", func() {
			transactionCursorResponse, _, err := Client().TransactionsApi.
				ListTransactions(TestContext(), "default").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(transactionCursorResponse.Cursor.Data).To(HaveLen(0))
		})
	})
})
