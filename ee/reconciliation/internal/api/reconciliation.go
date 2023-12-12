package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/reconciliation/internal/api/backend"
	"github.com/formancehq/reconciliation/internal/api/service"
	"github.com/formancehq/reconciliation/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/go-chi/chi/v5"
)

type reconciliationResponse struct {
	ID                   string              `json:"id"`
	PolicyID             string              `json:"policyID"`
	CreatedAt            time.Time           `json:"createdAt"`
	ReconciledAtLedger   time.Time           `json:"reconciledAtLedger"`
	ReconciledAtPayments time.Time           `json:"reconciledAtPayments"`
	Status               string              `json:"status"`
	PaymentsBalances     map[string]*big.Int `json:"paymentsBalances"`
	LedgerBalances       map[string]*big.Int `json:"ledgerBalances"`
	Error                string              `json:"error"`
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

		policyID := chi.URLParam(r, "policyID")
		if policyID == "" {
			api.BadRequest(w, ErrValidation, errors.New("missing policyID"))
			return
		}

		res, err := b.GetService().Reconciliation(r.Context(), policyID, &req)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		data := &reconciliationResponse{
			ID:                   res.ID.String(),
			PolicyID:             policyID,
			CreatedAt:            res.CreatedAt,
			ReconciledAtLedger:   res.ReconciledAtLedger,
			ReconciledAtPayments: res.ReconciledAtPayments,
			Status:               res.Status.String(),
			PaymentsBalances:     res.PaymentsBalances,
			LedgerBalances:       res.LedgerBalances,
			Error:                res.Error,
		}

		api.Ok(w, data)
	}
}

func getReconciliationHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "reconciliationID")

		res, err := b.GetService().GetReconciliation(r.Context(), id)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		data := &reconciliationResponse{
			ID:                   res.ID.String(),
			PolicyID:             res.PolicyID.String(),
			CreatedAt:            res.CreatedAt,
			ReconciledAtLedger:   res.ReconciledAtLedger,
			ReconciledAtPayments: res.ReconciledAtPayments,
			Status:               res.Status.String(),
			PaymentsBalances:     res.PaymentsBalances,
			LedgerBalances:       res.LedgerBalances,
			Error:                res.Error,
		}

		api.Ok(w, data)
	}
}

func listReconciliationsHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := storage.GetReconciliationsQuery{}

		if r.URL.Query().Get(QueryKeyCursor) != "" {
			err := bunpaginate.UnmarshalCursor(r.URL.Query().Get(QueryKeyCursor), &q)
			if err != nil {
				api.BadRequest(w, ErrValidation, fmt.Errorf("invalid '%s' query param", QueryKeyCursor))
				return
			}
		} else {
			options, err := getPaginatedQueryOptionsReconciliations(r)
			if err != nil {
				api.BadRequest(w, ErrValidation, err)
				return
			}
			q = storage.NewGetReconciliationsQuery(*options)
		}

		cursor, err := b.GetService().ListReconciliations(r.Context(), q)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.RenderCursor(w, *cursor)
	}
}
