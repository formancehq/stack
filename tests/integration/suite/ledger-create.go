package suite

import (
	"net/http"
	"strings"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/sdkerrors"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Ledger}, func() {
	When("creating a bulk on a ledger", func() {
		var m map[string]string
		JustBeforeEach(func() {
			response, err := Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
				Ledger: "default",
				V2CreateLedgerRequest: &shared.V2CreateLedgerRequest{
					Metadata: m,
				},
			})
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(http.StatusNoContent))
		})
		It("Should be ok", func() {})
		Then("trying to create another ledger with the same name", func() {
			JustBeforeEach(func() {
				_, err := Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
					Ledger: "default",
				})
				Expect(err).NotTo(BeNil())
				Expect(err.(*sdkerrors.V2ErrorResponse).ErrorCode).To(Equal(shared.V2ErrorsEnumValidation))
			})
			It("should fail", func() {})
		})
		Context("With metadata", func() {
			BeforeEach(func() {
				m = map[string]string{
					"foo": "bar",
				}
			})
			It("Should be ok", func() {
				ledger, err := Client().Ledger.V2.GetLedger(TestContext(), operations.V2GetLedgerRequest{
					Ledger: "default",
				})
				Expect(err).To(BeNil())
				Expect(ledger.V2GetLedgerResponse.Data.Metadata).To(Equal(m))
			})
		})
	})
	When("bucket naming convention depends on the database 63 bytes length (pg constraint)", func() {
		It("should fail with > 63 characters in ledger or bucket name", func() {
			_, err := Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
				V2CreateLedgerRequest: &shared.V2CreateLedgerRequest{
					Bucket: pointer.For(strings.Repeat("a", 64)),
				},
				Ledger: "default",
			})
			Expect(err).To(HaveOccurred())
		})
	})
})
