package apierrors

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/sqlstorage/sqlerrors"
	"github.com/pkg/errors"
)

const (
	ErrInternal                = "INTERNAL"
	ErrConflict                = "CONFLICT"
	ErrInsufficientFund        = "INSUFFICIENT_FUND"
	ErrValidation              = "VALIDATION"
	ErrContextCancelled        = "CONTEXT_CANCELLED"
	ErrStore                   = "STORE"
	ErrNotFound                = "NOT_FOUND"
	ErrScriptCompilationFailed = "COMPILATION_FAILED"
	ErrScriptNoScript          = "NO_SCRIPT"
	ErrScriptMetadataOverride  = "METADATA_OVERRIDE"
)

func ResponseError(w http.ResponseWriter, r *http.Request, err error) {
	status, code, details := coreErrorToErrorCode(err)

	w.WriteHeader(status)
	if status < 500 {
		err := json.NewEncoder(w).Encode(api.ErrorResponse{
			ErrorCode:    code,
			ErrorMessage: err.Error(),
			Details:      details,
		})
		if err != nil {
			panic(err)
		}
	} else {
		logging.FromContext(r.Context()).Errorf("internal server error: %s", err)
	}
}

func coreErrorToErrorCode(err error) (int, string, string) {
	switch {
	case IsConflictError(err):
		return http.StatusConflict, ErrConflict, ""
	case IsInsufficientFundError(err):
		return http.StatusBadRequest, ErrInsufficientFund, ""
	case IsValidationError(err):
		return http.StatusBadRequest, ErrValidation, ""
	case IsNotFoundError(err):
		return http.StatusNotFound, ErrNotFound, ""
	case IsScriptErrorWithCode(err, ErrScriptNoScript),
		IsScriptErrorWithCode(err, ErrInsufficientFund),
		IsScriptErrorWithCode(err, ErrScriptCompilationFailed),
		IsScriptErrorWithCode(err, ErrScriptMetadataOverride):
		scriptErr := err.(*ScriptError)
		return http.StatusBadRequest, scriptErr.Code, EncodeLink(scriptErr.Message)
	case errors.Is(err, context.Canceled):
		return http.StatusInternalServerError, ErrContextCancelled, ""
	case sqlerrors.IsError(err):
		return http.StatusServiceUnavailable, ErrStore, ""
	default:
		return http.StatusInternalServerError, ErrInternal, ""
	}
}

func EncodeLink(errStr string) string {
	if errStr == "" {
		return ""
	}

	errStr = strings.ReplaceAll(errStr, "\n", "\r\n")
	payload, err := json.Marshal(map[string]string{
		"error": errStr,
	})
	if err != nil {
		panic(err)
	}
	payloadB64 := base64.StdEncoding.EncodeToString(payload)
	return fmt.Sprintf("https://play.numscript.org/?payload=%v", payloadB64)
}

type InsufficientFundError struct {
	Asset string
}

func (e InsufficientFundError) Error() string {
	return fmt.Sprintf("balance.insufficient.%s", e.Asset)
}

func (e InsufficientFundError) Is(err error) bool {
	_, ok := err.(*InsufficientFundError)
	return ok
}

func IsInsufficientFundError(err error) bool {
	return errors.Is(err, &InsufficientFundError{})
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

type ConflictError struct{}

func (e ConflictError) Error() string {
	return "conflict error on reference"
}

func (e ConflictError) Is(err error) bool {
	_, ok := err.(*ConflictError)
	return ok
}

func NewConflictError() *ConflictError {
	return &ConflictError{}
}

func IsConflictError(err error) bool {
	return errors.Is(err, &ConflictError{})
}

const (
	ScriptErrorInsufficientFund  = "INSUFFICIENT_FUND"
	ScriptErrorCompilationFailed = "COMPILATION_FAILED"
	ScriptErrorNoScript          = "NO_SCRIPT"
	ScriptErrorMetadataOverride  = "METADATA_OVERRIDE"
)

type ScriptError struct {
	Code    string
	Message string
}

func (e ScriptError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e ScriptError) Is(err error) bool {
	eerr, ok := err.(*ScriptError)
	if !ok {
		return false
	}
	return e.Code == eerr.Code
}

func IsScriptErrorWithCode(err error, code string) bool {
	return errors.Is(err, &ScriptError{
		Code: code,
	})
}

func NewScriptError(code string, message string) *ScriptError {
	return &ScriptError{
		Code:    code,
		Message: message,
	}
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
