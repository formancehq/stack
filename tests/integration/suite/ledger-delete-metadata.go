package suite

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = WithModules([]*Module{modules.Ledger}, func() {
	When("creating a ledger with some metadata", func() {
		BeforeEach(func() {
			response, err := Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
				Ledger: "default",
				V2CreateLedgerRequest: &shared.V2CreateLedgerRequest{
					Metadata: map[string]string{
						"foo": "bar",
					},
				},
			})
			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(http.StatusNoContent))
		})
		Then("deleting metadata", func() {
			BeforeEach(func() {
				_, err := Client().Ledger.V2.DeleteLedgerMetadata(TestContext(), operations.V2DeleteLedgerMetadataRequest{
					Ledger: "default",
					Key:    "foo",
				})
				Expect(err).To(BeNil())
			})
			It("should be ok", func() {
				ledger, err := Client().Ledger.V2.GetLedger(TestContext(), operations.V2GetLedgerRequest{
					Ledger: "default",
				})
				Expect(err).To(BeNil())
				Expect(ledger.V2GetLedgerResponse.Data.Metadata).To(BeEmpty())
			})
		})
	})
})
