package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/formancehq/webhooks/pkg/server/apierrors"
)

func decodeJSONBody(r *http.Request, dst interface{}, allowEmpty bool) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&dst); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &apierrors.ValidationError{Msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &apierrors.ValidationError{Msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &apierrors.ValidationError{Msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &apierrors.ValidationError{Msg: msg}

		case errors.Is(err, io.EOF):
			if allowEmpty {
				return nil
			}
			msg := "Request body must not be empty"
			return &apierrors.ValidationError{Msg: msg}

		default:
			return fmt.Errorf("json.Decoder.Decode: %w", err)
		}
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		msg := "Request body must only contain a single JSON object"
		return &apierrors.ValidationError{Msg: msg}
	}

	if r.Header.Get("Content-Type") != "application/json" {
		msg := "Content-Type header should be application/json"
		return &apierrors.ValidationError{Msg: msg}
	}

	return nil
}
