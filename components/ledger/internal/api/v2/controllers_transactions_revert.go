package v2

import (
	"net/http"
	"strconv"

	"github.com/formancehq/go-libs/api"
	"github.com/formancehq/ledger/internal/api/common"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

func revertTransaction(w http.ResponseWriter, r *http.Request) {
	l := common.LedgerFromContext(r.Context())

	txId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		api.BadRequest(w, ErrValidation, err)
		return
	}

	tx, err := l.RevertTransaction(r.Context(), getCommandParameters(r), int(txId),
		api.QueryParamBool(r, "force"),
		api.QueryParamBool(r, "atEffectiveDate"),
	)
	if err != nil {
		switch {
		case errors.Is(err, &common.ErrInsufficientFunds{}):
			api.BadRequest(w, ErrInsufficientFund, err)
		case errors.Is(err, common.ErrAlreadyReverted{}):
			api.BadRequest(w, ErrAlreadyRevert, err)
		case errors.Is(err, common.ErrNotFound):
			api.NotFound(w, err)
		default:
			api.InternalServerError(w, r, err)
		}
		return
	}

	api.Created(w, tx)
}
