package v2

import (
	"fmt"
	"net/http"

	"github.com/formancehq/go-libs/api"
	"github.com/formancehq/ledger/internal/api/common"
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	"github.com/pkg/errors"
)

func countAccounts(w http.ResponseWriter, r *http.Request) {
	l := common.LedgerFromContext(r.Context())

	options, err := getPaginatedQueryOptionsOfPITFilterWithVolumes(r)
	if err != nil {
		api.BadRequest(w, ErrValidation, err)
		return
	}

	count, err := l.CountAccounts(r.Context(), ledgercontroller.NewListAccountsQuery(*options))
	if err != nil {
		switch {
		case errors.Is(err, common.ErrInvalidQuery{}) || errors.Is(err, common.ErrMissingFeature{}):
			api.BadRequest(w, ErrValidation, err)
		default:
			api.InternalServerError(w, r, err)
		}
		return
	}

	w.Header().Set("Count", fmt.Sprint(count))
	api.NoContent(w)
}
