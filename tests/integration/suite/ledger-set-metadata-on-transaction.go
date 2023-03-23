package suite_test

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
			Expect(err).To(BeNil())

			// Check existence on api
			_, _, err := Client().TransactionsApi.
				GetTransaction(TestContext(), "default", rsp.Data.Txid).
				Execute()
			Expect(err).To(BeNil())
		})
		Then("adding a metadata", func() {
			var (
				metadata = map[string]any{
					"foo": "bar",
				}
			)
			BeforeEach(func() {
				_, err := Client().TransactionsApi.
					AddMetadataOnTransaction(TestContext(), "default", rsp.Data.Txid).
					RequestBody(metadata).
					Execute()
				Expect(err).To(Succeed())
			})
			It("should eventually be available on api", func() {
				// Check existence on api
				transaction, _, err := Client().TransactionsApi.
					GetTransaction(TestContext(), "default", rsp.Data.Txid).
					Execute()
				Expect(err).To(BeNil())
				Expect(transaction.Data.Metadata).Should(Equal(metadata))
			})
		})
	})
})
