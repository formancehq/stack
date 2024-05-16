package v1

import (
	"errors"
	"net/http"

	"github.com/formancehq/reconciliation/internal/api/v1/backend"
	"github.com/formancehq/reconciliation/internal/api/v1/service"
	storageerrors "github.com/formancehq/reconciliation/internal/storage/errors"
	"github.com/formancehq/stack/libs/go-libs/api"
	"go.uber.org/fx"
)

const (
	ErrInvalidID            = "INVALID_ID"
	ErrMissingOrInvalidBody = "MISSING_OR_INVALID_BODY"
	ErrValidation           = "VALIDATION"
)

func HTTPModule() fx.Option {
	return fx.Options(
		fx.Provide(fx.Annotate(service.NewService, fx.As(new(backend.Service)))),
		fx.Provide(backend.NewDefaultBackend),
	)
}

func handleServiceErrors(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, service.ErrValidation):
		api.BadRequest(w, ErrValidation, err)
	case errors.Is(err, service.ErrInvalidID):
		api.BadRequest(w, ErrInvalidID, err)
	case errors.Is(storageerrors.ErrInvalidQuery, err):
		api.BadRequest(w, ErrValidation, err)
	case errors.Is(storageerrors.ErrNotFound, err):
		api.NotFound(w, err)
	default:
		api.InternalServerError(w, r, err)
	}
}
