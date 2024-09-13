package v3

import (
	"errors"
	"net/http"

	"github.com/formancehq/payments/internal/api/services"
	"github.com/formancehq/payments/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
)

const (
	ErrUniqueReference      = "CONFLICT"
	ErrNotFound             = "NOT_FOUND"
	ErrInvalidID            = "INVALID_ID"
	ErrMissingOrInvalidBody = "MISSING_OR_INVALID_BODY"
	ErrValidation           = "VALIDATION"
)

func handleServiceErrors(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, storage.ErrDuplicateKeyValue):
		api.BadRequest(w, ErrUniqueReference, err)
	case errors.Is(err, storage.ErrNotFound):
		api.NotFound(w, err)
	case errors.Is(err, storage.ErrValidation):
		api.BadRequest(w, ErrValidation, err)
	case errors.Is(err, services.ErrValidation):
		api.BadRequest(w, ErrValidation, err)
	case errors.Is(err, services.ErrNotFound):
		api.NotFound(w, err)
	default:
		api.InternalServerError(w, r, err)
	}
}
