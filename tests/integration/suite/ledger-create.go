package suite

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = WithModules([]*Module{modules.Ledger}, func() {
	When("creating a bulk on a ledger", func() {
		BeforeEach(func() {
			response, err := Client().Ledger.CreateLedger(TestContext(), operations.CreateLedgerRequest{
				Ledger: "default",
			})
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(http.StatusNoContent))
		})
		It("Should be ok", func() {})
		Then("trying to create another ledger with the same name", func() {
			var (
				response *operations.CreateLedgerResponse
			)
			BeforeEach(func() {
				var err error
				response, err = Client().Ledger.CreateLedger(TestContext(), operations.CreateLedgerRequest{
					Ledger: "default",
				})
				Expect(err).To(BeNil())
			})
			It("should fail", func() {
				Expect(response.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(response.ErrorResponse.ErrorCode).To(Equal(shared.ErrorsEnumValidation))
			})
		})
	})
})
