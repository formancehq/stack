package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/formancehq/reconciliation/internal/v2/api/backend"
	"github.com/formancehq/reconciliation/internal/v2/api/service"
	"github.com/formancehq/reconciliation/internal/v2/models"
	"github.com/formancehq/reconciliation/internal/v2/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/go-chi/chi/v5"
)

func createReconciliationHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req service.CreateReconciliationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		if err := req.Validate(); err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		reconciliation, err := b.GetService().CreateReconciliation(r.Context(), &req)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.Created(w, reconciliation)
	}
}

func getReconciliationHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "reconciliationID")

		reconciliation, err := b.GetService().GetReconciliation(r.Context(), id)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.Ok(w, reconciliation)
	}
}

func getReconciliationDetails(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "reconciliationID")

		reconciliation, err := b.GetService().GetReconciliation(r.Context(), id)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		switch reconciliation.PolicyType {
		case models.PolicyTypeAccountBased:
			details, err := b.GetService().GetAccountBasedReconciliationDetails(r.Context(), id)
			if err != nil {
				handleServiceErrors(w, r, err)
				return
			}

			api.Ok(w, details)
		case models.PolicyTypeTransactionBased:
			// TODO(polo): add transaction based details
		}
	}
}

func listReconciliationsHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := storage.ListReconciliationsQuery{}

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
			q = storage.NewListReconciliationsQuery(*options)
		}

		cursor, err := b.GetService().ListReconciliations(r.Context(), q)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.RenderCursor(w, *cursor)
	}
}
