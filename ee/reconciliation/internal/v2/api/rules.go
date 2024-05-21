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

func createRuleHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req service.CreateRuleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		rule, err := b.GetService().CreateRule(r.Context(), &req)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.Created(w, rule)
	}
}

func deleteRuleHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "ruleID")

		err := b.GetService().DeleteRule(r.Context(), id)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.NoContent(w)
	}
}

func getRuleHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "ruleID")

		rule, err := b.GetService().GetRule(r.Context(), id)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.Ok(w, rule)
	}
}

func listRulesHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := storage.ListRulesQuery{}

		if r.URL.Query().Get(QueryKeyCursor) != "" {
			err := bunpaginate.UnmarshalCursor(r.URL.Query().Get(QueryKeyCursor), &q)
			if err != nil {
				api.BadRequest(w, ErrValidation, fmt.Errorf("invalid '%s' query param", QueryKeyCursor))
				return
			}
		} else {
			options, err := getPaginatedQueryOptionsRules(r)
			if err != nil {
				api.BadRequest(w, ErrValidation, err)
				return
			}
			q = storage.NewListRulesQuery(*options)
		}

		cursor, err := b.GetService().ListRules(r.Context(), q)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		api.RenderCursor(w, *cursor)
	}
}
