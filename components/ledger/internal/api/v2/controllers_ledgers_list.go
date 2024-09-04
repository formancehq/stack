package v2

import (
	"net/http"

	"github.com/formancehq/go-libs/api"
	"github.com/formancehq/go-libs/bun/bunpaginate"
	"github.com/formancehq/go-libs/pointer"
	"github.com/formancehq/ledger/internal/api/common"
	"github.com/formancehq/ledger/internal/controller/ledger"
	"github.com/formancehq/ledger/internal/controller/system"
	"github.com/pkg/errors"
)

func listLedgers(b system.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		query, err := bunpaginate.Extract[ledger.ListLedgersQuery](r, func() (*ledger.ListLedgersQuery, error) {
			pageSize, err := bunpaginate.GetPageSize(r)
			if err != nil {
				return nil, err
			}

			return pointer.For(ledger.NewListLedgersQuery(pageSize)), nil
		})
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		ledgers, err := b.ListLedgers(r.Context(), *query)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrInvalidQuery{}) || errors.Is(err, common.ErrMissingFeature{}):
				api.BadRequest(w, ErrValidation, err)
			default:
				api.InternalServerError(w, r, err)
			}
			return
		}

		api.RenderCursor(w, *ledgers)
	}
}
