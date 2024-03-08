package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/pkg/errors"
)

type mangopayError struct {
	StatusCode int               `json:"-"`
	Message    string            `json:"Message"`
	Type       string            `json:"Type"`
	Errors     map[string]string `json:"Errors"`
	WithRetry  bool              `json:"-"`
}

func (me *mangopayError) Error() error {
	var errorMessage string
	if len(me.Errors) > 0 {
		for _, message := range me.Errors {
			errorMessage = message
			break
		}
	}

	var err error
	if errorMessage == "" {
		err = fmt.Errorf("unexpected status code: %d", me.StatusCode)
	} else {
		err = fmt.Errorf("%d: %s", me.StatusCode, errorMessage)
	}

	if me.WithRetry {
		return checkStatusCodeError(me.StatusCode, err)
	} else {
		return errors.Wrap(task.ErrNonRetryable, err.Error())
	}
}

func unmarshalErrorWithRetry(statusCode int, body io.ReadCloser) *mangopayError {
	var ce mangopayError
	_ = json.NewDecoder(body).Decode(&ce)

	ce.StatusCode = statusCode
	ce.WithRetry = true

	return &ce
}

func unmarshalErrorWithoutRetry(statusCode int, body io.ReadCloser) *mangopayError {
	var ce mangopayError
	_ = json.NewDecoder(body).Decode(&ce)

	ce.StatusCode = statusCode
	ce.WithRetry = false

	return &ce
}

func checkStatusCodeError(statusCode int, err error) error {
	switch statusCode {
	case http.StatusTooEarly, http.StatusRequestTimeout:
		return errors.Wrap(task.ErrRetryable, err.Error())
	case http.StatusTooManyRequests:
		// Retry rate limit errors
		// TODO(polo): add rate limit handling
		return errors.Wrap(task.ErrRetryable, err.Error())
	case http.StatusInternalServerError, http.StatusBadGateway,
		http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		// Retry internal errors
		return errors.Wrap(task.ErrRetryable, err.Error())
	default:
		return errors.Wrap(task.ErrNonRetryable, err.Error())
	}
}
