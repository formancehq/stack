package suite

import (
	"fmt"
	"sort"
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("some empty environment", func() {
	When("get balances and aggregate balances", func() {
		var (
			timestamp1 = time.Now().Add(-10 * time.Hour).Round(time.Second).UTC()
			timestamp2 = time.Now().Add(-9 * time.Hour).Round(time.Second).UTC()
		)
		BeforeEach(func() {
			response, err := Client().Ledger.CreateTransaction(
				TestContext(),
				operations.CreateTransactionRequest{
					PostTransaction: shared.PostTransaction{
						Metadata: map[string]string{},
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

			response, err = Client().Ledger.CreateTransaction(
				TestContext(),
				operations.CreateTransactionRequest{
					PostTransaction: shared.PostTransaction{
						Metadata: map[string]string{},
						Postings: []shared.Posting{
							{
								Amount:      200,
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
		})
		It("should be listed on api with GetBalances", func() {
			response, err := Client().Ledger.GetBalances(
				TestContext(),
				operations.GetBalancesRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
			Expect(response.BalancesCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(response.BalancesCursorResponse.Cursor.Data[0]).To(Equal(map[string]map[string]int64{
				"world": {
					"USD": -300,
				},
			}))
			Expect(response.BalancesCursorResponse.Cursor.Data[1]).To(Equal(map[string]map[string]int64{
				"foo:foo": {
					"USD": 100,
				},
			}))
			Expect(response.BalancesCursorResponse.Cursor.Data[2]).To(Equal(map[string]map[string]int64{
				"foo:bar": {
					"USD": 200,
				},
			}))
		})
		It("should be listed on api with GetBalances using accounts filter", func() {
			response, err := Client().Ledger.GetBalances(
				TestContext(),
				operations.GetBalancesRequest{
					Address: ptr("foo:.*"),
					Ledger:  "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			balancesCursorResponse := response.BalancesCursorResponse
			Expect(balancesCursorResponse.Cursor.Data).To(HaveLen(2))
			Expect(balancesCursorResponse.Cursor.Data[0]).To(Equal(map[string]map[string]int64{
				"foo:foo": {
					"USD": 100,
				},
			}))
			Expect(balancesCursorResponse.Cursor.Data[1]).To(Equal(map[string]map[string]int64{
				"foo:bar": {
					"USD": 200,
				},
			}))

			response, err = Client().Ledger.GetBalances(
				TestContext(),
				operations.GetBalancesRequest{
					Address: ptr(".*:foo"),
					Ledger:  "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			balancesCursorResponse = response.BalancesCursorResponse
			Expect(balancesCursorResponse.Cursor.Data).To(HaveLen(1))
			Expect(balancesCursorResponse.Cursor.Data[0]).To(Equal(map[string]map[string]int64{
				"foo:foo": {
					"USD": 100,
				},
			}))
		})
		It("should be listed on api with GetBalancesAggregated", func() {
			response, err := Client().Ledger.GetBalancesAggregated(
				TestContext(),
				operations.GetBalancesAggregatedRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			balances := response.AggregateBalancesResponse
			Expect(balances.Data).To(Equal(map[string]int64{
				"USD": 0,
			}))
		})
		It("should be listed on api with GetBalancesAggregated using accounts filter", func() {
			response, err := Client().Ledger.GetBalancesAggregated(
				TestContext(),
				operations.GetBalancesAggregatedRequest{
					Address: ptr("foo:.*"),
					Ledger:  "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			balances := response.AggregateBalancesResponse
			Expect(balances.Data).To(Equal(map[string]int64{
				// Should be the sum of the two accounts foo:foo and foo:bar
				"USD": 300,
			}))
		})
	})
})

var _ = Given("some empty environment", func() {
	When("get balances and aggregate balances emtpy", func() {
		It("should be listed on api with GetBalances even if empty", func() {
			response, err := Client().Ledger.GetBalances(
				TestContext(),
				operations.GetBalancesRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			balancesCursorResponse := response.BalancesCursorResponse
			Expect(balancesCursorResponse.Cursor.Data).To(HaveLen(0))
		})
		It("should be listed on api with GetBalancesAggregated", func() {
			response, err := Client().Ledger.GetBalancesAggregated(
				TestContext(),
				operations.GetBalancesAggregatedRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			balances := response.AggregateBalancesResponse
			Expect(balances.Data).To(Equal(map[string]int64{}))
		})
	})
})

var _ = Given("some environment with accounts and transactions", func() {
	const (
		pageSize          = int64(10)
		transactionCounts = 2 * pageSize
	)
	When("creating transactions", func() {
		var (
			balances []map[string]map[string]int64
		)
		BeforeEach(func() {
			for i := 0; i < int(transactionCounts); i++ {
				now := time.Now()
				asset := "USD"
				if i%2 == 0 {
					asset = "EUR"
				}

				_, err := Client().Ledger.CreateTransaction(
					TestContext(),
					operations.CreateTransactionRequest{
						PostTransaction: shared.PostTransaction{
							Metadata: map[string]string{},
							Postings: []shared.Posting{
								{
									Amount:      100,
									Asset:       asset,
									Source:      "world",
									Destination: fmt.Sprintf("foo:%d", i),
								},
							},
							Timestamp: &now,
						},
						Ledger: "default",
					},
				)
				Expect(err).ToNot(HaveOccurred())

				balances = append(balances, map[string]map[string]int64{
					fmt.Sprintf("foo:%d", i): {
						asset: 100,
					},
				})
			}

			sort.Slice(balances, func(i, j int) bool {
				name1 := ""
				for name := range balances[i] {
					name1 = name
					break
				}
				name2 := ""
				for name := range balances[j] {
					name2 = name
					break
				}
				return name1 > name2
			})

			balances = append([]map[string]map[string]int64{
				{
					"world": {
						"USD": -transactionCounts / 2 * 100,
						"EUR": -transactionCounts / 2 * 100,
					},
				},
			}, balances...)
		})
		AfterEach(func() {
			balances = nil
		})
		Then(fmt.Sprintf("listing balances using page size of %d", pageSize), func() {
			var (
				rsp *shared.BalancesCursorResponse
			)
			BeforeEach(func() {
				response, err := Client().Ledger.GetBalances(
					TestContext(),
					operations.GetBalancesRequest{
						Ledger:   "default",
						PageSize: ptr(pageSize),
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				rsp = response.BalancesCursorResponse
				Expect(rsp.Cursor.HasMore).To(BeTrue())
				Expect(rsp.Cursor.Previous).To(BeNil())
				Expect(rsp.Cursor.Next).NotTo(BeNil())
			})
			It("should return the first page", func() {
				Expect(rsp.Cursor.PageSize).To(Equal(pageSize))
				Expect(rsp.Cursor.Data).To(Equal(balances[:pageSize]))
			})
			Then("following next cursor", func() {
				BeforeEach(func() {
					response, err := Client().Ledger.GetBalances(
						TestContext(),
						operations.GetBalancesRequest{
							Cursor: rsp.Cursor.Next,
							Ledger: "default",
						},
					)
					Expect(err).ToNot(HaveOccurred())
					Expect(response.StatusCode).To(Equal(200))

					rsp = response.BalancesCursorResponse
				})
				It("should return next page", func() {
					Expect(rsp.Cursor.PageSize).To(Equal(pageSize))
					Expect(rsp.Cursor.Data).To(Equal(balances[pageSize : 2*pageSize]))
				})
				Then("following next cursor", func() {
					BeforeEach(func() {
						response, err := Client().Ledger.GetBalances(
							TestContext(),
							operations.GetBalancesRequest{
								Cursor: rsp.Cursor.Next,
								Ledger: "default",
							},
						)
						Expect(err).ToNot(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))

						rsp = response.BalancesCursorResponse
					})
					It("should return next page", func() {
						Expect(rsp.Cursor.PageSize).To(Equal(pageSize))
						Expect(rsp.Cursor.Data).To(Equal(balances[len(balances)-1:]))
						Expect(rsp.Cursor.Next).To(BeNil())
					})
					Then("following previous cursor", func() {
						BeforeEach(func() {
							response, err := Client().Ledger.GetBalances(
								TestContext(),
								operations.GetBalancesRequest{
									Cursor: rsp.Cursor.Previous,
									Ledger: "default",
								},
							)
							Expect(err).ToNot(HaveOccurred())
							Expect(response.StatusCode).To(Equal(200))

							rsp = response.BalancesCursorResponse
						})
						It("should return first page", func() {
							Expect(rsp.Cursor.PageSize).To(Equal(pageSize))
							Expect(rsp.Cursor.Data).To(Equal(balances[pageSize : 2*pageSize]))
						})
					})
				})

			})
		})
	})
})
