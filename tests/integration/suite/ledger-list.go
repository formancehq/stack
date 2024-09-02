package suite

import (
	"fmt"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = WithModules([]*Module{modules.Ledger}, func() {
	const count = 3
	When(fmt.Sprintf("creating %d ledger", count), func() {
		BeforeEach(func() {
			for i := 0; i < count; i++ {
				createLedgerResponse, err := Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
					Ledger: fmt.Sprintf("ledger%d", i),
				})
				Expect(err).To(BeNil())
				Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))
			}
		})
		It("should be listable on api", func() {
			response, err := Client().Ledger.V2.ListLedgers(TestContext(), operations.V2ListLedgersRequest{})
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(response.V2LedgerListResponse.Cursor.Data).To(HaveLen(count))
		})
	})
})
