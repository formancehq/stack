package v2

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/go-libs/api"
	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/api/common"
	"github.com/pkg/errors"
)

func createTransaction(w http.ResponseWriter, r *http.Request) {
	l := common.LedgerFromContext(r.Context())

	payload := ledger.TransactionRequest{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		api.BadRequest(w, ErrValidation, errors.New("invalid transaction format"))
		return
	}

	if len(payload.Postings) > 0 && payload.Script.Plain != "" {
		api.BadRequest(w, ErrValidation, errors.New("cannot pass postings and numscript in the same request"))
		return
	}

	if len(payload.Postings) == 0 && payload.Script.Plain == "" {
		api.BadRequest(w, ErrNoPostings, errors.New("you need to pass either a posting array or a numscript script"))
		return
	}

	res, err := l.CreateTransaction(r.Context(), getCommandParameters(r), *payload.ToRunScript(api.QueryParamBool(r, "force")))
	if err != nil {
		switch {
		case errors.Is(err, &common.ErrInsufficientFunds{}):
			api.BadRequest(w, ErrInsufficientFund, err)
		case errors.Is(err, &common.ErrInvalidVars{}) || errors.Is(err, &common.ErrCompileErrorList{}):
			api.BadRequest(w, ErrCompilationFailed, err)
		case errors.Is(err, &common.ErrMetadataOverride{}):
			api.BadRequest(w, ErrMetadataOverride, err)
		case errors.Is(err, common.ErrNoPostings):
			api.BadRequest(w, ErrNoPostings, err)
		case errors.Is(err, common.ErrConstraintsFailed{}):
			api.WriteErrorResponse(w, http.StatusConflict, ErrConflict, err)
		default:
			api.InternalServerError(w, r, err)
		}
		return
	}

	api.Ok(w, res)
}
