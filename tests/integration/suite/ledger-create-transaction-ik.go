package suite

import (
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/stack/tests/integration/internal/modules"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Search, modules.Ledger}, func() {
	BeforeEach(func() {
		createLedgerResponse, err := Client().Ledger.V2CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
			Ledger: "default",
		})
		Expect(err).To(BeNil())
		Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))
	})
	When("creating a transaction on a ledger", func() {
		var (
			timestamp = time.Now().Round(time.Second).UTC()
		)
		BeforeEach(func() {
			// Create a transaction
			response, err := Client().Ledger.V2CreateTransaction(
				TestContext(),
				operations.V2CreateTransactionRequest{
					IdempotencyKey: pointer.For("foo"),
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
						Timestamp: &timestamp,
						Reference: pointer.For("foo"),
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

		})
		Then("trying to commit a new transaction with same ik and different ledger", func() {
			BeforeEach(func() {
				// Create a transaction
				response, err := Client().Ledger.V2CreateTransaction(
					TestContext(),
					operations.V2CreateTransactionRequest{
						IdempotencyKey: pointer.For("foo"),
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
							Timestamp: &timestamp,
							Reference: pointer.For("foo"),
						},
						Ledger: "test",
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))
			})
		})
	})
})
