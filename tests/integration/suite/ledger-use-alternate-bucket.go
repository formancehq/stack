package suite

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/big"
	"net/http"
)

var _ = WithModules([]*Module{modules.Ledger}, func() {
	When("Creating a ledger on an alternate bucket", func() {
		var (
			ledger1 string
		)
		BeforeEach(func() {
			ledger1 = uuid.NewString()
			response, err := Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
				V2CreateLedgerRequest: &shared.V2CreateLedgerRequest{
					Bucket: pointer.For("bucket0"),
				},
				Ledger: ledger1,
			})
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(http.StatusNoContent))
		})
		Then("Creating a tx on this ledger", func() {
			BeforeEach(func() {
				// Create a transaction
				response, err := Client().Ledger.V2.CreateTransaction(
					TestContext(),
					operations.V2CreateTransactionRequest{
						V2PostTransaction: shared.V2PostTransaction{
							Metadata: map[string]string{},
							Postings: []shared.V2Posting{
								{
									Amount:      big.NewInt(100),
									Asset:       "USD",
									Source:      "world",
									Destination: "alice",
								},
							},
						},
						Ledger: ledger1,
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))
			})
			Then("creating another ledger on the same bucket", func() {
				var (
					ledger2 string
				)
				BeforeEach(func() {
					ledger2 = uuid.NewString()
					response, err := Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
						V2CreateLedgerRequest: &shared.V2CreateLedgerRequest{
							Bucket: pointer.For("bucket0"),
						},
						Ledger: ledger2,
					})
					Expect(err).To(BeNil())
					Expect(response.StatusCode).To(Equal(http.StatusNoContent))
				})
				Then("Creating another tx on this new ledger", func() {
					BeforeEach(func() {
						// Create a transaction
						response, err := Client().Ledger.V2.CreateTransaction(
							TestContext(),
							operations.V2CreateTransactionRequest{
								V2PostTransaction: shared.V2PostTransaction{
									Metadata: map[string]string{},
									Postings: []shared.V2Posting{
										{
											Amount:      big.NewInt(100),
											Asset:       "USD",
											Source:      "world",
											Destination: "alice",
										},
									},
								},
								Ledger: ledger2,
							},
						)
						Expect(err).ToNot(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))
					})
					It("should have one tx on both ledger", func() {
						response, err := Client().Ledger.V2.ListTransactions(TestContext(), operations.V2ListTransactionsRequest{
							Ledger: ledger1,
						})
						Expect(err).ToNot(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))
						Expect(response.V2TransactionsCursorResponse.Cursor.Data).To(HaveLen(1))

						response, err = Client().Ledger.V2.ListTransactions(TestContext(), operations.V2ListTransactionsRequest{
							Ledger: ledger2,
						})
						Expect(err).ToNot(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))
						Expect(response.V2TransactionsCursorResponse.Cursor.Data).To(HaveLen(1))
					})
				})
			})
		})
	})
})
