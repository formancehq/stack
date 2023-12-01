package suite

import (
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Ledger}, func() {
	BeforeEach(func() {
		createLedgerResponse, err := Client().Ledger.CreateLedger(TestContext(), operations.CreateLedgerRequest{
			Ledger: "default",
		})
		Expect(err).To(BeNil())
		Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))
	})
	When("creating a transaction on a ledger", func() {
		var (
			timestamp = time.Now().Round(time.Second).UTC()
			rsp       *shared.CreateTransactionResponse
		)
		BeforeEach(func() {
			// Create a transaction
			response, err := Client().Ledger.V2.CreateTransaction(
				TestContext(),
				operations.CreateTransactionRequest{
					PostTransaction: shared.PostTransaction{
						Metadata: map[string]string{},
						Postings: []shared.Posting{
							{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Source:      "world",
								Destination: "alice",
							},
						},
						Timestamp: &timestamp,
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			rsp = response.CreateTransactionResponse

			// Check existence on api
			getResponse, err := Client().Ledger.V2.GetTransaction(
				TestContext(),
				operations.GetTransactionRequest{
					Ledger: "default",
					ID:     rsp.Data.ID,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(getResponse.StatusCode).To(Equal(200))
		})
		It("should fail if the transaction does not exist", func() {
			metadata := map[string]string{
				"foo": "bar",
			}

			response, err := Client().Ledger.V2.AddMetadataOnTransaction(
				TestContext(),
				operations.AddMetadataOnTransactionRequest{
					RequestBody: metadata,
					Ledger:      "default",
					ID:          big.NewInt(666),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(404))
		})
		Then("adding a metadata", func() {
			metadata := map[string]string{
				"foo": "bar",
			}
			BeforeEach(func() {
				response, err := Client().Ledger.V2.AddMetadataOnTransaction(
					TestContext(),
					operations.AddMetadataOnTransactionRequest{
						RequestBody: metadata,
						Ledger:      "default",
						ID:          rsp.Data.ID,
					},
				)
				Expect(err).To(Succeed())
				Expect(response.StatusCode).To(Equal(204))
			})
			It("should be available on api", func() {
				// Check existence on api
				response, err := Client().Ledger.V2.GetTransaction(
					TestContext(),
					operations.GetTransactionRequest{
						Ledger: "default",
						ID:     rsp.Data.ID,
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				Expect(response.GetTransactionResponse.Data.Metadata).Should(Equal(metadata))
			})
		})
	})
})
