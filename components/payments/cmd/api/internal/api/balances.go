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

		pagination, err := getPagination(r)
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		accountID, err := models.AccountIDFromString(mux.Vars(r)["accountID"])
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		balanceQuery := storage.NewBalanceQuery(pagination).
			WithAccountID(accountID)

		balanceQuery, err = populateBalanceQueryFromRequest(r, balanceQuery)
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		ret, paginationDetails, err := b.GetService().ListBalances(r.Context(), balanceQuery)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

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
				PageSize: paginationDetails.PageSize,
				HasMore:  paginationDetails.HasMore,
				Previous: paginationDetails.PreviousPage,
				Next:     paginationDetails.NextPage,
				Data:     data,
			},
		})
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
	}
}

func populateBalanceQueryFromRequest(r *http.Request, balanceQuery storage.BalanceQuery) (storage.BalanceQuery, error) {
	balanceQuery = balanceQuery.WithCurrency(r.URL.Query().Get("asset"))

	var startTimeParsed, endTimeParsed time.Time
	var err error

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
	if r.URL.Query().Get("limit") != "" {
		limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
		if err != nil {
			return balanceQuery, err
		}

		if limit > 0 {
			balanceQuery = balanceQuery.WithLimit(int(limit))
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
