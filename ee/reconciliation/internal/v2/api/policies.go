package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/formancehq/reconciliation/internal/v2/api/backend"
	"github.com/formancehq/reconciliation/internal/v2/api/service"
	"github.com/formancehq/reconciliation/internal/v2/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/go-chi/chi/v5"
)

func createPolicyHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req service.CreatePolicyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		policy, err := b.GetService().CreatePolicy(r.Context(), &req)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.Created(w, policy)
	}
}

func udpatePolicyRulesHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req service.UpdatePolicyRulesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		id := chi.URLParam(r, "policyID")

		err := b.GetService().UpdatePolicyRules(r.Context(), id, &req)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.NoContent(w)
	}
}

func enablePolicyHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "policyID")

		err := b.GetService().EnablePolicy(r.Context(), id)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.NoContent(w)
	}
}

func disablePolicyHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "policyID")

		err := b.GetService().DisablePolicy(r.Context(), id)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.NoContent(w)
	}
}

func deletePolicyHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "policyID")

		err := b.GetService().DeletePolicy(r.Context(), id)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.NoContent(w)
	}
}

func getPolicyHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "policyID")

		policy, err := b.GetService().GetPolicy(r.Context(), id)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.Ok(w, policy)
	}
}

func listPoliciesHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := storage.ListPoliciesQuery{}

		if r.URL.Query().Get(QueryKeyCursor) != "" {
			err := bunpaginate.UnmarshalCursor(r.URL.Query().Get(QueryKeyCursor), &q)
			if err != nil {
				api.BadRequest(w, ErrValidation, fmt.Errorf("invalid '%s' query param", QueryKeyCursor))
				return
			}
		} else {
			options, err := getPaginatedQueryOptionsPolicies(r)
			if err != nil {
				api.BadRequest(w, ErrValidation, err)
				return
			}
			q = storage.NewListPoliciesQuery(*options)
		}

		cursor, err := b.GetService().ListPolicies(r.Context(), q)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.RenderCursor(w, *cursor)
	}
}
