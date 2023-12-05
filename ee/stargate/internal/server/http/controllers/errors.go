package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
)

const (
	errInternal         = "INTERNAL"
	errValidation       = "VALIDATION"
	errContextCancelled = "CONTEXT_CANCELLED"
)

var (
	ErrValidation   = errors.New("validation error")
	ErrNoResponders = errors.New("no responders")
)

func ResponseError(w http.ResponseWriter, r *http.Request, err error) int {
	status, code, details := coreErrorToErrorCode(err)

	baseError := errors.Cause(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if status < 500 {
		err := json.NewEncoder(w).Encode(api.ErrorResponse{
			ErrorCode:    code,
			ErrorMessage: baseError.Error(),
			Details:      details,
		})
		if err != nil {
			panic(err)
		}
	} else {
		logging.FromContext(r.Context()).Errorf("internal server error: %s", err)
	}

	return status
}

func coreErrorToErrorCode(err error) (int, string, string) {
	switch {
	case errors.Is(err, ErrValidation):
		return http.StatusBadRequest, errValidation, ""
	case errors.Is(err, context.Canceled):
		return http.StatusInternalServerError, errContextCancelled, ""
	case errors.Is(err, ErrNoResponders):
		return 524, errInternal, ""
	default:
		return http.StatusInternalServerError, errInternal, ""
	}
}
