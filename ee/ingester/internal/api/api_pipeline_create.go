package api

import (
	"net/http"

	"github.com/pkg/errors"

	ingester "github.com/formancehq/stack/ee/ingester/internal"
	"github.com/formancehq/stack/libs/go-libs/api"
)

func (a *API) createPipeline(w http.ResponseWriter, r *http.Request) {
	withBody[ingester.PipelineConfiguration](w, r, func(req ingester.PipelineConfiguration) {
		p, err := a.backend.CreatePipeline(r.Context(), req)
		if err != nil {
			switch {
			case errors.Is(err, ErrModuleNotAvailable("")) ||
				errors.Is(err, ErrConnectorNotFound("")) ||
				errors.Is(err, ErrPipelineAlreadyExists{}) ||
				errors.Is(err, ErrInUsePipeline("")):
				api.BadRequest(w, "VALIDATION", err)
			default:
				api.InternalServerError(w, r, err)
			}
			return
		}

		api.Created(w, p)
	})
}
