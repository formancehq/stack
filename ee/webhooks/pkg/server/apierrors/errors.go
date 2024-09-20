package apierrors

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/formancehq/go-libs/api"
	"github.com/formancehq/go-libs/logging"
	"github.com/pkg/errors"
)

const (
	ErrInternal         = "INTERNAL"
	ErrValidation       = "VALIDATION"
	ErrContextCancelled = "CONTEXT_CANCELLED"
	ErrNotFound         = "NOT_FOUND"
)

func ResponseError(w http.ResponseWriter, r *http.Request, err error) {
	status, code := coreErrorToErrorCode(err)
	w.WriteHeader(status)
	if status < 500 {
		err := json.NewEncoder(w).Encode(api.ErrorResponse{
			ErrorCode:    code,
			ErrorMessage: err.Error(),
		})
		if err != nil {
			panic(err)
		}
	} else {
		logging.FromContext(r.Context()).Errorf("internal server error: %s", err)
	}
}

func coreErrorToErrorCode(err error) (int, string) {
	switch {
	case IsValidationError(err):
		return http.StatusBadRequest, ErrValidation
	case IsNotFoundError(err):
		return http.StatusNotFound, ErrNotFound
	case errors.Is(err, context.Canceled):
		return http.StatusInternalServerError, ErrContextCancelled
	default:
		return http.StatusInternalServerError, ErrInternal
	}
}

type ValidationError struct {
	Msg string
}

func (v ValidationError) Error() string {
	return v.Msg
}

func (v ValidationError) Is(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

func NewValidationError(msg string) *ValidationError {
	return &ValidationError{
		Msg: msg,
	}
}

func IsValidationError(err error) bool {
	return errors.Is(err, &ValidationError{})
}

type NotFoundError struct {
	Msg string
}

func (v NotFoundError) Error() string {
	return v.Msg
}

func (v NotFoundError) Is(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{
		Msg: msg,
	}
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, &NotFoundError{})
}
