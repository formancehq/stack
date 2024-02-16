package api

import (
	"encoding/json"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/formancehq/payments/cmd/api/internal/api/backend"
	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/gorilla/mux"
)

type balancesResponse struct {
	AccountID     string    `json:"accountId"`
	CreatedAt     time.Time `json:"createdAt"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
	Currency      string    `json:"currency"` // Deprecated: should be removed soon
	Asset         string    `json:"asset"`
	Balance       *big.Int  `json:"balance"`
}

func listBalancesForAccount(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		balanceQuery, err := populateBalanceQueryFromRequest(r)
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		query, err := bunpaginate.Extract[storage.ListBalancesQuery](r, func() (*storage.ListBalancesQuery, error) {
			options, err := getPagination(r, balanceQuery)
			if err != nil {
				return nil, err
			}
			return pointer.For(storage.NewListBalancesQuery(*options)), nil
		})
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		// In order to support the legacy API, we need to check if the limit query parameter is set
		// and if so, we need to override the pageSize pagination option
		if r.URL.Query().Get("limit") != "" {
			limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
			if err != nil {
				api.BadRequest(w, ErrValidation, err)
				return
			}

			if limit > 0 {
				query.PageSize = uint64(limit)
				query.Options.PageSize = uint64(limit)
			}
		}

		cursor, err := b.GetService().ListBalances(r.Context(), *query)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		ret := cursor.Data
		data := make([]*balancesResponse, len(ret))

		for i := range ret {
			data[i] = &balancesResponse{
				AccountID:     ret[i].AccountID.String(),
				CreatedAt:     ret[i].CreatedAt,
				Currency:      ret[i].Asset.String(),
				Asset:         ret[i].Asset.String(),
				Balance:       ret[i].Balance,
				LastUpdatedAt: ret[i].LastUpdatedAt,
			}
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[*balancesResponse]{
			Cursor: &api.Cursor[*balancesResponse]{
				PageSize: cursor.PageSize,
				HasMore:  cursor.HasMore,
				Previous: cursor.Previous,
				Next:     cursor.Next,
				Data:     data,
			},
		})
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
	}
}

func populateBalanceQueryFromRequest(r *http.Request) (storage.BalanceQuery, error) {
	var balanceQuery storage.BalanceQuery

	balanceQuery = balanceQuery.WithCurrency(r.URL.Query().Get("asset"))

	accountID, err := models.AccountIDFromString(mux.Vars(r)["accountID"])
	if err != nil {
		return balanceQuery, err
	}
	balanceQuery = balanceQuery.WithAccountID(accountID)

	var startTimeParsed, endTimeParsed time.Time

	from, to := r.URL.Query().Get("from"), r.URL.Query().Get("to")
	if from != "" {
		startTimeParsed, err = time.Parse(time.RFC3339Nano, from)
		if err != nil {
			return balanceQuery, err
		}
	}
	if to != "" {
		endTimeParsed, err = time.Parse(time.RFC3339Nano, to)
		if err != nil {
			return balanceQuery, err
		}
	}

	switch {
	case startTimeParsed.IsZero() && endTimeParsed.IsZero():
		balanceQuery = balanceQuery.
			WithTo(time.Now())
	case !startTimeParsed.IsZero() && endTimeParsed.IsZero():
		balanceQuery = balanceQuery.
			WithFrom(startTimeParsed).
			WithTo(time.Now())
	case startTimeParsed.IsZero() && !endTimeParsed.IsZero():
		balanceQuery = balanceQuery.
			WithTo(endTimeParsed)
	default:
		balanceQuery = balanceQuery.
			WithFrom(startTimeParsed).
			WithTo(endTimeParsed)
	}

	return balanceQuery, nil
}
