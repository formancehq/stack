package suite

import (
	"fmt"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/stack/tests/integration/internal/modules"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Ledger}, func() {
	BeforeEach(func() {
		createLedgerResponse, err := Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
			Ledger: "default",
		})
		Expect(err).To(BeNil())
		Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))
	})
	const (
		pageSize = int64(10)
		txCount  = 2 * pageSize
	)
	When(fmt.Sprintf("creating %d transactions", txCount), func() {
		var (
			timestamp = time.Now().Round(time.Second).UTC()
		)
		BeforeEach(func() {
			for i := 0; i < int(txCount); i++ {
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
									Destination: fmt.Sprintf("account:%d", i),
								},
							},
							Timestamp: &timestamp,
						},
						Ledger: "default",
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))
			}
		})
		Then("Listing balances using v1 endpoint", func() {
			var (
				rsp *operations.GetBalancesResponse
				err error
			)
			BeforeEach(func() {
				rsp, err = Client().Ledger.V1.GetBalances(
					TestContext(),
					operations.GetBalancesRequest{
						Ledger:  "default",
						Address: pointer.For("world"),
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(rsp.StatusCode).To(Equal(200))
			})
			It("Should be return non empty balances", func() {
				Expect(rsp.BalancesCursorResponse.Cursor.Data).To(HaveLen(1))
				balances := rsp.BalancesCursorResponse.Cursor.Data[0]
				Expect(balances).To(HaveKey("world"))
				Expect(balances["world"]).To(HaveKey("USD"))
				Expect(balances["world"]["USD"]).To(Equal(int64(-2000)))
			})
		})
	})
})
