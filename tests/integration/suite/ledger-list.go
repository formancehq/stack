package suite

import (
	"fmt"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
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
				createLedgerResponse, err := Client().Ledger.CreateLedger(TestContext(), operations.CreateLedgerRequest{
					Ledger: fmt.Sprintf("ledger%d", i),
				})
				Expect(err).To(BeNil())
				Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))
			}
		})
		It("should be listable on api", func() {
			response, err := Client().Ledger.ListLedgers(TestContext(), operations.ListLedgersRequest{})
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(response.LedgerListResponse.Cursor.Data).To(HaveLen(count))
		})
	})
})
