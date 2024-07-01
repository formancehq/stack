package v2

import (
	"encoding/json"
	"io"
	"net/http"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/api/backend"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/pkg/errors"
)

func importLogs(w http.ResponseWriter, r *http.Request) {

	stream := make(chan *ledger.ChainedLog)
	errChan := make(chan error, 1)
	go func() {
		errChan <- backend.LedgerFromContext(r.Context()).Import(r.Context(), stream)
	}()
	dec := json.NewDecoder(r.Body)
	for {
		l := &ledger.ChainedLog{}
		if err := dec.Decode(l); err != nil {
			if errors.Is(err, io.EOF) {
				close(stream)
				break
			}
		}
		select {
		case stream <- l:
		case <-r.Context().Done():
			api.InternalServerError(w, r, r.Context().Err())
			return
		case err := <-errChan:
			api.InternalServerError(w, r, err)
			return
		}
	}
	select {
	case err := <-errChan:
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
	case <-r.Context().Done():
		api.InternalServerError(w, r, r.Context().Err())
		return
	}

	api.NoContent(w)
}
