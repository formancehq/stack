package suite

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/sdkerrors"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = WithModules([]*Module{modules.Ledger}, func() {
	When("creating a bulk on a ledger", func() {
		BeforeEach(func() {
			response, err := Client().Ledger.V2CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
				Ledger: "default",
			})
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(http.StatusNoContent))
		})
		It("Should be ok", func() {})
		Then("trying to create another ledger with the same name", func() {
			BeforeEach(func() {
				_, err := Client().Ledger.V2CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
					Ledger: "default",
				})
				Expect(err).NotTo(BeNil())
				Expect(err.(*sdkerrors.V2ErrorResponse).ErrorCode).To(Equal(sdkerrors.V2ErrorsEnumValidation))
			})
			It("should fail", func() {})
		})
	})
})
