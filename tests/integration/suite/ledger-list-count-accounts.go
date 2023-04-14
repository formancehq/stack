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
			_, err := Client().AccountsApi.
				AddMetadataToAccount(TestContext(), "default", "foo:foo").
				RequestBody(metadata1).
				Execute()
			Expect(err).ToNot(HaveOccurred())

			_, err = Client().AccountsApi.
				AddMetadataToAccount(TestContext(), "default", "foo:bar").
				RequestBody(metadata2).
				Execute()
			Expect(err).ToNot(HaveOccurred())

			_, _, err = Client().TransactionsApi.
				CreateTransaction(TestContext(), "default").
				PostTransaction(formance.PostTransaction{
					Timestamp: &timestamp,
					Postings: []formance.Posting{{
						Amount:      100,
						Asset:       "USD",
						Source:      "world",
						Destination: "foo:foo",
					}},
					Metadata: metadata.Metadata{},
				}).
				Execute()
			Expect(err).ToNot(HaveOccurred())
		})
		It("should be countable on api", func() {
			accountResponse, err := Client().AccountsApi.
				CountAccounts(TestContext(), "default").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(accountResponse.Header.Get("Count")).Should(Equal("3"))
		})
		It("should be listed on api", func() {
			accountsCursorResponse, _, err := Client().AccountsApi.
				ListAccounts(TestContext(), "default").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(accountsCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(accountsCursorResponse.Cursor.Data[0]).To(Equal(formance.Account{
				Address:  "foo:foo",
				Type:     nil,
				Metadata: metadata1,
			}))
			Expect(accountsCursorResponse.Cursor.Data[1]).To(Equal(formance.Account{
				Address:  "foo:bar",
				Type:     nil,
				Metadata: metadata2,
			}))

			Expect(accountsCursorResponse.Cursor.Data[2]).To(Equal(formance.Account{
				Address:  "world",
				Type:     nil,
				Metadata: metadata.Metadata{},
			}))
		})
		It("should be listed on api using address filters", func() {
			accountsCursorResponse, _, err := Client().AccountsApi.
				ListAccounts(TestContext(), "default").
				Address("foo:.*").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(accountsCursorResponse.Cursor.Data).To(HaveLen(2))
			Expect(accountsCursorResponse.Cursor.Data[0]).To(Equal(formance.Account{
				Address:  "foo:foo",
				Type:     nil,
				Metadata: metadata1,
			}))
			Expect(accountsCursorResponse.Cursor.Data[1]).To(Equal(formance.Account{
				Address:  "foo:bar",
				Type:     nil,
				Metadata: metadata2,
			}))

			accountsCursorResponse, _, err = Client().AccountsApi.
				ListAccounts(TestContext(), "default").
				Address("foo:f.*").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(accountsCursorResponse.Cursor.Data).To(HaveLen(1))
			Expect(accountsCursorResponse.Cursor.Data[0]).To(Equal(formance.Account{
				Address:  "foo:foo",
				Type:     nil,
				Metadata: metadata1,
			}))
		})
		It("should be listed on api using balance filters", func() {
			accountsCursorResponse, _, err := Client().AccountsApi.
				ListAccounts(TestContext(), "default").
				Balance(90).
				BalanceOperator("lte").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(accountsCursorResponse.Cursor.Data).To(HaveLen(2))
			Expect(accountsCursorResponse.Cursor.Data[0]).To(Equal(formance.Account{
				Address:  "foo:bar",
				Type:     nil,
				Metadata: metadata2,
			}))
			Expect(accountsCursorResponse.Cursor.Data[1]).To(Equal(formance.Account{
				Address:  "world",
				Type:     nil,
				Metadata: metadata.Metadata{},
			}))

			// Default operator should be gte
			accountsCursorResponse, _, err = Client().AccountsApi.
				ListAccounts(TestContext(), "default").
				Balance(90).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(accountsCursorResponse.Cursor.Data).To(HaveLen(1))
			Expect(accountsCursorResponse.Cursor.Data[0]).To(Equal(formance.Account{
				Address:  "foo:foo",
				Type:     nil,
				Metadata: metadata1,
			}))
		})
		It("should be listed on api using metadata filters", func() {
			accountsCursorResponse, _, err := Client().AccountsApi.
				ListAccounts(TestContext(), "default").
				Metadata(map[string]string{
					"clientType": "gold",
				}).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(accountsCursorResponse.Cursor.Data).To(HaveLen(1))
			Expect(accountsCursorResponse.Cursor.Data[0]).To(Equal(formance.Account{
				Address:  "foo:foo",
				Type:     nil,
				Metadata: metadata1,
			}))
		})
	})
})

var _ = Given("some empty environment", func() {
	When("counting and listing accounts empty", func() {
		It("should be countable on api even if empty", func() {
			accountResponse, err := Client().AccountsApi.
				CountAccounts(TestContext(), "default").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(accountResponse.Header.Get("Count")).Should(Equal("0"))
		})
		It("should be listed on api even if empty", func() {
			accountsCursorResponse, _, err := Client().AccountsApi.
				ListAccounts(TestContext(), "default").
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(accountsCursorResponse.Cursor.Data).To(HaveLen(0))
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
			accounts []formance.Account
		)
		BeforeEach(func() {
			for i := 0; i < int(accountCounts); i++ {
				m := map[string]string{
					"id": fmt.Sprintf("%d", i),
				}

				_, err := Client().AccountsApi.
					AddMetadataToAccount(TestContext(), "default", fmt.Sprintf("foo:%d", i)).
					RequestBody(m).
					Execute()
				Expect(err).ToNot(HaveOccurred())
				accounts = append(accounts, formance.Account{
					Address:  fmt.Sprintf("foo:%d", i),
					Type:     nil,
					Metadata: m,
				})
			}
		})
		AfterEach(func() {
			accounts = nil
		})
		Then(fmt.Sprintf("listing accounts using page size of %d", pageSize), func() {
			var (
				rsp *formance.AccountsCursorResponse
				err error
			)
			BeforeEach(func() {
				rsp, _, err = Client().AccountsApi.
					ListAccounts(TestContext(), "default").
					PageSize(pageSize).
					Execute()
				Expect(err).ToNot(HaveOccurred())
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
					rsp, _, err = Client().AccountsApi.
						ListAccounts(TestContext(), "default").
						Cursor(*rsp.Cursor.Next).
						Execute()
					Expect(err).ToNot(HaveOccurred())
				})
				It("should return next page", func() {
					Expect(rsp.Cursor.PageSize).To(Equal(pageSize))
					Expect(rsp.Cursor.Data).To(Equal(accounts[pageSize : 2*pageSize]))
					Expect(rsp.Cursor.Next).To(BeNil())
				})
				Then("following previous cursor", func() {
					BeforeEach(func() {
						rsp, _, err = Client().AccountsApi.
							ListAccounts(TestContext(), "default").
							Cursor(*rsp.Cursor.Previous).
							Execute()
						Expect(err).ToNot(HaveOccurred())
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
