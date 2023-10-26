package api

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type balancesRepository interface {
	ListBalances(ctx context.Context, query storage.BalanceQuery) ([]*models.Balance, storage.PaginationDetails, error)
}

type balancesResponse struct {
	AccountID     string    `json:"accountId"`
	CreatedAt     time.Time `json:"createdAt"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
	Currency      string    `json:"currency"` // Deprecated: should be removed soon
	Asset         string    `json:"asset"`
	Balance       *big.Int  `json:"balance"`
}

func listBalancesForAccount(repo balancesRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var sorter storage.Sorter

		if sortParams := r.URL.Query()["sort"]; sortParams != nil {
			for _, s := range sortParams {
				parts := strings.SplitN(s, ":", 2)

				var order storage.SortOrder

				if len(parts) > 1 {
					//nolint:goconst // allow duplicate string
					switch parts[1] {
					case "asc", "ASC":
						order = storage.SortOrderAsc
					case "dsc", "desc", "DSC", "DESC":
						order = storage.SortOrderDesc
					default:
						api.BadRequest(w, ErrValidation, errors.New("sort order not well specified, got "+parts[1]))
						return
					}
				}

				column := parts[0]

				sorter.Add(column, order)
			}
		}

		pageSize, err := pageSizeQueryParam(r)
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		pagination, err := storage.Paginate(pageSize, r.URL.Query().Get("cursor"), sorter, nil)
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
			WithAccountID(accountID).
			WithCurrency(r.URL.Query().Get("asset"))

		var startTimeParsed, endTimeParsed time.Time
		from, to := r.URL.Query().Get("from"), r.URL.Query().Get("to")
		if from != "" {
			startTimeParsed, err = time.Parse(time.RFC3339Nano, from)
			if err != nil {
				api.BadRequest(w, ErrValidation, err)
				return
			}
		}
		if to != "" {
			endTimeParsed, err = time.Parse(time.RFC3339Nano, to)
			if err != nil {
				api.BadRequest(w, ErrValidation, err)
				return
			}
		}
		if r.URL.Query().Get("limit") != "" {
			limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
			if err != nil {
				api.BadRequest(w, ErrValidation, err)
				return
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

		ret, paginationDetails, err := repo.ListBalances(r.Context(), balanceQuery)
		if err != nil {
			api.InternalServerError(w, r, err)
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
