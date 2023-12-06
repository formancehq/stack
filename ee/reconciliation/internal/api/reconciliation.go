package api

import (
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/formancehq/reconciliation/internal/api/backend"
	"github.com/formancehq/reconciliation/internal/api/service"
	"github.com/formancehq/stack/libs/go-libs/api"
)

type reconciliationResponse struct {
	Status          string              `json:"status"`
	PaymentBalances map[string]*big.Int `json:"paymentBalances"`
	LedgerBalances  map[string]*big.Int `json:"ledgerBalances"`
	Error           string              `json:"error"`
}

func reconciliationHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req service.ReconciliationRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		if err := req.Validate(); err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		res, err := b.GetService().Reconciliation(r.Context(), &req)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		data := &reconciliationResponse{
			Status:          res.Status.String(),
			PaymentBalances: res.PaymentBalances,
			LedgerBalances:  res.LedgerBalances,
			Error:           res.Error,
		}

		api.Ok(w, data)
	}
}
