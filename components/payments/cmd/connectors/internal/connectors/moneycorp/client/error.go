package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/pkg/errors"
)

type moneycorpErrors struct {
	Errors []*moneycorpError `json:"errors"`
}

type moneycorpError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Title      string `json:"title"`
	Detail     string `json:"detail"`
	WithRetry  bool   `json:"-"`
}

func (me *moneycorpError) Error() error {
	var err error
	if me.Detail == "" {
		err = fmt.Errorf("unexpected status code: %d", me.StatusCode)
	} else {
		err = fmt.Errorf("%d: %s", me.StatusCode, me.Detail)
	}

	if me.WithRetry {
		return checkStatusCodeError(me.StatusCode, err)
	}

	return errors.Wrap(task.ErrNonRetryable, err.Error())
}

func unmarshalError(statusCode int, body io.ReadCloser, withRetry bool) *moneycorpError {
	var ces moneycorpErrors
	_ = json.NewDecoder(body).Decode(&ces)

	if len(ces.Errors) == 0 {
		return &moneycorpError{
			StatusCode: statusCode,
			WithRetry:  withRetry,
		}
	}

	return &moneycorpError{
		StatusCode: statusCode,
		Code:       ces.Errors[0].Code,
		Title:      ces.Errors[0].Title,
		Detail:     ces.Errors[0].Detail,
		WithRetry:  withRetry,
	}
}

func unmarshalErrorWithoutRetry(statusCode int, body io.ReadCloser) *moneycorpError {
	return unmarshalError(statusCode, body, false)
}

func unmarshalErrorWithRetry(statusCode int, body io.ReadCloser) *moneycorpError {
	return unmarshalError(statusCode, body, true)
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
