package suite

import (
	"time"

	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("some empty environment", func() {
	When("creating a transaction on a ledger", func() {
		var (
			timestamp = time.Now().Round(time.Second).UTC()
			err       error
			rsp       *formance.CreateTransactionResponse
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
					Metadata: metadata.Metadata{},
				}).
				Execute()
			Expect(err).ToNot(HaveOccurred())

			// Check existence on api
			_, _, err := Client().TransactionsApi.
				GetTransaction(TestContext(), "default", rsp.Data.Txid).
				Execute()
			Expect(err).ToNot(HaveOccurred())
		})
		It("should fail if the transaction does not exist", func() {
			metadata := map[string]string{
				"foo": "bar",
			}

			resp, err := Client().TransactionsApi.
				AddMetadataOnTransaction(TestContext(), "default", 666).
				RequestBody(metadata).
				Execute()
			Expect(err).To(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(404))
		})
		Then("adding a metadata", func() {
			metadata := map[string]string{
				"foo": "bar",
			}
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
				Expect(err).ToNot(HaveOccurred())
				Expect(transaction.Data.Metadata).Should(Equal(metadata))
			})
		})
	})
})
