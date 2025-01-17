package client

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrStatusCodeClientError = errors.New("client error")
	ErrStatusCodeServerError = errors.New("server error")
)

func wrapError(err error, resp *http.Response) error {
	statusCode := resp.StatusCode

	if statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError {
		return fmt.Errorf("%w: %w", err, ErrStatusCodeClientError)
	}

	return fmt.Errorf("%w: %w", err, ErrStatusCodeServerError)
}
