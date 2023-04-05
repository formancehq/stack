package suite

import (
	"fmt"

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
			transactions []formance.Transaction
		)
		BeforeEach(func() {
			for i := 0; i < int(txCount); i++ {
				ret, _, err := Client().TransactionsApi.
					CreateTransaction(TestContext(), "default").
					PostTransaction(formance.PostTransaction{
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
				transactions = append([]formance.Transaction{
					ret.Data,
				}, transactions...)
			}
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
