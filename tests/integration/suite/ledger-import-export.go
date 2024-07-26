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
	"io"
	"math/big"
	"net/http"
)

var _ = WithModules([]*Module{modules.Ledger}, func() {
	const logsToWrite = 5
	When("creating a ledger with a bunch of transactions", func() {
		BeforeEach(func() {
			response, err := Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
				Ledger: "default",
			})
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(http.StatusNoContent))

			for i := 0; i < logsToWrite; i++ {
				_, err := Client().Ledger.V2.CreateTransaction(TestContext(), operations.V2CreateTransactionRequest{
					V2PostTransaction: shared.V2PostTransaction{
						Postings: []shared.V2Posting{{
							Amount:      big.NewInt(100),
							Asset:       "USD/2",
							Destination: "bank",
							Source:      "world",
						}},
						Metadata: map[string]string{},
					},
					Ledger: "default",
				})
				Expect(err).To(BeNil())
			}
		})
		Then("exporting the ledger", func() {
			var (
				data []byte
			)
			BeforeEach(func() {
				ret, err := Client().Ledger.V2.ExportLogs(TestContext(), operations.V2ExportLogsRequest{
					Ledger: "default",
				})
				Expect(err).To(BeNil())
				data, err = io.ReadAll(ret.RawResponse.Body)
				Expect(err).To(BeNil())
			})
			It("Should be ok", func() {
				Expect(len(data)).NotTo(Equal(0))
			})
			Then("Import the ledger under a new name", func() {
				var (
					newLedger string
				)
				BeforeEach(func() {
					newLedger = uuid.NewString()
					_, err := Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
						Ledger:                newLedger,
						V2CreateLedgerRequest: &shared.V2CreateLedgerRequest{},
					})
					Expect(err).To(BeNil())

					_, err = Client().Ledger.V2.ImportLogs(TestContext(), operations.V2ImportLogsRequest{
						RequestBody: pointer.For(string(data)),
						Ledger:      newLedger,
					})
					Expect(err).To(BeNil())
				})
				It("Should be ok", func() {
					logs, err := Client().Ledger.V2.ListLogs(TestContext(), operations.V2ListLogsRequest{
						Ledger:   newLedger,
						PageSize: pointer.For(int64(200)),
					})
					Expect(err).To(BeNil())
					Expect(logs.V2LogsCursorResponse.Cursor.Data).To(HaveLen(logsToWrite))
				})
				Then("retrying to import the same export", func() {
					var (
						err error
					)
					BeforeEach(func() {
						_, err = Client().Ledger.V2.ImportLogs(TestContext(), operations.V2ImportLogsRequest{
							RequestBody: pointer.For(string(data)),
							Ledger:      newLedger,
						})
					})
					It("Should trigger an error with IMPORT error code", func() {
						Expect(err).NotTo(BeNil())
					})
				})
			})
		})
	})
})
