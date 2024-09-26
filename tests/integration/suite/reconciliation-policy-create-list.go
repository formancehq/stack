package suite

import (
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Auth, modules.Ledger, modules.Payments, modules.Reconciliation}, func() {
	When("1 - reconciliation list policies", func() {
		var (
			policies []shared.Policy
		)
		JustBeforeEach(func() {
			policiesResponse, err := Client().Reconciliation.V1.ListPolicies(
				TestContext(),
				operations.ListPoliciesRequest{},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(policiesResponse.StatusCode).To(Equal(200))

			policies = policiesResponse.PoliciesCursorResponse.Cursor.Data
		})
		It("should respond with empty lists", func() {
			Expect(policies).To(BeEmpty())
		})
	})
	When("2 - reconciliation create 2 policies", func() {
		JustBeforeEach(func() {
			response, err := Client().Reconciliation.V1.CreatePolicy(
				TestContext(),
				shared.PolicyRequest{
					LedgerName:     "default",
					LedgerQuery:    map[string]interface{}{},
					Name:           "test 1",
					PaymentsPoolID: uuid.New().String(),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(201))

			response, err = Client().Reconciliation.V1.CreatePolicy(
				TestContext(),
				shared.PolicyRequest{
					LedgerName:     "default",
					LedgerQuery:    map[string]interface{}{},
					Name:           "test 2",
					PaymentsPoolID: uuid.New().String(),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(201))
		})
		Then("should list 2 policies", func() {
			var (
				policies []shared.Policy
			)
			JustBeforeEach(func() {
				policiesResponse, err := Client().Reconciliation.V1.ListPolicies(
					TestContext(),
					operations.ListPoliciesRequest{},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(policiesResponse.StatusCode).To(Equal(200))

				policies = policiesResponse.PoliciesCursorResponse.Cursor.Data
			})
			It("should return 2 items", func() {
				Expect(policies).To(HaveLen(2))
			})
		})
	})
})
