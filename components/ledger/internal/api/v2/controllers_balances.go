package v2

import (
	"net/http"

	"github.com/formancehq/ledger/internal/controller/ledger"
	"github.com/pkg/errors"

	"github.com/formancehq/go-libs/api"
	"github.com/formancehq/ledger/internal/api/common"
)

func readBalancesAggregated(w http.ResponseWriter, r *http.Request) {

	pitFilter, err := getPITFilter(r)
	if err != nil {
		api.BadRequest(w, ErrValidation, err)
		return
	}

	queryBuilder, err := getQueryBuilder(r)
	if err != nil {
		api.BadRequest(w, ErrValidation, err)
		return
	}

	balances, err := common.LedgerFromContext(r.Context()).
		GetAggregatedBalances(r.Context(), ledger.NewGetAggregatedBalancesQuery(
			*pitFilter, queryBuilder, api.QueryParamBool(r, "use_insertion_date") || api.QueryParamBool(r, "useInsertionDate")))
	if err != nil {
		switch {
		case errors.Is(err, common.ErrInvalidQuery{}) || errors.Is(err, common.ErrMissingFeature{}):
			api.BadRequest(w, ErrValidation, err)
		default:
			api.InternalServerError(w, r, err)
		}
		return
	}

	api.Ok(w, balances)
}
