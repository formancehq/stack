package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/pkg/errors"
)

type modulrError struct {
	StatusCode    int    `json:"-"`
	Field         string `json:"field"`
	Code          string `json:"code"`
	Message       string `json:"message"`
	ErrorCode     string `json:"errorCode"`
	SourceService string `json:"sourceService"`
	WithRetry     bool   `json:"-"`
}

func (me *modulrError) Error() error {
	var err error
	if me.Message == "" {
		err = fmt.Errorf("unexpected status code: %d", me.StatusCode)
	} else {
		err = fmt.Errorf("%d: %s", me.StatusCode, me.Message)
	}

	if me.WithRetry {
		return checkStatusCodeError(me.StatusCode, err)
	}

	return errors.Wrap(task.ErrNonRetryable, err.Error())
}

func unmarshalError(statusCode int, body io.ReadCloser, withRetry bool) *modulrError {
	var ces []modulrError
	_ = json.NewDecoder(body).Decode(&ces)

	if len(ces) == 0 {
		return &modulrError{
			StatusCode: statusCode,
			WithRetry:  withRetry,
		}
	}

	return &modulrError{
		StatusCode:    statusCode,
		Field:         ces[0].Field,
		Code:          ces[0].Code,
		Message:       ces[0].Message,
		ErrorCode:     ces[0].ErrorCode,
		SourceService: ces[0].SourceService,
		WithRetry:     withRetry,
	}
}

func unmarshalErrorWithRetry(statusCode int, body io.ReadCloser) *modulrError {
	return unmarshalError(statusCode, body, true)
}

func unmarshalErrorWithoutRetry(statusCode int, body io.ReadCloser) *modulrError {
	return unmarshalError(statusCode, body, false)
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
