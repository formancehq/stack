package v2

import (
	"net/http"

	"github.com/formancehq/ledger/internal/api/backend"
	"github.com/formancehq/ledger/internal/storage/ledgerstore"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"

	"github.com/formancehq/stack/libs/go-libs/pointer"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
)

func getVolumesWithBalances(w http.ResponseWriter, r *http.Request) {

	l := backend.LedgerFromContext(r.Context())

	query, err := bunpaginate.Extract[ledgerstore.GetVolumesWithBalancesQuery](r, func() (*ledgerstore.GetVolumesWithBalancesQuery, error) {
		options, err := getPaginatedQueryOptionsOfPITOOTFilterForVolumes(r)
		if err != nil {
			sharedapi.BadRequest(w, ErrValidation, err)
		}
		return pointer.For(ledgerstore.NewGetVolumesWithBalancesQuery(*options)), nil

	})

	if err != nil {
		sharedapi.BadRequest(w, ErrValidation, err)
		return
	}

	cursor, err := l.GetVolumesWithBalances(r.Context(), *query)

	if err != nil {
		switch {
		case ledgerstore.IsErrInvalidQuery(err):
			sharedapi.BadRequest(w, ErrValidation, err)
		default:
			sharedapi.InternalServerError(w, r, err)
		}
		return
	}

	sharedapi.RenderCursor(w, *cursor)

}
