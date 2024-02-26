package api

import (
	"encoding/json"
	"math/big"
	"net/http"
	"strings"

	"github.com/formancehq/payments/cmd/api/internal/api/backend"
	"github.com/formancehq/payments/cmd/api/internal/api/service"
	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
)

type poolResponse struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Accounts []string `json:"accounts"`
}

func createPoolHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "createPoolHandler")
		defer span.End()

		w.Header().Set("Content-Type", "application/json")

		var createPoolRequest service.CreatePoolRequest
		err := json.NewDecoder(r.Body).Decode(&createPoolRequest)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		span.SetAttributes(
			attribute.String("request.name", createPoolRequest.Name),
			attribute.String("request.accounts", strings.Join(createPoolRequest.AccountIDs, ",")),
		)

		if err := createPoolRequest.Validate(); err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		pool, err := b.GetService().CreatePool(ctx, &createPoolRequest)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		accounts := make([]string, len(pool.PoolAccounts))
		for i := range pool.PoolAccounts {
			accounts[i] = pool.PoolAccounts[i].AccountID.String()
		}

		data := &poolResponse{
			ID:       pool.ID.String(),
			Name:     pool.Name,
			Accounts: accounts,
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[poolResponse]{
			Data: data,
		})
		if err != nil {
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
		}
	}
}

func addAccountToPoolHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "addAccountToPoolHandler")
		defer span.End()

		poolID, ok := mux.Vars(r)["poolID"]
		if !ok {
			var err = errors.New("missing poolID")
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		span.SetAttributes(attribute.String("request.poolID", poolID))

		var addAccountToPoolRequest service.AddAccountToPoolRequest
		err := json.NewDecoder(r.Body).Decode(&addAccountToPoolRequest)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		span.SetAttributes(
			attribute.String("request.accountID", addAccountToPoolRequest.AccountID),
		)

		if err := addAccountToPoolRequest.Validate(); err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		err = b.GetService().AddAccountToPool(ctx, poolID, &addAccountToPoolRequest)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func removeAccountFromPoolHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "removeAccountFromPoolHandler")
		defer span.End()

		poolID, ok := mux.Vars(r)["poolID"]
		if !ok {
			var err = errors.New("missing poolID")
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		span.SetAttributes(attribute.String("request.poolID", poolID))

		accountID, ok := mux.Vars(r)["accountID"]
		if !ok {
			var err = errors.New("missing accountID")
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		span.SetAttributes(attribute.String("request.accountID", accountID))

		err := b.GetService().RemoveAccountFromPool(ctx, poolID, accountID)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func listPoolHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "listPoolHandler")
		defer span.End()

		w.Header().Set("Content-Type", "application/json")

		query, err := bunpaginate.Extract[storage.ListPoolsQuery](r, func() (*storage.ListPoolsQuery, error) {
			options, err := getPagination(r, storage.PoolQuery{})
			if err != nil {
				return nil, err
			}
			return pointer.For(storage.NewListPoolsQuery(*options)), nil
		})
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		cursor, err := b.GetService().ListPools(ctx, *query)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		ret := cursor.Data
		data := make([]*poolResponse, len(ret))

		for i := range ret {
			accounts := make([]string, len(ret[i].PoolAccounts))
			for j := range ret[i].PoolAccounts {
				accounts[j] = ret[i].PoolAccounts[j].AccountID.String()
			}

			data[i] = &poolResponse{
				ID:       ret[i].ID.String(),
				Name:     ret[i].Name,
				Accounts: accounts,
			}
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[*poolResponse]{
			Cursor: &api.Cursor[*poolResponse]{
				PageSize: cursor.PageSize,
				HasMore:  cursor.HasMore,
				Previous: cursor.Previous,
				Next:     cursor.Next,
				Data:     data,
			},
		})
		if err != nil {
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
		}
	}
}

func getPoolHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "getPoolHandler")
		defer span.End()

		poolID, ok := mux.Vars(r)["poolID"]
		if !ok {
			err := errors.New("missing poolID")
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		span.SetAttributes(attribute.String("request.poolID", poolID))

		pool, err := b.GetService().GetPool(ctx, poolID)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		accounts := make([]string, len(pool.PoolAccounts))
		for i := range pool.PoolAccounts {
			accounts[i] = pool.PoolAccounts[i].AccountID.String()
		}

		data := &poolResponse{
			ID:       pool.ID.String(),
			Name:     pool.Name,
			Accounts: accounts,
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[poolResponse]{
			Data: data,
		})
		if err != nil {
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
		}
	}
}

type poolBalancesResponse struct {
	Balances []*poolBalanceResponse `json:"balances"`
}

type poolBalanceResponse struct {
	Amount *big.Int `json:"amount"`
	Asset  string   `json:"asset"`
}

func getPoolBalances(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "getPoolBalances")
		defer span.End()

		poolID, ok := mux.Vars(r)["poolID"]
		if !ok {
			var err = errors.New("missing poolID")
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		span.SetAttributes(attribute.String("request.poolID", poolID))

		atTime := r.URL.Query().Get("at")
		if atTime == "" {
			var err = errors.New("missing atTime")
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		span.SetAttributes(attribute.String("request.atTime", atTime))

		balance, err := b.GetService().GetPoolBalance(ctx, poolID, atTime)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		data := &poolBalancesResponse{
			Balances: make([]*poolBalanceResponse, len(balance.Balances)),
		}

		for i := range balance.Balances {
			data.Balances[i] = &poolBalanceResponse{
				Amount: balance.Balances[i].Amount,
				Asset:  balance.Balances[i].Asset,
			}
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[poolBalancesResponse]{
			Data: data,
		})
		if err != nil {
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
		}
	}
}

func deletePoolHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "deletePoolHandler")
		defer span.End()

		poolID, ok := mux.Vars(r)["poolID"]
		if !ok {
			var err = errors.New("missing poolID")
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		span.SetAttributes(attribute.String("request.poolID", poolID))

		err := b.GetService().DeletePool(ctx, poolID)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
