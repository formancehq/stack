package suite

import (
	"time"

	"github.com/formancehq/formance-sdk-go"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("some empty environment", func() {
	When("creating a transaction on a ledger", func() {
		var (
			timestamp = time.Now().Round(time.Second).UTC()
			err       error
			rsp       *formance.TransactionResponse
		)
		BeforeEach(func() {
			// Create a transaction
			rsp, _, err = Client().TransactionsApi.
				CreateTransaction(TestContext(), "default").
				PostTransaction(formance.PostTransaction{
					Timestamp: &timestamp,
					Postings: []formance.Posting{{
						Amount:      100,
						Asset:       "USD",
						Source:      "world",
						Destination: "alice",
					}},
				}).
				Execute()
			Expect(err).NotTo(HaveOccurred())

			// Check existence on api
			Eventually(func(g Gomega) bool {
				_, _, err := Client().TransactionsApi.
					GetTransaction(TestContext(), "default", rsp.Data.Txid).
					Execute()
				if err != nil {
					return false
				}
				return true
			}).Should(BeTrue())
		})
		Then("adding a metadata", func() {
			metadata := map[string]any{
				"foo": "bar",
			}
			BeforeEach(func() {
				_, err := Client().TransactionsApi.
					AddMetadataOnTransaction(TestContext(), "default", rsp.Data.Txid).
					RequestBody(metadata).
					Execute()
				Expect(err).NotTo(HaveOccurred())
			})
			It("should eventually be available on api", func() {
				// Check existence on api
				Eventually(func(g Gomega) map[string]any {
					transaction, _, err := Client().TransactionsApi.
						GetTransaction(TestContext(), "default", rsp.Data.Txid).
						Execute()
					if err != nil {
						return map[string]any{}
					}
					return transaction.Data.Metadata
				}).Should(Equal(metadata))
			})
		})
	})
})
