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
	When("counting and listing accounts", func() {
		var (
			metadata1 = map[string]string{
				"clientType": "gold",
			}

			metadata2 = map[string]string{
				"clientType": "silver",
			}

			timestamp = time.Now().Round(time.Second).UTC()
		)
		BeforeEach(func() {
			// Subscribe to nats subject
			response, err := Client().Ledger.AddMetadataToAccount(
				TestContext(),
				operations.AddMetadataToAccountRequest{
					RequestBody: metadata1,
					Address:     "foo:foo",
					Ledger:      "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))

			response, err = Client().Ledger.AddMetadataToAccount(
				TestContext(),
				operations.AddMetadataToAccountRequest{
					RequestBody: metadata2,
					Address:     "foo:bar",
					Ledger:      "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))

			createTransactionResponse, err := Client().Ledger.CreateTransaction(
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
						Timestamp: &timestamp,
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(createTransactionResponse.StatusCode).To(Equal(200))
		})
		It("should be countable on api", func() {
			response, err := Client().Ledger.CountAccounts(
				TestContext(),
				operations.CountAccountsRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("3"))
		})
		It("should be listed on api", func() {
			response, err := Client().Ledger.ListAccounts(
				TestContext(),
				operations.ListAccountsRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			accountsCursorResponse := response.AccountsCursorResponse
			Expect(accountsCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(accountsCursorResponse.Cursor.Data[0]).To(Equal(shared.Account{
				Address:  "foo:foo",
				Metadata: metadata1,
			}))
			Expect(accountsCursorResponse.Cursor.Data[1]).To(Equal(shared.Account{
				Address:  "foo:bar",
				Metadata: metadata2,
			}))

			Expect(accountsCursorResponse.Cursor.Data[2]).To(Equal(shared.Account{
				Address:  "world",
				Metadata: metadata.Metadata{},
			}))
		})
		It("should be listed on api using address filters", func() {
			response, err := Client().Ledger.ListAccounts(
				TestContext(),
				operations.ListAccountsRequest{
					Address: ptr("foo:.*"),
					Ledger:  "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			accountsCursorResponse := response.AccountsCursorResponse
			Expect(accountsCursorResponse.Cursor.Data).To(HaveLen(2))
			Expect(accountsCursorResponse.Cursor.Data[0]).To(Equal(shared.Account{
				Address:  "foo:foo",
				Metadata: metadata1,
			}))
			Expect(accountsCursorResponse.Cursor.Data[1]).To(Equal(shared.Account{
				Address:  "foo:bar",
				Metadata: metadata2,
			}))

			response, err = Client().Ledger.ListAccounts(
				TestContext(),
				operations.ListAccountsRequest{
					Address: ptr(".*:foo"),
					Ledger:  "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			accountsCursorResponse = response.AccountsCursorResponse
			Expect(accountsCursorResponse.Cursor.Data).To(HaveLen(1))
			Expect(accountsCursorResponse.Cursor.Data[0]).To(Equal(shared.Account{
				Address:  "foo:foo",
				Metadata: metadata1,
			}))
		})
		It("should be listed on api using balance filters", func() {
			response, err := Client().Ledger.ListAccounts(
				TestContext(),
				operations.ListAccountsRequest{
					Balance:         ptr(int64(90)),
					BalanceOperator: ptr(operations.ListAccountsBalanceOperatorLte),
					Ledger:          "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			accountsCursorResponse := response.AccountsCursorResponse
			Expect(accountsCursorResponse.Cursor.Data).To(HaveLen(2))
			Expect(accountsCursorResponse.Cursor.Data[0]).To(Equal(shared.Account{
				Address:  "foo:bar",
				Metadata: metadata2,
			}))
			Expect(accountsCursorResponse.Cursor.Data[1]).To(Equal(shared.Account{
				Address:  "world",
				Metadata: metadata.Metadata{},
			}))

			// Default operator should be gte
			response, err = Client().Ledger.ListAccounts(
				TestContext(),
				operations.ListAccountsRequest{
					Balance: ptr(int64(90)),
					Ledger:  "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			accountsCursorResponse = response.AccountsCursorResponse
			Expect(accountsCursorResponse.Cursor.Data).To(HaveLen(1))
			Expect(accountsCursorResponse.Cursor.Data[0]).To(Equal(shared.Account{
				Address:  "foo:foo",
				Metadata: metadata1,
			}))
		})
		It("should be listed on api using metadata filters", func() {
			response, err := Client().Ledger.ListAccounts(
				TestContext(),
				operations.ListAccountsRequest{
					Ledger: "default",
					Metadata: map[string]string{
						"clientType": "gold",
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			accountsCursorResponse := response.AccountsCursorResponse
			Expect(accountsCursorResponse.Cursor.Data).To(HaveLen(1))
			Expect(accountsCursorResponse.Cursor.Data[0]).To(Equal(shared.Account{
				Address:  "foo:foo",
				Metadata: metadata1,
			}))
		})
	})
})

var _ = Given("some empty environment", func() {
	When("counting and listing accounts empty", func() {
		It("should be countable on api even if empty", func() {
			response, err := Client().Ledger.CountAccounts(
				TestContext(),
				operations.CountAccountsRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))

			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("0"))
		})
		It("should be listed on api even if empty", func() {
			response, err := Client().Ledger.ListAccounts(
				TestContext(),
				operations.ListAccountsRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			Expect(response.AccountsCursorResponse.Cursor.Data).To(HaveLen(0))
		})
	})
})

var _ = Given("some environment with accounts", func() {
	const (
		pageSize      = int64(10)
		accountCounts = 2 * pageSize
	)
	When("creating accounts", func() {
		var (
			accounts []shared.Account
		)
		BeforeEach(func() {
			for i := 0; i < int(accountCounts); i++ {
				m := map[string]string{
					"id": fmt.Sprintf("%d", i),
				}

				response, err := Client().Ledger.AddMetadataToAccount(
					TestContext(),
					operations.AddMetadataToAccountRequest{
						RequestBody: m,
						Address:     fmt.Sprintf("foo:%d", i),
						Ledger:      "default",
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(204))

				accounts = append(accounts, shared.Account{
					Address:  fmt.Sprintf("foo:%d", i),
					Metadata: m,
				})
			}
		})
		AfterEach(func() {
			accounts = nil
		})
		Then(fmt.Sprintf("listing accounts using page size of %d", pageSize), func() {
			var (
				rsp *shared.AccountsCursorResponse
			)
			BeforeEach(func() {
				response, err := Client().Ledger.ListAccounts(
					TestContext(),
					operations.ListAccountsRequest{
						Ledger:   "default",
						PageSize: ptr(pageSize),
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				rsp = response.AccountsCursorResponse
				Expect(rsp.Cursor.HasMore).To(BeTrue())
				Expect(rsp.Cursor.Previous).To(BeNil())
				Expect(rsp.Cursor.Next).NotTo(BeNil())
			})
			It("should return the first page", func() {
				Expect(rsp.Cursor.PageSize).To(Equal(pageSize))
				Expect(rsp.Cursor.Data).To(Equal(accounts[:pageSize]))
			})
			Then("following next cursor", func() {
				BeforeEach(func() {
					response, err := Client().Ledger.ListAccounts(
						TestContext(),
						operations.ListAccountsRequest{
							Cursor: rsp.Cursor.Next,
							Ledger: "default",
						},
					)
					Expect(err).ToNot(HaveOccurred())
					Expect(response.StatusCode).To(Equal(200))

					rsp = response.AccountsCursorResponse
				})
				It("should return next page", func() {
					Expect(rsp.Cursor.PageSize).To(Equal(pageSize))
					Expect(rsp.Cursor.Data).To(Equal(accounts[pageSize : 2*pageSize]))
					Expect(rsp.Cursor.Next).To(BeNil())
				})
				Then("following previous cursor", func() {
					BeforeEach(func() {
						response, err := Client().Ledger.ListAccounts(
							TestContext(),
							operations.ListAccountsRequest{
								Ledger: "default",
								Cursor: rsp.Cursor.Previous,
							},
						)
						Expect(err).ToNot(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))

						rsp = response.AccountsCursorResponse
					})
					It("should return first page", func() {
						Expect(rsp.Cursor.PageSize).To(Equal(pageSize))
						Expect(rsp.Cursor.Data).To(Equal(accounts[:pageSize]))
						Expect(rsp.Cursor.Previous).To(BeNil())
					})
				})
			})
		})
	})
})
