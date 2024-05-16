package v1

import (
	"github.com/formancehq/reconciliation/internal/api/v1/backend"
	"github.com/go-chi/chi/v5"
)

func NewRouter(
	b backend.Backend,
	r chi.Router,
) {
	r.Get("/reconciliations/{reconciliationID}", getReconciliationHandler(b))
	r.Get("/reconciliations", listReconciliationsHandler(b))

	r.Post("/policies", createPolicyHandler(b))
	r.Get("/policies", listPoliciesHandler(b))
	r.Delete("/policies/{policyID}", deletePolicyHandler(b))
	r.Get("/policies/{policyID}", getPolicyHandler(b))
	r.Post("/policies/{policyID}/reconciliation", reconciliationHandler(b))
}
